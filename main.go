package main

import (
	"bcv/internal/config"
	"bcv/internal/domain"
	"bcv/internal/modules/groq"
	"bcv/internal/platform/database"
	"bcv/internal/platform/providers/ia"
	"bcv/internal/platform/scraper"
	"bcv/internal/platform/server"
	"bcv/internal/worker"
	"log/slog"
)

func main() {
	config.Setup()
	port, conString, tokenTelegram, chatIDTelegram := config.Load()

	db, err := database.ConnectDB(conString)
	if err != nil {
		slog.Error("Error to connect to database", "error", err)
		return
	}

	if err = database.InitializeDB(db); err != nil {
		slog.Error("Error to initialize the database", "error", err)
		return
	}

	initIa, err := ia.ClientIA()
	if err != nil {
		slog.Error("Error initializing IA client", "error", err)
	}

	app := server.NewApp(db, domain.AuthTelegram{
		Token:  tokenTelegram,
		ChatID: chatIDTelegram,
	}, port, initIa)

	err = groq.GetTradeSignal(app)
	if err != nil {
		slog.Error("Error getting trade signal", "error", err)
	} else {
		slog.Info("Trade signal generated and saved successfully")
	}

	err = scraper.ScrapeLatestRates(app)

	if err != nil {
		slog.Error("Error during initial scrape", "error", err)
	} else {
		slog.Info("Initial scrape completed successfully")
	}

	worker.StartCron(app)

	server.StartServer(app)
}
