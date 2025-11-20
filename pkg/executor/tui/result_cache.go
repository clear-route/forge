// Package tui provides result caching for tool outputs
package tui

import (
	"sync"
	"time"

	"github.com/entrhq/forge/pkg/executor/tui/types"
)

// resultCache stores tool results for later viewing in overlays
type resultCache struct {
	mu      sync.RWMutex
	results map[string]*types.CachedResult // toolCallID -> cached result
	order   []string                       // LRU order (oldest first)
	maxSize int                            // Maximum number of results to cache
}

// newResultCache creates a new result cache with the specified size
func newResultCache(maxSize int) *resultCache {
	if maxSize <= 0 {
		maxSize = 20 // Default
	}
	return &resultCache{
		results: make(map[string]*types.CachedResult),
		order:   make([]string, 0, maxSize),
		maxSize: maxSize,
	}
}

// store adds a result to the cache, evicting the oldest if necessary
func (rc *resultCache) store(id string, toolName string, result string, summary string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// If already exists, update it and move to end
	if _, exists := rc.results[id]; exists {
		rc.remove(id) // Remove from current position
	}

	// Evict oldest if at capacity
	if len(rc.results) >= rc.maxSize {
		oldest := rc.order[0]
		delete(rc.results, oldest)
		rc.order = rc.order[1:]
	}

	// Add new result
	rc.results[id] = &types.CachedResult{
		ID:        id,
		ToolName:  toolName,
		Result:    result,
		Timestamp: time.Now(),
		Summary:   summary,
	}
	rc.order = append(rc.order, id)
}

// get retrieves a result from the cache
func (rc *resultCache) get(id string) (*types.CachedResult, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	result, exists := rc.results[id]
	return result, exists
}

// getLast retrieves the most recently added result
// getAll retrieves all cached results in reverse chronological order (newest first)
func (rc *resultCache) getAll() []*types.CachedResult {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	results := make([]*types.CachedResult, 0, len(rc.order))
	// Iterate in reverse order (newest first)
	for i := len(rc.order) - 1; i >= 0; i-- {
		id := rc.order[i]
		if result, exists := rc.results[id]; exists {
			results = append(results, result)
		}
	}
	return results
}

// remove removes a result from the cache (internal, assumes lock held)
func (rc *resultCache) remove(id string) {
	delete(rc.results, id)

	// Remove from order slice
	for i, oid := range rc.order {
		if oid == id {
			rc.order = append(rc.order[:i], rc.order[i+1:]...)
			break
		}
	}
}
