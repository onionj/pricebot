package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Telegram struct {
	botToken        string
	chatID          string
	LastMessageId   int
	LastMessageTime int64
}

type sendMessageResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
	Description string `json:"description"`
	ErrCode     int    `json:"error_code"`
}

func NewTelegram(botToken, chatID string) *Telegram {
	// TODO: save and reload LastMessageTime and LastMessageId from file/db
	return &Telegram{botToken: botToken, chatID: chatID, LastMessageTime: 0}
}

func (t *Telegram) SendMessage(msg string) error {
	payload := map[string]string{
		"chat_id": t.chatID,
		"text":    msg,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken),
		"application/json", bytes.NewBuffer(body))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response sendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if !response.OK {
		return fmt.Errorf("failed to send message: code:%d, description:%s", response.ErrCode, response.Description)
	}

	t.LastMessageId = response.Result.MessageID
	t.LastMessageTime = time.Now().Unix()
	return nil
}

func (t *Telegram) UpdateMessage(msg string, messageId int) error {
	payload := map[string]interface{}{
		"chat_id":    t.chatID,
		"message_id": messageId,
		"text":       msg,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText", t.botToken),
		"application/json", bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var response sendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if !response.OK {
		return fmt.Errorf("failed to update message: code:%d, description:%s", response.ErrCode, response.Description)
	}

	return nil
}
