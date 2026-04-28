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
		return fmt.Errorf("error creating payload: %v", err)
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Telegram response error: %s", resp.Status)
	}

	return nil
}

func BuildMessage(data CurrencyRatesData) string {
	var sb strings.Builder
	sb.Grow(30 + (len(data.List) * 40))

	sb.WriteString("📅 Date: ")
	sb.WriteString(data.Date.String())
	sb.WriteByte('\n')

	for _, rate := range data.List {
		sb.WriteString("✅ ")
		sb.WriteString(rate.Symbol)
		sb.WriteString(": ")
		sb.WriteString(rate.Price.StringFixed(2))
		sb.WriteByte('\n')
	}
	return sb.String()
}
