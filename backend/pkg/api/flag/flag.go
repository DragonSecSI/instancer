package flag

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

type FlagApi struct {
	DB     *gorm.DB
	Config config.ConfigApp
	Logger zerolog.Logger
}

func (rs FlagApi) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.AuthAdminMiddleware(rs.Config.Initializer.AdminPassword))

	r.With(httpin.NewInput(FlagSubmitRequest{})).Post("/submit", rs.FlagSubmit)

	return r
}

func (rs FlagApi) FlagSubmit(w http.ResponseWriter, r *http.Request) {
	req := r.Context().Value(httpin.Input).(*FlagSubmitRequest)

	metrics.FlagsSubmittedCounter.Inc()

	instance, err := models.InstanceGetByFlag(rs.DB, req.Body.Flag)
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to get instance by flag")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if instance == nil {
		metrics.FlagsIncorrectCounter.Inc()
		rs.Logger.Warn().Str("flag", req.Body.Flag).Msg("Flag not found")
		res := FlagSubmitResponse{
			Correct:        false,
			ActiveInstance: false,
			WrongTeam:      false,
		}
		helpers.Api.Response.Json(w, &rs.Logger, res)
		return
	}

	team, err := models.TeamGetByID(rs.DB, instance.TeamID)
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to get team by ID")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if team == nil {
		rs.Logger.Error().Uint("team_id", instance.TeamID).Msg("Team not found for instance")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Team not found", http.StatusInternalServerError)
		return
	}

	metrics.FlagsCorrectCounter.Inc()
	if team.RemoteID != req.Body.RemoteID {
		metrics.FlagsWrongCounter.Inc()
	}

	res := FlagSubmitResponse{
		Correct:        true,
		ActiveInstance: instance.Active,
		WrongTeam:      team.RemoteID != req.Body.RemoteID,
	}
	helpers.Api.Response.Json(w, &rs.Logger, res)
}
