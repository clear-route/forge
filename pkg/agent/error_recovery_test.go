package agent

import (
	"testing"

	"github.com/entrhq/forge/pkg/agent/prompts"
)

func TestErrorTracking(t *testing.T) {
	agent := &DefaultAgent{}

	t.Run("TracksSingleError", func(t *testing.T) {
		agent.resetErrorTracking()

		shouldBreak := agent.trackError("error 1")
		if shouldBreak {
			t.Error("should not trigger circuit breaker on first error")
		}
	})

	t.Run("TracksMultipleDifferentErrors", func(t *testing.T) {
		agent.resetErrorTracking()

		agent.trackError("error 1")
		agent.trackError("error 2")
		agent.trackError("error 3")
		shouldBreak := agent.trackError("error 4")

		if shouldBreak {
			t.Error("should not trigger circuit breaker on different errors")
		}
	})

	t.Run("TriggersCircuitBreakerOn5IdenticalErrors", func(t *testing.T) {
		agent.resetErrorTracking()

		agent.trackError("same error")
		agent.trackError("same error")
		agent.trackError("same error")
		agent.trackError("same error")
		shouldBreak := agent.trackError("same error")

		if !shouldBreak {
			t.Error("should trigger circuit breaker after 5 identical errors")
		}
	})

	t.Run("ResetsAfterSuccessfulIteration", func(t *testing.T) {
		agent.resetErrorTracking()

		agent.trackError("error")
		agent.trackError("error")
		agent.resetErrorTracking()

		// After reset, should not have any errors tracked
		if agent.lastErrors[0] != "" {
			t.Error("should reset error buffer")
		}
		if agent.errorIndex != 0 {
			t.Error("should reset error index")
		}
	})

	t.Run("CircuitBreakerRequiresAllFiveSlotsFilled", func(t *testing.T) {
		agent.resetErrorTracking()

		// Only 4 identical errors
		agent.trackError("error")
		agent.trackError("error")
		agent.trackError("error")
		shouldBreak := agent.trackError("error")

		if shouldBreak {
			t.Error("should not trigger with only 4 errors")
		}
	})
}

func TestErrorMessageBuilders(t *testing.T) {
	t.Run("BuildNoToolCallError", func(t *testing.T) {
		msg := prompts.BuildErrorRecoveryMessage(prompts.ErrorRecoveryContext{
			Type: prompts.ErrorTypeNoToolCall,
		})

		if msg == "" {
			t.Error("should return non-empty message")
		}

		// Check for key instructions
		if len(msg) < 50 {
			t.Error("error message should be detailed")
		}
	})

	t.Run("BuildMissingToolNameError", func(t *testing.T) {
		msg := prompts.BuildErrorRecoveryMessage(prompts.ErrorRecoveryContext{
			Type: prompts.ErrorTypeMissingToolName,
		})

		if msg == "" {
			t.Error("should return non-empty message")
		}
	})

	t.Run("BuildToolExecutionError", func(t *testing.T) {
		testErr := &testError{msg: "execution failed"}
		msg := prompts.BuildErrorRecoveryMessage(prompts.ErrorRecoveryContext{
			Type:     prompts.ErrorTypeToolExecution,
			ToolName: "calculator",
			Error:    testErr,
		})

		if msg == "" {
			t.Error("should return non-empty message")
		}

		// Should include tool name
		if len(msg) < 30 {
			t.Error("error message should be detailed")
		}
	})
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
