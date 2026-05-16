package database

import (
	"bcv/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(conString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(conString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return db, nil
}

func InitializeDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.ScrapeReport{}, &models.CurrencyRate{}, &models.TradeSignal{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	return nil
}
