package agent

import (
	"context"
	"testing"
	"time"

	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/types"
)

func TestApprovalSystem_RequestApproval(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name             string
		sendResponse     bool
		responseDecision types.ApprovalDecision
		timeout          time.Duration
		expectApproved   bool
		expectTimedOut   bool
	}{
		{
			name:             "approval granted",
			sendResponse:     true,
			responseDecision: types.ApprovalGranted,
			timeout:          1 * time.Second,
			expectApproved:   true,
			expectTimedOut:   false,
		},
		{
			name:             "approval rejected",
			sendResponse:     true,
			responseDecision: types.ApprovalRejected,
			timeout:          1 * time.Second,
			expectApproved:   false,
			expectTimedOut:   false,
		},
		{
			name:           "approval timeout",
			sendResponse:   false,
			timeout:        100 * time.Millisecond,
			expectApproved: false,
			expectTimedOut: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &DefaultAgent{
				approvalTimeout: tt.timeout,
				channels:        types.NewAgentChannels(10),
			}

			toolCall := tools.ToolCall{
				ServerName:   "local",
				ToolName: "test_tool",
				Arguments: tools.ArgumentsBlock{
					InnerXML: []byte(`<arg>value</arg>`),
				},
			}

			preview := &tools.ToolPreview{
				Type:    tools.PreviewTypeDiff,
				Title:   "Test preview",
				Content: "preview content",
			}

			if tt.sendResponse {
				go func() {
					time.Sleep(50 * time.Millisecond)
					agent.approvalMu.Lock()
					var approvalID string
					if agent.pendingApproval != nil {
						approvalID = agent.pendingApproval.approvalID
					}
					agent.approvalMu.Unlock()

					if approvalID != "" {
						response := types.NewApprovalResponse(approvalID, tt.responseDecision)
						agent.handleApprovalResponse(response)
					}
				}()
			}

			approved, timedOut := agent.requestApproval(ctx, toolCall, preview)

			if approved != tt.expectApproved {
				t.Errorf("approved = %v, want %v", approved, tt.expectApproved)
			}

			if timedOut != tt.expectTimedOut {
				t.Errorf("timedOut = %v, want %v", timedOut, tt.expectTimedOut)
			}

			agent.approvalMu.Lock()
			if agent.pendingApproval != nil {
				t.Error("pending approval should be nil after request completes")
			}
			agent.approvalMu.Unlock()
		})
	}
}
