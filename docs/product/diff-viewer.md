# Product Requirements: Diff Viewer

**Feature:** Interactive Code Change Visualization  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Core Team  
**Last Updated:** December 2024

---

## Overview

The Diff Viewer provides a rich, interactive interface for reviewing code changes before they're applied. It displays file modifications in either unified or side-by-side format with syntax highlighting, making it easy for developers to understand exactly what the agent wants to change and approve or deny modifications confidently.

---

## Problem Statement

Developers reviewing AI-generated code changes face several challenges:

1. **Blind Approval:** Hard to understand changes without seeing diffs
2. **Context Loss:** Changes shown without surrounding code context
3. **Syntax Confusion:** Plain text diffs are hard to parse for complex code
4. **Large Changes:** Difficult to review multi-line or multi-file changes
5. **Navigation Issues:** Can't easily jump between changed sections
6. **Approval Fatigue:** Poor diff UX leads to rubber-stamp approvals

Without a quality diff viewer, users either:
- Approve changes blindly (risky)
- Manually compare files in external tools (slow)
- Deny legitimate changes out of caution (inefficient)

---

## Goals

### Primary Goals

1. **Clear Visualization:** Show exactly what's changing in readable format
2. **Context Preservation:** Display surrounding code for understanding
3. **Syntax Awareness:** Highlight code syntax for easier comprehension
4. **Multiple Formats:** Support both unified and side-by-side views
5. **Easy Navigation:** Quick movement between changes
6. **Confident Approval:** Enable informed approval decisions

### Non-Goals

1. **Inline Editing:** This is NOT an editor (view-only)
2. **3-Way Merge:** Does NOT support complex merge scenarios
3. **Patch Generation:** Does NOT create patch files
4. **Version Control:** Does NOT integrate with git directly
5. **Change Suggestions:** Does NOT propose alternative changes

---

## User Personas

### Primary: Code-Quality-Focused Developer
- **Background:** Experienced developer who reviews all changes carefully
- **Workflow:** Wants to understand every modification before applying
- **Pain Points:** Needs clear visualization to make informed decisions
- **Goals:** Catch errors and understand impact of changes

### Secondary: Fast-Moving Developer
- **Background:** Developer who values speed but wants safety net
- **Workflow:** Quick reviews with ability to spot major issues
- **Pain Points:** Wants efficient review without missing important changes
- **Goals:** Fast but safe approval workflow

### Tertiary: Learning Developer
- **Background:** Less experienced, learning from AI suggestions
- **Workflow:** Studies changes to understand best practices
- **Pain Points:** Needs extra context and clarity
- **Goals:** Learn while reviewing changes

---

## Requirements

### Functional Requirements

#### FR1: Diff Display Modes
- **R1.1:** Unified diff mode (traditional format)
- **R1.2:** Side-by-side diff mode (split view)
- **R1.3:** Toggle between modes with keyboard shortcut
- **R1.4:** Remember user's preferred mode
- **R1.5:** Adapt to terminal width (collapse to unified if too narrow)

#### FR2: Syntax Highlighting
- **R2.1:** Detect file language from extension
- **R2.2:** Apply syntax highlighting to both old and new code
- **R2.3:** Support major languages (Go, Python, JavaScript, TypeScript, etc.)
- **R2.4:** Fall back to plain text for unknown languages
- **R2.5:** Highlight changed lines with different background color

#### FR3: Change Indicators
- **R3.1:** Red background for removed lines (-)
- **R3.2:** Green background for added lines (+)
- **R3.3:** Yellow/orange for modified lines (± in some formats)
- **R3.4:** Line numbers on both sides (old and new)
- **R3.5:** Clear visual separator between old and new code

#### FR4: Context Lines
- **R4.1:** Show 3 lines of context before changes
- **R4.2:** Show 3 lines of context after changes
- **R4.3:** Configurable context amount (1-10 lines)
- **R4.4:** Collapse unchanged sections if >10 lines
- **R4.5:** Expand collapsed sections on demand

#### FR5: Navigation
- **R5.1:** Scroll up/down with arrow keys or mouse
- **R5.2:** Jump to next change with 'n' key
- **R5.3:** Jump to previous change with 'p' key
- **R5.4:** Jump to top with 'g' or Home
- **R5.5:** Jump to bottom with 'G' or End
- **R5.6:** Page up/down for large diffs

#### FR6: File Information
- **R6.1:** Display file path prominently
- **R6.2:** Show file size (old vs new)
- **R6.3:** Indicate operation type (create, modify, delete)
- **R6.4:** Display change summary (e.g., "+15 -8 lines")
- **R6.5:** Show language/file type

#### FR7: Multi-File Support
- **R7.1:** Support reviewing multiple file changes in sequence
- **R7.2:** Navigate between files with Tab/Shift+Tab
- **R7.3:** Show file list with change counts
- **R7.4:** Approve/deny per file or all at once
- **R7.5:** Indicate which files have been reviewed

#### FR8: Diff Overlay Integration
- **R8.1:** Display in modal overlay (full screen)
- **R8.2:** Triggered from tool approval workflow
- **R8.3:** Close with Esc or explicit approve/deny
- **R8.4:** Preserve scroll position if reopened
- **R8.5:** Keyboard-accessible controls

#### FR9: Large Diff Handling
- **R9.1:** Warn if diff exceeds 1000 lines
- **R9.2:** Virtual scrolling for large files
- **R9.3:** Show summary stats for huge changes
- **R9.4:** Option to view in external tool
- **R9.5:** Performance optimization (lazy rendering)

#### FR10: Special Cases
- **R10.1:** New file creation → show all content as added
- **R10.2:** File deletion → show all content as removed
- **R10.3:** File rename → show old and new paths
- **R10.4:** Binary files → indicate "binary file changed"
- **R10.5:** Permission changes → show mode diff

### Non-Functional Requirements

#### NFR1: Performance
- **N1.1:** Render diff under 200ms for typical file (< 500 lines)
- **N1.2:** Smooth scrolling (60 FPS)
- **N1.3:** Syntax highlighting incremental (stream as user scrolls)
- **N1.4:** Handle files up to 5000 lines without lag
- **N1.5:** Memory efficient (< 50MB for typical diff)

#### NFR2: Visual Quality
- **N2.1:** Clear color distinction between additions/deletions
- **N2.2:** Readable on both dark and light terminal themes
- **N2.3:** Proper alignment in side-by-side mode
- **N2.4:** No text wrapping issues with long lines
- **N2.5:** Professional, polished appearance

#### NFR3: Usability
- **N3.1:** Intuitive navigation (familiar to git users)
- **N3.2:** Keyboard shortcuts discoverable
- **N3.3:** Clear indication of current position in diff
- **N3.4:** Helpful error messages for invalid diffs
- **N3.5:** Consistent behavior across file types

#### NFR4: Reliability
- **N4.1:** Never crash on malformed diffs
- **N4.2:** Graceful handling of encoding issues
- **N4.3:** Consistent rendering across terminal emulators
- **N4.4:** Safe handling of very long lines
- **N4.5:** Recovery from syntax highlighting errors

---

## User Experience

### Core Workflows

#### Workflow 1: Simple File Modification Review
1. Agent wants to modify file with `apply_diff`
2. Approval dialog appears
3. User presses 'v' to view diff
4. Diff viewer opens in unified mode
5. Shows 3 added lines, 2 removed lines
6. User reviews changes
7. Looks good, user presses 'a' to approve
8. Changes applied

**Success Criteria:** User reviews and approves in under 10 seconds

#### Workflow 2: Side-by-Side Comparison
1. Agent proposes complex refactoring
2. User opens diff viewer
3. Sees unified diff initially
4. Presses 's' to switch to side-by-side
5. Left shows old code, right shows new code
6. Easier to compare logic changes
7. User approves after thorough review

**Success Criteria:** Side-by-side mode makes comparison clearer

#### Workflow 3: Large Multi-File Change
1. Agent wants to modify 5 files
2. Diff viewer shows file list
3. User tabs through each file
4. Reviews changes one by one
5. File 3 has issue, user denies
6. Provides feedback to agent
7. Agent adjusts and resubmits

**Success Criteria:** User can review each file independently

#### Workflow 4: Navigating Large Diff
1. Agent modifies 300-line file
2. Diff viewer shows changes throughout file
3. User presses 'n' to jump to next change
4. Reviews first change
5. Presses 'n' again to jump to second change
6. Continues through all changes
7. Returns to top with 'g'
8. Approves all changes

**Success Criteria:** Easy navigation through scattered changes

#### Workflow 5: Context Expansion
1. Agent changes middle of function
2. Diff shows 3 lines context
3. User wants more context
4. Presses 'e' to expand context
5. Now sees 10 lines before/after
6. Understands change better
7. Approves with confidence

**Success Criteria:** User can get more context when needed

---

## Technical Architecture

### Component Structure

```
Diff Viewer System
├── Diff Parser
│   ├── Unified Parser
│   ├── Side-by-Side Builder
│   ├── Change Detector
│   └── Line Matcher
├── Syntax Highlighter
│   ├── Language Detector
│   ├── Chroma Integration
│   ├── Color Scheme Manager
│   └── Highlight Cache
├── Diff Renderer
│   ├── Unified Renderer
│   ├── Side-by-Side Renderer
│   ├── Line Number Formatter
│   └── Color Applicator
├── Navigation Controller
│   ├── Change Tracker
│   ├── Jump Handler
│   ├── Scroll Manager
│   └── Position Keeper
└── Diff Overlay (TUI)
    ├── Viewport
    ├── Header/Footer
    ├── Keyboard Handler
    └── State Manager
```

### Data Model

```go
type Diff struct {
    OldFile      FileInfo
    NewFile      FileInfo
    Hunks        []DiffHunk
    Stats        DiffStats
    IsBinary     bool
    Operation    FileOperation
}

type FileInfo struct {
    Path         string
    Content      []string
    Language     string
    LineCount    int
}

type DiffHunk struct {
    OldStart     int
    OldLines     int
    NewStart     int
    NewLines     int
    Changes      []LineChange
    Context      []string
}

type LineChange struct {
    Type         ChangeType  // Add, Delete, Modify, Context
    OldLineNum   int
    NewLineNum   int
    OldContent   string
    NewContent   string
    Highlighted  string
}

type ChangeType int
const (
    ChangeAdd ChangeType = iota
    ChangeDelete
    ChangeModify
    ChangeContext
)

type DiffStats struct {
    FilesChanged int
    Insertions   int
    Deletions    int
}
```

### Rendering Flow

```
Input: Old & New File Content
    ↓
Generate Unified Diff
    ↓
Parse into Hunks
    ↓
For Each Hunk:
  ├─ Detect changed lines
  ├─ Add context lines
  └─ Calculate line numbers
    ↓
Apply Syntax Highlighting
  ├─ Detect language
  ├─ Highlight old content
  └─ Highlight new content
    ↓
Render Based on Mode
  ├─ Unified: Single column with +/-
  └─ Side-by-side: Two columns
    ↓
Apply Color Coding
  ├─ Green background for additions
  ├─ Red background for deletions
  └─ Gray for context
    ↓
Display in Viewport
```

---

## Design Decisions

### Why Support Both Unified and Side-by-Side?
**Rationale:**
- **User preference:** Different developers prefer different formats
- **Use case specific:** Unified better for small changes, side-by-side for refactoring
- **Industry standard:** Both common in git tools
- **Terminal width:** Unified works on narrow terminals

**Decision:** Support both, unified as default

### Why 3 Lines of Context?
**Rationale:**
- **Balance:** Enough context to understand change, not too much clutter
- **Standard:** Git default is 3 lines
- **Screen space:** Fits well on typical terminals
- **Configurable:** Users can adjust if needed

**Testing showed:** 3 lines sufficient for 85% of reviews

### Why Syntax Highlighting in Diffs?
**Rationale:**
- **Readability:** Much easier to parse highlighted code
- **Error detection:** Syntax errors stand out
- **Professional:** Matches modern git UI tools
- **Adoption:** Users expect it from GitHub, GitLab, etc.

**User testing:** 90% preferred highlighted diffs

### Why Modal Overlay vs Inline?
**Rationale:**
- **Focus:** Full-screen diff is less distracting
- **Space:** More room for side-by-side comparison
- **Workflow:** Clear approve/deny decision point
- **Familiar:** Matches git commit UI patterns

---

## Display Formats

### Unified Diff Format

```
src/handler.go
Modified • +3 -2 lines

12 │ func handleRequest(req *Request) error {
13 │     if req == nil {
14 -         return errors.New("nil request")
15 +         return fmt.Errorf("request cannot be nil")
16 │     }
17 │
18 +     // Validate request fields
19 +     if err := req.Validate(); err != nil {
20 +         return err
21 +     }
22 │
23 │     return processRequest(req)
```

Key:
- Gray: Context lines (unchanged)
- Red background (-): Deleted lines
- Green background (+): Added lines
- Line numbers on left

---

### Side-by-Side Format

```
src/handler.go
Modified • +3 -2 lines

Old (14-16)                           New (14-20)
────────────────────────────────────  ────────────────────────────────────
14 │     if req == nil {              14 │     if req == nil {
15 │         return errors.New(       15 │         return fmt.Errorf(
     "nil request")                         "request cannot be nil")
16 │     }                            16 │     }
                                      17 │
                                      18 │     // Validate request fields
                                      19 │     if err := req.Validate(); 
                                              err != nil {
                                      20 │         return err
                                      21 │     }
```

Key:
- Left column: Old file
- Right column: New file  
- Red lines: Removed
- Green lines: Added
- Aligned context

---

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| ↑/↓ | Scroll up/down |
| PgUp/PgDn | Page up/down |
| n | Next change |
| p | Previous change |
| g/Home | Jump to top |
| G/End | Jump to bottom |
| s | Toggle side-by-side mode |
| e | Expand context |
| c | Collapse context |
| a | Approve changes |
| d | Deny changes |
| ? | Show help |
| Esc | Close viewer |

---

## Success Metrics

### Usage Metrics
- **Viewer access:** >90% of file change approvals use diff viewer
- **Mode preference:** Unified 60%, Side-by-side 40%
- **Navigation usage:** >50% use jump keys (n/p)
- **Context expansion:** >20% expand context at least once

### Effectiveness Metrics
- **Error detection:** >95% of problematic changes caught in review
- **Approval confidence:** >90% of users confident in approval decisions
- **Review time:** Average 15 seconds per file change
- **Denial accuracy:** <5% false denials (legitimate changes denied)

### Quality Metrics
- **Render performance:** p95 under 300ms
- **Syntax highlighting:** 100% coverage for supported languages
- **Visual clarity:** >95% of users find diffs readable
- **Navigation accuracy:** 100% of jump commands work correctly

---

## Dependencies

### External Dependencies
- Syntax highlighting library (Chroma)
- Diff generation utilities
- Text processing libraries

### Internal Dependencies
- TUI framework (Bubble Tea)
- Tool approval system
- Settings system (for preferences)
- File system access

### Platform Requirements
- Terminal with ANSI color support
- Unicode support (for diff symbols)
- Minimum 80 column width (120+ recommended for side-by-side)

---

## Risks & Mitigations

### Risk 1: Performance with Large Files
**Impact:** Medium  
**Probability:** Medium  
**Mitigation:**
- Virtual scrolling for large diffs
- Lazy syntax highlighting
- Warn before opening huge diffs
- Option to view in external tool
- Optimize rendering pipeline

### Risk 2: Syntax Highlighting Errors
**Impact:** Low  
**Probability:** Medium  
**Mitigation:**
- Fallback to plain text on error
- Catch and log highlighting exceptions
- Test with diverse code samples
- Graceful degradation
- User can disable highlighting

### Risk 3: Terminal Compatibility
**Impact:** Medium  
**Probability:** Low  
**Mitigation:**
- Test on major terminal emulators
- Fallback to simpler format if needed
- Adaptive width detection
- Clear documentation of requirements
- Support for minimal color modes

### Risk 4: Confusing Diffs
**Impact:** Medium  
**Probability:** Low  
**Mitigation:**
- Clear visual indicators
- Help text available
- Intuitive keyboard shortcuts
- User testing for clarity
- Examples in documentation

---

## Future Enhancements

### Phase 2 Ideas
- **Inline Comments:** Add notes to specific lines
- **Change Suggestions:** Propose modifications to agent's changes
- **Diff Export:** Save diffs to file
- **Split View Options:** Horizontal vs vertical split
- **Word-Level Diff:** Highlight changed words within lines

### Phase 3 Ideas
- **Interactive Editing:** Edit changes directly in viewer
- **3-Way Merge:** Support complex merge scenarios
- **Git Integration:** Show git-style patches
- **Diff History:** Compare multiple versions
- **AI Explanations:** Agent explains each change

---

## Open Questions

1. **Should we support inline editing of changes?**
   - Pro: Fix small issues without denying
   - Con: Complexity, potential for errors
   - Decision: Phase 3 feature if high demand

2. **Should we show word-level diffs?**
   - Pro: Highlights exact changes in modified lines
   - Con: More complex rendering
   - Decision: Phase 2 experiment

3. **Should we integrate with git diff?**
   - Pro: Familiar format for git users
   - Con: Adds dependency, not all workspaces have git
   - Decision: Keep internal diff format, consider git export

4. **Should we support external diff tools?**
   - Pro: Some users prefer their own tools
   - Con: Less integrated experience
   - Decision: Add "open in $DIFFTOOL" option in Phase 2

---

## Related Documentation

- [ADR-0019: Diff Viewer in Approval Flow](../adr/0019-diff-viewer-approval.md)
- [Tool Approval System PRD](tool-approval-system.md)
- [TUI Executor PRD](tui-executor.md)
- [How-to: Use TUI Interface - Diff Viewer](../how-to/use-tui-interface.md#diff-viewer)

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2024-12 | 1.0 | Initial PRD creation |
