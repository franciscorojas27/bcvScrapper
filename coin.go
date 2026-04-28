package main

import (
	"time"

	"github.com/shopspring/decimal"
)

type CurrencyRatesData struct {
	List []CurrencyRate `json:"list"`
	Date time.Time      `json:"date"`
}
type CurrencyRate struct {
	Symbol    string          `json:"symbol"`
	Price     decimal.Decimal `json:"price"`
	ChangePct decimal.Decimal `json:"change_pct"`
}
