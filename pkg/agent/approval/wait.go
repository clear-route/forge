package approval

import (
	"context"
	"time"

	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/types"
)

// waitForResponse waits for the user's approval response
func (m *Manager) waitForResponse(ctx context.Context, approvalID string, toolCall tools.ToolCall, responseChannel chan *types.ApprovalResponse, preview *tools.ToolPreview) (bool, bool) {
	timeout := time.NewTimer(m.timeout)
	defer timeout.Stop()

	select {
	case <-ctx.Done():
		return false, false

	case <-timeout.C:
		m.emitEvent(types.NewToolApprovalTimeoutEvent(approvalID, toolCall.ToolName))
		return false, true

	case approval := <-m.approvalChannel:
		return m.handleDirectApproval(ctx, approval, approvalID, toolCall, preview)

	case response := <-responseChannel:
		return m.handleChannelResponse(approvalID, toolCall, response)
	}
}

// handleDirectApproval handles approval received directly from executor channel
func (m *Manager) handleDirectApproval(ctx context.Context, approval *types.ApprovalResponse, approvalID string, toolCall tools.ToolCall, preview *tools.ToolPreview) (bool, bool) {
	if approval == nil {
		return false, false
	}

	// Verify it's for this approval request
	if approval.ApprovalID != approvalID {
		// Put it back for the right handler
		select {
		case m.approvalChannel <- approval:
		default:
		}
		// Continue waiting
		return m.RequestApproval(ctx, toolCall, preview)
	}

	// Process the approval
	if approval.IsGranted() {
		m.emitEvent(types.NewToolApprovalGrantedEvent(approvalID, toolCall.ToolName))
		return true, false
	}

	m.emitEvent(types.NewToolApprovalRejectedEvent(approvalID, toolCall.ToolName))
	return false, false
}

// handleChannelResponse handles response from internal response channel
// Returns true if approved, false if rejected
func (m *Manager) handleChannelResponse(approvalID string, toolCall tools.ToolCall, response *types.ApprovalResponse) bool {
	if response.IsGranted() {
		m.emitEvent(types.NewToolApprovalGrantedEvent(approvalID, toolCall.ToolName))
		return true
	}

	m.emitEvent(types.NewToolApprovalRejectedEvent(approvalID, toolCall.ToolName))
	return false
}
