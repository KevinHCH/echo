package telegram

import (
	"bytes"
	utils "echo/internal"
	"encoding/json"
	"fmt"
	"net/http"
)

type TelegramBot struct {
	Token   string
	ChatID  string
	BaseURL string
}

func NewTelegramBot() (*TelegramBot, error) {

	botToken, err := utils.GetEnv("TELEGRAM_BOT_TOKEN")
	if err != nil {
		return nil, err
	}

	chatID, err := utils.GetEnv("TELEGRAM_CHAT_ID")
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		Token:   botToken,
		ChatID:  chatID,
		BaseURL: "https://api.telegram.org/bot%s/sendMessage",
	}, nil
}

func (bot *TelegramBot) SendMessage(message string) error {
	payload := map[string]interface{}{
		"chat_id":                  bot.ChatID,
		"text":                     message,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	// Create the POST request
	url := fmt.Sprintf(bot.BaseURL, bot.Token)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message: %s", resp.Status)
	}

	return nil
}
