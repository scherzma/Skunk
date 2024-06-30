package test

import (
	"testing"
	"time"

	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/frontend"
)

// TestCreateMessage tests the createMessage function.
func TestCreateMessage(t *testing.T) {
	msg := frontend.CreateMessage("1", "Alice", "Bob", "Hello", "", "", frontend.SEND_MESSAGE)
	if msg.ChatID != "1" {
		t.Errorf("Expected ChatID '1', got %s", msg.ChatID)
	}
	if msg.SenderID != "Alice" {
		t.Errorf("Expected SenderID 'Alice', got %s", msg.SenderID)
	}
	if msg.ReceiverID != "Bob" {
		t.Errorf("Expected ReceiverID 'Bob', got %s", msg.ReceiverID)
	}
	if msg.Content != "Hello" {
		t.Errorf("Expected Content 'Hello', got %s", msg.Content)
	}
	if msg.Operation != frontend.SEND_MESSAGE {
		t.Errorf("Expected Operation 'SEND_MESSAGE', got %d", msg.Operation)
	}
}

// TestIsValidUsername tests the isValidUsername function.
func TestIsValidUsername(t *testing.T) {
	validUsernames := []string{"Alice", "Bob123"}
	invalidUsernames := []string{"Charlie_456", "", " ", "   ", "invalid username", "verylongusernamethatexceeds20chars"}

	for _, username := range validUsernames {
		if !frontend.IsValidUsername(username) {
			t.Errorf("Expected '%s' to be valid", username)
		}
	}

	for _, username := range invalidUsernames {
		if frontend.IsValidUsername(username) {
			t.Errorf("Expected '%s' to be invalid", username)
		}
	}
}

// TestSanitizeInput tests the sanitizeInput function.
func TestSanitizeInput(t *testing.T) {
	input := "  Hello, World!  "
	expected := "Hello, World!"
	sanitized := frontend.SanitizeInput(input)
	if sanitized != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sanitized)
	}
}

// TestValidateInput tests the validateInput function.
func TestValidateInput(t *testing.T) {
	validInput := "Hello"
	invalidInput := "Hello, World! This input is way too long and should be rejected."

	_, err := frontend.ValidateInput(validInput, 10)
	if err != nil {
		t.Errorf("Expected no error for valid input, got %v", err)
	}

	_, err = frontend.ValidateInput(invalidInput, 10)
	if err == nil {
		t.Errorf("Expected error for invalid input, got none")
	}
}

// TestHandleCommand tests the handleCommand function.
func TestHandleCommand(t *testing.T) {
	m := frontend.InitialModel()
	m.Usernames["1"] = "Alice"
	m.CurrentChat = "1"
	m.Chats["1"] = []frontend.FrontendMessage{} // Initialize the chat messages

	_, cmd := m.HandleCommand("/invite", "Bob")
	if cmd == nil {
		t.Errorf("Expected command, got nil")
	}

	if len(m.Chats["1"]) != 1 {
		t.Errorf("Expected 1 message, got %d", len(m.Chats["1"]))
	} else if m.Chats["1"][0].ReceiverID != "Bob" {
		t.Errorf("Expected ReceiverID 'Bob', got %s", m.Chats["1"][0].ReceiverID)
	}
}

// TestClearTempMessage tests the clearTempMessage function.
func TestClearTempMessage(t *testing.T) {
	m := frontend.InitialModel()
	m.TempMessage = "Test message"
	m.TempMessageExpire = time.Now().Add(1 * time.Second)

	// Wait for the temporary message to expire
	time.Sleep(2 * time.Second)

	// Call the update function with tempMsgTimeoutMsg
	modelInterface, _ := m.Update(frontend.TempMsgTimeoutMsg{})
	m, _ = modelInterface.(frontend.Model)

	if m.TempMessage != "" {
		t.Errorf("Expected message to be empty, got %s", m.TempMessage)
	}
}

// TestTestConnection tests the testConnection function.
func TestTestConnection(t *testing.T) {
	if !frontend.TestConnection("validOnionID") {
		t.Error("Expected connection to be successful")
	}
	if frontend.TestConnection("invalidOnionID") {
		t.Error("Expected connection to fail")
	}
}

// TestToggleFocus tests the toggleFocus function.
func TestToggleFocus(t *testing.T) {
	if frontend.ToggleFocus("Chats") != "invites" {
		t.Error("Expected focus to switch to 'invites'")
	}
	if frontend.ToggleFocus("invites") != "Chats" {
		t.Error("Expected focus to switch to 'Chats'")
	}
}

// TestGetSortedChatIDs tests the getSortedChatIDs function.
func TestGetSortedChatIDs(t *testing.T) {
	chatNames := map[string]string{
		"1": "Chat One",
		"2": "Chat Two",
		"3": "Chat Three",
	}
	expected := []string{"1", "2", "3"}
	sorted := frontend.GetSortedChatIDs(chatNames)

	for i, id := range sorted {
		if id != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], id)
		}
	}
}
