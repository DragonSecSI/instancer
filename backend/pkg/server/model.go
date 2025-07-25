package server

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/DragonSecSI/instancer/backend/pkg/instancer"
)

type Server struct {
	Config        *config.Config
	Logger        *zerolog.Logger
	DB            *gorm.DB
	Instancer     *instancer.Instancer
	CronScheduler gocron.Scheduler
}
