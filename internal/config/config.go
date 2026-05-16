package config

import (
	"os"

	"log/slog"
)

var (
	conString      = os.Getenv("DB_STRING")
	port           = os.Getenv("PORT")
	tokenTelegram  = os.Getenv("TELEGRAM_TOKEN")
	chatIDTelegram = os.Getenv("TELEGRAM_CHAT_ID")
)

func Setup() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
func Load() (string, string, string, string) {
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
	return port, conString, tokenTelegram, chatIDTelegram
}
