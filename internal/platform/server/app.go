package server

import (
	"bcv/internal/domain"

	"github.com/tmc/langchaingo/llms/openai"
	"gorm.io/gorm"
)

type App struct {
	DB   *gorm.DB
	Auth domain.AuthTelegram
	Port string
	IA   *openai.LLM
}

func NewApp(db *gorm.DB, auth domain.AuthTelegram, port string, ia *openai.LLM) *App {
	return &App{
		DB:   db,
		Auth: auth,
		Port: port,
		IA:   ia,
	}
}
