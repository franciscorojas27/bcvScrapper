package server

import (
	"bcv/internal/domain"

	"gorm.io/gorm"
)

type App struct {
	DB   *gorm.DB
	Auth domain.AuthTelegram
	Port string
}

func NewApp(db *gorm.DB, auth domain.AuthTelegram, port string) *App {
	return &App{
		DB:   db,
		Auth: auth,
		Port: port,
	}
}
