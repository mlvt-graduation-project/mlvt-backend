package notify_service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"mlvt/internal/infra/env"
	"mlvt/internal/infra/zap-logging/log"
)

// Telegram API structs
type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

type GetUpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type TelegramResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
}

type NotifyService interface {
	GetChatID() error
	SendTelegram(message string) error
	SendNotification(message string) error
	SendFormattedNotification(message string) error
}

type notifyServiceImpl struct{}

func NewNotifyService() NotifyService {
	return &notifyServiceImpl{}
}

// GetChatID retrieves recent messages and displays chat IDs
func (s *notifyServiceImpl) GetChatID() error {
	cfg := env.EnvConfig

	if cfg.TelegramBotToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
	}

	// Get updates from Telegram
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", cfg.TelegramBotToken)
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("error getting updates: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	var updates GetUpdatesResponse
	if err := json.Unmarshal(body, &updates); err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	log.Info("Recent messages to your bot:")
	log.Info("================================")

	if len(updates.Result) == 0 {
		log.Info("No messages found. Please send a message to your bot first!")
		log.Info("1. Open Telegram")
		log.Info("2. Search for your bot")
		log.Info("3. Send any message to it")
		log.Info("4. Run this command again")
		return nil
	}

	for _, update := range updates.Result {
		log.Infof("Chat ID: %d", update.Message.Chat.ID)
		log.Infof("From: %s (@%s)", update.Message.From.FirstName, update.Message.From.Username)
		log.Infof("Message: %s", update.Message.Text)
		log.Info("---")
	}

	// Get the most recent chat ID
	if len(updates.Result) > 0 {
		latestChatID := updates.Result[len(updates.Result)-1].Message.Chat.ID
		log.Infof("🎯 Use this Chat ID in your .env file:")
		log.Infof("TELEGRAM_CHAT_ID=%d", latestChatID)
	}

	return nil
}

// SendTelegram sends a message to the configured Telegram chat
func (s *notifyServiceImpl) SendTelegram(message string) error {
	cfg := env.EnvConfig

	if cfg.TelegramBotToken == "" || cfg.TelegramChatID == "" {
		log.Warn("Telegram configuration not set: missing bot token or chat ID")
		return fmt.Errorf("telegram configuration not set: missing bot token or chat ID")
	}

	log.Debugf("Sending Telegram notification with bot token: %s", maskToken(cfg.TelegramBotToken))
	log.Debugf("Sending to chat ID: %s", cfg.TelegramChatID)

	// apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.TelegramBotToken)

	// data := url.Values{}
	// data.Set("chat_id", cfg.TelegramChatID)
	// data.Set("text", message)
	// data.Set("parse_mode", "HTML") // Enable HTML formatting

	// resp, err := http.PostForm(apiURL, data)
	// if err != nil {
	// 	log.Errorf("Failed to send HTTP request to Telegram: %v", err)
	// 	return fmt.Errorf("failed to send HTTP request: %v", err)
	// }
	// defer resp.Body.Close()

	// // Read response body for better error handling
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Errorf("Error reading Telegram API response: %v", err)
	// 	return fmt.Errorf("error reading response: %v", err)
	// }

	// if resp.StatusCode != 200 {
	// 	log.Errorf("Telegram API returned non-200 status: %s, body: %s", resp.Status, string(body))
	// 	return fmt.Errorf("telegram API returned non-200 status: %s", resp.Status)
	// }

	// // Parse response to check if message was sent successfully
	// var telegramResp TelegramResponse
	// if err := json.Unmarshal(body, &telegramResp); err != nil {
	// 	log.Errorf("Error parsing Telegram API response: %v", err)
	// 	return fmt.Errorf("error parsing response: %v", err)
	// }

	// if !telegramResp.OK {
	// 	log.Errorf("Telegram API returned error: %s", string(body))
	// 	return fmt.Errorf("telegram API returned error")
	// }

	// log.Infof("Telegram message sent successfully with message ID: %d", telegramResp.Result.MessageID)
	return nil
}

// SendNotification sends a notification (currently only supports Telegram)
func (s *notifyServiceImpl) SendNotification(message string) error {
	return s.SendTelegram(message)
}

// SendFormattedNotification sends a notification with MLVT header format
func (s *notifyServiceImpl) SendFormattedNotification(message string) error {
	formattedMessage := s.formatNotificationMessage(message)
	return s.SendTelegram(formattedMessage)
}

// Helper function to format notification message with header
func (s *notifyServiceImpl) formatNotificationMessage(message string) string {
	cfg := env.EnvConfig
	appEnv := cfg.AppEnv
	if appEnv == "" {
		appEnv = "unknown"
	}

	header := fmt.Sprintf("🤖 <b>MLVT Backend [%s]</b>\n"+
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n",
		strings.ToUpper(appEnv))

	footer := "\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		fmt.Sprintf("⏰ <i>%s</i>", time.Now().Format("2006-01-02 15:04:05"))

	return header + message + footer
}

// maskToken masks the bot token for logging purposes
func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}
