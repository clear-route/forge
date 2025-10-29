package types

import (
	"testing"
)

func TestNewAgentChannels(t *testing.T) {
	bufferSize := 10
	channels := NewAgentChannels(bufferSize)

	if channels.Input == nil {
		t.Error("NewAgentChannels Input channel should be initialized")
	}
	if channels.Event == nil {
		t.Error("NewAgentChannels Event channel should be initialized")
	}
	if channels.Shutdown == nil {
		t.Error("NewAgentChannels Shutdown channel should be initialized")
	}
	if channels.Done == nil {
		t.Error("NewAgentChannels Done channel should be initialized")
	}

	// Check buffer sizes
	if cap(channels.Input) != bufferSize {
		t.Errorf("Input channel buffer = %v, want %v", cap(channels.Input), bufferSize)
	}
	if cap(channels.Event) != bufferSize {
		t.Errorf("Event channel buffer = %v, want %v", cap(channels.Event), bufferSize)
	}
}

func TestAgentChannelsBuffering(t *testing.T) {
	bufferSize := 2
	channels := NewAgentChannels(bufferSize)

	// Should be able to send bufferSize inputs without blocking
	input1 := NewUserInput("test1")
	input2 := NewUserInput("test2")

	channels.Input <- input1
	channels.Input <- input2

	// Verify inputs can be received
	received1 := <-channels.Input
	received2 := <-channels.Input

	if received1.Content != "test1" {
		t.Errorf("Received input 1 content = %v, want test1", received1.Content)
	}
	if received2.Content != "test2" {
		t.Errorf("Received input 2 content = %v, want test2", received2.Content)
	}
}

func TestAgentChannelsClose(t *testing.T) {
	channels := NewAgentChannels(10)

	// Close all channels
	channels.Close()

	// Verify channels are closed by checking if receive returns zero value
	// Note: We can't directly check if a channel is closed, but we can verify
	// that receiving from it doesn't block and returns the zero value

	select {
	case _, ok := <-channels.Event:
		if ok {
			t.Error("Event channel should be closed")
		}
	default:
		t.Error("Event channel should be closed and readable")
	}

	select {
	case _, ok := <-channels.Done:
		if ok {
			t.Error("Done channel should be closed")
		}
	default:
		t.Error("Done channel should be closed and readable")
	}
}

func TestAgentChannelsShutdownSignal(t *testing.T) {
	channels := NewAgentChannels(10)

	// Close shutdown channel to signal shutdown
	close(channels.Shutdown)

	// Verify we can detect the shutdown signal
	select {
	case <-channels.Shutdown:
		// Successfully received shutdown signal
	default:
		t.Error("Should be able to receive from closed Shutdown channel")
	}
}

func TestAgentChannelsDifferentBufferSizes(t *testing.T) {
	tests := []struct {
		name       string
		bufferSize int
	}{
		{"small buffer", 1},
		{"medium buffer", 10},
		{"large buffer", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channels := NewAgentChannels(tt.bufferSize)

			if cap(channels.Input) != tt.bufferSize {
				t.Errorf("Input channel buffer = %v, want %v", cap(channels.Input), tt.bufferSize)
			}
			if cap(channels.Event) != tt.bufferSize {
				t.Errorf("Event channel buffer = %v, want %v", cap(channels.Event), tt.bufferSize)
			}
		})
	}
}
