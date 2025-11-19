package types

import "time"

// CachedResult represents a single cached tool result with metadata
type CachedResult struct {
	ID        string    // Unique identifier
	ToolName  string    // Name of the tool that produced this result
	Result    string    // The actual result content
	Timestamp time.Time // When this result was created
	Summary   string    // Brief summary of the result
}
