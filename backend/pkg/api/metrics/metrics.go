package metrics

import (
	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/DragonSecSI/instancer/backend/pkg/server/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

type MetricsApi struct {
	Logger zerolog.Logger
	Config config.ConfigApp
}

func (rs MetricsApi) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.AuthAdminMiddleware(rs.Config.Initializer.AdminPassword))

	r.Handle("/", promhttp.Handler())

	return r
}
