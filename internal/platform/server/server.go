package server

import (
	"bcv/internal/modules/trade"
	"bcv/internal/platform/database"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/fiber/v3/middleware/compress"
)

func StartServer(app *App) {
	server := fiber.New()

	server.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	server.Use("/binance", cache.New(cache.Config{
		Expiration:  10 * time.Second,
		CacheHeader: "Cache-Control",
		KeyGenerator: func(c fiber.Ctx) string {
			return c.OriginalURL()
		},
	}))

	server.Get("/health", func(c fiber.Ctx) error {
		slog.Info("Request: /health", "ip:", c.IP())
		return c.JSON(fiber.Map{"status": "ok"})
	})

	server.Get("/binance", func(c fiber.Ctx) error {
		rates, err := trade.FetchBinanceRates()
		if err != nil {
			slog.Error("Error fetching Binance rates", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error fetching Binance rates: %v", err)})
		}

		slog.Info("Request: /binance", "ip:", c.IP())
		return c.JSON(rates)
	})

	server.Get("/rates", func(c fiber.Ctx) error {
		date, err := database.GetLatestRates(app.DB)
		slog.Info("Request: /rates", "ip:", c.IP())
		if err != nil {
			slog.Error("Error fetching latest rates", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error fetching latest rates: %v", err)})
		}
		return c.JSON(date)
	})

	server.Get("/trade-signal", func(c fiber.Ctx) error {
		tradeSignal, err := database.GetLatestTradeSignal(app.DB)
		if err != nil {
			slog.Error("Error fetching latest trade signal", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error fetching latest trade signal: %v", err)})
		}
		slog.Info("Request: /trade-signal", "ip:", c.IP())
		return c.JSON(tradeSignal)
	})

	server.Get("/rate-list", func(c fiber.Ctx) error {
		rates, err := database.GetListOfLatestReports(app.DB)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error fetching list of latest rates: %v", err)})
		}
		slog.Info("Request: /rate-list", "ip:", c.IP())
		return c.JSON(rates)
	})

	slog.Info("Starting server", "port", app.Port)
	server.Listen(":" + app.Port)
}
