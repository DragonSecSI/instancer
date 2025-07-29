package instance

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/DragonSecSI/instancer/backend/pkg/database/models"
	"github.com/DragonSecSI/instancer/backend/pkg/helpers"
	"github.com/DragonSecSI/instancer/backend/pkg/instancer"
	"github.com/DragonSecSI/instancer/backend/pkg/server/middleware"
)

type InstanceApi struct {
	DB        *gorm.DB
	Config    config.ConfigApp
	Logger    zerolog.Logger
	Instancer *instancer.Instancer
}

func (rs InstanceApi) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.AuthUserMiddleware(rs.DB))

	r.Get("/", rs.ListInstances)
	r.With(httpin.NewInput(InstanceRequest{})).Get("/{id}", rs.GetInstance)
	r.With(httpin.NewInput(InstanceRequest{})).Delete("/{id}", rs.DeleteInstance)
	r.With(httpin.NewInput(InstanceRequest{})).Get("/new/{id}", rs.NewInstance)

	return r
}

func (rs InstanceApi) ListInstances(w http.ResponseWriter, r *http.Request) {
	team, ok := r.Context().Value("team").(*models.Team)
	if !ok {
		rs.Logger.Error().Msg("Failed to get team from request context")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	instances, err := models.InstanceGetByTeamID(rs.DB, team.ID)
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to get instances for team")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	insts := newInstanceResponseList(instances)
	helpers.Api.Response.Json(w, &rs.Logger, insts)
}

func (rs InstanceApi) GetInstance(w http.ResponseWriter, r *http.Request) {
	team, ok := r.Context().Value("team").(*models.Team)
	if !ok {
		rs.Logger.Error().Msg("Failed to get team from request context")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	req := r.Context().Value(httpin.Input).(*InstanceRequest)

	instance, err := models.InstanceGetByID(rs.DB, req.ID)
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to get instance by ID")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if instance == nil {
		rs.Logger.Warn().Msg("Instance not found")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Not Found", http.StatusNotFound)
		return
	}

	if instance.TeamID != team.ID {
		rs.Logger.Warn().Str("team", team.Name).Uint("teamid", team.ID).Uint("instance", instance.ID).Msg("Instance does not belong to the team")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Forbidden", http.StatusForbidden)
		return
	}

	inst := newInstanceResponse(*instance)
	helpers.Api.Response.Json(w, &rs.Logger, inst)
}

func (rs InstanceApi) DeleteInstance(w http.ResponseWriter, r *http.Request) {
	team, ok := r.Context().Value("team").(*models.Team)
	if !ok {
		rs.Logger.Error().Msg("Failed to get team from request context")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	req := r.Context().Value(httpin.Input).(*InstanceRequest)

	instance, err := models.InstanceGetByID(rs.DB, req.ID)
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to get instance by ID")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if instance == nil {
		rs.Logger.Warn().Msg("Instance not found")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Not Found", http.StatusNotFound)
		return
	}

	if instance.TeamID != team.ID {
		rs.Logger.Warn().Str("team", team.Name).Uint("teamid", team.ID).Uint("instance", instance.ID).Msg("Instance does not belong to the team")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Forbidden", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs InstanceApi) NewInstance(w http.ResponseWriter, r *http.Request) {
	team, ok := r.Context().Value("team").(*models.Team)
	if !ok {
		rs.Logger.Error().Msg("Failed to get team from request context")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	instances, err := models.InstanceGetByTeamID(rs.DB, team.ID)
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to get instances for team")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	req := r.Context().Value(httpin.Input).(*InstanceRequest)

	duplicate := false
	for _, inst := range instances {
		if inst.Active && inst.ChallengeID == req.ID {
			duplicate = true
			break
		}
	}
	if duplicate {
		rs.Logger.Warn().Str("team", team.Name).Uint("teamid", team.ID).Uint("challenge", req.ID).Msg("Instance already exists for this challenge")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Instance already exists for this challenge", http.StatusConflict)
		return
	}

	challenge, err := models.ChallengeGetByID(rs.DB, req.ID)
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to get challenge by ID")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if challenge == nil {
		rs.Logger.Warn().Msg("Challenge not found")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Not Found", http.StatusNotFound)
		return
	}

	name, err := rs.Instancer.GenerateName()
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to generate instance name")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	flag := helpers.Flag.Process(challenge.Flag, challenge.FlagType)
	values := strings.Split(strings.TrimSpace(challenge.Values), "\n")
	if values[len(values)-1] == "" {
		values[len(values)-1] = "flag.flag=" + flag
	} else {
		values = append(values, "flag.flag="+flag)
	}

	err = rs.Instancer.NewInstance(instancer.InstancerConfig{
		Name:       name,
		Repository: fmt.Sprintf("%s/%s", challenge.Repository, challenge.Chart),
		Version:    challenge.ChartVersion,
		Values:     values,
	})
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to create new instance")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	instance := models.Instance{
		Name:          name,
		Flag:          flag,
		TeamID:        team.ID,
		ChallengeID:   challenge.ID,
		ChallengeType: challenge.Type,
		Active:        true,
	}

	err = models.InstanceCreate(rs.DB, &instance)
	if err != nil {
		rs.Logger.Error().Err(err).Msg("Failed to create instance in database")
		helpers.Api.Response.JsonError(w, &rs.Logger, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	helpers.Api.Response.Json(w, &rs.Logger, InstanceNewResponse{
		Name: instance.Name,
	})
}
