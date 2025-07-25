package health

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/DragonSecSI/instancer/backend/pkg/helpers"
)

type HealthApi struct {
	DB     *gorm.DB
	Config config.ConfigApp
	Logger zerolog.Logger
}

func (rs HealthApi) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/ready", rs.GetReady)

	return r
}

func (rs HealthApi) GetReady(w http.ResponseWriter, r *http.Request) {
	helpers.Api.Response.Json(w, &rs.Logger, HealthReadyResponse{
		Ready: true,
	})
}
