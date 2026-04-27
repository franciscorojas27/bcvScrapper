package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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
