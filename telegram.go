package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AuthTelegram struct {
	Token  string `json:"token"`
	ChatID string `json:"chat_id"`
}

func sendMessage(auth AuthTelegram, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", auth.Token)
	payload := map[string]string{
		"chat_id": auth.ChatID,
		"text":    message,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al crear el payload: %v", err)
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return fmt.Errorf("error al enviar el mensaje: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta de Telegram: %s", resp.Status)
	}

	return nil
}
