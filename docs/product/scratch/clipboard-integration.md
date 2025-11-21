# Feature Idea: Clipboard Integration

**Status:** Draft  
**Priority:** Medium Impact, Near-Term  
**Last Updated:** November 2025

---

## Overview

Seamless bidirectional clipboard integration that allows Forge to read from and write to the system clipboard, enabling frictionless data transfer between Forge, other applications, and the terminal. Makes it easy to paste code, share agent outputs, and integrate with external tools.

---

## Problem Statement

Moving data between Forge and other tools is clunky:
- Must manually copy/paste from terminal output
- Cannot easily share agent responses
- Difficult to paste code from browser/IDE into Forge
- No quick way to copy file contents
- Tool output buried in conversation history
- Cannot easily export results

This leads to:
- Retyping code instead of pasting
- Screenshot sharing instead of text
- Lost productivity from manual copying
- Inability to integrate with external tools
- Frustration with terminal limitations

---

## Key Capabilities

### Read from Clipboard
- Paste content directly into chat
- Auto-detect and format code blocks
- Handle multi-line pastes gracefully
- Paste into specific contexts (file content, command args)
- Keyboard shortcut for paste (Cmd/Ctrl+V)

### Write to Clipboard
- Copy agent responses to clipboard
- Copy code blocks with one keystroke
- Copy file contents without opening
- Copy command outputs
- Copy diff hunks

### Smart Clipboard Actions
- Auto-format pasted code
- Detect language from clipboard content
- Offer actions based on clipboard content
- Strip ANSI codes from terminal output
- Handle rich text gracefully (extract plain text)

### Clipboard History
- Access recent clipboard entries
- Search clipboard history
- Pin important clipboard items
- Clear clipboard history
- Export clipboard history

### Cross-Platform Support
- macOS: pbcopy/pbpaste
- Linux: xclip, xsel, wl-clipboard (Wayland)
- Windows: clip.exe, PowerShell
- SSH/Remote: OSC 52 escape sequences
- Fallback: internal clipboard buffer

---

## Example Use Cases

### 1. Paste Code from Browser

```
User: [Copies code from Stack Overflow]
User: [Presses Cmd+V in Forge]

Agent: I see you pasted some code! 

[Detected language: JavaScript]
[Auto-formatted and syntax-highlighted]

function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}

Would you like me to:
1. Add this to a file
2. Explain how it works
3. Adapt it for your codebase
4. Create tests for it
```

### 2. Copy Agent Response

```
Agent: Here's the SQL query you need:

SELECT u.id, u.name, COUNT(o.id) as order_count
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.created_at > NOW() - INTERVAL '30 days'
GROUP BY u.id, u.name
HAVING COUNT(o.id) > 5
ORDER BY order_count DESC;

[User hovers over code block]
[Copy button appears]
[User clicks or presses 'c']

✓ Copied to clipboard!

[User switches to database tool]
[Pastes and runs query]
```

### 3. Quick File Content Copy

```
User: /clipboard copy src/config.json

Agent: Copied file contents to clipboard!

File: src/config.json (234 lines, 5.2 KB)

✓ Ready to paste into another application

[User pastes into Slack to share with team]
```

### 4. Paste Multi-File Context

```
User: [Copies 3 files from IDE]
User: [Pastes into Forge]

Agent: I see you pasted multiple files:

1. auth.go (145 lines)
2. auth_test.go (89 lines)  
3. middleware.go (67 lines)

What would you like me to do with these?
- Review for issues
- Add documentation
- Refactor for better structure
- Create additional tests
```

### 5. Copy Diff for Code Review

```
Agent: Here are the changes I made:

[Shows diff with 3 files changed]

User: /clipboard copy-diff

Agent: Copied diff to clipboard!

Format: Unified diff (ready for `git apply`)
Stats: 3 files, +127 -45 lines

✓ You can now:
- Paste into code review tool
- Email to teammate
- Apply manually with `git apply`
- Share in Slack/Discord
```

---

## Technical Approach

### Platform Detection

**macOS:**
```go
// Write to clipboard
exec.Command("pbcopy").Stdin = strings.NewReader(content)

// Read from clipboard
output, _ := exec.Command("pbpaste").Output()
```

**Linux (X11):**
```go
// Write
exec.Command("xclip", "-selection", "clipboard").Stdin = strings.NewReader(content)

// Read
output, _ := exec.Command("xclip", "-selection", "clipboard", "-o").Output()
```

**Linux (Wayland):**
```go
// Write
exec.Command("wl-copy").Stdin = strings.NewReader(content)

// Read
output, _ := exec.Command("wl-paste").Output()
```

**Windows:**
```go
// Write
exec.Command("clip.exe").Stdin = strings.NewReader(content)

// Read (via PowerShell)
output, _ := exec.Command("powershell", "-command", "Get-Clipboard").Output()
```

### SSH/Remote Support (OSC 52)

**OSC 52 Escape Sequences:**
```go
// Write to clipboard over SSH
func writeClipboardOSC52(content string) {
    encoded := base64.StdEncoding.EncodeToString([]byte(content))
    fmt.Printf("\033]52;c;%s\007", encoded)
}

// Requires terminal emulator support:
// - iTerm2 (macOS)
// - tmux with clipboard enabled
// - Modern terminal emulators
```

**Fallback:**
- Display content for manual copy
- Save to temporary file
- Offer download link (if web TUI in future)

### Content Detection

**Auto-detect Language:**
```go
func detectLanguage(content string) string {
    // Check for shebangs
    if strings.HasPrefix(content, "#!/usr/bin/env python") {
        return "python"
    }
    
    // Check for language-specific syntax
    if strings.Contains(content, "func main()") {
        return "go"
    }
    
    // Check file extensions if pasted with path
    // Use heuristics (keywords, patterns)
    // Fallback to "text"
}
```

**Smart Formatting:**
- Preserve indentation
- Auto-detect code blocks
- Strip markdown formatting if pasted
- Handle mixed content (code + explanation)

### TUI Integration

**Input Handling:**
```go
// Bubble Tea key binding
case key.Matches(msg, m.keyMap.Paste):
    content := getClipboard()
    m.input.SetValue(m.input.Value() + content)
    return m, nil
```

**Copy Actions:**
```go
// Add copy button to code blocks
case key.Matches(msg, m.keyMap.Copy):
    setClipboard(m.currentCodeBlock.Content)
    m.showNotification("✓ Copied to clipboard")
    return m, nil
```

**Clipboard History:**
```go
type ClipboardHistory struct {
    entries []ClipboardEntry
    maxSize int
}

type ClipboardEntry struct {
    content   string
    timestamp time.Time
    contentType string // "code", "text", "diff"
    language  string
    pinned    bool
}
```

---

## Value Propositions

### For All Users
- Seamless data transfer
- Copy/paste just works
- Share outputs easily
- Integration with other tools
- Faster workflows

### For SSH/Remote Users
- Clipboard works over SSH (OSC 52)
- No manual file transfer needed
- Copy/paste between local and remote
- Terminal-native solution

### For Collaboration
- Share code snippets quickly
- Copy diffs for review
- Export results for reports
- Paste context from team members

---

## Implementation Phases

### Phase 1: Basic Copy/Paste (1 week)
- Platform detection
- Read/write clipboard
- Paste into input
- Copy agent responses
- Keyboard shortcuts

### Phase 2: Smart Actions (1 week)
- Auto-detect language
- Format pasted code
- Copy code blocks
- Copy file contents
- Copy diffs

### Phase 3: Clipboard History (1 week)
- Store clipboard entries
- Search history
- Pin important items
- Quick access in TUI
- Export history

### Phase 4: Remote Support (1 week)
- OSC 52 implementation
- Terminal capability detection
- Fallback strategies
- SSH clipboard sync

---

## Open Questions

1. **History Size:** How many clipboard entries to keep?
2. **Privacy:** Should clipboard history be encrypted?
3. **Large Content:** Limit clipboard content size?
4. **Binary Data:** Support binary clipboard content (images)?
5. **Sync:** Sync clipboard across Forge instances?

---

## Related Features

**Synergies with:**
- **TUI** - Keyboard shortcuts and visual feedback
- **File Operations** - Quick copy file contents
- **Diff Viewer** - Copy diffs easily
- **Command Execution** - Copy command output

---

## Success Metrics

**Adoption:**
- 80%+ users paste code into Forge
- 70%+ copy agent responses
- 50%+ use copy-file feature
- 40%+ access clipboard history

**Quality:**
- 95%+ paste success rate
- 90%+ correct language detection
- OSC 52 works on 80%+ supported terminals

**Satisfaction:**
- 4.5+ rating for clipboard features
- "Just works as expected" feedback
- "Love the quick copy buttons" comments

---

## Risks and Mitigations

### Risk: Platform Compatibility
**Impact:** Clipboard doesn't work on some systems  
**Mitigation:**
- Test on all major platforms
- Multiple fallback mechanisms
- Clear error messages
- Manual copy option always available

### Risk: SSH Limitations
**Impact:** OSC 52 not supported by all terminals  
**Mitigation:**
- Detect terminal capabilities
- Fallback to file-based transfer
- Document supported terminals
- Alternative methods (scp, tmux)

### Risk: Large Content
**Impact:** Copying huge files freezes UI  
**Mitigation:**
- Size limits on clipboard operations
- Async clipboard operations
- Progress indicators
- Chunked transfer for large content

---

## Next Steps

1. **Research** - Survey terminal clipboard support
2. **Prototype** - Basic copy/paste (1 week)
3. **Testing** - Test on all platforms
4. **OSC 52** - Implement and test remote support
5. **Polish** - Add smart features and history
6. **Document** - Write user guide for clipboard features
