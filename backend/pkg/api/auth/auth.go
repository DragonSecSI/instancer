package auth

import (
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/DragonSecSI/instancer/backend/pkg/database/models"
	"github.com/DragonSecSI/instancer/backend/pkg/helpers"
	"github.com/DragonSecSI/instancer/backend/pkg/metrics"
	"github.com/DragonSecSI/instancer/backend/pkg/server/middleware"
)

type AuthApi struct {
	DB     *gorm.DB
	Config config.ConfigApp
	Logger zerolog.Logger
}

func (rs AuthApi) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.AuthAdminMiddleware(rs.Config.Initializer.AdminPassword))

	r.With(httpin.NewInput(TeamRegisterRequest{})).Post("/team/register", rs.TeamRegister)

	return r
}

func (rs AuthApi) TeamRegister(w http.ResponseWriter, r *http.Request) {
	req := r.Context().Value(httpin.Input).(*TeamRegisterRequest)
	team, err := models.TeamGetByRemoteID(rs.DB, req.Team.RemoteID)
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to get team by remote ID")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if team != nil {
		rs.Logger.Warn().Str("remote_id", req.Team.RemoteID).Msg("Team already registered")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Team already registered", http.StatusConflict)
		return
	}

	token, err := helpers.Auth.Token.GenerateToken()
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to generate token")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	team = &models.Team{
		Name:     req.Team.Name,
		RemoteID: req.Team.RemoteID,
		Token:    token,
	}
	if err := rs.DB.Create(team).Error; err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to create team")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	metrics.TeamsCreatedCounter.Inc()

	teamres := TeamRegisterResponse{
		Token: token,
	}
	helpers.Api.Response.Json(w, &rs.Logger, teamres, http.StatusCreated)
}
