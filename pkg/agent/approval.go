package agent

import (
	"context"

	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/types"
)

// handleApprovalResponse processes an approval response from the user
func (a *DefaultAgent) handleApprovalResponse(response *types.ApprovalResponse) {
	a.approvalManager.HandleResponse(response)
}

// handleCommandCancellation processes a command cancellation request
func (a *DefaultAgent) handleCommandCancellation(req *types.CancellationRequest) {
	// Look up the cancel function for this execution ID
	if cancelFunc, ok := a.activeCommands.Load(req.ExecutionID); ok {
		// Cancel the context (cancellation never returns an error)
		if cf, ok := cancelFunc.(context.CancelFunc); ok {
			cf()
		}
		// Remove from active commands
		a.activeCommands.Delete(req.ExecutionID)
	}
}

// requestApproval sends an approval request and waits for user response
// Returns (approved, timedOut) where:
//   - approved: true if user approved, false if rejected
//   - timedOut: true if the request timed out waiting for response
func (a *DefaultAgent) requestApproval(ctx context.Context, toolCall tools.ToolCall, preview *tools.ToolPreview) (bool, bool) {
	// Delegate all approval logic to the approval manager
	return a.approvalManager.RequestApproval(ctx, toolCall, preview)
}
