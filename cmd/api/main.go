package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vlladoff/young_astrologer/internal/config"
	"github.com/vlladoff/young_astrologer/internal/http-server/handlers/astro"
	"github.com/vlladoff/young_astrologer/internal/lib/logger/sl"
	"github.com/vlladoff/young_astrologer/internal/storage/postgresql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	log.Info("starting young astrologer api")

	storage, err := postgresql.New(cfg.StorageDataSource)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	defer postgresql.Close(storage)

	router := initRouter(log, storage)

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}

func initRouter(log *slog.Logger, storage *postgresql.Storage) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/api/get_all", func(r chi.Router) {
		r.Post("/", astro.GetAll(log, storage))
	})

	router.Route("/api/get", func(r chi.Router) {
		r.Post("/", astro.GetByDay(log, storage))
	})

	return router
}
