package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
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

	server.Get("/binance", func(c fiber.Ctx) error {
		var sellPrice, buyPrice decimal.Decimal
		g, _ := errgroup.WithContext(c.Context())

		g.Go(func() error {
			var err error
			sellPrice, err = GetBinanceRates("BUY", 50000)
			return err
		})

		g.Go(func() error {
			var err error
			buyPrice, err = GetBinanceRates("SELL", 50000)
			return err
		})

		if err := g.Wait(); err != nil {
			slog.Error("Error fetching Binance rates", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error fetching Binance rates: %v", err)})
		}

		rates := map[string]decimal.Decimal{
			"sellPrice": sellPrice,
			"buyPrice":  buyPrice,
		}

		slog.Info("Request: /binance", "ip:", c.IP())
		return c.JSON(rates)
	})

	server.Get("/rates", func(c fiber.Ctx) error {
		date, err := GetLatestRates(app.DB)
		slog.Info("Request: /rates", "ip:", c.IP())
		if err != nil {
			slog.Error("Error fetching latest rates", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error fetching latest rates: %v", err)})
		}
		return c.JSON(date)
	})

	slog.Info("Starting server", "port", app.Port)
	server.Listen(":" + app.Port)
}
