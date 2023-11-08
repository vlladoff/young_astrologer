package main

import (
	"github.com/vlladoff/young_astrologer/internal/config"
	"github.com/vlladoff/young_astrologer/internal/lib/logger/sl"
	"github.com/vlladoff/young_astrologer/internal/storage/postgresql"
	"github.com/vlladoff/young_astrologer/internal/worker"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	storage, err := postgresql.New(cfg.StorageDataSource)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	defer postgresql.Close(storage)

	wrkr := worker.NewWorker(cfg, storage, log)
	wrkr.Start()
}
