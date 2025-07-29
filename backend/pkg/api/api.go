package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/DragonSecSI/instancer/backend/pkg/api/auth"
	"github.com/DragonSecSI/instancer/backend/pkg/api/challenge"
	"github.com/DragonSecSI/instancer/backend/pkg/api/flag"
	"github.com/DragonSecSI/instancer/backend/pkg/api/health"
	"github.com/DragonSecSI/instancer/backend/pkg/api/instance"
	"github.com/DragonSecSI/instancer/backend/pkg/api/meta"
	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/DragonSecSI/instancer/backend/pkg/instancer"
)

type Api struct {
	DB        *gorm.DB
	Config    config.ConfigApp
	Logger    zerolog.Logger
	Instancer *instancer.Instancer
}

func (rs Api) Routes() chi.Router {
	r := chi.NewRouter()

	r.Mount("/health", health.HealthApi{
		DB:     rs.DB,
		Config: rs.Config,
		Logger: rs.Logger.With().Str("module", "api/health").Logger(),
	}.Routes())

	r.Mount("/auth", auth.AuthApi{
		DB:     rs.DB,
		Config: rs.Config,
		Logger: rs.Logger.With().Str("module", "api/auth").Logger(),
	}.Routes())

	r.Mount("/instance", instance.InstanceApi{
		DB:        rs.DB,
		Config:    rs.Config,
		Logger:    rs.Logger.With().Str("module", "api/instance").Logger(),
		Instancer: rs.Instancer,
	}.Routes())

	r.Mount("/challenge", challenge.ChallengeApi{
		DB:     rs.DB,
		Config: rs.Config,
		Logger: rs.Logger.With().Str("module", "api/challenge").Logger(),
	}.Routes())

	r.Mount("/flag", flag.FlagApi{
		DB:     rs.DB,
		Config: rs.Config,
		Logger: rs.Logger.With().Str("module", "api/flag").Logger(),
	}.Routes())

	r.Mount("/meta", meta.MetaApi{
		DB:     rs.DB,
		Config: rs.Config,
		Logger: rs.Logger.With().Str("module", "api/meta").Logger(),
	}.Routes())

	return r
}
