package telegram

import (
	"testing"
	"time"

	"bcv/models"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestBuildMessage(t *testing.T) {
	ts := time.Date(2026, time.May, 16, 10, 30, 0, 0, time.FixedZone("UTC-4", -4*60*60))
	message := BuildMessage(models.ScrapeReport{
		BcvDate: ts,
		Rates: []models.CurrencyRate{
			{Symbol: "USD", Price: decimal.RequireFromString("3.5")},
			{Symbol: "EUR", Price: decimal.RequireFromString("4.25")},
		},
	})

	assert.Contains(t, message, ts.Local().String())
	assert.Contains(t, message, "✅ USD: 3.50")
	assert.Contains(t, message, "✅ EUR: 4.25")
}
