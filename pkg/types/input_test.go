package types

import (
	"testing"
)

func TestInputType(t *testing.T) {
	tests := []struct {
		name      string
		inputType InputType
		expected  string
	}{
		{
			name:      "cancel type",
			inputType: InputTypeCancel,
			expected:  "cancel",
		},
		{
			name:      "user_input type",
			inputType: InputTypeUserInput,
			expected:  "user_input",
		},
		{
			name:      "form_input type",
			inputType: InputTypeFormInput,
			expected:  "form_input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.inputType) != tt.expected {
				t.Errorf("InputType = %v, want %v", tt.inputType, tt.expected)
			}
		})
	}
}

func TestNewCancelInput(t *testing.T) {
	input := NewCancelInput()

	if input.Type != InputTypeCancel {
		t.Errorf("NewCancelInput type = %v, want %v", input.Type, InputTypeCancel)
	}
	if input.Content != "" {
		t.Error("NewCancelInput should not have content")
	}
	if input.FormData != nil {
		t.Error("NewCancelInput should not have form data")
	}
	if input.Metadata == nil {
		t.Error("NewCancelInput metadata should be initialized")
	}
}

func TestNewUserInput(t *testing.T) {
	content := "Hello, world!"
	input := NewUserInput(content)

	if input.Type != InputTypeUserInput {
		t.Errorf("NewUserInput type = %v, want %v", input.Type, InputTypeUserInput)
	}
	if input.Content != content {
		t.Errorf("NewUserInput content = %v, want %v", input.Content, content)
	}
	if input.FormData != nil {
		t.Error("NewUserInput should not have form data")
	}
	if input.Metadata == nil {
		t.Error("NewUserInput metadata should be initialized")
	}
}

func TestNewFormInput(t *testing.T) {
	formData := map[string]string{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   "30",
	}
	input := NewFormInput(formData)

	if input.Type != InputTypeFormInput {
		t.Errorf("NewFormInput type = %v, want %v", input.Type, InputTypeFormInput)
	}
	if input.Content != "" {
		t.Error("NewFormInput should not have content")
	}
	if input.FormData == nil {
		t.Fatal("NewFormInput should have form data")
	}
	if len(input.FormData) != len(formData) {
		t.Errorf("NewFormInput form data length = %v, want %v", len(input.FormData), len(formData))
	}
	for key, value := range formData {
		if input.FormData[key] != value {
			t.Errorf("NewFormInput form data[%s] = %v, want %v", key, input.FormData[key], value)
		}
	}
	if input.Metadata == nil {
		t.Error("NewFormInput metadata should be initialized")
	}
}

func TestInputWithMetadata(t *testing.T) {
	input := NewUserInput("test")
	key := "test_key"
	value := "test_value"

	result := input.WithMetadata(key, value)

	if result != input {
		t.Error("WithMetadata should return the same input for chaining")
	}
	if input.Metadata[key] != value {
		t.Errorf("WithMetadata did not set metadata correctly, got %v, want %v", input.Metadata[key], value)
	}
}

func TestInputHelpers(t *testing.T) {
	tests := []struct {
		input       *Input
		name        string
		isCancel    bool
		isUserInput bool
		isFormInput bool
	}{
		{
			name:        "cancel input",
			input:       NewCancelInput(),
			isCancel:    true,
			isUserInput: false,
			isFormInput: false,
		},
		{
			name:        "user input",
			input:       NewUserInput("test"),
			isCancel:    false,
			isUserInput: true,
			isFormInput: false,
		},
		{
			name:        "form input",
			input:       NewFormInput(map[string]string{"key": "value"}),
			isCancel:    false,
			isUserInput: false,
			isFormInput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input.IsCancel() != tt.isCancel {
				t.Errorf("IsCancel() = %v, want %v", tt.input.IsCancel(), tt.isCancel)
			}
			if tt.input.IsUserInput() != tt.isUserInput {
				t.Errorf("IsUserInput() = %v, want %v", tt.input.IsUserInput(), tt.isUserInput)
			}
			if tt.input.IsFormInput() != tt.isFormInput {
				t.Errorf("IsFormInput() = %v, want %v", tt.input.IsFormInput(), tt.isFormInput)
			}
		})
	}
}

func TestInputWithNilMetadata(t *testing.T) {
	input := &Input{
		Type:    InputTypeUserInput,
		Content: "test",
	}

	// Metadata is nil initially
	if input.Metadata != nil {
		t.Error("Metadata should be nil initially")
	}

	// Adding metadata should initialize it
	input.WithMetadata("key", "value")

	if input.Metadata == nil {
		t.Error("WithMetadata should initialize Metadata if nil")
	}
	if input.Metadata["key"] != "value" {
		t.Error("WithMetadata should set the value correctly")
	}
}

func TestFormInputEmpty(t *testing.T) {
	emptyForm := make(map[string]string)
	input := NewFormInput(emptyForm)

	if input.Type != InputTypeFormInput {
		t.Errorf("NewFormInput type = %v, want %v", input.Type, InputTypeFormInput)
	}
	if input.FormData == nil {
		t.Error("NewFormInput should have form data even if empty")
	}
	if len(input.FormData) != 0 {
		t.Errorf("NewFormInput with empty map should have 0 entries, got %v", len(input.FormData))
	}
}
