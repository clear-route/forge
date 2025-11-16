package config

import (
	"fmt"
	"sync"
)

// Manager coordinates all configuration sections and handles persistence.
type Manager struct {
	store    Store
	sections map[string]Section
	order    []string
	mu       sync.RWMutex
}

// NewManager creates a new configuration manager.
func NewManager(store Store) *Manager {
	return &Manager{
		store:    store,
		sections: make(map[string]Section),
		order:    make([]string, 0),
	}
}

// RegisterSection adds a new configuration section.
// Sections are displayed in the order they are registered.
func (m *Manager) RegisterSection(section Section) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := section.ID()
	if _, exists := m.sections[id]; exists {
		return fmt.Errorf("section '%s' is already registered", id)
	}

	m.sections[id] = section
	m.order = append(m.order, id)
	return nil
}

// GetSection retrieves a section by ID.
func (m *Manager) GetSection(id string) (Section, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	section, exists := m.sections[id]
	return section, exists
}

// GetSections returns all registered sections in registration order.
func (m *Manager) GetSections() []Section {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sections := make([]Section, 0, len(m.order))
	for _, id := range m.order {
		if section, exists := m.sections[id]; exists {
			sections = append(sections, section)
		}
	}
	return sections
}

// LoadAll loads configuration for all registered sections from the store.
func (m *Manager) LoadAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// First, load data from store
	if err := m.store.Load(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Then, populate each section with its data
	for id, section := range m.sections {
		data, err := m.store.GetSection(id)
		if err != nil {
			return fmt.Errorf("failed to get section '%s': %w", id, err)
		}

		if err := section.SetData(data); err != nil {
			return fmt.Errorf("failed to set data for section '%s': %w", id, err)
		}
	}

	return nil
}

// SaveAll saves configuration from all registered sections to the store.
func (m *Manager) SaveAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Validate all sections first
	for id, section := range m.sections {
		if err := section.Validate(); err != nil {
			return fmt.Errorf("validation failed for section '%s': %w", id, err)
		}
	}

	// Save each section's data to the store
	for id, section := range m.sections {
		data := section.Data()
		if err := m.store.SetSection(id, data); err != nil {
			return fmt.Errorf("failed to save section '%s': %w", id, err)
		}
	}

	// Persist to disk
	if err := m.store.Save(); err != nil {
		return fmt.Errorf("failed to save configuration to disk: %w", err)
	}

	return nil
}

// ResetAll resets all sections to their default configuration.
func (m *Manager) ResetAll() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, section := range m.sections {
		section.Reset()
	}
}

// Store returns the underlying configuration store.
func (m *Manager) Store() Store {
	return m.store
}