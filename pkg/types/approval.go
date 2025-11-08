package types

import "time"

// ApprovalDecision represents a user's decision on a tool approval request.
type ApprovalDecision string

const (
	// ApprovalGranted indicates the user approved the tool execution
	ApprovalGranted ApprovalDecision = "granted"

	// ApprovalRejected indicates the user rejected the tool execution
	ApprovalRejected ApprovalDecision = "rejected"
)

// ApprovalResponse represents a response to a tool approval request.
type ApprovalResponse struct {
	// ApprovalID matches the ID from the approval request
	ApprovalID string

	// Decision is the user's approval decision
	Decision ApprovalDecision

	// Timestamp when the decision was made
	Timestamp time.Time
}

// NewApprovalResponse creates a new approval response.
func NewApprovalResponse(approvalID string, decision ApprovalDecision) *ApprovalResponse {
	return &ApprovalResponse{
		ApprovalID: approvalID,
		Decision:   decision,
		Timestamp:  time.Now(),
	}
}

// IsGranted returns true if the approval was granted.
func (r *ApprovalResponse) IsGranted() bool {
	return r.Decision == ApprovalGranted
}

// IsRejected returns true if the approval was rejected.
func (r *ApprovalResponse) IsRejected() bool {
	return r.Decision == ApprovalRejected
}
