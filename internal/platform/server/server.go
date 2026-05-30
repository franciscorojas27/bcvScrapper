package server

import (
	"bcv/internal/domain"
	"bcv/internal/modules/news"
	"bcv/internal/platform/database"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func StartServer(app *App) {
	app.Hub = domain.NewHub()
	go app.Hub.Run()

	server := fiber.New()

	server.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	server.Use(cors.New(cors.Config{
		AllowOrigins:     app.CORS.AllowOrigins,
		AllowMethods:     app.CORS.AllowMethods,
		AllowHeaders:     app.CORS.AllowHeaders,
		ExposeHeaders:    app.CORS.ExposeHeaders,
		AllowCredentials: app.CORS.AllowCredentials,
		MaxAge:           app.CORS.MaxAge,
	}))

	server.Use("/ws", func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	api := server.Group("/api")

	api.Get("/ws", websocket.New(func(c *websocket.Conn) {
		app.Hub.Mu.Lock()
		app.Hub.Clients[c] = true
		app.Hub.Mu.Unlock()

		slog.Info("WebSocket: Cliente conectado", "ip", c.RemoteAddr())

		defer func() {
			app.Hub.Mu.Lock()
			delete(app.Hub.Clients, c)
			app.Hub.Mu.Unlock()
			c.Close()
			slog.Info("WebSocket: Cliente desconectado")
		}()

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))

	api.Get("/binance", func(c fiber.Ctx) error {
		rates, err := database.GetLatestBinanceRates(app.DB)
		if err != nil {
			slog.Error("Error fetching latest binance rates from DB", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(rates)
	})

	api.Get("/binance-list", func(c fiber.Ctx) error {
		rates, err := database.GetlistOfLatestBinanceRates(app.DB)
		if err != nil {
			slog.Error("Error fetching list of latest binance rates from DB", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(rates)
	})

	api.Get("/health", func(c fiber.Ctx) error {
		slog.Info("Request: /health", "ip:", c.IP())
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api.Get("/rates", func(c fiber.Ctx) error {
		date, err := database.GetLatestRates(app.DB)
		slog.Info("Request: /rates", "ip:", c.IP())
		if err != nil {
			slog.Error("Error fetching latest rates", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error fetching latest rates: %v", err)})
		}
		return c.JSON(date)
	})

	api.Get("/trade-signal", func(c fiber.Ctx) error {
		tradeSignal, err := database.GetLatestTradeSignal(app.DB)
		if err != nil {
			slog.Error("Error fetching latest trade signal", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error fetching latest trade signal: %v", err)})
		}
		slog.Info("Request: /trade-signal", "ip:", c.IP())
		return c.JSON(tradeSignal)
	})

	api.Get("/news", cache.New(cache.Config{
		Expiration: 4 * time.Hour,
	}), func(c fiber.Ctx) error {
		news, err := news.FetchNewsTitles()
		if err != nil {
			slog.Error("Error fetching latest news", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error fetching latest news: %v", err)})
		}
		slog.Info("Request: /news", "ip:", c.IP())
		return c.JSON(news)
	})

	api.Get("/rate-list", func(c fiber.Ctx) error {
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
