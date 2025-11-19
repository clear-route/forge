package overlay

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/executor/tui/types"
)

// resultListItem represents a single item in the result list
type resultListItem struct {
	result *types.CachedResult
}

func (i resultListItem) FilterValue() string {
	return i.result.ToolName
}

func (i resultListItem) Title() string {
	return fmt.Sprintf("%s - %s", i.result.ToolName, i.result.Timestamp.Format("15:04:05"))
}

func (i resultListItem) Description() string {
	// Truncate summary if too long (77 chars + 3 for "..." = 80 total)
	summary := i.result.Summary
	if len(summary) > 77 {
		summary = summary[:77] + "..."
	}
	return summary
}

// resultListDelegate is a custom delegate for rendering result list items
type resultListDelegate struct {
	list.DefaultDelegate
}

func newResultListDelegate() resultListDelegate {
	d := list.NewDefaultDelegate()

	// Customize styles to match Forge theme using types constants
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(types.SalmonPink).
		BorderForeground(types.SalmonPink)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(types.MutedGray).
		BorderForeground(types.SalmonPink)

	return resultListDelegate{DefaultDelegate: d}
}

// ResultListModel represents the state of the result list overlay
type ResultListModel struct {
	list     list.Model
	width    int
	height   int
	active   bool
	quitting bool
}

// NewResultListModel creates a new result list model
func NewResultListModel() ResultListModel {
	delegate := newResultListDelegate()

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Tool Result History"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)

	l.Styles.Title = lipgloss.NewStyle().
		Foreground(types.SalmonPink).
		Bold(true).
		Padding(0, 1)

	// Add custom key bindings
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "view result"),
			),
			key.NewBinding(
				key.WithKeys("esc", "q"),
				key.WithHelp("esc/q", "close"),
			),
		}
	}

	return ResultListModel{
		list:   l,
		active: false,
	}
}

// Activate shows the result list with the given cached results
func (m *ResultListModel) Activate(results []*types.CachedResult, width, height int) {
	m.active = true
	m.width = width
	m.height = height

	// Convert cached results to list items
	items := make([]list.Item, len(results))
	for i, result := range results {
		items[i] = resultListItem{result: result}
	}

	m.list.SetItems(items)
	m.list.SetSize(width-4, height-4)
}

// Deactivate hides the result list
func (m *ResultListModel) Deactivate() {
	m.active = false
	m.quitting = false
}

// IsActive returns whether the result list is active
func (m *ResultListModel) IsActive() bool {
	return m.active
}

// Init implements tea.Model
func (m *ResultListModel) Init() tea.Cmd {
	return nil
}

func (m *ResultListModel) Update(msg tea.Msg, state types.StateProvider, actions types.ActionHandler) (types.Overlay, tea.Cmd) {
	if !m.active {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", keyEsc:
			m.Deactivate()
			return m, nil
		case keyEnter:
			// Return the selected result
			if item, ok := m.list.SelectedItem().(resultListItem); ok {
				// Signal that we want to view this result
				m.quitting = true
				return m, func() tea.Msg {
					return types.ViewResultMsg{ResultID: item.result.ID}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-4, msg.Height-4)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the result list
func (m *ResultListModel) View() string {
	if !m.active {
		return ""
	}

	// Create a bordered container for the list
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(types.SalmonPink).
		Padding(1, 2).
		Width(m.width - 4).
		Height(m.height - 4)

	return boxStyle.Render(m.list.View())
}

// Focused returns whether the result list should handle input
func (m *ResultListModel) Focused() bool {
	return m.active
}

// Width returns the width of the result list
func (m *ResultListModel) Width() int {
	return m.width
}

// Height returns the height of the result list
func (m *ResultListModel) Height() int {
	return m.height
}
