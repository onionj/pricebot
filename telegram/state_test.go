package telegram

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestLoadState_Success(t *testing.T) {
	// Create a temporary test file
	stateFile = "test_state.json"
	defer os.Remove(stateFile)

	// Create test data
	testData := &Telegram{
		botToken:        "test_token",
		chatID:          "test_chat",
		LastMessageId:   100,
		LastMessageTime: time.Now().Unix(),
	}

	// Write test data to file
	data, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create new telegram instance and load state
	telegram := NewTelegram("test_token", "test_chat")

	// Verify loaded state
	if telegram.LastMessageId != testData.LastMessageId {
		t.Errorf("Expected LastMessageId %d, got %d", testData.LastMessageId, telegram.LastMessageId)
	}
	if telegram.LastMessageTime != testData.LastMessageTime {
		t.Errorf("Expected LastMessageTime %d, got %d", testData.LastMessageTime, telegram.LastMessageTime)
	}
}

func TestLoadState_FileNotExists(t *testing.T) {
	// Set non-existent file
	stateFile = "nonexistent_state.json"

	// Create new telegram instance - should not panic
	telegram := NewTelegram("test_token", "test_chat")

	// Verify default values
	if telegram.LastMessageId != 0 {
		t.Errorf("Expected LastMessageId 0, got %d", telegram.LastMessageId)
	}
	if telegram.LastMessageTime != 0 {
		t.Errorf("Expected LastMessageTime 0, got %d", telegram.LastMessageTime)
	}
}

func TestLoadState_InvalidJSON(t *testing.T) {
	// Create a temporary test file with invalid JSON
	tmpStateFile := "test_invalid_state.json"
	defer os.Remove(tmpStateFile)

	// Write invalid JSON to file
	invalidJSON := []byte(`{"last_message_id": "not_a_number"}`)
	if err := os.WriteFile(tmpStateFile, invalidJSON, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Set up test environment
	stateFile = tmpStateFile

	// Create new telegram instance - should not panic
	telegram := NewTelegram("test_token", "test_chat")

	// Verify default values are used
	if telegram.LastMessageId != 0 {
		t.Errorf("Expected LastMessageId 0, got %d", telegram.LastMessageId)
	}
	if telegram.LastMessageTime != 0 {
		t.Errorf("Expected LastMessageTime 0, got %d", telegram.LastMessageTime)
	}
}

func TestLoadState_EmptyFile(t *testing.T) {
	// Create a temporary empty test file
	tmpStateFile := "test_empty_state.json"
	defer os.Remove(tmpStateFile)

	// Create empty file
	if err := os.WriteFile(tmpStateFile, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Set up test environment
	stateFile = tmpStateFile

	// Create new telegram instance - should not panic
	telegram := NewTelegram("test_token", "test_chat")

	// Verify default values are used
	if telegram.LastMessageId != 0 {
		t.Errorf("Expected LastMessageId 0, got %d", telegram.LastMessageId)
	}
	if telegram.LastMessageTime != 0 {
		t.Errorf("Expected LastMessageTime 0, got %d", telegram.LastMessageTime)
	}
}

func TestSaveState_Success(t *testing.T) {
	// Create a temporary test file
	tmpStateFile := "test_save_state.json"
	defer os.Remove(tmpStateFile)

	// Set up test environment
	stateFile = tmpStateFile

	// Create and populate telegram instance
	telegram := NewTelegram("test_token", "test_chat")
	telegram.LastMessageId = 200
	telegram.LastMessageTime = time.Now().Unix()

	// Save state
	if err := telegram.saveState(); err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}

	// Read saved file
	data, err := os.ReadFile(tmpStateFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	// Parse saved data
	var savedState Telegram
	if err := json.Unmarshal(data, &savedState); err != nil {
		t.Fatalf("Failed to parse saved state: %v", err)
	}

	// Verify saved state
	if savedState.LastMessageId != telegram.LastMessageId {
		t.Errorf("Expected saved LastMessageId %d, got %d", telegram.LastMessageId, savedState.LastMessageId)
	}
	if savedState.LastMessageTime != telegram.LastMessageTime {
		t.Errorf("Expected saved LastMessageTime %d, got %d", telegram.LastMessageTime, savedState.LastMessageTime)
	}
}
