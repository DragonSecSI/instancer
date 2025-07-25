package challenge

import (
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/DragonSecSI/instancer/backend/pkg/database/models"
	"github.com/DragonSecSI/instancer/backend/pkg/helpers"
	"github.com/DragonSecSI/instancer/backend/pkg/server/middleware"
)

type ChallengeApi struct {
	DB     *gorm.DB
	Config config.ConfigApp
	Logger zerolog.Logger
}

func (rs ChallengeApi) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.AuthAnyMiddleware(rs.DB, rs.Config.Initializer.AdminPassword))

	r.With(httpin.NewInput(ChallengeListRequest{})).Get("/", rs.ListChallenges)
	r.With(httpin.NewInput(ChallengeRequest{})).Get("/{id}", rs.GetChallenge)
	r.With(httpin.NewInput(ChallengeNewRequest{})).Post("/", rs.CreateChallenge)

	return r
}

func (rs ChallengeApi) ListChallenges(w http.ResponseWriter, r *http.Request) {
	req := r.Context().Value(httpin.Input).(*ChallengeListRequest)

	offset := (req.Page - 1) * req.Pagesize
	limit := req.Pagesize
	challenges, err := models.ChallengeGetList(rs.DB, offset, limit)
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to get challenge list")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	challs := newChallengeResponseList(challenges)
	helpers.Api.Response.Json(w, &rs.Logger, challs)
}

func (rs ChallengeApi) GetChallenge(w http.ResponseWriter, r *http.Request) {
	req := r.Context().Value(httpin.Input).(*ChallengeRequest)

	challenge, err := models.ChallengeGetByID(rs.DB, req.ID)
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to get challenge by ID")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if challenge == nil {
		helpers.Api.Response.JsonError(w, &rs.Logger, "Challenge not found", http.StatusNotFound)
		return
	}

	chall := newChallengeResponse(*challenge)
	helpers.Api.Response.Json(w, &rs.Logger, chall)
}

func (rs ChallengeApi) CreateChallenge(w http.ResponseWriter, r *http.Request) {
	if isAdmin, ok := r.Context().Value("admin").(bool); !ok || !isAdmin {
		helpers.Api.Response.JsonError(w, &rs.Logger, "Forbidden", http.StatusForbidden)
		return
	}

	req := r.Context().Value(httpin.Input).(*ChallengeNewRequest)
	challenge := req.Challenge.ToModel()

	if err := models.ChallengeCreate(rs.DB, &challenge); err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to create challenge")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	chall := newChallengeResponse(challenge)
	helpers.Api.Response.Json(w, &rs.Logger, chall, http.StatusCreated)
}
