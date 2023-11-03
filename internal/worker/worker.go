package worker

import (
	"github.com/go-co-op/gocron"
	"github.com/vlladoff/young_astrologer/internal/config"
	"github.com/vlladoff/young_astrologer/internal/lib/logger/sl"
	"github.com/vlladoff/young_astrologer/internal/storage/postgresql"
	"log/slog"
	"time"
)

type Worker struct {
	cfg     *config.Config
	storage *postgresql.Storage
	log     *slog.Logger
}

func NewWorker(cfg *config.Config, storage *postgresql.Storage, log *slog.Logger) *Worker {
	return &Worker{
		cfg:     cfg,
		storage: storage,
		log:     log,
	}
}

func (w *Worker) Start() {
	s := gocron.NewScheduler(time.UTC)
	job, err := s.Every(1).Day().Do(w.fetchAPOD)

	w.log.Info("young astrologer worker started")

	if err != nil {
		w.log.Error("Error in job", job.GetName(), sl.Err(err))
	}

	job.RegisterEventListeners(
		gocron.BeforeJobRuns(func(jobName string) {
			w.log.Info("job started", jobName)
		}),
		gocron.WhenJobReturnsError(func(jobName string, err error) {
			w.log.Error("Error in job", jobName, sl.Err(err))
		}),
		gocron.WhenJobReturnsNoError(func(jobName string) {
			w.log.Info("job done", jobName)
		}),
	)

	s.StartBlocking()
}
