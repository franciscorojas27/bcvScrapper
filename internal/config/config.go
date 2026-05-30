package config

import (
	"bcv/internal/domain"
	"os"
	"strconv"
	"strings"

	"log/slog"
)

var (
	port                 = os.Getenv("PORT")
	tokenTelegram        = os.Getenv("TELEGRAM_TOKEN")
	chatIDTelegram       = os.Getenv("TELEGRAM_CHAT_ID")
	corsAllowOrigins     = os.Getenv("CORS_ALLOW_ORIGINS")
	corsAllowMethods     = os.Getenv("CORS_ALLOW_METHODS")
	corsAllowHeaders     = os.Getenv("CORS_ALLOW_HEADERS")
	corsExposeHeaders    = os.Getenv("CORS_EXPOSE_HEADERS")
	corsAllowCredentials = os.Getenv("CORS_ALLOW_CREDENTIALS")
	corsMaxAge           = os.Getenv("CORS_MAX_AGE")
)

func Setup() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
func Load() (string, string, string) {
	if port == "" {
		port = "8080"
	}
	if tokenTelegram == "" || chatIDTelegram == "" {
		slog.Warn("the environment variables TELEGRAM_TOKEN o TELEGRAM_CHAT_ID not set.")
		os.Exit(1)
	}
	return port, tokenTelegram, chatIDTelegram
}

func LoadCORS() domain.CORSConfig {
	cfg := domain.CORSConfig{
		AllowOrigins:     splitList(corsAllowOrigins),
		AllowMethods:     splitList(corsAllowMethods),
		AllowHeaders:     splitList(corsAllowHeaders),
		ExposeHeaders:    splitList(corsExposeHeaders),
		AllowCredentials: false,
		MaxAge:           0,
	}

	if len(cfg.AllowOrigins) == 0 {
		cfg.AllowOrigins = []string{"*"}
	}
	if len(cfg.AllowMethods) == 0 {
		cfg.AllowMethods = []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"}
	}
	if corsAllowCredentials != "" {
		if val, err := strconv.ParseBool(corsAllowCredentials); err == nil {
			cfg.AllowCredentials = val
		} else {
			slog.Warn("invalid CORS_ALLOW_CREDENTIALS value", "error", err)
		}
	}
	if corsMaxAge != "" {
		if val, err := strconv.Atoi(corsMaxAge); err == nil {
			cfg.MaxAge = val
		} else {
			slog.Warn("invalid CORS_MAX_AGE value", "error", err)
		}
	}

	if cfg.AllowCredentials && containsWildcard(cfg.AllowOrigins) {
		if len(cfg.AllowOrigins) > 1 {
			cfg.AllowOrigins = removeWildcard(cfg.AllowOrigins)
			slog.Warn("CORS_ALLOW_CREDENTIALS true; wildcard removed from CORS_ALLOW_ORIGINS")
		} else {
			cfg.AllowCredentials = false
			slog.Warn("CORS_ALLOW_CREDENTIALS true with wildcard origins; disabling credentials")
		}
	}

	return cfg
}

func splitList(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}

func containsWildcard(values []string) bool {
	for _, value := range values {
		if value == "*" {
			return true
		}
	}
	return false
}

func removeWildcard(values []string) []string {
	items := make([]string, 0, len(values))
	for _, value := range values {
		if value != "*" {
			items = append(items, value)
		}
	}
	return items
}
