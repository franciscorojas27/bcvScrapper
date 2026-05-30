package database

import (
	"bcv/models"
	"fmt"
	"net/url"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func BuildConnString() string {
	u := &url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
		),
		Host: fmt.Sprintf(
			"%s:%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
		),
	}

	u = u.JoinPath(os.Getenv("DB_NAME"))

	return u.String()
}

func ConnectDB() (*gorm.DB, error) {
	conString := BuildConnString()
	db, err := gorm.Open(postgres.Open(conString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return db, nil
}

func InitializeDB(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.ScrapeReport{},
		&models.CurrencyRate{},
		&models.TradeSignal{},
		&models.Gap{},
		&models.BinanceRate{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	return nil
}
