package server

import (
	"net/http"
	"strconv"
	"time"

	httpin_integration "github.com/ggicci/httpin/integration"
	"github.com/go-chi/chi/v5"

	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/DragonSecSI/instancer/backend/pkg/database"
	"github.com/DragonSecSI/instancer/backend/pkg/instancer"
	"github.com/DragonSecSI/instancer/backend/pkg/logs"
)

func NewServer() (*Server, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	logger, err := logs.NewLogger(&config.Logs)
	if err != nil {
		return nil, err
	}

	db, err := database.NewDatabase(&config.Database)
	if err != nil {
		return nil, err
	}

	instancer := instancer.Instancer{
		Logger:           logger.With().Str("component", "instancer").Logger(),
		HelmConfig:       config.App.Helm,
		KubernetesConfig: config.App.Kubernetes,

		DB:              db,
		Prefix:          "inst-",
		CleanupDuration: time.Minute * 5,
	}

	return &Server{
		Config:    config,
		Logger:    logger,
		DB:        db,
		Instancer: &instancer,
	}, nil
}

func (s *Server) Start() {
	httpin_integration.UseGochiURLParam("path", chi.URLParam)

	router := s.NewRouter()

	bindAddr := s.Config.Server.Address + ":" + strconv.Itoa(s.Config.Server.Port)
	s.Logger.Info().Msgf("Starting Instancer on %s", bindAddr)

	err := http.ListenAndServe(bindAddr, router)
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
