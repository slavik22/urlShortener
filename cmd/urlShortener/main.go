package main

import (
	"github.com/slavik22/urlShortener/internal/config"
	"github.com/slavik22/urlShortener/internal/lib/logger/sl"
	"github.com/slavik22/urlShortener/internal/storage/sqlite"
	"golang.org/x/exp/slog"
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

	_, _ = storage.SaveURL("www.google.com", "google")

	_, err = storage.SaveURL("www.google.com", "google")
	if err != nil {
		log.Error("Failed to save url")
		os.Exit(1)
	}
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
