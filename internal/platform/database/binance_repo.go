package database

import (
	"bcv/models"

	"gorm.io/gorm"
)

func SaveBinanceRates(db *gorm.DB, rates []models.BinanceRate) error {
	for i := range rates {
		var lastRate models.BinanceRate
		err := db.Where("type_value = ?", rates[i].TypeValue).Order("id desc").First(&lastRate).Error

		if err == nil {
			if lastRate.Price.Equal(rates[i].Price) {
				rates[i].ID = lastRate.ID
				rates[i].CreatedAt = lastRate.CreatedAt
				continue
			}
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := db.Create(&rates[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetLatestBinanceRates(db *gorm.DB) ([]models.BinanceRate, error) {
	var buyRate models.BinanceRate
	var sellRate models.BinanceRate

	err := db.Where("type_value = ?", "buy").Order("id desc").First(&buyRate).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []models.BinanceRate{}, nil
		}
		return nil, err
	}

	err = db.Where("type_value = ?", "sell").Order("id desc").First(&sellRate).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []models.BinanceRate{buyRate}, nil
		}
		return nil, err
	}

	return []models.BinanceRate{buyRate, sellRate}, nil
}

type BinanceRateList struct {
	SellList []models.BinanceRate `json:"sell_list"`
	BuyList  []models.BinanceRate `json:"buy_list"`
}

func GetlistOfLatestBinanceRates(db *gorm.DB) (BinanceRateList, error) {
	var buyRates []models.BinanceRate
	var sellRates []models.BinanceRate
	err := db.Where("type_value = ?", "buy").Order("id desc").Limit(15).Find(&buyRates).Error
	if err != nil {
		return BinanceRateList{}, err
	}
	err = db.Where("type_value = ?", "sell").Order("id desc").Limit(15).Find(&sellRates).Error
	if err != nil {
		return BinanceRateList{}, err
	}
	return BinanceRateList{
		SellList: sellRates,
		BuyList:  buyRates,
	}, nil
}
