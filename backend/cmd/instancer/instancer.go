package main

import (
	"github.com/DragonSecSI/instancer/backend/pkg/server"
)

func main() {
	srv, err := server.NewServer()
	if err != nil {
		panic(err)
	}

	srv.InitCron()
	srv.CronScheduler.Start()
	defer func() {
		if err := srv.CronScheduler.Shutdown(); err != nil {
			srv.Logger.Error().Err(err).Msg("Failed to stop cron scheduler")
		}
	}()

	srv.Start()
}
