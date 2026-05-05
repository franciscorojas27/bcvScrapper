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
}
