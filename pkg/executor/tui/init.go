package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// initialModel returns the initial state of the TUI.
// It creates and configures all components needed for the interactive interface.
func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Type your message..."
	ta.Focus()
	ta.Prompt = "> "
	ta.CharLimit = 0
	ta.SetHeight(1)
	ta.MaxHeight = 10 // Allow up to 10 lines
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false) // Disable default Enter behavior
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(salmonPink)
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(brightWhite)

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().Padding(0, 2)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(salmonPink)

	return model{
		viewport:         vp,
		textarea:         ta,
		content:          &strings.Builder{},
		thinkingBuffer:   &strings.Builder{},
		messageBuffer:    &strings.Builder{},
		overlay:          newOverlayState(),
		commandPalette:   newCommandPalette(),
		summarization:    &summarizationStatus{},
		toast:            &toastNotification{},
		spinner:          s,
		agentBusy:        false,
		resultClassifier: NewToolResultClassifier(),
		resultSummarizer: NewToolResultSummarizer(),
		resultCache:      newResultCache(20),
		resultList:       newResultListModel(),
	}
}

// Init is the first function that will be called by Bubble Tea.
// It returns commands to start the textarea blink animation and spinner.
func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.spinner.Tick)
}
