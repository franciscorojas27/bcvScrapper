package main

import (
	"fmt"
	"os"

	"log/slog"
)

func Setup() {
	f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		os.Exit(1)
	}
	logger := slog.New(slog.NewJSONHandler(f, nil))
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
