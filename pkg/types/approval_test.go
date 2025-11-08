package types

import (
	"testing"
	"time"
)

func TestNewApprovalResponse(t *testing.T) {
	tests := []struct {
		name       string
		approvalID string
		decision   ApprovalDecision
	}{
		{
			name:       "granted approval",
			approvalID: "test-123",
			decision:   ApprovalGranted,
		},
		{
			name:       "rejected approval",
			approvalID: "test-456",
			decision:   ApprovalRejected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			resp := NewApprovalResponse(tt.approvalID, tt.decision)
			after := time.Now()

			if resp.ApprovalID != tt.approvalID {
				t.Errorf("ApprovalID = %v, want %v", resp.ApprovalID, tt.approvalID)
			}

			if resp.Decision != tt.decision {
				t.Errorf("Decision = %v, want %v", resp.Decision, tt.decision)
			}

			if resp.Timestamp.Before(before) || resp.Timestamp.After(after) {
				t.Errorf("Timestamp %v not between %v and %v", resp.Timestamp, before, after)
			}
		})
	}
}

func TestApprovalResponse_IsGranted(t *testing.T) {
	tests := []struct {
		name     string
		decision ApprovalDecision
		want     bool
	}{
		{"granted", ApprovalGranted, true},
		{"rejected", ApprovalRejected, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &ApprovalResponse{Decision: tt.decision}
			if got := resp.IsGranted(); got != tt.want {
				t.Errorf("IsGranted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApprovalResponse_IsRejected(t *testing.T) {
	tests := []struct {
		name     string
		decision ApprovalDecision
		want     bool
	}{
		{"granted", ApprovalGranted, false},
		{"rejected", ApprovalRejected, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &ApprovalResponse{Decision: tt.decision}
			if got := resp.IsRejected(); got != tt.want {
				t.Errorf("IsRejected() = %v, want %v", got, tt.want)
			}
		})
	}
}
