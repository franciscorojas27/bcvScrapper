package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CurrencyRatesData struct {
	List []CurrencyRate `json:"list"`
	Date string         `json:"date"`
}
type CurrencyRate struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	ChangePct float64 `json:"change_pct"`
}

func normalizeToUTCDate(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func (c CurrencyRatesData) InsertRates(app *App) error {

	if validate, err := c.ValidateNewData(app.DB); err != nil {
		return err
	} else if !validate {
		return nil
	}

	ctx := context.Background()

	dateParsed, err := time.Parse(time.RFC3339, c.Date)
	if err != nil {
		return err
	}
	dateToInsert := normalizeToUTCDate(dateParsed)
	prevDate := dateToInsert.AddDate(0, 0, -1)

	tx, err := app.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var reportID int
	err = tx.QueryRow(ctx, "INSERT INTO scrape_reports (bcv_date) VALUES ($1) RETURNING id", dateToInsert).Scan(&reportID)
	if err != nil {
		return err
	}

	for _, rate := range c.List {
		var prevPrice float64
		err := tx.QueryRow(ctx, "SELECT cr.price FROM currency_rates cr JOIN scrape_reports sr ON cr.report_id = sr.id WHERE sr.bcv_date = $1 AND cr.symbol = $2", prevDate, rate.Symbol).Scan(&prevPrice)
		if err != nil {
			if err == pgx.ErrNoRows {
				prevPrice = 0
			} else {
				return err
			}
		}

		var changePct float64
		if prevPrice == 0 {
			changePct = 0
		} else {
			changePct = ((rate.Price - prevPrice) / prevPrice) * 100
		}

		_, err = tx.Exec(ctx, "INSERT INTO currency_rates (report_id, symbol, price, change_pct) VALUES ($1, $2, $3, $4)", reportID, rate.Symbol, rate.Price, changePct)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (c CurrencyRatesData) ValidateNewData(db *pgxpool.Pool) (bool, error) {
	ctx := context.Background()

	dateParsed, err := time.Parse(time.RFC3339, c.Date)
	if err != nil {
		return false, err
	}
	dateToCheck := normalizeToUTCDate(dateParsed)

	var found time.Time
	err = db.QueryRow(ctx, "SELECT bcv_date FROM scrape_reports WHERE bcv_date = $1", dateToCheck).Scan(&found)
	if err != nil {
		if err == pgx.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	// If we get a row back, it means this date already exists in DB
	return false, nil
}

func GetLatestRates(app *App) (CurrencyRatesData, error) {
	ctx := context.Background()

	var data CurrencyRatesData

	var latestDate time.Time
	err := app.DB.QueryRow(ctx, "SELECT bcv_date FROM scrape_reports ORDER BY bcv_date DESC LIMIT 1").Scan(&latestDate)
	if err != nil {
		return data, err
	}

	// Ensure we treat the stored date as UTC date (midnight) and then convert
	// to the requested UTC-4 zone for returning.
	latestDate = normalizeToUTCDate(latestDate)
	loc := time.FixedZone("UTC-4", -4*3600)
	data.Date = latestDate.In(loc).Format(time.RFC3339)

	rows, err := app.DB.Query(ctx, "SELECT cr.symbol, cr.price, cr.change_pct FROM currency_rates cr JOIN scrape_reports sr ON cr.report_id = sr.id WHERE sr.bcv_date = $1", latestDate)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var rate CurrencyRate
		if err := rows.Scan(&rate.Symbol, &rate.Price, &rate.ChangePct); err != nil {
			return data, err
		}
		data.List = append(data.List, rate)
	}

	if err := rows.Err(); err != nil {
		return data, err
	}

	return data, nil
}
