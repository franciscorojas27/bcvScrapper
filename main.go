package main

import (
	"log/slog"
	"os"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type App struct {
	DB   *gorm.DB
	Auth AuthTelegram
	Port string
}

var (
	conString      = os.Getenv("DB_STRING")
	port           = os.Getenv("PORT")
	tokenTelegram  = os.Getenv("TELEGRAM_TOKEN")
	chatIDTelegram = os.Getenv("TELEGRAM_CHAT_ID")
)

func main() {
	Setup()
	Load()
	db, err := ConnectDB(conString)
	if err != nil {
		slog.Error("Error to connect to database", "error", err)
		return
	}

	if err = InitializeDB(db); err != nil {
		slog.Error("Error to initialize the database", "error", err)
		return
	}
	app := &App{DB: db, Auth: AuthTelegram{
		Token:  tokenTelegram,
		ChatID: chatIDTelegram,
	},
		Port: port,
	}

	c := cron.New()

	c.AddFunc("@daily", func() {
		slog.Info("Starting scheduled scrape")
		if err := ScrapeLatestRates(app); err != nil {
			slog.Error("Error during scheduled scrape", "error", err)
		} else {
			slog.Info("Scheduled scrape completed successfully")
		}
	})
	c.Start()

	err = ScrapeLatestRates(app)

	if err != nil {
		slog.Error("Error during initial scrape", "error", err)
	} else {
		slog.Info("Initial scrape completed successfully")
	}

	StartServer(app)
}
