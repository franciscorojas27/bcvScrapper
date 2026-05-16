package domain

type AuthTelegram struct {
	Token  string `json:"token"`
	ChatID string `json:"chat_id"`
}
