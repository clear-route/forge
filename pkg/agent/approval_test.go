package agent

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/entrhq/forge/pkg/agent/approval"
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
			channels := types.NewAgentChannels(10)

			// Track events for verification - use mutex to prevent data race
			var lastApprovalID string
			var approvalIDMutex sync.Mutex

			emitEvent := func(event *types.AgentEvent) {
				if event.Type == types.EventTypeToolApprovalRequest {
					approvalIDMutex.Lock()
					lastApprovalID = event.ApprovalID
					approvalIDMutex.Unlock()
				}
				channels.Event <- event
			}

			agent := &DefaultAgent{
				channels:        channels,
				approvalManager: approval.NewManager(tt.timeout, emitEvent),
			}

			toolCall := tools.ToolCall{
				ServerName: "local",
				ToolName:   "test_tool",
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
					// Wait for approval request event
					time.Sleep(50 * time.Millisecond)

					approvalIDMutex.Lock()
					id := lastApprovalID
					approvalIDMutex.Unlock()

					if id != "" {
						response := types.NewApprovalResponse(id, tt.responseDecision)
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
		})
	}
}
