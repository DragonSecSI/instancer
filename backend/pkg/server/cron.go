package server

import (
	"github.com/go-co-op/gocron/v2"
)

func (s *Server) InitCron() {
	sch, err := gocron.NewScheduler()
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("Failed to create cron scheduler")
	}

	_, err = sch.NewJob(
		gocron.CronJob("0 * * * * *", true),
		gocron.NewTask(
			func() {
				err := s.Instancer.CronRun()
				if err != nil {
					s.Logger.Error().Err(err).Msg("Cron job failed")
				}
			},
		),
	)
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("Failed to create cron job")
	}

	s.CronScheduler = sch
}
