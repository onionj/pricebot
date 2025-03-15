package telegram

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Define request body struct for better type safety
type updateMessageRequest struct {
	ChatID    string `json:"chat_id"`
	MessageID int    `json:"message_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type sendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func TestTelegram_SendMessage(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot123456:ABC-DEF/sendMessage" {
			t.Errorf("Expected path '/bot123456:ABC-DEF/sendMessage', got %s", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}

		// Check request body using the struct
		var requestBody sendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if requestBody.ChatID != "test_chat_id" {
			t.Errorf("Expected chat_id 'test_chat_id', got '%s'", requestBody.ChatID)
		}

		if requestBody.ParseMode != "Markdown" {
			t.Errorf("Expected parse_mode 'Markdown', got '%s'", requestBody.ParseMode)
		}

		// Send mock response
		response := messageResponse{
			OK: true,
			Result: struct {
				MessageID int `json:"message_id"`
			}{
				MessageID: 123,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create temporary state file
	tmpStateFile := "test_state.json"
	defer os.Remove(tmpStateFile)

	// Set test values
	stateFile = tmpStateFile
	httpClient = server.Client()
	baseURL = server.URL + "/bot%s%s"

	telegram := NewTelegram("123456:ABC-DEF", "test_chat_id")

	// Test SendMessage
	err := telegram.SendMessage("Test message")
	if err != nil {
		t.Errorf("SendMessage failed: %v", err)
	}

	if telegram.LastMessageId != 123 {
		t.Errorf("Expected LastMessageId to be 123, got %d", telegram.LastMessageId)
	}
}

func TestTelegram_UpdateMessage(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot123456:ABC-DEF/editMessageText" {
			t.Errorf("Expected path '/bot123456:ABC-DEF/editMessageText', got %s", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}

		// Check request body using the struct
		var requestBody updateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if requestBody.ChatID != "test_chat_id" {
			t.Errorf("Expected chat_id 'test_chat_id', got '%s'", requestBody.ChatID)
		}

		if requestBody.MessageID != 123 {
			t.Errorf("Expected message_id 123, got %d", requestBody.MessageID)
		}

		if requestBody.ParseMode != "Markdown" {
			t.Errorf("Expected parse_mode 'Markdown', got '%s'", requestBody.ParseMode)
		}

		// Send mock response
		response := messageResponse{
			OK: true,
			Result: struct {
				MessageID int `json:"message_id"`
			}{
				MessageID: 123,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create temporary state file
	tmpStateFile := "test_state.json"
	defer os.Remove(tmpStateFile)

	// Set test values
	stateFile = tmpStateFile
	httpClient = server.Client()
	baseURL = server.URL + "/bot%s%s"

	telegram := NewTelegram("123456:ABC-DEF", "test_chat_id")

	// Test UpdateMessage
	err := telegram.UpdateMessage("Updated message", 123)
	if err != nil {
		t.Errorf("UpdateMessage failed: %v", err)
	}
}

func TestTelegram_StateManagement(t *testing.T) {
	// Create temporary state file
	tmpStateFile := "test_state.json"
	defer os.Remove(tmpStateFile)

	// Set up test state
	testState := &Telegram{
		LastMessageId:   456,
		LastMessageTime: 1234567890,
	}

	// Save test state to file
	stateData, _ := json.Marshal(testState)
	if err := os.WriteFile(tmpStateFile, stateData, 0644); err != nil {
		t.Fatalf("Failed to write test state file: %v", err)
	}

	// Save original state file path and restore after test
	originalStateFile := stateFile
	defer func() { stateFile = originalStateFile }()

	// Set test state file
	stateFile = tmpStateFile

	telegram := NewTelegram("test_token", "test_chat_id")

	// Verify state was loaded correctly
	if telegram.LastMessageId != 456 {
		t.Errorf("Expected LastMessageId to be 456, got %d", telegram.LastMessageId)
	}
	if telegram.LastMessageTime != 1234567890 {
		t.Errorf("Expected LastMessageTime to be 1234567890, got %d", telegram.LastMessageTime)
	}
}
