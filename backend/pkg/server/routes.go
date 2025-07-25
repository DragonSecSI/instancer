package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/DragonSecSI/instancer/backend/pkg/api"
	"github.com/DragonSecSI/instancer/backend/pkg/server/middleware"
)

func (s *Server) NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger(s.Logger))

	r.Mount("/api/v1", api.Api{
		DB:        s.DB,
		Config:    s.Config.App,
		Logger:    s.Logger.With().Str("module", "api").Logger(),
		Instancer: s.Instancer,
	}.Routes())

	return r
}
