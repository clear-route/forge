// Package git provides Git integration functionality for Forge agents,
// including file modification tracking, commit creation, and PR management.
package git

import (
	"sync"
	"time"
)

// FileModification represents a file that was modified during the session.
type FileModification struct {
	Path      string    // Relative path from workspace root
	Operation string    // "write", "diff", "delete"
	Timestamp time.Time // When the modification occurred
}

// ModificationTracker tracks file modifications made during an agent session.
type ModificationTracker struct {
	mu            sync.RWMutex
	modifications map[string]*FileModification
}

// NewModificationTracker creates a new file modification tracker.
func NewModificationTracker() *ModificationTracker {
	return &ModificationTracker{
		modifications: make(map[string]*FileModification),
	}
}

// Track records a file modification.
func (t *ModificationTracker) Track(path, operation string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.modifications[path] = &FileModification{
		Path:      path,
		Operation: operation,
		Timestamp: time.Now(),
	}
}

// GetModified returns all modified file paths.
func (t *ModificationTracker) GetModified() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	paths := make([]string, 0, len(t.modifications))
	for path := range t.modifications {
		paths = append(paths, path)
	}
	return paths
}

// GetModifications returns all file modifications with their metadata.
func (t *ModificationTracker) GetModifications() []*FileModification {
	t.mu.RLock()
	defer t.mu.RUnlock()

	mods := make([]*FileModification, 0, len(t.modifications))
	for _, mod := range t.modifications {
		mods = append(mods, mod)
	}
	return mods
}

// Clear resets the tracker, removing all tracked modifications.
func (t *ModificationTracker) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.modifications = make(map[string]*FileModification)
}

// Count returns the number of tracked modifications.
func (t *ModificationTracker) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return len(t.modifications)
}
