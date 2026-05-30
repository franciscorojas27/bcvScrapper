package database

import (
	"bcv/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func SaveScrapeReport(db *gorm.DB, data models.ScrapeReport) (bool, error) {
	saved := false
	err := db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.ScrapeReport{}).Where("bcv_date = ?", data.BcvDate).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return nil
		}

		var latestBinance models.BinanceRate
		errBinance := tx.Where("type_value = ?", "sell").Order("id desc").First(&latestBinance).Error

		report := models.ScrapeReport{
			BcvDate: data.BcvDate,
		}

		var ratesToInsert []models.CurrencyRate
		for _, r := range data.Rates {
			var prevRate models.CurrencyRate

			err := tx.Joins("JOIN scrape_reports ON scrape_reports.id = currency_rates.report_id").
				Where("currency_rates.symbol = ?", r.Symbol).
				Order("scrape_reports.bcv_date DESC").
				First(&prevRate).Error

			changePct := decimal.Zero
			if err == nil && !prevRate.Price.IsZero() {
				changePct = r.Price.Sub(prevRate.Price).Div(prevRate.Price).Mul(decimal.NewFromInt(100))
			}

			currentRate := models.CurrencyRate{
				Symbol:    r.Symbol,
				Price:     r.Price,
				ChangePct: changePct,
			}

			if errBinance == nil && !currentRate.Price.IsZero() {
				gapValue := latestBinance.Price.Sub(currentRate.Price)

				gapPct := gapValue.Div(currentRate.Price).Mul(decimal.NewFromInt(100))

				currentRate.Gap = &models.Gap{
					Value:           gapValue.Round(2),
					ValuePorcentual: gapPct.Round(2),
					BinanceRateID:   latestBinance.ID,
				}
			}

			ratesToInsert = append(ratesToInsert, currentRate)
		}

		report.Rates = ratesToInsert

		if err := tx.Create(&report).Error; err != nil {
			return err
		}
		saved = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return saved, nil
}

func GetLatestRates(db *gorm.DB) (models.ScrapeReport, error) {
	var report models.ScrapeReport

	err := db.Preload("Rates").
		Preload("Rates.Gap.BinanceRate").
		Order("bcv_date DESC").
		First(&report).Error

	if err != nil {
		return models.ScrapeReport{}, err
	}

	return report, nil
}

func GetListOfLatestRates(db *gorm.DB) ([]models.CurrencyRate, error) {
	var reports []models.ScrapeReport
	if err := db.Preload("Rates").Preload("Rates.Gap.BinanceRate").Order("bcv_date DESC").Limit(15).Find(&reports).Error; err != nil {
		return nil, err
	}

	var rates []models.CurrencyRate
	for _, rpt := range reports {
		for _, r := range rpt.Rates {
			rates = append(rates, r)
		}
	}

	return rates, nil
}

func GetListOfLatestReports(db *gorm.DB) ([]models.ScrapeReport, error) {
	var reports []models.ScrapeReport
	if err := db.Preload("Rates").Preload("Rates.Gap.BinanceRate").Order("bcv_date DESC").Limit(15).Find(&reports).Error; err != nil {
		return nil, err
	}

	return reports, nil
}
