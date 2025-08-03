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
		rs.Logger.Warn().
			Str("flag", req.Body.Flag).
			Str("remote_id", req.Body.RemoteID).
			Str("remote_challenge_id", req.Body.RemoteChallengeID).
			Msg("Flag not found")
		res := FlagSubmitResponse{
			Correct:        false,
			ActiveInstance: false,
			WrongTeam:      false,
		}
		helpers.Api.Response.Json(w, &rs.Logger, res)
		return
	}

	if instance.Challenge.RemoteID != req.Body.RemoteChallengeID {
		metrics.FlagsIncorrectCounter.Inc()
		rs.Logger.Warn().
			Str("flag", req.Body.Flag).
			Str("remote_id", req.Body.RemoteID).
			Str("remote_challenge_id", req.Body.RemoteChallengeID).
			Str("instance", instance.Name).
			Msg("Remote challenge ID does not match")
		res := FlagSubmitResponse{
			Correct:        false,
			ActiveInstance: instance.Active,
			WrongTeam:      instance.Team.RemoteID != req.Body.RemoteID,
		}
		helpers.Api.Response.Json(w, &rs.Logger, res)
		return
	}

	metrics.FlagsCorrectCounter.Inc()
	if instance.Team.RemoteID != req.Body.RemoteID {
		metrics.FlagsWrongCounter.Inc()
	}

	res := FlagSubmitResponse{
		Correct:        true,
		ActiveInstance: instance.Active,
		WrongTeam:      instance.Team.RemoteID != req.Body.RemoteID,
	}
	helpers.Api.Response.Json(w, &rs.Logger, res)
}
