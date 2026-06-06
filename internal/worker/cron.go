package worker

import (
	"bcv/internal/modules/trade"
	"bcv/internal/platform/database"
	"bcv/internal/platform/scraper"
	"bcv/internal/platform/server"
	"bcv/models"
	"log/slog"

	"github.com/robfig/cron/v3"
)

func StartCron(app *server.App) {
	c := cron.New(cron.WithSeconds())

	c.AddFunc("*/15 * * * * *", func() {
		rates, err := trade.FetchBinanceRates()
		if err != nil {
			slog.Error("Cron Error: Failed to fetch rates from Binance API", "error", err)
			return
		}

		dbRates := []models.BinanceRate{
			{Price: rates.BuyPrice, TypeValue: "buy"},
			{Price: rates.SellPrice, TypeValue: "sell"},
		}

		if err := database.SaveBinanceRates(app.DB, dbRates); err != nil {
			slog.Error("Cron Error: Failed to persist Binance rates to database", "error", err)
			return
		}

		if app.Hub != nil {
			app.Hub.Broadcast <- dbRates
		}
	})

	c.AddFunc("0 0 16-22 * * *", func() {
		slog.Info("Starting scheduled scrape")
		if err := scraper.ScrapeLatestRates(app.DB, app.Auth); err != nil {
			slog.Error("Error during scheduled scrape", "error", err)
		} else {
			slog.Info("Scheduled scrape completed successfully")
		}
	})

	c.Start()
}
