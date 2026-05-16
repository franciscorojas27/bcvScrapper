package worker

import (
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

	c.Run()
}
