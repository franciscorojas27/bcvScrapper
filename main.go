package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

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
		logger.Error("Error to connect to database", "error", err)
		return
	}
	if err = initializeDB(pool); err != nil {
		logger.Error("Error to initialize the database", "error", err)
		return
	}
	if err = SaveScrapeReport(pool, FinalData); err != nil {
		logger.Error("Error to save scrape report", "error", err)
		return
	}

	defer pool.Close()

	// app := &App{DB: pool}

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

	mux := http.NewServeMux()
	mux.HandleFunc("/rates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		date, err := GetLatestRates(pool)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching latest rates: %v", err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(date)
	})

	handler := cors.Default().Handler(mux)
	fmt.Printf("Server http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, handler)
}
