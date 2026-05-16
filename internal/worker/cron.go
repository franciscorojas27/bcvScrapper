package worker

import (
	"bcv/internal/modules/groq"
	"bcv/internal/platform/scraper"
	"bcv/internal/platform/server"
	"log/slog"

	"github.com/robfig/cron/v3"
)

func StartCron(app *server.App) {
	c := cron.New()

	c.AddFunc("@daily", func() {
		slog.Info("Starting scheduled scrape")
		if err := scraper.ScrapeLatestRates(app); err != nil {
			slog.Error("Error during scheduled scrape", "error", err)
		} else {
			slog.Info("Scheduled scrape completed successfully")
		}
	})

	runScrape := func() {
		slog.Info("Starting scheduled trade signal generation")
		if err := groq.GetTradeSignal(app); err != nil {
			slog.Error("Error during scheduled trade signal generation", "error", err)
		} else {
			slog.Info("Scheduled trade signal generation completed successfully")
		}
	}

	c.AddFunc("0 9 * * *", runScrape)
	c.AddFunc("0 16 * * *", runScrape)

	c.Start()
}
