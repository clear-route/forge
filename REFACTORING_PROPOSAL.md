# Forge Refactoring Proposal

## Overview

This document provides a detailed, actionable plan for refactoring the Forge codebase based on the comprehensive code review. The refactoring is designed to improve maintainability, reduce complexity, and establish consistent patterns throughout the codebase.

---

## Guiding Principles

1. **Incremental Change** - Small, safe steps with testing at each stage
2. **Backward Compatibility** - Maintain existing API contracts
3. **Test-Driven** - Add tests before refactoring, verify after
4. **Documentation** - Update docs alongside code changes
5. **Team Consensus** - Review and approve each phase

---

## Phase 1: Critical Cleanup (Week 1)

**Goal:** Remove dead code and split monolithic files  
**Effort:** 16 hours  
**Risk:** Low

### Task 1.1: Remove Empty Packages (1 hour)

**Files to Remove:**
```
internal/core/core.go
internal/utils/utils.go
```

**Steps:**
1. Search for any imports of these packages
2. Delete the files
3. Remove directories if empty
4. Update documentation

**Verification:**
```bash
go build ./...
go test ./...
```

---

### Task 1.2: Split TUI Executor (8 hours)

**Current:** `pkg/executor/tui/executor.go` (1,400+ lines)

**Target Structure:**
```
pkg/executor/tui/
├── executor.go           # Main executor (150 lines)
├── model.go              # Model struct and initialization (200 lines)
├── init.go               # Initialization logic (100 lines)
├── update.go             # Bubble Tea Update method (300 lines)
├── view.go               # Bubble Tea View method (200 lines)
├── events.go             # Event handling (300 lines)
├── rendering.go          # Rendering utilities (150 lines)
└── formatting.go         # Text formatting helpers (100 lines)
```

**Detailed Breakdown:**

#### executor.go (Main)
```go
package tui

import (
    "context"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/entrhq/forge/pkg/agent"
    "github.com/entrhq/forge/pkg/llm"
)

// Executor is a TUI-based executor
type Executor struct {
    agent        agent.Agent
    program      *tea.Program
    provider     llm.Provider
    workspaceDir string
}

// NewExecutor creates a new TUI executor
func NewExecutor(agent agent.Agent, provider llm.Provider, workspaceDir string) *Executor {
    return &Executor{
        agent:        agent,
        provider:     provider,
        workspaceDir: workspaceDir,
    }
}

// Run starts the TUI executor
func (e *Executor) Run(ctx context.Context) error {
    if err := e.agent.Start(ctx); err != nil {
        return fmt.Errorf("failed to start agent: %w", err)
    }

    m := newModel(e.agent, e.provider, e.workspaceDir)
    
    e.program = tea.NewProgram(
        m,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )

    go e.forwardEvents(m)

    if _, err := e.program.Run(); err != nil {
        return fmt.Errorf("failed to run TUI program: %w", err)
    }

    return nil
}

func (e *Executor) forwardEvents(m *model) {
    for event := range m.channels.Event {
        e.program.Send(event)
    }
}
```

#### model.go (State)
```go
package tui

import (
    "strings"
    "time"
    
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/bubbles/viewport"
    "github.com/entrhq/forge/pkg/agent"
    "github.com/entrhq/forge/pkg/types"
)

// model represents the state of the TUI application
type model struct {
    // Core components
    viewport   viewport.Model
    textarea   textarea.Model
    agent      agent.Agent
    channels   *types.AgentChannels
    
    // Content buffers
    content        *strings.Builder
    thinkingBuffer *strings.Builder
    messageBuffer  *strings.Builder
    
    // UI state
    overlay         *overlayState
    commandPalette  *CommandPalette
    summarization   *summarizationStatus
    toast           *toastNotification
    spinner         spinner.Model
    
    // Flags
    isThinking               bool
    agentBusy                bool
    bashMode                 bool
    toolNameDisplayed        bool
    ready                    bool
    hasMessageContentStarted bool
    
    // Dimensions
    width  int
    height int
    
    // Token tracking
    totalPromptTokens     int
    totalCompletionTokens int
    totalTokens           int
    currentContextTokens  int
    maxContextTokens      int
    
    // Tool results
    resultClassifier *ToolResultClassifier
    resultSummarizer *ToolResultSummarizer
    resultCache      *resultCache
    resultList       resultListModel
    lastToolCallID   string
    lastToolName     string
    
    // Git operations
    workspaceDir string
    slashHandler *slash.Handler
    commitGen    *git.CommitMessageGenerator
    prGen        *git.PRGenerator
}

// newModel creates a new model with default values
func newModel(agent agent.Agent, provider llm.Provider, workspaceDir string) *model {
    m := &model{
        agent:        agent,
        channels:     agent.GetChannels(),
        workspaceDir: workspaceDir,
    }
    
    m.initializeComponents()
    m.initializeGitComponents(provider)
    
    return m
}
```

#### init.go (Initialization)
```go
package tui

import (
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/bubbles/viewport"
    "github.com/charmbracelet/lipgloss"
)

func (m *model) initializeComponents() {
    m.initializeTextArea()
    m.initializeViewport()
    m.initializeSpinner()
    m.initializeBuffers()
    m.initializeToolComponents()
}

func (m *model) initializeTextArea() {
    ta := textarea.New()
    ta.Placeholder = "Type your message..."
    ta.Focus()
    ta.Prompt = "> "
    ta.CharLimit = 0
    ta.SetHeight(1)
    ta.MaxHeight = 10
    ta.ShowLineNumbers = false
    ta.KeyMap.InsertNewline.SetEnabled(false)
    ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
    ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(salmonPink)
    ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(brightWhite)
    
    m.textarea = ta
}

func (m *model) initializeViewport() {
    vp := viewport.New(80, 20)
    vp.Style = lipgloss.NewStyle().Padding(0, 2)
    m.viewport = vp
}

func (m *model) initializeSpinner() {
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(salmonPink)
    m.spinner = s
}

// ... more initialization methods
```

#### update.go (Update Logic)
```go
package tui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/entrhq/forge/pkg/types"
)

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle overlay updates first
    if m.overlay.active {
        return m.updateOverlay(msg)
    }
    
    // Handle command palette
    if m.commandPalette.active {
        return m.updateCommandPalette(msg)
    }
    
    // Route to appropriate handler
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case tea.WindowSizeMsg:
        return m.handleWindowResize(msg)
    case *types.AgentEvent:
        return m.handleAgentEvent(msg)
    case spinner.TickMsg:
        return m.handleSpinnerTick(msg)
    default:
        return m.handleOtherMessages(msg)
    }
}

func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // Delegate to specific handlers
    switch msg.String() {
    case "ctrl+c":
        return m.handleQuit()
    case "ctrl+p":
        return m.handleCommandPalette()
    case "ctrl+l":
        return m.handleClearScreen()
    default:
        return m.handleTextInput(msg)
    }
}
```

#### events.go (Event Handling)
```go
package tui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/entrhq/forge/pkg/types"
)

func (m *model) handleAgentEvent(event *types.AgentEvent) (tea.Model, tea.Cmd) {
    switch event.Type {
    case types.EventTypeThinkingStart:
        return m.handleThinkingStart(event)
    case types.EventTypeThinkingContent:
        return m.handleThinkingContent(event)
    case types.EventTypeThinkingEnd:
        return m.handleThinkingEnd(event)
    case types.EventTypeMessageStart:
        return m.handleMessageStart(event)
    case types.EventTypeMessageContent:
        return m.handleMessageContent(event)
    // ... handle all event types
    default:
        return m, nil
    }
}

func (m *model) handleThinkingStart(event *types.AgentEvent) (tea.Model, tea.Cmd) {
    m.isThinking = true
    m.thinkingBuffer.Reset()
    return m, nil
}

// ... more event handlers
```

#### view.go (Rendering)
```go
package tui

import (
    "fmt"
    "strings"
)

// View renders the TUI
func (m model) View() string {
    if !m.ready {
        return "Initializing..."
    }
    
    // Build view components
    var parts []string
    
    // Main viewport
    parts = append(parts, m.renderViewport())
    
    // Summarization status if active
    if m.summarization.active {
        parts = append(parts, m.renderSummarizationStatus())
    }
    
    // Toast notification if active
    if m.toast.active && time.Now().Before(m.toast.showUntil) {
        parts = append(parts, m.renderToast())
    }
    
    // Input area
    parts = append(parts, m.renderInputArea())
    
    // Status bar
    parts = append(parts, m.renderStatusBar())
    
    // Overlay on top if active
    if m.overlay.active {
        return m.renderWithOverlay(strings.Join(parts, "\n"))
    }
    
    // Command palette on top if active
    if m.commandPalette.active {
        return m.renderWithCommandPalette(strings.Join(parts, "\n"))
    }
    
    return strings.Join(parts, "\n")
}

func (m model) renderViewport() string {
    return m.viewport.View()
}

func (m model) renderInputArea() string {
    if m.agentBusy {
        return m.renderBusyIndicator()
    }
    return m.textarea.View()
}

func (m model) renderBusyIndicator() string {
    return fmt.Sprintf("%s %s", m.spinner.View(), m.currentLoadingMessage)
}
```

**Migration Steps:**

1. Create new files with extracted code
2. Update imports in each file
3. Run tests after each file creation
4. Remove code from original file
5. Final verification

**Tests:**
```bash
go test ./pkg/executor/tui/...
```

---

### Task 1.3: Extract Approval Manager (6 hours)

**Current:** Approval logic scattered across 10+ functions in `default.go`

**New Structure:**
```
pkg/agent/approval/
├── manager.go        # Main approval manager
├── auto_approval.go  # Auto-approval logic
└── approval_test.go  # Tests
```

#### manager.go
```go
package approval

import (
    "context"
    "sync"
    "time"
    
    "github.com/entrhq/forge/pkg/agent/tools"
    "github.com/entrhq/forge/pkg/config"
    "github.com/entrhq/forge/pkg/types"
    "github.com/google/uuid"
)

// Manager handles tool approval workflows
type Manager struct {
    timeout         time.Duration
    pendingApproval *pendingApproval
    approvalMu      sync.Mutex
    eventEmitter    EventEmitter
}

type EventEmitter func(*types.AgentEvent)

type pendingApproval struct {
    approvalID string
    toolName   string
    toolCall   tools.ToolCall
    response   chan *types.ApprovalResponse
}

// NewManager creates a new approval manager
func NewManager(timeout time.Duration, emitter EventEmitter) *Manager {
    return &Manager{
        timeout:      timeout,
        eventEmitter: emitter,
    }
}

// RequestApproval requests approval for a tool call
func (m *Manager) RequestApproval(
    ctx context.Context,
    toolCall tools.ToolCall,
    preview *tools.ToolPreview,
    approvalChan <-chan *types.ApprovalResponse,
) (approved bool, timedOut bool) {
    
    // Check auto-approval first
    if m.checkAutoApproval(toolCall) {
        return true, false
    }
    
    // Generate approval ID
    approvalID := uuid.New().String()
    
    // Create response channel
    responseChan := make(chan *types.ApprovalResponse, 1)
    
    // Setup pending approval
    m.setupPending(approvalID, toolCall, responseChan)
    defer m.cleanup(responseChan)
    
    // Emit approval request
    argsMap := parseToolArguments(toolCall)
    m.eventEmitter(types.NewToolApprovalRequestEvent(
        approvalID,
        toolCall.ToolName,
        argsMap,
        preview,
    ))
    
    // Wait for response
    return m.waitForResponse(ctx, approvalID, toolCall, responseChan, approvalChan)
}

// HandleResponse processes an approval response
func (m *Manager) HandleResponse(response *types.ApprovalResponse) {
    m.approvalMu.Lock()
    defer m.approvalMu.Unlock()
    
    if m.pendingApproval == nil || m.pendingApproval.approvalID != response.ApprovalID {
        return
    }
    
    select {
    case m.pendingApproval.response <- response:
    default:
    }
}

func (m *Manager) checkAutoApproval(toolCall tools.ToolCall) bool {
    // Special handling for execute_command
    if toolCall.ToolName == "execute_command" {
        return checkCommandWhitelist(toolCall)
    }
    
    // Check if tool is auto-approved
    return config.IsToolAutoApproved(toolCall.ToolName)
}

// ... more helper methods
```

**Integration into DefaultAgent:**
```go
// In pkg/agent/default.go

import "github.com/entrhq/forge/pkg/agent/approval"

type DefaultAgent struct {
    // ... existing fields
    approvalManager *approval.Manager
}

func NewDefaultAgent(provider llm.Provider, opts ...AgentOption) *DefaultAgent {
    a := &DefaultAgent{
        // ... existing initialization
        approvalManager: approval.NewManager(
            5*time.Minute,
            a.emitEvent, // Pass emitter function
        ),
    }
    // ...
}

func (a *DefaultAgent) executeTool(ctx context.Context, toolCall tools.ToolCall) (bool, string) {
    // ... existing code
    
    // Simplified approval check
    if previewable, ok := tool.(tools.Previewable); ok {
        preview, err := previewable.GeneratePreview(ctx, toolCall.GetArgumentsXML())
        if err != nil {
            a.emitEvent(types.NewErrorEvent(err))
        } else {
            approved, timedOut := a.approvalManager.RequestApproval(
                ctx,
                toolCall,
                preview,
                a.channels.Approval,
            )
            
            if timedOut {
                errMsg := fmt.Sprintf("Tool approval timed out")
                a.memory.Add(types.NewUserMessage(errMsg))
                return true, ""
            }
            
            if !approved {
                errMsg := fmt.Sprintf("Tool '%s' rejected", toolCall.ToolName)
                a.memory.Add(types.NewUserMessage(errMsg))
                return true, ""
            }
        }
    }
    
    // ... continue with execution
}
```

**Benefits:**
- 300+ lines removed from `default.go`
- Clear separation of concerns
- Easier to test approval logic in isolation
- Reusable across different agent implementations

---

## Phase 2: Core Refactoring (Week 2-3)

**Goal:** Simplify complex logic and standardize patterns  
**Effort:** 24 hours  
**Risk:** Medium

### Task 2.1: Simplify Agent Loop Methods (8 hours)

**Target:** Break down complex methods in `default.go`

#### Before:
```go
func (a *DefaultAgent) executeIteration(ctx context.Context, errorContext string) (bool, string) {
    // 107 lines of complex logic
    // Mixed concerns: prompt building, LLM calls, token counting, summarization, tool processing
}
```

#### After:
```go
func (a *DefaultAgent) executeIteration(ctx context.Context, errorContext string) (bool, string) {
    // Build and send prompt
    messages := a.buildIterationMessages(errorContext)
    
    // Evaluate context management
    if err := a.evaluateContextManagement(ctx, messages); err != nil {
        a.emitEvent(types.NewErrorEvent(err))
    }
    
    // Call LLM
    response, err := a.callLLM(ctx, messages)
    if err != nil {
        return a.handleLLMError(ctx, err)
    }
    
    // Track token usage
    a.trackTokenUsage(response)
    
    // Process response
    return a.processResponse(ctx, response)
}

func (a *DefaultAgent) buildIterationMessages(errorContext string) []*types.Message {
    systemPrompt := a.buildSystemPrompt()
    history := a.memory.GetAll()
    return prompts.BuildMessages(systemPrompt, history, "", errorContext)
}

func (a *DefaultAgent) evaluateContextManagement(ctx context.Context, messages []*types.Message) error {
    if a.contextManager == nil || a.tokenizer == nil {
        return nil
    }
    
    convMem, ok := a.memory.(*memory.ConversationMemory)
    if !ok {
        return nil
    }
    
    promptTokens := a.tokenizer.CountMessagesTokens(messages)
    _, err := a.contextManager.EvaluateAndSummarize(ctx, convMem, promptTokens)
    return err
}

func (a *DefaultAgent) callLLM(ctx context.Context, messages []*types.Message) (*llmResponse, error) {
    // Emit API call start
    promptTokens := 0
    if a.tokenizer != nil {
        promptTokens = a.tokenizer.CountMessagesTokens(messages)
    }
    
    maxTokens := 0
    if a.contextManager != nil {
        maxTokens = a.contextManager.GetMaxTokens()
    }
    a.emitEvent(types.NewApiCallStartEvent("llm", promptTokens, maxTokens))
    
    // Stream completion
    stream, err := a.provider.StreamCompletion(ctx, messages)
    if err != nil {
        return nil, err
    }
    
    // Process stream
    var assistantContent, toolCallContent string
    core.ProcessStream(stream, a.emitEvent, func(content, thinking, toolCall, role string) {
        assistantContent = content
        toolCallContent = toolCall
    })
    
    return &llmResponse{
        assistantContent: assistantContent,
        toolCallContent:  toolCallContent,
    }, nil
}
```

---

### Task 2.2: Standardize Error Handling (6 hours)

**Create:** `pkg/agent/errors/errors.go`

```go
package errors

import "fmt"

// ErrorType categorizes different error scenarios
type ErrorType int

const (
    ErrorTypeUnknown ErrorType = iota
    ErrorTypeNoToolCall
    ErrorTypeInvalidXML
    ErrorTypeMissingToolName
    ErrorTypeUnknownTool
    ErrorTypeToolExecution
    ErrorTypeCircuitBreaker
    ErrorTypeLLMFailure
    ErrorTypeContextCanceled
)

// String returns the error type as a string
func (e ErrorType) String() string {
    switch e {
    case ErrorTypeNoToolCall:
        return "NO_TOOL_CALL"
    case ErrorTypeInvalidXML:
        return "INVALID_XML"
    case ErrorTypeMissingToolName:
        return "MISSING_TOOL_NAME"
    case ErrorTypeUnknownTool:
        return "UNKNOWN_TOOL"
    case ErrorTypeToolExecution:
        return "TOOL_EXECUTION"
    case ErrorTypeCircuitBreaker:
        return "CIRCUIT_BREAKER"
    case ErrorTypeLLMFailure:
        return "LLM_FAILURE"
    case ErrorTypeContextCanceled:
        return "CONTEXT_CANCELED"
    default:
        return "UNKNOWN"
    }
}

// AgentError represents a structured error with type and context
type AgentError struct {
    Type    ErrorType
    Message string
    Cause   error
    Context map[string]interface{}
}

// Error implements the error interface
func (e *AgentError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Cause)
    }
    return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *AgentError) Unwrap() error {
    return e.Cause
}

// WithContext adds context to the error
func (e *AgentError) WithContext(key string, value interface{}) *AgentError {
    if e.Context == nil {
        e.Context = make(map[string]interface{})
    }
    e.Context[key] = value
    return e
}

// New creates a new AgentError
func New(errType ErrorType, message string) *AgentError {
    return &AgentError{
        Type:    errType,
        Message: message,
    }
}

// Wrap wraps an existing error with additional context
func Wrap(errType ErrorType, message string, cause error) *AgentError {
    return &AgentError{
        Type:    errType,
        Message: message,
        Cause:   cause,
    }
}
```

**Usage in DefaultAgent:**
```go
import agenterrors "github.com/entrhq/forge/pkg/agent/errors"

// Before
func (a *DefaultAgent) trackError(errMsg string) bool {
    a.lastErrors[a.errorIndex] = errMsg
    // String comparison
}

// After
func (a *DefaultAgent) trackError(err *agenterrors.AgentError) bool {
    a.lastErrors[a.errorIndex] = err
    a.errorIndex = (a.errorIndex + 1) % 5
    
    // Check circuit breaker by error type
    if a.lastErrors[0] == nil {
        return false
    }
    
    firstType := a.lastErrors[0].Type
    for i := 1; i < 5; i++ {
        if a.lastErrors[i] == nil || a.lastErrors[i].Type != firstType {
            return false
        }
    }
    
    return true
}
```

---

### Task 2.3: Consolidate Overlay Components (4 hours)

**Create:** `pkg/executor/tui/overlay/base.go`

```go
package overlay

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Base provides common overlay functionality
type Base struct {
    Title   string
    Width   int
    Height  int
    Style   lipgloss.Style
}

// NewBase creates a new base overlay
func NewBase(title string, width, height int) *Base {
    return &Base{
        Title:  title,
        Width:  width,
        Height: height,
        Style:  defaultOverlayStyle(),
    }
}

// RenderFrame renders a framed overlay with title
func (b *Base) RenderFrame(content string) string {
    // Create border style
    borderStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(salmonPink).
        Padding(1, 2).
        Width(b.Width).
        Height(b.Height)
    
    // Add title
    titleStyle := lipgloss.NewStyle().
        Foreground(salmonPink).
        Bold(true)
    
    header := titleStyle.Render(b.Title)
    
    // Combine
    body := header + "\n\n" + content
    
    return borderStyle.Render(body)
}

// Center centers the overlay on screen
func (b *Base) Center(screenWidth, screenHeight int) lipgloss.Style {
    return lipgloss.NewStyle().
        Position(lipgloss.Position{}).
        AlignHorizontal(lipgloss.Center).
        AlignVertical(lipgloss.Center)
}

func defaultOverlayStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Background(lipgloss.Color("#1a1a1a")).
        Foreground(lipgloss.Color("#ffffff"))
}
```

**Refactor Overlays:**
```go
// pkg/executor/tui/help_overlay.go
type HelpOverlay struct {
    *overlay.Base
    commands []CommandHelp
}

func NewHelpOverlay(width, height int) *HelpOverlay {
    return &HelpOverlay{
        Base:     overlay.NewBase("Help", width, height),
        commands: getHelpCommands(),
    }
}

func (h *HelpOverlay) View() string {
    content := h.renderCommands()
    return h.RenderFrame(content)
}
```

---

### Task 2.4: Implement Structured Logging (6 hours)

**Create:** `pkg/logging/logger.go`

```go
package logging

import (
    "io"
    "log/slog"
    "os"
)

// Logger provides structured logging
type Logger struct {
    *slog.Logger
}

// Config holds logger configuration
type Config struct {
    Level      slog.Level
    Output     io.Writer
    AddSource  bool
    JSONFormat bool
}

// NewLogger creates a new structured logger
func NewLogger(cfg Config) *Logger {
    opts := &slog.HandlerOptions{
        Level:     cfg.Level,
        AddSource: cfg.AddSource,
    }
    
    var handler slog.Handler
    if cfg.JSONFormat {
        handler = slog.NewJSONHandler(cfg.Output, opts)
    } else {
        handler = slog.NewTextHandler(cfg.Output, opts)
    }
    
    return &Logger{
        Logger: slog.New(handler),
    }
}

// NewDefaultLogger creates a logger with sensible defaults
func NewDefaultLogger() *Logger {
    return NewLogger(Config{
        Level:      slog.LevelInfo,
        Output:     os.Stderr,
        AddSource:  false,
        JSONFormat: false,
    })
}

// NewDebugLogger creates a debug logger that writes to a file
func NewDebugLogger(filepath string) (*Logger, error) {
    f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }
    
    return NewLogger(Config{
        Level:      slog.LevelDebug,
        Output:     f,
        AddSource:  true,
        JSONFormat: false,
    }), nil
}
```

**Integration:**
```go
// pkg/agent/default.go

type DefaultAgent struct {
    // ... existing fields
    logger *logging.Logger
}

func WithLogger(logger *logging.Logger) AgentOption {
    return func(a *DefaultAgent) {
        a.logger = logger
    }
}

// Usage
func (a *DefaultAgent) executeIteration(ctx context.Context, errorContext string) (bool, string) {
    a.logger.Debug("starting iteration",
        slog.String("error_context", errorContext),
        slog.Int("message_count", len(a.memory.GetAll())),
    )
    
    // ... rest of the method
}
```

---

## Phase 3: Polish and Documentation (Week 4)

**Goal:** Improve code documentation and add missing tests  
**Effort:** 16 hours  
**Risk:** Low

### Task 3.1: Document Magic Numbers (2 hours)

**Example:**
```go
// Before
const (
    defaultMaxTokens        = 100000
    defaultThresholdPercent = 80.0
    defaultToolCallAge      = 20
)

// After
const (
    // defaultMaxTokens is the conservative context window limit with headroom
    // for 128K context models. This leaves ~28K tokens for system prompt,
    // tools, and completion buffer while maintaining stable performance.
    defaultMaxTokens = 100000
    
    // defaultThresholdPercent triggers context summarization at 80% capacity.
    // This threshold balances:
    // - Early intervention to prevent sudden context exhaustion
    // - Batch efficiency by collecting multiple items before summarizing
    // - Model performance (models perform better with more context)
    defaultThresholdPercent = 80.0
    
    // defaultToolCallAge defines how many messages back a tool call must be
    // before entering the summarization buffer. This ensures recent tool
    // operations remain in full context for immediate reference.
    defaultToolCallAge = 20
    
    // defaultMinToolCalls is the minimum buffer size before triggering
    // batch summarization. This reduces LLM API calls by processing
    // multiple tool results in a single summarization request.
    defaultMinToolCalls = 10
    
    // defaultMaxToolCallDist is the maximum age (in messages) before
    // forcing summarization regardless of buffer size. This prevents
    // very old tool calls from consuming context indefinitely.
    defaultMaxToolCallDist = 40
)
```

### Task 3.2: Add Integration Tests (8 hours)

**Create:** `pkg/agent/integration_test.go`

```go
package agent_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/entrhq/forge/pkg/agent"
    "github.com/entrhq/forge/pkg/agent/tools"
    "github.com/entrhq/forge/pkg/types"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAgentLoop_ToolExecution_Integration(t *testing.T) {
    // Setup mock provider
    provider := newMockProvider(t)
    
    // Create agent
    ag := agent.NewDefaultAgent(provider)
    
    // Register test tool
    testTool := newTestTool()
    require.NoError(t, ag.RegisterTool(testTool))
    
    // Start agent
    ctx := context.Background()
    require.NoError(t, ag.Start(ctx))
    defer ag.Shutdown(ctx)
    
    // Send user input
    channels := ag.GetChannels()
    channels.Input <- types.NewUserInput("Test message")
    
    // Collect events
    var events []*types.AgentEvent
    timeout := time.After(5 * time.Second)
    
eventLoop:
    for {
        select {
        case event := <-channels.Event:
            events = append(events, event)
            if event.Type == types.EventTypeTurnEnd {
                break eventLoop
            }
        case <-timeout:
            t.Fatal("timeout waiting for turn end")
        }
    }
    
    // Verify event sequence
    assert.Contains(t, eventTypes(events), types.EventTypeToolCall)
    assert.Contains(t, eventTypes(events), types.EventTypeToolResult)
}

func TestApprovalWorkflow_Integration(t *testing.T) {
    provider := newMockProvider(t)
    ag := agent.NewDefaultAgent(provider)
    
    // Register previewable tool
    previewableTool := newPreviewableTool()
    require.NoError(t, ag.RegisterTool(previewableTool))
    
    ctx := context.Background()
    require.NoError(t, ag.Start(ctx))
    defer ag.Shutdown(ctx)
    
    channels := ag.GetChannels()
    channels.Input <- types.NewUserInput("Execute previewed operation")
    
    // Wait for approval request
    var approvalID string
    select {
    case event := <-channels.Event:
        if event.Type == types.EventTypeToolApprovalRequest {
            approvalID = event.ApprovalRequest.ApprovalID
        }
    case <-time.After(2 * time.Second):
        t.Fatal("timeout waiting for approval request")
    }
    
    // Grant approval
    channels.Approval <- &types.ApprovalResponse{
        ApprovalID: approvalID,
        Action:     types.ApprovalActionApprove,
    }
    
    // Verify tool executed
    var toolExecuted bool
    timeout := time.After(3 * time.Second)
    
    for {
        select {
        case event := <-channels.Event:
            if event.Type == types.EventTypeToolResult {
                toolExecuted = true
            }
            if event.Type == types.EventTypeTurnEnd {
                assert.True(t, toolExecuted, "tool should have executed after approval")
                return
            }
        case <-timeout:
            t.Fatal("timeout waiting for tool execution")
        }
    }
}
```

### Task 3.3: Update ADRs (2 hours)

**Create:** `docs/adr/0025-refactoring-2024.md`

```markdown
# ADR-0025: Codebase Refactoring 2024

## Status
Accepted

## Context
After comprehensive code review, several areas were identified for improvement:
- Large, complex files (executor.go 1400+ lines)
- Approval logic scattered across 10+ functions
- Inconsistent error handling patterns
- Empty placeholder packages

## Decision
Implement phased refactoring over 4 weeks:

### Phase 1: Critical Cleanup
- Remove empty packages
- Split TUI executor into focused files
- Extract approval manager

### Phase 2: Core Refactoring
- Simplify agent loop methods
- Standardize error handling
- Consolidate overlay components
- Implement structured logging

### Phase 3: Polish
- Document magic numbers
- Add integration tests
- Update documentation

## Consequences

### Positive
- Improved maintainability
- Better testability
- Clearer separation of concerns
- Consistent patterns

### Negative
- Initial time investment (~40 hours)
- Learning curve for new structure
- Risk of introducing bugs during refactoring

### Mitigation
- Incremental changes with testing
- Comprehensive test coverage
- Code review for each phase
```

### Task 3.4: Create Developer Guide (4 hours)

**Create:** `docs/CONTRIBUTING_CODE.md`

```markdown
# Code Contribution Guide

## Code Organization

### Package Structure
```
pkg/
├── agent/          # Agent implementations
│   ├── approval/   # Approval workflow management
│   ├── context/    # Context management strategies
│   ├── memory/     # Conversation memory
│   └── tools/      # Tool system
├── executor/       # Execution environments
│   ├── cli/        # CLI executor
│   └── tui/        # TUI executor
├── llm/            # LLM provider abstractions
├── tools/          # Built-in tools
│   └── coding/     # Coding assistance tools
└── types/          # Shared types and interfaces
```

## Coding Standards

### File Size
- Keep files under 500 lines
- Split large files by concern
- One primary type per file

### Function Complexity
- Maximum cyclomatic complexity: 10
- Extract helper functions liberally
- Use early returns to reduce nesting

### Error Handling
```go
// Use typed errors
import agenterrors "github.com/entrhq/forge/pkg/agent/errors"

err := agenterrors.New(
    agenterrors.ErrorTypeToolExecution,
    "failed to execute tool",
).WithContext("tool_name", toolName)
```

### Logging
```go
// Use structured logging
logger.Debug("operation completed",
    slog.String("tool", toolName),
    slog.Int("duration_ms", duration),
)
```

### Testing
- Write tests before refactoring
- Maintain >80% coverage
- Use table-driven tests
- Add integration tests for workflows

## Pull Request Process

1. Create feature branch from main
2. Make changes in small, logical commits
3. Add/update tests
4. Update documentation
5. Run `make all` (lint, test, build)
6. Create PR with description
7. Address review feedback
8. Squash and merge
```

---

## Migration Checklist

### Pre-Refactoring
- [ ] Full test suite passing
- [ ] Create backup branch
- [ ] Document current behavior
- [ ] Set up test coverage baseline

### Phase 1
- [ ] Remove `internal/core/core.go`
- [ ] Remove `internal/utils/utils.go`
- [ ] Create `pkg/executor/tui/model.go`
- [ ] Create `pkg/executor/tui/init.go`
- [ ] Create `pkg/executor/tui/update.go`
- [ ] Create `pkg/executor/tui/view.go`
- [ ] Create `pkg/executor/tui/events.go`
- [ ] Refactor `executor.go` to use new files
- [ ] Create `pkg/agent/approval/manager.go`
- [ ] Update `DefaultAgent` to use approval manager
- [ ] Run full test suite
- [ ] Update documentation

### Phase 2
- [ ] Simplify `executeIteration()`
- [ ] Simplify `executeTool()`
- [ ] Simplify `processToolCall()`
- [ ] Create `pkg/agent/errors/errors.go`
- [ ] Update error handling throughout
- [ ] Create `pkg/executor/tui/overlay/base.go`
- [ ] Refactor overlay components
- [ ] Create `pkg/logging/logger.go`
- [ ] Replace debug logging
- [ ] Run full test suite
- [ ] Verify no regressions

### Phase 3
- [ ] Document all constants
- [ ] Add integration tests
- [ ] Create ADR-0025
- [ ] Update CONTRIBUTING_CODE.md
- [ ] Update README if needed
- [ ] Final code review
- [ ] Merge to main

---

## Success Metrics

### Before Refactoring
- Largest file: 1,400 lines
- Average file: ~250 lines
- Test files: 36
- Code duplication: ~15%

### After Refactoring (Target)
- Largest file: <500 lines
- Average file: ~200 lines
- Test files: 40+
- Code duplication: <5%
- Test coverage: >85%

### Quality Gates
- [ ] All tests passing
- [ ] No linter warnings
- [ ] Test coverage maintained/improved
- [ ] Documentation updated
- [ ] Code review approved
- [ ] No performance regressions

---

## Rollback Plan

If issues arise during refactoring:

1. **Minor Issues**
   - Fix forward with additional commits
   - Add regression tests

2. **Major Issues**
   - Revert problematic commits
   - Create hotfix branch
   - Re-attempt with more caution

3. **Critical Issues**
   - Revert entire phase
   - Reassess approach
   - Add more tests before retry

---

## Timeline

| Week | Phase | Tasks | Hours |
|------|-------|-------|-------|
| 1 | Phase 1 | Critical cleanup | 16 |
| 2-3 | Phase 2 | Core refactoring | 24 |
| 4 | Phase 3 | Polish & docs | 16 |
| **Total** | | | **56 hours** |

---

## Conclusion

This refactoring proposal provides a structured, low-risk approach to improving the Forge codebase. By following this plan incrementally with proper testing and documentation, we can achieve significant improvements in code quality and maintainability without disrupting existing functionality.

**Next Steps:**
1. Review and approve this proposal
2. Create GitHub issues for each phase
3. Begin Phase 1 implementation
4. Regular progress reviews

**Questions or Concerns:**
- Open GitHub discussion
- Tag relevant team members
- Update proposal as needed
