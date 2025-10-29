package types

import (
	"errors"
	"testing"
)

func TestErrorCode(t *testing.T) {
	tests := []struct {
		name string
		code ErrorCode
		want string
	}{
		{"llm failure", ErrorCodeLLMFailure, "llm_failure"},
		{"shutdown", ErrorCodeShutdown, "shutdown"},
		{"invalid input", ErrorCodeInvalidInput, "invalid_input"},
		{"timeout", ErrorCodeTimeout, "timeout"},
		{"canceled", ErrorCodeCanceled, "canceled"},
		{"internal", ErrorCodeInternal, "internal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.code) != tt.want {
				t.Errorf("ErrorCode = %v, want %v", tt.code, tt.want)
			}
		})
	}
}

func TestNewAgentError(t *testing.T) {
	code := ErrorCodeLLMFailure
	message := "Failed to call LLM"
	err := NewAgentError(code, message)

	if err.Code != code {
		t.Errorf("NewAgentError code = %v, want %v", err.Code, code)
	}
	if err.Message != message {
		t.Errorf("NewAgentError message = %v, want %v", err.Message, message)
	}
	if err.Cause != nil {
		t.Error("NewAgentError should not have a cause")
	}
	if err.Metadata == nil {
		t.Error("NewAgentError metadata should be initialized")
	}
}

func TestNewAgentErrorWithCause(t *testing.T) {
	code := ErrorCodeLLMFailure
	message := "Failed to call LLM"
	cause := errors.New("network error")
	err := NewAgentErrorWithCause(code, message, cause)

	if err.Code != code {
		t.Errorf("NewAgentErrorWithCause code = %v, want %v", err.Code, code)
	}
	if err.Message != message {
		t.Errorf("NewAgentErrorWithCause message = %v, want %v", err.Message, message)
	}
	if err.Cause != cause {
		t.Errorf("NewAgentErrorWithCause cause = %v, want %v", err.Cause, cause)
	}
	if err.Metadata == nil {
		t.Error("NewAgentErrorWithCause metadata should be initialized")
	}
}

func TestAgentErrorError(t *testing.T) {
	tests := []struct {
		name string
		err  *AgentError
		want string
	}{
		{
			name: "error without cause",
			err:  NewAgentError(ErrorCodeLLMFailure, "test error"),
			want: "llm_failure: test error",
		},
		{
			name: "error with cause",
			err:  NewAgentErrorWithCause(ErrorCodeLLMFailure, "test error", errors.New("network error")),
			want: "llm_failure: test error (caused by: network error)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("AgentError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgentErrorUnwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewAgentErrorWithCause(ErrorCodeInternal, "wrapper error", cause)

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("AgentError.Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestAgentErrorWithMetadata(t *testing.T) {
	err := NewAgentError(ErrorCodeInternal, "test")
	key := "test_key"
	value := "test_value"

	result := err.WithMetadata(key, value)

	if result != err {
		t.Error("WithMetadata should return the same error for chaining")
	}
	if err.Metadata[key] != value {
		t.Errorf("WithMetadata did not set metadata correctly, got %v, want %v", err.Metadata[key], value)
	}
}

func TestIsAgentError(t *testing.T) {
	tests := []struct {
		err       error
		name      string
		wantAgent bool
		wantErr   bool
	}{
		{
			name:      "nil error",
			err:       nil,
			wantErr:   false,
			wantAgent: false,
		},
		{
			name:      "agent error",
			err:       NewAgentError(ErrorCodeInternal, "test"),
			wantErr:   true,
			wantAgent: true,
		},
		{
			name:      "standard error",
			err:       errors.New("standard error"),
			wantErr:   true,
			wantAgent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agentErr, isAgent := IsAgentError(tt.err)

			if isAgent != tt.wantAgent {
				t.Errorf("IsAgentError() isAgent = %v, want %v", isAgent, tt.wantAgent)
			}

			if tt.wantAgent && agentErr == nil {
				t.Error("IsAgentError() returned true but agentErr is nil")
			}

			if !tt.wantAgent && agentErr != nil {
				t.Error("IsAgentError() returned false but agentErr is not nil")
			}
		})
	}
}
