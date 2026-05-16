package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type JSONB []string

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan type for JSONB")
	}
	return json.Unmarshal(bytes, j)
}

type TradeSignal struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Action       string    `gorm:"type:varchar(10);not null;index" json:"action"`
	Rationale    string    `gorm:"type:text;not null" json:"rationale"`
	KeyFactors   JSONB     `gorm:"type:jsonb" json:"key_factors"`
	WinPoints    float64   `gorm:"type:numeric(10,2);not null" json:"win_points"`
	AccuracyRate float64   `gorm:"type:numeric(5,2);not null" json:"accuracy_rate"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp;index" json:"created_at"`
}
