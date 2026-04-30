package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	DB *pgxpool.Pool
}

var (
	conString      = os.Getenv("DB_STRING")
	port           = os.Getenv("PORT")
	tokenTelegram  = os.Getenv("TELEGRAM_TOKEN")
	chatIDTelegram = os.Getenv("TELEGRAM_CHAT_ID")
	FinalData      CurrencyRatesData
	logger         = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	app            *App
)

func init() {
	if conString == "" {
		logger.Warn("the environment variable DB_STRING is not set.")
		panic(0)
	}
	if port == "" {
		port = "8080"
	}
	if tokenTelegram == "" || chatIDTelegram == "" {
		logger.Warn("the environment variables TELEGRAM_TOKEN o TELEGRAM_CHAT_ID not set.")
		panic(0)
	}

	pool, err := connectDB(conString)
	if err != nil {
		logger.Error("Error to connect to database", "error", err)
		return
	}
	defer pool.Close()

	if err = initializeDB(pool); err != nil {
		logger.Error("Error to initialize the database", "error", err)
		return
	}
	app = &App{DB: pool}
}

func main() {
	FinalData = scrapeBCV()
	if err := SaveScrapeReport(app.DB, FinalData); err != nil {
		logger.Error("Error to save scrape report", "error", err)
		return
	}

	// authTelegram := AuthTelegram{
	// 	Token:  tokenTelegram,
	// 	ChatID: chatIDTelegram,
	// }

	// message := BuildMessage(FinalData)

	// if err := sendMessage(authTelegram, message); err != nil {
	// 	logger.Error("Error sending message to telegram", "error", err)
	// } else {
	// 	logger.Info("Message sent to Telegram successfully")
	// }

	server := fiber.New()

	server.Get("/rates", func(c fiber.Ctx) error {
		date, err := GetLatestRates(app.DB)
		logger.Info("Request: /rates", "ip:", c.IP())
		if err != nil {
			logger.Error("Error fetching latest rates", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error fetching latest rates: %v", err)})
		}
		return c.JSON(date)
	})
	logger.Info("Starting server", "port", port)
	server.Listen(":" + port)
}
