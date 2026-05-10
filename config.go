package main

import (
	"os"

	"log/slog"
)

func Setup() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
func Load() {
	if conString == "" {
		slog.Warn("the environment variable DB_STRING is not set.")
		os.Exit(1)
	}
	if port == "" {
		port = "8080"
	}
	if tokenTelegram == "" || chatIDTelegram == "" {
		slog.Warn("the environment variables TELEGRAM_TOKEN o TELEGRAM_CHAT_ID not set.")
		os.Exit(1)
	}
}
