package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
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
}

func main() {
	FinalData = scrapeBCV()

	pool, err := connectDB(conString)
	if err != nil {
		logger.Error("Error al conectar a la base de datos", "error", err)
		return
	}
	if err = initializeDB(pool); err != nil {
		logger.Error("Error al inicializar la base de datos", "error", err)
		return
	}

	defer pool.Close()

	app := &App{DB: pool}

	FinalData.InsertRates(app)
	authTelegram := AuthTelegram{
		Token:  tokenTelegram,
		ChatID: chatIDTelegram,
	}
	var sb strings.Builder

	fmt.Fprintf(&sb, "📅 Fecha: %s\n", FinalData.Date)

	for _, rate := range FinalData.List {
		fmt.Fprintf(&sb, "✅ %s: %f\n", rate.Symbol, rate.Price)
	}

	message := sb.String()
	if err := sendMessage(authTelegram, message); err != nil {
		logger.Error("Error sending message to telegram", "error", err)
	} else {
		logger.Info("Message sent to Telegram successfully")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/rates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if date, err := time.Parse(time.RFC3339, FinalData.Date); err == nil {
			FinalData.Date = date.UTC().String()
		} else {
			logger.Error("Error al formatear la fecha para la respuesta HTTP", "error", err)
		}
		json.NewEncoder(w).Encode(FinalData)
	})

	handler := cors.Default().Handler(mux)
	fmt.Printf("Servidor en http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, handler)
}
