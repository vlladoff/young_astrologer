package main

import (
	"github.com/go-co-op/gocron"
	"github.com/vlladoff/young_astrologer/internal/config"
	"github.com/vlladoff/young_astrologer/internal/lib/logger/sl"
	"log/slog"
	"os"
	"time"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	s := gocron.NewScheduler(time.UTC)
	job, err := s.Every(1).Day().At("00:00").Do(fetchAPOD(log, cfg))

	log.Info("young astrologer worker started")

	if err != nil {
		log.Error("Error in job ", job.GetName(), sl.Err(err))
	}
}

func fetchAPOD(log *slog.Logger, cfg *config.Config) error {
	log.Info("fetched apod data successful")

	return nil
}
