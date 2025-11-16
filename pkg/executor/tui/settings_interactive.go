package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/config"
)

// InteractiveSettingsOverlay provides a full interactive settings editor
type InteractiveSettingsOverlay struct {
	width   int
	height  int
	focused bool

	// Navigation state
	selectedSection int
	selectedItem    int
	sections        []settingsSection

	// Edit state
	hasChanges bool
	editMode   bool

	// Scroll state
	scrollOffset int

	// Dialog state
	activeDialog      *inputDialog
	confirmDialog     *confirmDialog
}

// settingsSection represents a section with its items
type settingsSection struct {
	id          string
	title       string
	description string
	items       []settingsItem
}

// settingsItem represents an editable item
type settingsItem struct {
	key         string
	displayName string
	value       interface{}
	itemType    itemType
	modified    bool
}

// itemType defines the type of setting item
type itemType int

const (
	itemTypeToggle itemType = iota
	itemTypeText
	itemTypeList
)

// inputDialog represents a modal dialog for text input
type inputDialog struct {
	title         string
	fields        []inputField
	selectedField int
	onConfirm     func(values map[string]string) error
	onCancel      func()
}

// inputField represents a single input field in a dialog
type inputField struct {
	label      string
	key        string
	value      string
	fieldType  fieldType
	options    []string // For radio buttons
	maxLength  int
	validator  func(string) error
	errorMsg   string
}

// fieldType defines the type of input field
type fieldType int

const (
	fieldTypeText fieldType = iota
	fieldTypeRadio
)

// confirmDialog represents a confirmation dialog
type confirmDialog struct {
	title   string
	message string
	details []string
	onYes   func()
	onNo    func()
}

// NewInteractiveSettingsOverlay creates a new interactive settings overlay
func NewInteractiveSettingsOverlay(width, height int) *InteractiveSettingsOverlay {
	overlay := &InteractiveSettingsOverlay{
		width:           width,
		height:          height,
		focused:         true,
		selectedSection: 0,
		selectedItem:    0,
		hasChanges:      false,
		editMode:        false,
		scrollOffset:    0,
	}

	overlay.loadSettings()
	return overlay
}

// loadSettings loads settings from config into editable sections
func (s *InteractiveSettingsOverlay) loadSettings() {
	if !config.IsInitialized() {
		return
	}

	manager := config.Global()
	configSections := manager.GetSections()

	s.sections = make([]settingsSection, 0, len(configSections))

	for _, sec := range configSections {
		section := settingsSection{
			id:          sec.ID(),
			title:       sec.Title(),
			description: sec.Description(),
			items:       make([]settingsItem, 0),
		}

		data := sec.Data()

		switch sec.ID() {
		case "auto_approval":
			// Create toggle items for each tool
			keys := make([]string, 0, len(data))
			for k := range data {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, key := range keys {
				value := data[key]
				item := settingsItem{
					key:         key,
					displayName: key,
					value:       value,
					itemType:    itemTypeToggle,
					modified:    false,
				}
				section.items = append(section.items, item)
			}

		case "command_whitelist":
			// Show patterns as list items
			if patterns, ok := data["patterns"].([]interface{}); ok {
				for i, p := range patterns {
					if patternMap, ok := p.(map[string]interface{}); ok {
						pattern := patternMap["pattern"]
						desc := patternMap["description"]
						displayName := fmt.Sprintf("%v", pattern)
						if desc != nil && desc != "" {
							displayName = fmt.Sprintf("%v - %v", pattern, desc)
						}

						item := settingsItem{
							key:         fmt.Sprintf("pattern_%d", i),
							displayName: displayName,
							value:       patternMap,
							itemType:    itemTypeList,
							modified:    false,
						}
						section.items = append(section.items, item)
					}
				}
			}
		}

		s.sections = append(s.sections, section)
	}
}

// Update handles messages for the interactive settings overlay
func (s *InteractiveSettingsOverlay) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	// Handle active dialog input first
	if s.activeDialog != nil {
		return s.handleDialogInput(msg)
	}

	// Handle confirmation dialog
	if s.confirmDialog != nil {
		return s.handleConfirmInput(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			if s.hasChanges {
				// Show confirmation dialog
				s.showUnsavedChangesDialog()
				return s, nil
			}
			return nil, nil

		case "ctrl+s":
			// Save changes
			if s.hasChanges {
				if err := s.saveSettings(); err == nil {
					s.hasChanges = false
				}
			}
			return s, nil

		case "up", "k":
			s.navigateUp()
			return s, nil

		case "down", "j":
			s.navigateDown()
			return s, nil

		case "left", "h":
			s.navigateLeft()
			return s, nil

		case "right", "l":
			s.navigateRight()
			return s, nil

		case " ", "enter":
			s.toggleCurrent()
			return s, nil

		case "tab":
			s.nextSection()
			return s, nil

		case "shift+tab":
			s.prevSection()
			return s, nil

		case "a":
			// Add new pattern (only in whitelist section)
			if s.isInWhitelistSection() {
				s.showAddPatternDialog()
			}
			return s, nil

		case "e":
			// Edit selected pattern
			if s.isInWhitelistSection() && s.isPatternSelected() {
				s.showEditPatternDialog()
			}
			return s, nil

		case "d":
			// Delete selected pattern
			if s.isInWhitelistSection() && s.isPatternSelected() {
				s.showDeleteConfirmation()
			}
			return s, nil
		}
	}

	return s, nil
}

// navigateUp moves selection up
func (s *InteractiveSettingsOverlay) navigateUp() {
	if len(s.sections) == 0 {
		return
	}

	if s.selectedItem > 0 {
		s.selectedItem--
	} else if s.selectedSection > 0 {
		s.selectedSection--
		if len(s.sections[s.selectedSection].items) > 0 {
			s.selectedItem = len(s.sections[s.selectedSection].items) - 1
		}
	}
}

// navigateDown moves selection down
func (s *InteractiveSettingsOverlay) navigateDown() {
	if len(s.sections) == 0 {
		return
	}

	currentSection := s.sections[s.selectedSection]
	if s.selectedItem < len(currentSection.items)-1 {
		s.selectedItem++
	} else if s.selectedSection < len(s.sections)-1 {
		s.selectedSection++
		s.selectedItem = 0
	}
}

// navigateLeft moves to previous section
func (s *InteractiveSettingsOverlay) navigateLeft() {
	s.prevSection()
}

// navigateRight moves to next section
func (s *InteractiveSettingsOverlay) navigateRight() {
	s.nextSection()
}

// nextSection moves to the next section
func (s *InteractiveSettingsOverlay) nextSection() {
	if s.selectedSection < len(s.sections)-1 {
		s.selectedSection++
		s.selectedItem = 0
	}
}

// prevSection moves to the previous section
func (s *InteractiveSettingsOverlay) prevSection() {
	if s.selectedSection > 0 {
		s.selectedSection--
		s.selectedItem = 0
	}
}

// toggleCurrent toggles the current item
func (s *InteractiveSettingsOverlay) toggleCurrent() {
	if len(s.sections) == 0 {
		return
	}

	section := &s.sections[s.selectedSection]
	if s.selectedItem >= len(section.items) {
		return
	}

	item := &section.items[s.selectedItem]
	if item.itemType == itemTypeToggle {
		if boolVal, ok := item.value.(bool); ok {
			item.value = !boolVal
			item.modified = true
			s.hasChanges = true
		}
	}
}

// saveSettings saves changes back to config
func (s *InteractiveSettingsOverlay) saveSettings() error {
	if !config.IsInitialized() {
		return fmt.Errorf("config not initialized")
	}

	manager := config.Global()

	for _, section := range s.sections {
		configSection, exists := manager.GetSection(section.id)
		if !exists {
			continue
		}

		// Build updated data map
		data := make(map[string]interface{})

		switch section.id {
		case "auto_approval":
			for _, item := range section.items {
				data[item.key] = item.value
			}

		case "command_whitelist":
			// Reconstruct patterns array
			patterns := make([]interface{}, 0)
			for _, item := range section.items {
				if item.itemType == itemTypeList {
					patterns = append(patterns, item.value)
				}
			}
			data["patterns"] = patterns
		}

		// Update section
		if err := configSection.SetData(data); err != nil {
			return fmt.Errorf("failed to update section %s: %w", section.id, err)
		}
	}

	// Save all changes
	return manager.SaveAll()
}

// View renders the interactive settings overlay
func (s *InteractiveSettingsOverlay) View() string {
	if !config.IsInitialized() {
		return s.renderError("Configuration not initialized")
	}

	// If dialog is active, render it on top
	if s.activeDialog != nil {
		return s.renderWithDialog()
	}

	// If confirmation dialog is active, render it on top
	if s.confirmDialog != nil {
		return s.renderWithConfirmation()
	}

	var content strings.Builder

	// Title
	title := OverlayTitleStyle.Render("⚙️  Settings")
	content.WriteString(title)
	content.WriteString("\n\n")

	// Help text
	helpText := s.buildHelpText()
	content.WriteString(OverlaySubtitleStyle.Render(helpText))
	content.WriteString("\n\n")

	// Render sections
	for i, section := range s.sections {
		if i > 0 {
			content.WriteString("\n")
		}
		content.WriteString(s.renderSection(section, i == s.selectedSection))
	}

	// Status bar
	if s.hasChanges {
		content.WriteString("\n\n")
		saveHint := lipgloss.NewStyle().
			Foreground(salmonPink).
			Bold(true).
			Render("● Unsaved changes - Press Ctrl+S to save")
		content.WriteString(saveHint)
	}

	// Create bordered box
	boxStyle := CreateOverlayContainerStyle(s.width - 4).Height(s.height - 4)

	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		boxStyle.Render(content.String()),
	)
}

// buildHelpText creates the help text based on current state
func (s *InteractiveSettingsOverlay) buildHelpText() string {
	shortcuts := []string{
		"↑↓/jk: Navigate",
		"Tab/←→/hl: Switch section",
		"Space/Enter: Toggle",
	}

	// Add whitelist-specific shortcuts if in that section
	if s.isInWhitelistSection() {
		shortcuts = append(shortcuts, "a: Add", "e: Edit", "d: Delete")
	}

	shortcuts = append(shortcuts, "Ctrl+S: Save", "Esc/q: Close")
	return strings.Join(shortcuts, " • ")
}

// renderSection renders a settings section
func (s *InteractiveSettingsOverlay) renderSection(section settingsSection, isSelected bool) string {
	var out strings.Builder

	// Section title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(mintGreen)

	if isSelected {
		titleStyle = titleStyle.Foreground(salmonPink)
	}

	out.WriteString(titleStyle.Render("▸ " + section.title))
	out.WriteString("\n")

	// Section description
	if section.description != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(mutedGray).
			Italic(true)
		out.WriteString("  ")
		out.WriteString(descStyle.Render(section.description))
		out.WriteString("\n")
	}

	// Render items
	if len(section.items) == 0 {
		out.WriteString("  ")
		out.WriteString(lipgloss.NewStyle().Foreground(mutedGray).Render("(no settings)"))
		out.WriteString("\n")
	} else {
		for i, item := range section.items {
			isItemSelected := isSelected && i == s.selectedItem
			out.WriteString(s.renderItem(item, isItemSelected))
		}
	}

	return out.String()
}

// renderItem renders a single settings item
func (s *InteractiveSettingsOverlay) renderItem(item settingsItem, isSelected bool) string {
	var out strings.Builder
	out.WriteString("  ")

	// Selection indicator
	if isSelected {
		out.WriteString(lipgloss.NewStyle().Foreground(salmonPink).Render("▸ "))
	} else {
		out.WriteString("  ")
	}

	// Render based on type
	switch item.itemType {
	case itemTypeToggle:
		out.WriteString(s.renderToggle(item, isSelected))
	case itemTypeList:
		out.WriteString(s.renderListItem(item, isSelected))
	default:
		out.WriteString(item.displayName)
	}

	// Modified indicator
	if item.modified {
		out.WriteString(" ")
		out.WriteString(lipgloss.NewStyle().Foreground(salmonPink).Render("*"))
	}

	out.WriteString("\n")
	return out.String()
}

// Helper methods for dialog and pattern management

// isInWhitelistSection returns true if currently in command whitelist section
func (s *InteractiveSettingsOverlay) isInWhitelistSection() bool {
	if len(s.sections) == 0 || s.selectedSection >= len(s.sections) {
		return false
	}
	return s.sections[s.selectedSection].id == "command_whitelist"
}

// isPatternSelected returns true if a pattern item is currently selected
func (s *InteractiveSettingsOverlay) isPatternSelected() bool {
	if len(s.sections) == 0 || s.selectedSection >= len(s.sections) {
		return false
	}
	section := s.sections[s.selectedSection]
	if s.selectedItem >= len(section.items) {
		return false
	}
	return section.items[s.selectedItem].itemType == itemTypeList
}

// showAddPatternDialog displays the dialog for adding a new command pattern
func (s *InteractiveSettingsOverlay) showAddPatternDialog() {
	s.activeDialog = &inputDialog{
		title: "Add New Command Pattern",
		fields: []inputField{
			{
				label:     "Pattern (command or prefix):",
				key:       "pattern",
				value:     "",
				fieldType: fieldTypeText,
				maxLength: 100,
				validator: func(v string) error {
					if strings.TrimSpace(v) == "" {
						return fmt.Errorf("pattern cannot be empty")
					}
					return nil
				},
			},
			{
				label:     "Description:",
				key:       "description",
				value:     "",
				fieldType: fieldTypeText,
				maxLength: 100,
			},
			{
				label:     "Pattern type:",
				key:       "type",
				value:     "prefix",
				fieldType: fieldTypeRadio,
				options:   []string{"prefix", "exact"},
			},
		},
		selectedField: 0,
		onConfirm: func(values map[string]string) error {
			return s.addPattern(values)
		},
		onCancel: func() {
			s.activeDialog = nil
		},
	}
}

// showEditPatternDialog displays the dialog for editing an existing pattern
func (s *InteractiveSettingsOverlay) showEditPatternDialog() {
	if !s.isPatternSelected() {
		return
	}

	section := &s.sections[s.selectedSection]
	item := section.items[s.selectedItem]
	
	patternMap, ok := item.value.(map[string]interface{})
	if !ok {
		return
	}

	pattern := fmt.Sprintf("%v", patternMap["pattern"])
	description := fmt.Sprintf("%v", patternMap["description"])
	patternType := "prefix"
	if t, ok := patternMap["type"].(string); ok && t == "exact" {
		patternType = "exact"
	}

	s.activeDialog = &inputDialog{
		title: "Edit Command Pattern",
		fields: []inputField{
			{
				label:     "Pattern (command or prefix):",
				key:       "pattern",
				value:     pattern,
				fieldType: fieldTypeText,
				maxLength: 100,
				validator: func(v string) error {
					if strings.TrimSpace(v) == "" {
						return fmt.Errorf("pattern cannot be empty")
					}
					return nil
				},
			},
			{
				label:     "Description:",
				key:       "description",
				value:     description,
				fieldType: fieldTypeText,
				maxLength: 100,
			},
			{
				label:     "Pattern type:",
				key:       "type",
				value:     patternType,
				fieldType: fieldTypeRadio,
				options:   []string{"prefix", "exact"},
			},
		},
		selectedField: 0,
		onConfirm: func(values map[string]string) error {
			return s.updatePattern(s.selectedItem, values)
		},
		onCancel: func() {
			s.activeDialog = nil
		},
	}
}

// showDeleteConfirmation displays confirmation dialog for deleting a pattern
func (s *InteractiveSettingsOverlay) showDeleteConfirmation() {
	if !s.isPatternSelected() {
		return
	}

	section := &s.sections[s.selectedSection]
	item := section.items[s.selectedItem]
	
	patternMap, ok := item.value.(map[string]interface{})
	if !ok {
		return
	}

	pattern := fmt.Sprintf("%v", patternMap["pattern"])
	description := fmt.Sprintf("%v", patternMap["description"])

	s.confirmDialog = &confirmDialog{
		title:   "Confirm Delete",
		message: "⚠️  Are you sure you want to delete this pattern?",
		details: []string{
			fmt.Sprintf("Pattern: %s", pattern),
			fmt.Sprintf("Description: %s", description),
			"",
			"This command will require manual approval after deletion.",
		},
		onYes: func() {
			s.deletePattern(s.selectedItem)
			s.confirmDialog = nil
		},
		onNo: func() {
			s.confirmDialog = nil
		},
	}
}

// showUnsavedChangesDialog displays confirmation for closing with unsaved changes
func (s *InteractiveSettingsOverlay) showUnsavedChangesDialog() {
	s.confirmDialog = &confirmDialog{
		title:   "Unsaved Changes",
		message: "⚠️  You have unsaved changes.",
		details: []string{
			"",
			"Do you want to save before closing?",
		},
		onYes: func() {
			if err := s.saveSettings(); err == nil {
				s.hasChanges = false
			}
			s.confirmDialog = nil
			// Close the overlay by returning nil in next update
		},
		onNo: func() {
			s.confirmDialog = nil
			// Discard changes and close
		},
	}
}

// addPattern adds a new command pattern
func (s *InteractiveSettingsOverlay) addPattern(values map[string]string) error {
	pattern := strings.TrimSpace(values["pattern"])
	description := strings.TrimSpace(values["description"])
	patternType := values["type"]

	if pattern == "" {
		return fmt.Errorf("pattern cannot be empty")
	}

	// Create new pattern map
	newPattern := map[string]interface{}{
		"pattern":     pattern,
		"description": description,
		"type":        patternType,
	}

	// Add to whitelist section
	section := &s.sections[s.selectedSection]
	displayName := pattern
	if description != "" {
		displayName = fmt.Sprintf("%s - %s", pattern, description)
	}

	newItem := settingsItem{
		key:         fmt.Sprintf("pattern_%d", len(section.items)),
		displayName: displayName,
		value:       newPattern,
		itemType:    itemTypeList,
		modified:    true,
	}

	section.items = append(section.items, newItem)
	s.hasChanges = true
	s.activeDialog = nil

	return nil
}

// updatePattern updates an existing command pattern
func (s *InteractiveSettingsOverlay) updatePattern(index int, values map[string]string) error {
	pattern := strings.TrimSpace(values["pattern"])
	description := strings.TrimSpace(values["description"])
	patternType := values["type"]

	if pattern == "" {
		return fmt.Errorf("pattern cannot be empty")
	}

	section := &s.sections[s.selectedSection]
	if index >= len(section.items) {
		return fmt.Errorf("invalid pattern index")
	}

	// Update pattern map
	updatedPattern := map[string]interface{}{
		"pattern":     pattern,
		"description": description,
		"type":        patternType,
	}

	displayName := pattern
	if description != "" {
		displayName = fmt.Sprintf("%s - %s", pattern, description)
	}

	section.items[index].value = updatedPattern
	section.items[index].displayName = displayName
	section.items[index].modified = true

	s.hasChanges = true
	s.activeDialog = nil

	return nil
}

// deletePattern removes a command pattern
func (s *InteractiveSettingsOverlay) deletePattern(index int) {
	section := &s.sections[s.selectedSection]
	if index >= len(section.items) {
		return
	}

	// Remove the item
	section.items = append(section.items[:index], section.items[index+1:]...)
	
	// Adjust selected item if needed
	if s.selectedItem >= len(section.items) && s.selectedItem > 0 {
		s.selectedItem = len(section.items) - 1
	}

	s.hasChanges = true
}

// handleDialogInput handles keyboard input for the active dialog
func (s *InteractiveSettingsOverlay) handleDialogInput(msg tea.Msg) (Overlay, tea.Cmd) {
	if s.activeDialog == nil {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if s.activeDialog.onCancel != nil {
				s.activeDialog.onCancel()
			}
			return s, nil

		case "enter":
			// Validate and confirm
			values := make(map[string]string)
			for _, field := range s.activeDialog.fields {
				values[field.key] = field.value
				
				// Run validator if present
				if field.validator != nil {
					if err := field.validator(field.value); err != nil {
						field.errorMsg = err.Error()
						return s, nil
					}
				}
			}

			// Call onConfirm
			if s.activeDialog.onConfirm != nil {
				if err := s.activeDialog.onConfirm(values); err != nil {
					// Show error
					return s, nil
				}
			}
			return s, nil

		case "tab", "down":
			// Move to next field
			s.activeDialog.selectedField++
			if s.activeDialog.selectedField >= len(s.activeDialog.fields) {
				s.activeDialog.selectedField = 0
			}
			return s, nil

		case "shift+tab", "up":
			// Move to previous field
			s.activeDialog.selectedField--
			if s.activeDialog.selectedField < 0 {
				s.activeDialog.selectedField = len(s.activeDialog.fields) - 1
			}
			return s, nil

		case " ":
			// Handle space for both radio buttons and text fields
			field := &s.activeDialog.fields[s.activeDialog.selectedField]
			if field.fieldType == fieldTypeRadio && len(field.options) > 0 {
				// Toggle radio button
				currentIdx := 0
				for i, opt := range field.options {
					if opt == field.value {
						currentIdx = i
						break
					}
				}
				nextIdx := (currentIdx + 1) % len(field.options)
				field.value = field.options[nextIdx]
			} else if field.fieldType == fieldTypeText {
				// Add space to text field
				if field.maxLength == 0 || len(field.value) < field.maxLength {
					field.value += " "
					field.errorMsg = ""
				}
			}
			return s, nil

		case "backspace":
			// Delete character from text field
			field := &s.activeDialog.fields[s.activeDialog.selectedField]
			if field.fieldType == fieldTypeText && len(field.value) > 0 {
				field.value = field.value[:len(field.value)-1]
				field.errorMsg = ""
			}
			return s, nil

		default:
			// Add character to text field
			field := &s.activeDialog.fields[s.activeDialog.selectedField]
			if field.fieldType == fieldTypeText {
				if len(msg.String()) == 1 && (field.maxLength == 0 || len(field.value) < field.maxLength) {
					field.value += msg.String()
					field.errorMsg = ""
				}
			}
			return s, nil
		}
	}

	return s, nil
}

// handleConfirmInput handles keyboard input for confirmation dialogs
func (s *InteractiveSettingsOverlay) handleConfirmInput(msg tea.Msg) (Overlay, tea.Cmd) {
	if s.confirmDialog == nil {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			if s.confirmDialog.onYes != nil {
				s.confirmDialog.onYes()
			}
			// If this was unsaved changes dialog and user said yes, close overlay
			if s.confirmDialog == nil && !s.hasChanges {
				return nil, nil
			}
			return s, nil

		case "n", "N", "esc":
			if s.confirmDialog.onNo != nil {
				s.confirmDialog.onNo()
			}
			return s, nil
		}
	}

	return s, nil
}

// renderWithDialog renders the settings view with an input dialog overlay
func (s *InteractiveSettingsOverlay) renderWithDialog() string {
	// Render base settings view (dimmed)
	baseView := s.renderBaseView(true)

	// Render dialog on top
	dialogView := s.renderInputDialog()

	// Layer dialog over base view
	return s.layerDialogOver(baseView, dialogView)
}

// renderWithConfirmation renders the settings view with a confirmation dialog overlay
func (s *InteractiveSettingsOverlay) renderWithConfirmation() string {
	// Render base settings view (dimmed)
	baseView := s.renderBaseView(true)

	// Render confirmation dialog on top
	dialogView := s.renderConfirmDialog()

	// Layer dialog over base view
	return s.layerDialogOver(baseView, dialogView)
}

// renderBaseView renders the main settings view
func (s *InteractiveSettingsOverlay) renderBaseView(dimmed bool) string {
	var content strings.Builder

	// Title
	titleStyle := OverlayTitleStyle
	if dimmed {
		titleStyle = titleStyle.Foreground(mutedGray)
	}
	title := titleStyle.Render("⚙️  Settings")
	content.WriteString(title)
	content.WriteString("\n\n")

	// Help text
	helpText := s.buildHelpText()
	helpStyle := OverlaySubtitleStyle
	if dimmed {
		helpStyle = helpStyle.Foreground(mutedGray)
	}
	content.WriteString(helpStyle.Render(helpText))
	content.WriteString("\n\n")

	// Render sections
	for i, section := range s.sections {
		if i > 0 {
			content.WriteString("\n")
		}
		sectionView := s.renderSection(section, i == s.selectedSection)
		if dimmed {
			// Apply dimmed style to section
			sectionView = lipgloss.NewStyle().Foreground(mutedGray).Render(sectionView)
		}
		content.WriteString(sectionView)
	}

	// Status bar
	if s.hasChanges {
		content.WriteString("\n\n")
		saveHint := lipgloss.NewStyle().
			Foreground(salmonPink).
			Bold(true).
			Render("● Unsaved changes - Press Ctrl+S to save")
		content.WriteString(saveHint)
	}

	// Create bordered box
	boxStyle := CreateOverlayContainerStyle(s.width - 4).Height(s.height - 4)

	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		boxStyle.Render(content.String()),
	)
}

// renderInputDialog renders an input dialog
func (s *InteractiveSettingsOverlay) renderInputDialog() string {
	if s.activeDialog == nil {
		return ""
	}

	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(salmonPink).
		Bold(true)
	content.WriteString(titleStyle.Render(s.activeDialog.title))
	content.WriteString("\n\n")

	// Fields
	for i, field := range s.activeDialog.fields {
		isSelected := i == s.activeDialog.selectedField

		// Label
		labelStyle := lipgloss.NewStyle().Foreground(brightWhite)
		content.WriteString(labelStyle.Render(field.label))
		content.WriteString("\n")

		// Field content based on type
		switch field.fieldType {
		case fieldTypeText:
			// Text input field
			fieldStyle := lipgloss.NewStyle().
				Foreground(brightWhite).
				Background(darkBg).
				Padding(0, 1)

			if isSelected {
				fieldStyle = fieldStyle.Border(lipgloss.RoundedBorder()).
					BorderForeground(salmonPink)
			}

			value := field.value
			if isSelected {
				value += "▸" // Cursor
			}

			content.WriteString(fieldStyle.Render(value))
			content.WriteString("\n")

			// Character count or error
			if field.errorMsg != "" {
				errorStyle := lipgloss.NewStyle().Foreground(salmonPink)
				content.WriteString(errorStyle.Render(field.errorMsg))
				content.WriteString("\n")
			} else if field.maxLength > 0 {
				countStyle := lipgloss.NewStyle().Foreground(mutedGray)
				count := fmt.Sprintf("[%d/%d]", len(field.value), field.maxLength)
				content.WriteString(countStyle.Render(count))
				content.WriteString("\n")
			}

		case fieldTypeRadio:
			// Radio buttons
			for j, option := range field.options {
				radioStyle := lipgloss.NewStyle().Foreground(brightWhite)
				if isSelected {
					radioStyle = radioStyle.Bold(true)
				}

				bullet := "○"
				if option == field.value {
					bullet = "●"
					radioStyle = radioStyle.Foreground(mintGreen)
				}

				optionText := fmt.Sprintf("%s %s", bullet, option)
				if j > 0 {
					content.WriteString("    ") // Indent
				}
				content.WriteString(radioStyle.Render(optionText))
				if j < len(field.options)-1 {
					content.WriteString("    ")
				}
			}
			content.WriteString("\n")
		}

		content.WriteString("\n")
	}

	// Buttons
	buttonRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		"[Enter to Add] [Esc to Cancel]",
	)
	buttonStyle := lipgloss.NewStyle().Foreground(mutedGray)
	content.WriteString(buttonStyle.Render(buttonRow))

	// Create dialog box
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(salmonPink).
		Background(darkBg).
		Padding(1, 2).
		Width(60)

	return dialogStyle.Render(content.String())
}

// renderConfirmDialog renders a confirmation dialog
func (s *InteractiveSettingsOverlay) renderConfirmDialog() string {
	if s.confirmDialog == nil {
		return ""
	}

	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(salmonPink).
		Bold(true)
	content.WriteString(titleStyle.Render(s.confirmDialog.title))
	content.WriteString("\n\n")

	// Message
	messageStyle := lipgloss.NewStyle().Foreground(brightWhite)
	content.WriteString(messageStyle.Render(s.confirmDialog.message))
	content.WriteString("\n\n")

	// Details
	detailStyle := lipgloss.NewStyle().Foreground(mutedGray)
	for _, detail := range s.confirmDialog.details {
		content.WriteString(detailStyle.Render(detail))
		content.WriteString("\n")
	}

	content.WriteString("\n")

	// Buttons
	buttonRow := "[y] Yes, delete    [n] No, cancel"
	buttonStyle := lipgloss.NewStyle().Foreground(mutedGray)
	content.WriteString(buttonStyle.Render(buttonRow))

	// Create dialog box
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(salmonPink).
		Background(darkBg).
		Padding(1, 2).
		Width(60)

	return dialogStyle.Render(content.String())
}

// layerDialogOver layers a dialog over the base view
func (s *InteractiveSettingsOverlay) layerDialogOver(baseView, dialogView string) string {
	// Place dialog in center
	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		dialogView,
		lipgloss.WithWhitespaceChars(""),
		lipgloss.WithWhitespaceForeground(lipgloss.NoColor{}),
	)
}

// renderToggle renders a toggle item
func (s *InteractiveSettingsOverlay) renderToggle(item settingsItem, isSelected bool) string {
	boolVal, ok := item.value.(bool)
	if !ok {
		boolVal = false
	}

	// Toggle indicator
	var toggle string
	if boolVal {
		toggle = lipgloss.NewStyle().Foreground(mintGreen).Render("[✓]")
	} else {
		toggle = lipgloss.NewStyle().Foreground(mutedGray).Render("[ ]")
	}

	// Item name
	nameStyle := lipgloss.NewStyle().Foreground(brightWhite)
	if isSelected {
		nameStyle = nameStyle.Bold(true)
	}

	return fmt.Sprintf("%s %s", toggle, nameStyle.Render(item.displayName))
}

// renderListItem renders a list item
func (s *InteractiveSettingsOverlay) renderListItem(item settingsItem, isSelected bool) string {
	nameStyle := lipgloss.NewStyle().Foreground(brightWhite)
	if isSelected {
		nameStyle = nameStyle.Bold(true)
	}

	bullet := lipgloss.NewStyle().Foreground(mintGreen).Render("✓")
	return fmt.Sprintf("%s %s", bullet, nameStyle.Render(item.displayName))
}

// renderError renders an error message
func (s *InteractiveSettingsOverlay) renderError(message string) string {
	errorStyle := lipgloss.NewStyle().
		Foreground(salmonPink).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(salmonPink).
		Background(darkBg).
		Padding(1, 2).
		Width(s.width - 4).
		Height(s.height - 4)

	content := errorStyle.Render("Error: ") + message

	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		boxStyle.Render(content),
	)
}

// Focused returns whether this overlay should handle input
func (s *InteractiveSettingsOverlay) Focused() bool {
	return s.focused
}

// Width returns the overlay width
func (s *InteractiveSettingsOverlay) Width() int {
	return s.width
}

// Height returns the overlay height
func (s *InteractiveSettingsOverlay) Height() int {
	return s.height
}