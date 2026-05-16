package database

import (
	"bcv/models"

	"gorm.io/gorm"
)

func SaveTradeSignal(db *gorm.DB, signal *models.TradeSignal) error {
	if err := db.Create(signal).Error; err != nil {
		return err
	}
	return nil
}

func GetLatestTradeSignal(db *gorm.DB) (*models.TradeSignal, error) {
	var signal models.TradeSignal
	if err := db.Last(&signal).Error; err != nil {
		return nil, err
	}
	return &signal, nil
}
