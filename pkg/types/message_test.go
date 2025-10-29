package types

import (
	"testing"
	"time"
)

func TestMessageRole(t *testing.T) {
	tests := []struct {
		name string
		role MessageRole
		want string
	}{
		{"system role", RoleSystem, "system"},
		{"user role", RoleUser, "user"},
		{"assistant role", RoleAssistant, "assistant"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.role) != tt.want {
				t.Errorf("MessageRole = %v, want %v", tt.role, tt.want)
			}
		})
	}
}

func TestNewMessage(t *testing.T) {
	content := "Hello, world!"
	msg := NewMessage(RoleUser, content)

	if msg.Role != RoleUser {
		t.Errorf("NewMessage role = %v, want %v", msg.Role, RoleUser)
	}
	if msg.Content != content {
		t.Errorf("NewMessage content = %v, want %v", msg.Content, content)
	}
	if msg.Timestamp.IsZero() {
		t.Error("NewMessage timestamp should not be zero")
	}
	if msg.Metadata == nil {
		t.Error("NewMessage metadata should be initialized")
	}
}

func TestNewSystemMessage(t *testing.T) {
	content := "You are a helpful assistant"
	msg := NewSystemMessage(content)

	if msg.Role != RoleSystem {
		t.Errorf("NewSystemMessage role = %v, want %v", msg.Role, RoleSystem)
	}
	if msg.Content != content {
		t.Errorf("NewSystemMessage content = %v, want %v", msg.Content, content)
	}
}

func TestNewUserMessage(t *testing.T) {
	content := "What is the weather?"
	msg := NewUserMessage(content)

	if msg.Role != RoleUser {
		t.Errorf("NewUserMessage role = %v, want %v", msg.Role, RoleUser)
	}
	if msg.Content != content {
		t.Errorf("NewUserMessage content = %v, want %v", msg.Content, content)
	}
}

func TestNewAssistantMessage(t *testing.T) {
	content := "The weather is sunny"
	msg := NewAssistantMessage(content)

	if msg.Role != RoleAssistant {
		t.Errorf("NewAssistantMessage role = %v, want %v", msg.Role, RoleAssistant)
	}
	if msg.Content != content {
		t.Errorf("NewAssistantMessage content = %v, want %v", msg.Content, content)
	}
}

func TestMessageWithMetadata(t *testing.T) {
	msg := NewUserMessage("test")
	key := "test_key"
	value := "test_value"

	result := msg.WithMetadata(key, value)

	if result != msg {
		t.Error("WithMetadata should return the same message for chaining")
	}
	if msg.Metadata[key] != value {
		t.Errorf("WithMetadata did not set metadata correctly, got %v, want %v", msg.Metadata[key], value)
	}
}

func TestMessageTimestamp(t *testing.T) {
	before := time.Now()
	msg := NewMessage(RoleUser, "test")
	after := time.Now()

	if msg.Timestamp.Before(before) || msg.Timestamp.After(after) {
		t.Errorf("Message timestamp %v should be between %v and %v", msg.Timestamp, before, after)
	}
}
