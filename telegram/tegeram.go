package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Make these package variables so they can be modified in tests
var (
	stateFile  = "telegram_state.json"
	httpClient = &http.Client{}
	baseURL    = "https://api.telegram.org/bot%s%s"
)

type Telegram struct {
	botToken        string
	chatID          string
	LastMessageId   int   `json:"last_message_id"`
	LastMessageTime int64 `json:"last_message_time"`
}

type messageResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
	Description string `json:"description"`
	ErrCode     int    `json:"error_code"`
}

// NewTelegram initializes a Telegram bot and loads state from a file
func NewTelegram(botToken, chatID string) *Telegram {
	t := &Telegram{botToken: botToken, chatID: chatID}

	// Load state from file
	if err := t.loadState(); err != nil {
		fmt.Println("Warning: Could not load state,", err)
	}

	return t
}

// Save state to file
func (t *Telegram) saveState() error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(stateFile, data, 0644)
}

// Load state from file
func (t *Telegram) loadState() error {
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, t); err != nil {
		return err
	}

	return nil
}

func (t *Telegram) SendMessage(msg string) error {
	payload := map[string]string{
		"chat_id":    t.chatID,
		"text":       msg,
		"parse_mode": "Markdown",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := httpClient.Post(
		fmt.Sprintf(baseURL, t.botToken, "/sendMessage"),
		"application/json", bytes.NewBuffer(body))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response messageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if !response.OK {
		return fmt.Errorf("failed to send message: code:%d, description:%s", response.ErrCode, response.Description)
	}

	t.LastMessageId = response.Result.MessageID
	t.LastMessageTime = time.Now().Unix()
	return t.saveState()
}

func (t *Telegram) UpdateMessage(msg string, messageId int) error {
	payload := map[string]interface{}{
		"chat_id":    t.chatID,
		"message_id": messageId,
		"text":       msg,
		"parse_mode": "Markdown",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := httpClient.Post(
		fmt.Sprintf(baseURL, t.botToken, "/editMessageText"),
		"application/json", bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var response messageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if !response.OK {
		return fmt.Errorf("failed to update message: code:%d, description:%s", response.ErrCode, response.Description)
	}

	return nil
}
