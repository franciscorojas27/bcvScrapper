package database

import (
	"bcv/models"

	"gorm.io/gorm"
)

func SaveGap(db *gorm.DB, gap *models.Gap) error {
	return db.Create(gap).Error
}

func GetLatestGap(db *gorm.DB) (*models.Gap, error) {
	var gap models.Gap
	if err := db.Order("id desc").First(&gap).Error; err != nil {
		return nil, err
	}
	return &gap, nil
}
