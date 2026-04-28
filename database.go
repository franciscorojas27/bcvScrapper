package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

func connectDB(conString string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, conString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return dbpool, nil
}

func initializeDB(pool *pgxpool.Pool) error {
	ctx := context.Background()

	createReports := `CREATE TABLE IF NOT EXISTS scrape_reports (
		id SERIAL PRIMARY KEY,
		bcv_date TIMESTAMPTZ UNIQUE NOT NULL,
		fetched_at TIMESTAMPTZ DEFAULT NOW()
	);`

	if _, err := pool.Exec(ctx, createReports); err != nil {
		return fmt.Errorf("failed to create scrape_reports table: %w", err)
	}

	createCurrencyRates := `CREATE TABLE IF NOT EXISTS currency_rates (
		id SERIAL PRIMARY KEY,
		report_id INTEGER NOT NULL REFERENCES scrape_reports(id) ON DELETE CASCADE,
		symbol VARCHAR(10) NOT NULL CHECK (symbol IN ('USD', 'EUR', 'CNY', 'TRY', 'RUB')),
		price NUMERIC(20,8) NOT NULL,
		change_pct NUMERIC(10,4) DEFAULT 0.00,
		CONSTRAINT unique_report_currency UNIQUE(report_id, symbol)
	);`

	if _, err := pool.Exec(ctx, createCurrencyRates); err != nil {
		return fmt.Errorf("failed to create currency_rates table: %w", err)
	}

	if _, err := pool.Exec(ctx, "CREATE INDEX IF NOT EXISTS idx_rates_symbol_lookup ON currency_rates(symbol);"); err != nil {
		return fmt.Errorf("failed to create index idx_rates_symbol_lookup: %w", err)
	}

	return nil
}

func SaveScrapeReport(pool *pgxpool.Pool, data CurrencyRatesData) error {
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	isNewReport := `
		SELECT bcv_date FROM scrape_reports WHERE bcv_date = $1
	`
	var existingDate time.Time
	err = tx.QueryRow(ctx, isNewReport, data.Date).Scan(&existingDate)
	if err != nil && err != pgx.ErrNoRows {
		return fmt.Errorf("failed to check existing report: %w", err)
	}

	if err == nil {
		if existingDate.Equal(data.Date) {
			return nil
		}
	}

	var reportID int
	err = tx.QueryRow(ctx, "INSERT INTO scrape_reports (bcv_date) VALUES ($1) RETURNING id", data.Date).Scan(&reportID)
	if err != nil {
		return fmt.Errorf("failed to insert scrape report: %w", err)
	}

	insertRate := `INSERT INTO currency_rates (report_id, symbol, price, change_pct) VALUES ($1, $2, $3, $4)`
	prevQuery := `
		SELECT cr.price FROM currency_rates cr
		JOIN scrape_reports sr ON cr.report_id = sr.id
		WHERE cr.symbol = $1
		ORDER BY sr.bcv_date DESC LIMIT 1
	`

	for _, rate := range data.List {
		var changePct decimal.Decimal
		var prevPriceStr string
		errPrev := tx.QueryRow(ctx, prevQuery, rate.Symbol).Scan(&prevPriceStr)
		if errPrev != nil && errPrev != pgx.ErrNoRows {
			return fmt.Errorf("failed to fetch previous price for %s: %w", rate.Symbol, errPrev)
		}

		if errPrev == pgx.ErrNoRows || prevPriceStr == "" {
			changePct = decimal.Zero
		} else {
			prevPriceDec, err := decimal.NewFromString(prevPriceStr)
			if err != nil {
				return fmt.Errorf("failed to parse previous price for %s: %w", rate.Symbol, err)
			}
			currPriceDec := rate.Price
			if prevPriceDec.IsZero() {
				changePct = decimal.Zero
			} else {
				changePct = currPriceDec.Sub(prevPriceDec).Div(prevPriceDec).Mul(decimal.NewFromInt(100))
			}
		}

		if _, err := tx.Exec(ctx, insertRate, reportID, rate.Symbol, rate.Price.String(), changePct.String()); err != nil {
			return fmt.Errorf("failed to insert currency rate for %s: %w", rate.Symbol, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func GetLatestRates(pool *pgxpool.Pool) (CurrencyRatesData, error) {
	ctx := context.Background()
	query := `
		SELECT sr.bcv_date, cr.symbol, cr.price, cr.change_pct
		FROM scrape_reports sr
		JOIN currency_rates cr ON cr.report_id = sr.id
		WHERE sr.bcv_date = (SELECT MAX(bcv_date) FROM scrape_reports)
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return CurrencyRatesData{}, fmt.Errorf("failed to query latest rates: %w", err)
	}
	defer rows.Close()

	var data CurrencyRatesData
	data.List = []CurrencyRate{}
	for rows.Next() {
		var rate CurrencyRate
		var bcvDate time.Time
		var priceStr, changePctStr string
		if err := rows.Scan(&bcvDate, &rate.Symbol, &priceStr, &changePctStr); err != nil {
			return CurrencyRatesData{}, fmt.Errorf("failed to scan row: %w", err)
		}
		rate.Price, err = decimal.NewFromString(priceStr)
		if err != nil {
			return CurrencyRatesData{}, fmt.Errorf("failed to parse price for %s: %w", rate.Symbol, err)
		}
		rate.ChangePct, err = decimal.NewFromString(changePctStr)
		if err != nil {
			return CurrencyRatesData{}, fmt.Errorf("failed to parse change_pct for %s: %w", rate.Symbol, err)
		}
		data.List = append(data.List, rate)
		data.Date = bcvDate
	}
	return data, nil
}
