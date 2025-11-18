# TUI Refactor Recovery Report

## Executive Summary

This document tracks the recovery effort for the TUI executor refactor. The refactor successfully split a 1,377-line monolithic file into focused modules, but introduced several regressions that broke core functionality.

**Branch:** `refactor/split-tui-executor`  
**Comparison Base:** `main`  
**Recovery Status:** In Progress

---

## Issues Found and Fixed

### 1. ✅ FIXED: Token Count Display Regression
**File:** `pkg/executor/tui/helpers.go`  
**Severity:** Low  
**Impact:** Token counts over 1M displayed incorrectly

**Problem:**
```go
// Missing million formatting
if count >= 1000 {
    return fmt.Sprintf("%.1fK", float64(count)/1000)
}
```

**Solution:**
```go
if count >= 1000000 {
    return fmt.Sprintf("%.1fM", float64(count)/1000000)
}
if count >= 1000 {
    return fmt.Sprintf("%.1fK", float64(count)/1000)
}
```

---

### 2. ✅ FIXED: Command Execution Overlay Missing
**File:** `pkg/executor/tui/events.go`  
**Severity:** Critical  
**Impact:** Commands executed without interactive overlay, no cancel support

**Problem:**
The `handleCommandExecution()` function was empty - no overlay creation.

**Solution:**
Restored overlay creation logic:
```go
func (m *model) handleCommandExecution(event *types.AgentEvent) {
    overlay := NewCommandExecutionOverlay(
        event.CommandExecution.Command,
        event.CommandExecution.WorkingDir,
        event.CommandExecution.ExecutionID,
        m.channels.Cancel,
    )
    m.overlay.activate(OverlayModeCommandOutput, overlay)
}
```

---

### 3. ✅ FIXED: Approval Request Handler Missing
**File:** `pkg/executor/tui/update.go`  
**Severity:** Critical  
**Impact:** Slash command approvals wouldn't display UI

**Problem:**
Missing case in `Update()` for `approvalRequestMsg` and missing handler function.

**Solution:**
Added message handler:
```go
case approvalRequestMsg:
    return m.handleApprovalRequest(msg)

func (m model) handleApprovalRequest(msg approvalRequestMsg) (tea.Model, tea.Cmd) {
    overlay := NewGenericApprovalOverlay(msg.request, m.width, m.height)
    m.overlay.activate(OverlayModeSlashCommandPreview, overlay)
    return m, nil
}
```

---

### 4. ✅ FIXED: Streaming Content Not Working
**File:** `pkg/executor/tui/update.go`  
**Severity:** Critical  
**Impact:** Thinking and message content not streaming in real-time

**Root Cause:**
Event processing order was reversed in refactor:
- **Main (correct):** `viewport.Update(msg)` → `handleAgentEvent(msg)`
- **Refactor (wrong):** `handleAgentEvent(msg)` → `viewport.Update(msg)`

The viewport update was overwriting the streaming updates made by handlers.

**Solution:**
Restored correct order and added command execution overlay forwarding:
```go
case *types.AgentEvent:
    // If overlay is active and it's a command execution event, forward to overlay
    if m.overlay.isActive() && msg.IsCommandExecutionEvent() {
        var overlayCmd tea.Cmd
        m.overlay.overlay, overlayCmd = m.overlay.overlay.Update(msg)
        // Still handle the event in the main model too
        m.handleAgentEvent(msg)
        return m, tea.Batch(tiCmd, vpCmd, overlayCmd, spinnerCmd)
    }
    
    // Update viewport BEFORE handling event (critical for streaming)
    m.viewport, vpCmd = m.viewport.Update(msg)
    m.handleAgentEvent(msg)
    return m, tea.Batch(tiCmd, vpCmd, spinnerCmd)
```

**Why This Matters:**
1. Viewport processes raw message first
2. Handlers append/modify viewport content
3. If viewport runs after handlers, it overwrites streaming updates

---

## Remaining Tasks

### High Priority
- [ ] Check each overlay type for WindowSizeMsg handling
- [ ] Verify approval workflow end-to-end integration
- [ ] Manual testing of all TUI functionality

### Medium Priority
- [ ] Add nil checks and defensive programming to overlays
- [ ] Create integration tests for key workflows

### Low Priority
- [ ] Commit fixes with descriptive message
- [ ] Update refactor todo list

---

## Files Modified During Recovery

1. `pkg/executor/tui/helpers.go` - Token count formatting
2. `pkg/executor/tui/events.go` - Command execution overlay
3. `pkg/executor/tui/update.go` - Approval requests + event order fix

---

### 10. ✅ FIXED: Command Palette Activation Missing
**File:** `pkg/executor/tui/update.go`  
**Severity:** Critical  
**Impact:** Typing `/` did not activate command palette at all

**Problem:**
The textarea update logic was missing the check for `/` character to activate/update the command palette. This entire block was deleted during refactor.

**Solution:**
Restored command palette activation logic after textarea update:
```go
// Check if we should activate/deactivate command palette based on input
value := m.textarea.Value()

// Handle command palette activation/deactivation based on input
switch {
case value == "/" && !m.commandPalette.active:
    // Only activate palette if input is exactly "/" as first character
    m.commandPalette.activate()
    m.commandPalette.updateFilter("")
case strings.HasPrefix(value, "/") && m.commandPalette.active:
    // Update filter if palette is already active
    filter := strings.TrimPrefix(value, "/")
    m.commandPalette.updateFilter(filter)
case !strings.HasPrefix(value, "/") && m.commandPalette.active:
    // Deactivate palette if input no longer starts with /
    m.commandPalette.deactivate()
}
```


### 11. ✅ FIXED: Command Palette Enter Key Submits Instead of Autocompleting
**File:** `pkg/executor/tui/update.go`  
**Severity:** Critical  
**Impact:** Pressing Enter in command palette submitted `/` to chat instead of autocompleting

**Problem:**
The command palette keyboard handling was happening AFTER textarea update in the event flow:
1. User types `/` and presses Enter
2. Textarea processes Enter key, adds newline
3. Command palette handler tries to autocomplete
4. But textarea already submitted the message

**Solution:**
Moved command palette keyboard handling to BEFORE textarea update in `Update()` function:
```go
// Handle command palette keyboard input BEFORE updating textarea
// This prevents Enter from being processed by textarea when palette is active
if keyMsg, ok := msg.(tea.KeyMsg); ok && m.commandPalette.active {
    switch keyMsg.Type {
    case tea.KeyEnter:
        // Autocomplete with the selected command and close the palette
        selected := m.commandPalette.getSelected()
        if selected != nil {
            m.textarea.SetValue("/" + selected.Name + " ")
            m.textarea.CursorEnd()
        }
        m.commandPalette.deactivate()
        return m, tea.Batch(tiCmd, vpCmd, spinnerCmd)
    // ... other keys
    }
}

// Only THEN update textarea
if !m.overlay.isActive() && !m.resultList.active {
    m.textarea, tiCmd = m.textarea.Update(msg)
    // ...

### 12. ✅ FIXED: Slash Commands Displayed in Chat History
**File:** `pkg/executor/tui/update.go`  
**Severity:** High  
**Impact:** Slash commands like `/context` were being displayed in chat history instead of executing silently

**Problem:**
The `handleSlashCommand()` function was incorrectly displaying slash commands in the chat viewport:
```go
// Display the slash command
formatted := formatEntry("You: ", input, userStyle, m.width, true)
m.content.WriteString(formatted + "\n\n")  // Wrong! Adds to chat
```

**Solution:**
Removed the display logic - slash commands execute silently:
```go
// Do NOT display slash commands in chat history - they are executed silently

// Clear the input area
m.textarea.Reset()

// Parse and execute slash command
commandName, args, ok := parseSlashCommand(input)
// ...
```

### 13. ✅ FIXED: Missing Textarea Auto-Height Function
**File:** `pkg/executor/tui/helpers.go`, `pkg/executor/tui/update.go`  
**Severity:** Medium  
**Impact:** Textarea did not auto-expand when typing multi-line input

**Problem:**
The `updateTextAreaHeight()` function was completely missing from the refactored code. This function dynamically adjusts textarea height based on content, accounting for line wrapping.

**Solution:**
Added the function to helpers.go and called it in two places:
1. After Alt+Enter inserts newline
2. After any textarea content change

```go
// In helpers.go
func (m *model) updateTextAreaHeight() {
    value := m.textarea.Value()
    // Calculate visual lines accounting for wrapping
    // Clamp between 1 and MaxHeight
    // Only update if height changed
}

// In update.go - after Alt+Enter
m.textarea.InsertString("\n")
m.updateTextAreaHeight()

// In update.go - after command palette logic
m.updateTextAreaHeight()
```

---

### 14. ✅ FIXED: Missing Mouse Event Handling
**File:** `pkg/executor/tui/update.go`  
**Severity:** Medium  
**Impact:** Mouse wheel scrolling and clicks did not work

**Problem:**
The `tea.MouseMsg` case was completely missing from the Update() switch statement, so mouse events were being ignored.

**Solution:**
Added mouse event handling that:
- Forwards mouse events to active overlays
- Routes mouse events to viewport for scrolling when no overlay active

```go
case tea.MouseMsg:
    // If overlay is active, forward mouse events to it
    if m.overlay.isActive() {
        updatedOverlay, overlayCmd := m.overlay.overlay.Update(msg)
        if updatedOverlay == nil {
            m.overlay.deactivate()
            return m, overlayCmd
        }
        m.overlay.overlay = updatedOverlay
        return m, overlayCmd
    }
    
    // Route mouse events to viewport for scrolling
    m.viewport, vpCmd = m.viewport.Update(msg)
    return m, tea.Batch(tiCmd, vpCmd, spinnerCmd)

### 15. ✅ FIXED: Command Output Formatting/Indentation Lost
**File:** `pkg/executor/tui/events.go`  
**Severity:** Medium  
**Impact:** Command output lost indentation and formatting due to style rendering

**Problem:**
The `handleCommandExecutionOutput()` function was applying `commandOutputStyle.Render()` to command output, which reformats/wraps text and destroys original formatting:
```go
m.content.WriteString(commandOutputStyle.Render(event.CommandExecution.Output))
```

**Solution:**
Write command output directly without styling to preserve original formatting:
```go
// Write output directly without styling to preserve formatting/indentation
if event.CommandExecution != nil && event.CommandExecution.Output != "" {
    m.content.WriteString(event.CommandExecution.Output)
}
```

This preserves:
- Original indentation (spaces/tabs)
- Line breaks and formatting
- Code structure in command output
- Table formatting, etc.

---
```

---

---
}
```

This ensures command palette intercepts special keys (Enter/Esc/Up/Down/Tab) before textarea can process them.

---
---

## Testing Notes

### Compilation
✅ `go build` passes in `pkg/executor/tui`

### Unit Tests
Status: Not yet run against fixes

### Integration Tests
Status: Manual testing required

---

## Key Learnings

1. **Event Processing Order Matters**: Viewport updates must happen before handlers for streaming to work
2. **Overlay Lifecycle**: Command execution events need forwarding to active overlays
3. **Message Routing**: All message types in main need corresponding handlers in refactored code
4. **Business Logic Preservation**: Critical to compare diffs line-by-line to catch missing logic

---

## Next Steps

1. Run comprehensive manual testing
2. Verify all overlay interactions
3. Test approval workflows end-to-end
4. Consider adding integration tests to prevent future regressions
5. Document any additional findings

---

---

### 16. ✅ FIXED: Summarization Progress Display
**File:** `pkg/executor/tui/view.go`  
**Severity:** Medium  
**Impact:** Progress bar showed only percentage, not item counts

**Problem:**
```go
progressLine := fmt.Sprintf("%s %.0f%%", bar, m.summarization.progressPercent)
```

The progress bar only showed percentage. User feedback indicated:
- Missing "X/Y items" display
- Percentage calculation was already cumulative (based on ItemsProcessed/TotalItems)
- Display format needed improvement for better visibility

**Solution:**
```go
// Show both item count and percentage
if m.summarization.totalItems > 0 {
    progressLine := fmt.Sprintf("%s %d/%d items (%.0f%%)", 
        bar, m.summarization.itemsProcessed, m.summarization.totalItems, m.summarization.progressPercent)
    content.WriteString(progressLine)
} else {
    progressLine := fmt.Sprintf("%s %.0f%%", bar, m.summarization.progressPercent)
    content.WriteString(progressLine)
}
```

**Notes:**
- Progress calculation in `events.go` was already correct (cumulative)
- Only the display format needed updating
- Fallback to percentage-only when totalItems is unavailable

---

### 17. ✅ FIXED: Unused Style Definition
**File:** `pkg/executor/tui/styles.go`  
**Severity:** Low  
**Impact:** Compilation error due to undefined color

**Problem:**
```go
commandOutputStyle = lipgloss.NewStyle().
    Foreground(softGray)  // softGray color constant doesn't exist
```

---

### 18. ✅ FIXED: Bash Mode Exit Not Working
**File:** `pkg/executor/tui/update.go`  
**Severity:** High  
**Impact:** Users trapped in bash mode - Escape and Ctrl+C didn't restore normal prompt

**Problem:**
When entering bash mode with `/bash`, pressing Escape or Ctrl+C would set `bashMode = false` but not restore the normal prompt. The `bash>` prompt remained, confusing users.

```go
// Missing updatePrompt() call
func (m model) handleCtrlC() (tea.Model, tea.Cmd) {
    if m.bashMode {
        m.bashMode = false
        m.textarea.Reset()
        m.recalculateLayout()  // Prompt not updated!
        return m, nil
    }
    return m, tea.Quit
}
```

Also, Escape key had no handler for exiting bash mode.

**Solution:**
```go
// Added Escape key handler
case tea.KeyEsc:
    // Escape exits bash mode if active
    if m.bashMode {
        m.bashMode = false
        m.textarea.Reset()
        m.updatePrompt()  // Restore normal prompt
        m.recalculateLayout()
        return m, nil
    }

// Fixed Ctrl+C handler
func (m model) handleCtrlC() (tea.Model, tea.Cmd) {
    if m.bashMode {
        m.bashMode = false
        m.textarea.Reset()
        m.updatePrompt()  // Restore normal prompt
        m.recalculateLayout()
        return m, nil
    }
    return m, tea.Quit
}
```

**Notes:**
- Both Escape and Ctrl+C now properly exit bash mode
- `updatePrompt()` restores the normal "> " prompt
- Consistent with user expectations for modal interfaces


During refactor, `commandOutputStyle` was defined but never used, and referenced a non-existent `softGray` color.

**Solution:**
Removed the unused style definition entirely. Command output formatting was already fixed in regression #15 to use direct string writing without styling.


*Last Updated: 2025-11-18T03:55:00Z*