package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/slavik22/urlShortener/internal/config"
	remove "github.com/slavik22/urlShortener/internal/http-server/handlers/url/delete"
	"github.com/slavik22/urlShortener/internal/http-server/handlers/url/redirect"
	"github.com/slavik22/urlShortener/internal/http-server/handlers/url/save"
	"github.com/slavik22/urlShortener/internal/http-server/middleware/logger"
	"github.com/slavik22/urlShortener/internal/lib/logger/sl"
	"github.com/slavik22/urlShortener/internal/storage/sqlite"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting logging info")
	log.Debug("Starting logging debug")

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{cfg.HTTPServer.User: cfg.HTTPServer.Password}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", remove.New(log, storage))
	})

	router.Get("/url/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
