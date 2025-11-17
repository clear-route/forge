// Package tui provides result caching for tool outputs
package tui

import (
	"sync"
)

// resultCache stores tool results for later viewing in overlays
type resultCache struct {
	mu      sync.RWMutex
	results map[string]string // toolCallID -> result
	order   []string          // LRU order (oldest first)
	maxSize int               // Maximum number of results to cache
}

// newResultCache creates a new result cache with the specified size
func newResultCache(maxSize int) *resultCache {
	if maxSize <= 0 {
		maxSize = 20 // Default
	}
	return &resultCache{
		results: make(map[string]string),
		order:   make([]string, 0, maxSize),
		maxSize: maxSize,
	}
}

// store adds a result to the cache, evicting the oldest if necessary
func (rc *resultCache) store(id string, result string) {
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
	rc.results[id] = result
	rc.order = append(rc.order, id)
}

// get retrieves a result from the cache
func (rc *resultCache) get(id string) (string, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	result, exists := rc.results[id]
	return result, exists
}

// getLast retrieves the most recently added result
func (rc *resultCache) getLast() (string, string, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	if len(rc.order) == 0 {
		return "", "", false
	}

	lastID := rc.order[len(rc.order)-1]
	result := rc.results[lastID]
	return lastID, result, true
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

// clear removes all results from the cache
func (rc *resultCache) clear() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.results = make(map[string]string)
	rc.order = make([]string, 0, rc.maxSize)
}

// size returns the current number of cached results
func (rc *resultCache) size() int {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	return len(rc.results)
}
