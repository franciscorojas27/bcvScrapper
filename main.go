package main

import (
	"bcv/config"
	"bcv/internal/domain"
	pricebank "bcv/internal/modules/price_bank"
	"bcv/internal/platform/database"
	"bcv/internal/platform/scraper"
	"bcv/internal/platform/server"
	"bcv/internal/worker"
	"log/slog"
)

func main() {
	config.Setup()
	port, conString, tokenTelegram, chatIDTelegram := config.Load()

	data, err := pricebank.FetchNewsTitles()
	if err != nil {
		slog.Error("Error fetching news titles", "error", err)
	} else {
		slog.Info("Fetched news titles successfully", "count", data)
	}

	db, err := database.ConnectDB(conString)
	if err != nil {
		slog.Error("Error to connect to database", "error", err)
		return
	}

	if err = database.InitializeDB(db); err != nil {
		slog.Error("Error to initialize the database", "error", err)
		return
	}
	app := server.NewApp(db, domain.AuthTelegram{
		Token:  tokenTelegram,
		ChatID: chatIDTelegram,
	}, port)

	err = scraper.ScrapeLatestRates(app)

	if err != nil {
		slog.Error("Error during initial scrape", "error", err)
	} else {
		slog.Info("Initial scrape completed successfully")
	}

	go worker.StartCron(app)

	server.StartServer(app)
}
