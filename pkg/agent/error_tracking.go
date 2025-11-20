package agent

// trackError adds an error to the ring buffer and checks if we've hit the circuit breaker
// Returns true if the circuit breaker should trigger (5 identical consecutive errors)
func (a *DefaultAgent) trackError(errMsg string) bool {
	// Add to ring buffer
	a.lastErrors[a.errorIndex] = errMsg
	a.errorIndex = (a.errorIndex + 1) % 5

	// Check if all 5 are identical and non-empty
	if a.lastErrors[0] == "" {
		return false // Not enough errors yet
	}

	first := a.lastErrors[0]
	for i := 1; i < 5; i++ {
		if a.lastErrors[i] != first {
			return false
		}
	}

	return true // All 5 errors are identical
}

// resetErrorTracking clears the error ring buffer after a successful iteration
func (a *DefaultAgent) resetErrorTracking() {
	for i := range a.lastErrors {
		a.lastErrors[i] = ""
	}
	a.errorIndex = 0
}
