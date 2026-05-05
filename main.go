package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type App struct {
	DB   *gorm.DB
	Auth AuthTelegram
}

var (
	conString      = os.Getenv("DB_STRING")
	port           = os.Getenv("PORT")
	tokenTelegram  = os.Getenv("TELEGRAM_TOKEN")
	chatIDTelegram = os.Getenv("TELEGRAM_CHAT_ID")
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
func ScrapeLatestRates(app *App) error {
	data := scrapeBCV()
	if err := SaveScrapeReport(app.DB, data); err != nil {
		return fmt.Errorf("Error to save scrape report: %w", err)
	}
	message := BuildMessage(data)

	if err := sendMessage(app.Auth, message); err != nil {
		return fmt.Errorf("Error sending message to telegram: %w ", err)
	}
	return nil
}

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
	}}

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

	server := fiber.New()

	server.Get("/rates", func(c fiber.Ctx) error {
		date, err := GetLatestRates(app.DB)
		slog.Info("Request: /rates", "ip:", c.IP())
		if err != nil {
			slog.Error("Error fetching latest rates", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error fetching latest rates: %v", err)})
		}
		return c.JSON(date)
	})
	slog.Info("Starting server", "port", port)
	server.Listen(":" + port)
}
