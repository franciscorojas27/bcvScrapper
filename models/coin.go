package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ScrapeReport struct {
	gorm.Model
	BcvDate   time.Time      `gorm:"uniqueIndex;not null" json:"bcv_date"`
	FetchedAt time.Time      `gorm:"autoCreateTime" json:"fetched_at"`
	Rates     []CurrencyRate `gorm:"foreignKey:ReportID;constraint:OnDelete:CASCADE;" json:"list"`
}

type CurrencyRate struct {
	gorm.Model
	ReportID  uint            `gorm:"uniqueIndex:idx_report_symbol" json:"report_id"`
	Symbol    string          `gorm:"type:varchar(10);not null;index:idx_rates_symbol_lookup;uniqueIndex:idx_report_symbol" json:"symbol"`
	Price     decimal.Decimal `gorm:"type:numeric(20,8);not null" json:"price"`
	ChangePct decimal.Decimal `gorm:"type:numeric(10,4);default:0.00" json:"change_pct"`
	Gap       *Gap            `gorm:"foreignKey:CurrencyRateID;constraint:OnDelete:CASCADE;" json:"gap,omitempty"`
}

type Gap struct {
	gorm.Model
	CurrencyRateID  uint            `gorm:"uniqueIndex;not null" json:"currency_rate_id"`
	Value           decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"value"`
	ValuePorcentual decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"value_porcentual"`
	BinanceRateID   uint            `gorm:"not null" json:"binance_rate_id"`
	BinanceRate     BinanceRate     `gorm:"foreignKey:BinanceRateID" json:"binance_rate"`
}

type BinanceRate struct {
	gorm.Model
	Price     decimal.Decimal `gorm:"type:numeric(20,8);not null" json:"price"`
	TypeValue string          `gorm:"type:varchar(20);not null" json:"type_value"`
}
