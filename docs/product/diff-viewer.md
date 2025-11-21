# Product Requirements: Diff Viewer

**Feature:** Interactive Code Change Visualization  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Core Team  
**Last Updated:** December 2024

---

## Product Vision

Empower developers to review AI-generated code changes with complete confidence and clarity. The diff viewer transforms blind approval into informed decision-making by showing exactly what's changing in a beautiful, readable format—making code review feel natural, not tedious.

**Strategic Alignment:** Trust through transparency—when developers can see and understand every change, they gain confidence in AI assistance. This transforms Forge from a black box into a collaborative coding partner.

---

## Problem Statement

Developers reviewing AI-generated code changes face critical challenges that undermine both productivity and code quality:

1. **Blind Approval Risk:** Without seeing diffs, developers either rubber-stamp dangerous changes or waste time manually comparing files
2. **Context Loss:** Changes shown in isolation without surrounding code make impact impossible to assess
3. **Syntax Confusion:** Plain text diffs of complex code are mentally exhausting to parse
4. **Large Change Overwhelm:** Multi-line or multi-file modifications become review nightmares
5. **Navigation Friction:** Can't efficiently jump between scattered changes in large files
6. **Approval Fatigue:** Poor review UX leads to either dangerous rubber-stamping or overly cautious denials

**Current Workarounds (All Problematic):**
- **Approve blindly** → Risky, can introduce bugs or security issues
- **Open external diff tool** → Slow, breaks workflow, context switching
- **Manually compare files** → Time-consuming, error-prone
- **Deny legitimate changes** → Inefficient, slows development

**Real-World Impact:**
- Developer misses subtle bug in AI change because diff was hard to read → Production issue
- Team wastes 10 minutes per review opening files in VSCode to compare → 30+ minutes daily
- Junior developer approves dangerous change without understanding it → Security vulnerability

---

## Key Value Propositions

### For Quality-Focused Developers
- **Confident Decisions:** See exactly what's changing before approving
- **Bug Prevention:** Catch subtle issues through clear visualization
- **Context Understanding:** Sufficient surrounding code to assess impact
- **Professional Review:** Syntax highlighting makes code instantly readable

### For Fast-Moving Developers
- **Efficient Review:** Quick scan with ability to spot major issues
- **Workflow Integration:** Seamless review without leaving Forge
- **Navigation Speed:** Jump between changes, skip unchanged sections
- **Time Savings:** 15-second reviews vs. 2-minute manual comparisons

### For Learning Developers
- **Educational:** Study AI suggestions to learn best practices
- **Comprehension:** Syntax highlighting aids understanding
- **Safe Exploration:** Review without commitment, learn from proposals
- **Pattern Recognition:** See refactoring patterns in action

---

## Target Users & Use Cases

### Primary: Code-Quality-Focused Developer

**Profile:**
- Experienced engineer who reviews all changes carefully
- Values correctness and maintainability
- Comfortable with command-line tools
- Familiar with git diff workflows

**Key Use Cases:**
- Reviewing refactoring changes line-by-line
- Catching subtle logic errors before approval
- Understanding impact of multi-line modifications
- Verifying edge case handling in AI suggestions

**Pain Points Addressed:**
- Can't confidently approve without seeing full context
- Plain text diffs are hard to parse for complex code
- Need to verify AI didn't introduce bugs

**Success Story:**
"I caught a subtle off-by-one error in the diff viewer that I would have missed in plain text. The syntax highlighting made the logic flow clear. I denied the change, explained the issue, and the agent fixed it. Perfect workflow."

---

### Secondary: Fast-Moving Developer

**Profile:**
- Moves quickly but wants safety net
- Trusts AI for most changes but wants sanity check
- Uses keyboard shortcuts extensively
- Values speed without sacrificing quality

**Key Use Cases:**
- Quick review of straightforward changes
- Spotting obviously wrong modifications
- Approving safe changes in seconds
- Jumping between multiple file modifications

**Pain Points Addressed:**
- Manual file comparison breaks flow
- Need fast but safe approval process
- Want to spot major issues without deep review

**Success Story:**
"I can review 5 files in under a minute with the diff viewer. Jump through changes with 'n', approve with 'a'. Fast, safe, and I still catch the important stuff."

---

### Tertiary: Learning Developer

**Profile:**
- Less experienced, learning from AI
- Wants to understand what's changing and why
- Building mental models of good code
- Cautious about approving changes

**Key Use Cases:**
- Studying AI refactoring suggestions
- Learning code patterns from examples
- Understanding best practices through changes
- Building confidence in AI assistance

**Pain Points Addressed:**
- Hard to learn when can't see full context
- Need clear explanation of what's changing
- Want to understand before approving

**Success Story:**
"The diff viewer is like having a mentor show me better ways to write code. I can see the before and after, understand the improvement, and learn the pattern. It's educational."

---

## Product Requirements

### Priority 0 (Must Have)

#### P0-1: Unified Diff Display
**Description:** Traditional git-style diff format with +/- indicators

**User Stories:**
- As a developer, I want to see changes in familiar unified diff format
- As a user, I want clear visual distinction between additions and deletions

**Acceptance Criteria:**
- Display removed lines with red background and - prefix
- Display added lines with green background and + prefix
- Show 3 lines of context before and after changes
- Include line numbers for old and new file
- Syntax highlighting applied to all code

**Example:**
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
```

---

#### P0-2: Side-by-Side Comparison
**Description:** Split-screen view showing old and new code side by side

**User Stories:**
- As a developer, I want to compare old and new code directly
- As a user reviewing refactoring, I want side-by-side view for clarity

**Acceptance Criteria:**
- Left column shows old file content
- Right column shows new file content
- Proper alignment of corresponding lines
- Clear visual separator between columns
- Toggle between unified and side-by-side with 's' key
- Auto-switch to unified on narrow terminals (<120 cols)

**Example:**
```
src/handler.go
Modified • +3 -2 lines

Old (14-16)                    New (14-20)
─────────────────────────────  ─────────────────────────────
14 │     if req == nil {        14 │     if req == nil {
15 │         return errors      15 │         return fmt.Errorf(
        .New("nil request")              "request cannot be nil")
16 │     }                      16 │     }
                                17 │
                                18 │     // Validate fields
                                19 │     if err := req.Validate();
                                        err != nil {
                                20 │         return err
                                21 │     }
```

---

#### P0-3: Syntax Highlighting
**Description:** Language-aware code highlighting in diffs

**User Stories:**
- As a developer, I want syntax highlighting to read code easily
- As a user, I want to quickly identify syntax errors in changes

**Acceptance Criteria:**
- Detect file language from extension (.go, .py, .js, .ts, etc.)
- Apply appropriate syntax highlighting to both old and new code
- Support major languages (Go, Python, JavaScript, TypeScript, Rust, etc.)
- Fall back to plain text for unknown file types
- Highlight change backgrounds (red/green) without obscuring syntax colors
- Performance: Apply highlighting incrementally for large files

**Visual:**
- Keywords highlighted (func, return, if, etc.)
- Strings in distinct color
- Comments grayed out
- Types and functions distinguished
- All layered on top of diff colors (red/green backgrounds)

---

#### P0-4: Efficient Navigation
**Description:** Keyboard shortcuts for quick movement through diffs

**User Stories:**
- As a developer, I want to jump between changes quickly
- As a user reviewing large files, I want to skip unchanged sections

**Acceptance Criteria:**
- 'n' key jumps to next change
- 'p' key jumps to previous change
- 'g' or Home jumps to top
- 'G' or End jumps to bottom
- Arrow keys scroll line by line
- Page Up/Down for larger jumps
- Current change highlighted or indicated

**Keyboard Shortcuts:**
| Key | Action |
|-----|--------|
| n | Next change |
| p | Previous change |
| g / Home | Jump to top |
| G / End | Jump to bottom |
| ↑ / ↓ | Scroll up/down |
| PgUp / PgDn | Page up/down |
| s | Toggle side-by-side |
| a | Approve |
| d | Deny |
| Esc | Close viewer |

---

#### P0-5: File Change Information
**Description:** Clear metadata about what's being changed

**User Stories:**
- As a developer, I want to know which file is being modified
- As a user, I want to see change statistics at a glance

**Acceptance Criteria:**
- File path displayed prominently at top
- Operation type shown (Modified, Created, Deleted)
- Change summary (e.g., "+15 -8 lines")
- File language/type indicated
- Old vs. new file size comparison (for new/deleted files)

**Header Example:**
```
┌─ src/auth/handler.go ─────────────────────────────────────┐
│ Modified • Go • +12 -5 lines                              │
└───────────────────────────────────────────────────────────┘
```

---

#### P0-6: Context Lines
**Description:** Show surrounding code for change comprehension

**User Stories:**
- As a developer, I want to see code around changes to understand context
- As a user, I want enough context to assess impact

**Acceptance Criteria:**
- Default 3 lines of context before each change
- Default 3 lines of context after each change
- Context lines shown in gray/dimmed color
- Clear visual distinction from changed lines
- Collapse large unchanged sections (>10 lines)
- Show "... X lines unchanged ..." for collapsed sections

**Visual:**
```
45 │     // Context before
46 │     existingCode()
47 │     moreContext()
48 -     oldCode()        ← Red: deleted
49 +     newCode()        ← Green: added
50 │     contextAfter()
51 │     moreContextAfter()
```

---

### Priority 1 (Should Have)

#### P1-1: Expandable Context
**Description:** Allow users to see more surrounding code on demand

**User Stories:**
- As a developer, I sometimes need more context to understand a change
- As a user, I want to expand context without leaving diff viewer

**Acceptance Criteria:**
- 'e' key expands context to 10 lines before/after
- 'c' key collapses back to default 3 lines
- Context expansion per-hunk or global setting
- Visual indication of current context amount
- Preserve expansion state while navigating

**Workflow:**
```
User sees change with 3 lines context
    ↓
Not enough context to understand
    ↓
Presses 'e' to expand
    ↓
Now sees 10 lines before/after
    ↓
Understands full function context
    ↓
Makes informed approval decision
```

---

#### P1-2: Multi-File Support
**Description:** Review multiple file changes in single session

**User Stories:**
- As a developer, I want to review all changes together
- As a user, I want to approve/deny files independently

**Acceptance Criteria:**
- Show file list with change counts
- Tab key navigates to next file
- Shift+Tab navigates to previous file
- Approve/deny per file or all at once
- Visual indicator of which files reviewed
- File index (e.g., "File 2 of 5") in header

**Multi-File Interface:**
```
┌─ Review Changes: 3 files ─────────────────────────────────┐
│ [1/3] src/handler.go          Modified • +8 -3            │
│ [2/3] test/handler_test.go    Modified • +12 -0           │
│ [3/3] docs/api.md              Modified • +2 -1            │
│                                                            │
│ Currently viewing: [1/3] src/handler.go                   │
│                                                            │
│ [Tab] Next file  [Shift+Tab] Prev file  [A] Approve All   │
└────────────────────────────────────────────────────────────┘
```

---

#### P1-3: Large Diff Handling
**Description:** Graceful handling of large file modifications

**User Stories:**
- As a developer, I want warnings about huge diffs
- As a user, I want performance even with large files

**Acceptance Criteria:**
- Warn if diff exceeds 1000 lines
- Show summary statistics for huge changes
- Virtual scrolling for smooth performance
- Option to view in external tool for massive diffs
- Progress indicator for syntax highlighting
- "View summary" option for 5000+ line diffs

**Large Diff Warning:**
```
⚠️  Large Diff Warning

This change affects 2,847 lines across multiple functions.

Consider:
• Review in smaller incremental changes
• Use side-by-side mode for clarity
• Open in external diff tool ($DIFFTOOL)

[Continue Review]  [Open in $DIFFTOOL]  [Request Smaller Changes]
```

---

#### P1-4: Special File Operations
**Description:** Handle file creation, deletion, and special cases

**User Stories:**
- As a developer, I want to see new file content clearly
- As a user, I want to understand file deletions and renames

**Acceptance Criteria:**
- **New files:** Show all content as additions (green)
- **Deleted files:** Show all content as deletions (red)
- **Renamed files:** Show old and new paths with rename indicator
- **Binary files:** Display "Binary file changed" message with size
- **Permission changes:** Show mode diff (755 → 644)

**Examples:**

**New File:**
```
┌─ src/auth/validator.go ───────────────────────────────────┐
│ Created • Go • 47 lines                                    │
└────────────────────────────────────────────────────────────┘

1 + package auth
2 +
3 + import "fmt"
4 +
5 + func Validate(token string) error {
...all green...
```

**Deleted File:**
```
┌─ src/deprecated/old.go ───────────────────────────────────┐
│ Deleted • Go • 23 lines removed                            │
└────────────────────────────────────────────────────────────┘

1 - package deprecated
2 -
3 - // Old implementation
...all red...
```

**Renamed File:**
```
┌─ Renamed: auth.go → validator.go ─────────────────────────┐
│ Modified • Go • +5 -2 lines                                │
└────────────────────────────────────────────────────────────┘
```

---

#### P1-5: Help Overlay
**Description:** Discoverable keyboard shortcuts and commands

**User Stories:**
- As a new user, I want to discover available shortcuts
- As a developer, I want quick reference without leaving viewer

**Acceptance Criteria:**
- '?' key shows help overlay
- Lists all keyboard shortcuts with descriptions
- Organized by category (Navigation, Actions, View)
- Dismissable with Esc or any key
- Help button visible in footer

**Help Overlay:**
```
┌─ Diff Viewer Help ────────────────────────────────────────┐
│                                                            │
│ Navigation                                                 │
│   n           Jump to next change                          │
│   p           Jump to previous change                      │
│   g / Home    Jump to top                                  │
│   G / End     Jump to bottom                               │
│   ↑ / ↓       Scroll up/down                               │
│   PgUp/PgDn   Page up/down                                 │
│                                                            │
│ View Options                                               │
│   s           Toggle side-by-side mode                     │
│   e           Expand context                               │
│   c           Collapse context                             │
│                                                            │
│ Actions                                                    │
│   a           Approve changes                              │
│   d           Deny changes                                 │
│   Esc         Close viewer                                 │
│                                                            │
│ Files (multi-file mode)                                    │
│   Tab         Next file                                    │
│   Shift+Tab   Previous file                                │
│   A           Approve all files                            │
│                                                            │
│ Press any key to close help                                │
└────────────────────────────────────────────────────────────┘
```

---

### Priority 2 (Nice to Have)

#### P2-1: Word-Level Diff Highlighting
**Description:** Highlight specific word changes within modified lines

**User Stories:**
- As a developer, I want to see exactly which words changed
- As a user reviewing string changes, I want precise highlighting

**Acceptance Criteria:**
- Highlight changed words within modified lines
- Different color for word-level changes (darker shade)
- Works in both unified and side-by-side modes
- Toggle on/off with setting
- Performance optimized (word diff is expensive)

**Example:**
```
Before word-level diff:
14 - return errors.New("nil request")
15 + return fmt.Errorf("request cannot be nil")

With word-level diff:
14 - return [errors.New]("[nil request]")
15 + return [fmt.Errorf]("[request cannot be nil]")
     └─ Darker background on changed words
```

---

#### P2-2: Diff Export
**Description:** Save diff to file for external review or sharing

**User Stories:**
- As a developer, I want to save diffs for later review
- As a team lead, I want to share diffs with team members

**Acceptance Criteria:**
- Export to unified diff format (.diff, .patch)
- Export to HTML with syntax highlighting
- Export to markdown for documentation
- Configurable output location
- Keyboard shortcut for quick export

---

#### P2-3: Configurable Themes
**Description:** Customizable color schemes for diffs

**User Stories:**
- As a user, I want diff colors that match my terminal theme
- As a developer with accessibility needs, I want high-contrast options

**Acceptance Criteria:**
- Preset themes (GitHub, GitLab, Monokai, Solarized)
- Custom color configuration in settings
- High-contrast mode for accessibility
- Preview theme before applying
- Per-user theme preference persistence

---

## User Experience Flows

### First-Time Review Flow

```
Agent proposes file modification
    ↓
Approval dialog appears
    ↓
User sees: "Press 'v' to view diff"
    ↓
User presses 'v'
    ↓
Diff viewer opens in full screen
    ↓
Shows unified diff with syntax highlighting
    ↓
User sees familiar git-style format
    ↓
"This looks good!"
    ↓
User presses 'a' to approve
    ↓
Changes applied
    ↓
Toast: "Changes approved and applied"
```

**Experience:** Familiar, confidence-building, quick

---

### Complex Refactoring Review Flow

```
Agent proposes refactoring with 127 line changes
    ↓
User opens diff viewer
    ↓
Sees unified diff initially—hard to compare
    ↓
User presses 's' to switch to side-by-side
    ↓
Left: Old implementation | Right: New implementation
    ↓
Much clearer comparison of logic flow
    ↓
User presses 'n' to jump to next change
    ↓
Reviews each change section by section
    ↓
In one section, needs more context
    ↓
Presses 'e' to expand context to 10 lines
    ↓
Now sees full function—understands impact
    ↓
All changes look good
    ↓
Presses 'a' to approve
    ↓
Refactoring applied successfully
```

**Experience:** Flexible, comprehensive, controlled

---

### Multi-File Change Review Flow

```
Agent wants to modify 5 files
    ↓
Diff viewer shows file list
    ↓
┌─ Review Changes: 5 files ────────┐
│ [1/5] handler.go      +12 -3    │
│ [2/5] handler_test.go +25 -0    │
│ [3/5] validator.go    +8 -2     │
│ [4/5] types.go        +3 -1     │
│ [5/5] README.md       +2 -0     │
└──────────────────────────────────┘
    ↓
Starts with file 1: handler.go
    ↓
Reviews changes—looks good
    ↓
Presses Tab to go to next file
    ↓
File 2: handler_test.go
    ↓
New test cases—excellent
    ↓
Presses Tab again
    ↓
File 3: validator.go
    ↓
Wait—this logic looks wrong
    ↓
User presses 'd' to deny this file
    ↓
Continues reviewing other files
    ↓
Approves files 1, 2, 4, 5
    ↓
Denies file 3 with feedback:
"Validator logic incorrect—should check for empty string before length"
    ↓
Agent receives feedback, corrects file 3
    ↓
User reviews updated file 3
    ↓
Now correct—approves
```

**Experience:** Granular control, selective approval, feedback loop

---

### Large Diff Navigation Flow

```
Agent modifies 500-line file with changes throughout
    ↓
Diff viewer opens
    ↓
User sees first change at line 45
    ↓
Presses 'n' to jump to next change
    ↓
Jump to line 127—second change
    ↓
Reviews, presses 'n' again
    ↓
Jump to line 289—third change
    ↓
Continues jumping through changes
    ↓
Reaches last change at line 456
    ↓
Presses 'g' to return to top
    ↓
Quick second pass through all changes
    ↓
Everything checks out
    ↓
Approves
```

**Experience:** Efficient, no wasted scrolling, focused review

---

## User Interface Design

### Unified Diff View

```
┌─ src/auth/handler.go ─────────────────────────────────────┐
│ Modified • Go • +12 -5 lines                    [1/3]     │
├───────────────────────────────────────────────────────────┤
│                                                            │
│  42 │     // Existing function context                    │
│  43 │     existingCode()                                  │
│  44 │     moreExistingCode()                              │
│  45 -     if err != nil {                    ← Red bg     │
│  46 -         return err                                   │
│  47 -     }                                                │
│  48 +     if err := validate(req); err != nil { ← Green   │
│  49 +         return fmt.Errorf("validation: %w", err)    │
│  50 +     }                                                │
│  51 │     processRequest()                                │
│  52 │     moreCode()                                      │
│                                                            │
│     ... 45 lines unchanged ...                            │
│                                                            │
│  97 │     anotherFunction()                               │
│  98 -     oldImplementation()                ← Red bg     │
│  99 +     newImplementation()                ← Green bg   │
│ 100 │     continueProcessing()                            │
│                                                            │
├───────────────────────────────────────────────────────────┤
│ [n]ext  [p]rev  [s]ide-by-side  [a]pprove  [d]eny  [?]help│
└────────────────────────────────────────────────────────────┘
```

**Design Elements:**
- File path and metadata in header
- Line numbers aligned on left
- +/- indicators clear
- Red/green backgrounds for changes
- Gray text for context
- Collapsed unchanged sections
- Keyboard shortcuts in footer
- File counter for multi-file mode

---

### Side-by-Side View

```
┌─ src/auth/handler.go ─────────────────────────────────────┐
│ Modified • Go • +12 -5 lines                    [1/3]     │
├───────────────────────────────────────────────────────────┤
│                                                            │
│ Old (45-47)                │ New (45-50)                  │
│ ───────────────────────────┼───────────────────────────── │
│ 45 │     if err != nil {    │ 45 │     if err :=          │
│      Red background         │              validate(req);  │
│ 46 │         return err     │              err != nil {    │
│                             │      Green background        │
│ 47 │     }                  │ 46 │         return         │
│                             │              fmt.Errorf(     │
│                             │              "validation:    │
│                             │              %w", err)       │
│                             │ 47 │     }                  │
│                             │                              │
│ Context lines in gray       │ Context lines in gray       │
│                             │                              │
├───────────────────────────────────────────────────────────┤
│ [n]ext  [p]rev  [s]ide-by-side  [a]pprove  [d]eny  [?]help│
└────────────────────────────────────────────────────────────┘
```

**Design Elements:**
- Split screen with clear divider
- Aligned corresponding lines
- Old code on left, new on right
- Line numbers for both versions
- Color coding maintained
- Wrapping handled gracefully
- Column width adapts to terminal

---

### Multi-File Selector

```
┌─ Review Changes ──────────────────────────────────────────┐
│                                                            │
│ Select file to review:                                     │
│                                                            │
│ ▸ [1] src/auth/handler.go           Modified  +12 -5     │
│   [2] src/auth/validator.go         Modified  +8 -2      │
│   [3] test/handler_test.go          Modified  +25 -0     │
│   [4] internal/types.go              Modified  +3 -1      │
│   [5] docs/API.md                    Modified  +2 -0      │
│                                                            │
│ ✓ = Reviewed    ✗ = Denied    ○ = Not reviewed           │
│                                                            │
│ [Enter] View selected  [A] Approve all  [D] Deny all      │
└────────────────────────────────────────────────────────────┘
```

---

## Success Metrics

### Adoption Metrics

**Diff Viewer Usage:**
- Target: >90% of file change approvals use diff viewer
- Measure: Approval events with vs. without diff viewing

**Mode Preference:**
- Target: 60% unified, 40% side-by-side
- Measure: Mode selection frequency

**Navigation Usage:**
- Target: >50% use jump keys (n/p) in large diffs
- Measure: Navigation command usage

**Context Expansion:**
- Target: >20% expand context at least once per review
- Measure: Context expansion command usage

---

### Effectiveness Metrics

**Error Detection:**
- Target: >95% of problematic changes caught in review
- Measure: Denied changes that would have caused issues

**Approval Confidence:**
- Target: >90% of users feel confident in approval decisions
- Measure: User survey after review sessions

**Review Time:**
- Target: Average 15 seconds per file change (<50 lines)
- Measure: Time from diff open to approve/deny

**Denial Accuracy:**
- Target: <5% false denials (legitimate changes denied)
- Measure: User reports of denied-then-reapproved changes

---

### Quality Metrics

**Render Performance:**
- Target: p95 under 300ms for typical file (<500 lines)
- Measure: Time from diff requested to fully rendered

**Syntax Highlighting:**
- Target: 100% coverage for supported languages (Go, Python, JS, TS, Rust, etc.)
- Measure: Language detection and highlighting success rate

**Visual Clarity:**
- Target: >95% of users find diffs readable and clear
- Measure: User survey on visual quality

**Navigation Accuracy:**
- Target: 100% of jump commands work correctly
- Measure: Navigation command success rate

---

### User Satisfaction

**Overall Satisfaction:**
- Target: >4.5/5.0 rating for diff viewer experience
- Measure: Post-session user rating

**Feature Discovery:**
- Target: >70% discover side-by-side mode within first 10 reviews
- Measure: First usage of mode toggle

**Productivity Impact:**
- Target: 50% reduction in time spent on code review
- Measure: Compare review time before/after diff viewer

---

## User Enablement

### Discoverability

**First-Time Experience:**
- Approval dialog shows: "Press 'v' to view diff"
- Tutorial overlay on first diff view
- Help button visible in footer
- Keyboard shortcuts in footer bar

**Progressive Disclosure:**
1. **Beginner:** Basic unified view, approve/deny
2. **Intermediate:** Discover side-by-side, navigation shortcuts
3. **Advanced:** Context expansion, multi-file workflow, export

---

### Learning Path

**Session 1-5 (Beginner):**
1. Learn to open diff viewer ('v' from approval)
2. Understand red (removed) vs. green (added)
3. Use basic scrolling
4. Approve/deny with 'a' and 'd'

**Session 6-20 (Intermediate):**
1. Discover side-by-side mode ('s')
2. Learn navigation shortcuts ('n'/'p' for changes)
3. Use jump to top/bottom ('g'/'G')
4. Expand context when needed ('e')

**Session 21+ (Advanced):**
1. Efficiently review multi-file changes
2. Quick navigation through large diffs
3. Use all keyboard shortcuts fluently
4. Customize view preferences

---

### Support Materials

**Documentation:**
1. "Understanding Diffs" - Concept overview
2. "Diff Viewer Quick Start" - 2-minute guide
3. "Keyboard Shortcuts Reference" - Complete list
4. "Best Practices for Code Review" - Review guidelines

**In-App Help:**
- '?' key shows complete keyboard shortcut reference
- Footer shows most common actions
- Tooltips on hover (where applicable)
- Tutorial on first use (dismissable)

**Video Tutorials:**
1. "Your First Code Review" (1 min)
2. "Side-by-Side vs. Unified" (2 min)
3. "Navigating Large Diffs" (3 min)
4. "Multi-File Review Workflow" (4 min)

---

## Risk & Mitigation

### Risk 1: Performance with Large Files
**Impact:** Medium - Slow rendering frustrates users  
**Probability:** Medium  
**User Impact:** Laggy scrolling, delayed syntax highlighting

**Mitigation:**
- Virtual scrolling for large diffs (only render visible lines)
- Lazy syntax highlighting (highlight as user scrolls)
- Warn before opening diffs >1000 lines
- Option to view in external tool for massive files
- Performance benchmarking during development
- Optimize rendering pipeline

**User Communication:**
"Large diffs (>1000 lines) may take a moment to render. Consider requesting smaller, incremental changes from the agent for easier review."

---

### Risk 2: Syntax Highlighting Errors
**Impact:** Low - Falls back to plain text  
**Probability:** Medium  
**User Impact:** Less readable but functional

**Mitigation:**
- Graceful fallback to plain text on error
- Catch and log highlighting exceptions
- Test with diverse code samples
- User can disable highlighting in settings
- Report highlighting issues for language improvement

**Monitoring:**
Track highlighting error rate by language, prioritize fixes for common languages.

---

### Risk 3: Terminal Compatibility Issues
**Impact:** Medium - Rendering broken in some terminals  
**Probability:** Low  
**User Impact:** Garbled display, unreadable diffs

**Mitigation:**
- Test on major terminal emulators (iTerm2, Terminal.app, Windows Terminal, Alacritty, etc.)
- Fallback to simpler format if advanced features fail
- Adaptive width detection (collapse to unified on narrow terminals)
- Clear documentation of terminal requirements
- Support for minimal color modes (16 colors)

**Requirements:**
- ANSI color support
- Unicode support
- Minimum 80 column width (120+ for side-by-side)

---

### Risk 4: Confusing Diff Interpretation
**Impact:** Medium - Users misunderstand changes  
**Probability:** Low  
**User Impact:** Wrong approval decisions

**Mitigation:**
- Clear visual indicators (red = removed, green = added)
- Help text always available ('?' key)
- Intuitive keyboard shortcuts (familiar from git)
- User testing for clarity
- Examples in documentation
- Tutorial on first use

**User Education:**
Video showing how to read diffs, common patterns, what to look for in reviews.

---

### Risk 5: Over-Reliance on Visual Review
**Impact:** Low - Users might miss logical errors  
**Probability:** Medium  
**User Impact:** Bugs slip through despite review

**Mitigation:**
- Encourage understanding, not just visual scanning
- Promote context expansion for deeper review
- Suggest testing after approval
- Educational content on effective code review
- Highlight that diff viewer aids review but doesn't replace critical thinking

**User Guidance:**
"The diff viewer makes it easy to see what's changing, but always think about the logic, edge cases, and potential issues. When in doubt, expand context or ask the agent to explain."

---

## Competitive Analysis

### GitHub Web Diff Viewer
**Approach:** Web-based, rich UI with color coding  
**Strengths:** Beautiful, intuitive, familiar to millions  
**Weaknesses:** Web-only, not integrated into workflow  
**Differentiation:** We bring GitHub-quality diffs to the terminal

### VSCode Diff Editor
**Approach:** Full IDE diff with inline editing  
**Strengths:** Powerful, integrated with editor  
**Weaknesses:** Heavy, requires opening files in VSCode  
**Differentiation:** Lightweight, workflow-integrated, keyboard-first

### Git Diff (Command Line)
**Approach:** Traditional unified diff, plain text  
**Strengths:** Universal, fast, simple  
**Weaknesses:** No syntax highlighting, hard to read  
**Differentiation:** Modern UI with highlighting while keeping speed

### Cursor AI Diff
**Approach:** Inline diffs in editor  
**Strengths:** Immediate, in-context  
**Weaknesses:** Less suitable for large changes  
**Differentiation:** Dedicated full-screen review for thorough examination

### Aider Diff Display
**Approach:** Search/replace block visualization  
**Strengths:** Shows exact blocks being changed  
**Weaknesses:** Non-standard format, less familiar  
**Differentiation:** Industry-standard unified/side-by-side formats

---

## Go-to-Market Considerations

### Positioning

**Primary Message:**  
"Review AI code changes with confidence—Forge's diff viewer brings GitHub-quality diffs to your terminal with syntax highlighting, side-by-side comparison, and keyboard-driven navigation."

**Key Differentiators:**
- Syntax-highlighted diffs (unlike plain git diff)
- Side-by-side comparison (unlike most terminal tools)
- Integrated workflow (unlike VSCode or web tools)
- Keyboard-first navigation (unlike GUI-only tools)
- Multi-file support (review entire changesets)

---

### Target Segments

**Early Adopters:**
- Developers frustrated with plain text diffs
- Teams wanting better AI code review tools
- Power users who live in the terminal

**Value Propositions by Segment:**
- **Terminal Power Users:** "Git-quality diffs without leaving terminal"
- **Code Quality Teams:** "Thorough review UI prevents bugs"
- **Learning Developers:** "Understand AI suggestions through clear visualization"

---

### Documentation Needs

**Essential Documentation:**
1. "Diff Viewer Quick Start" - 5-minute guide
2. "Keyboard Shortcuts Reference" - Complete list
3. "Effective Code Review" - Best practices
4. "Troubleshooting Diff Display" - Common issues

**FAQ Topics:**
- "How do I switch between unified and side-by-side?"
- "What do the colors mean?"
- "How do I jump between changes?"
- "Can I view diffs in my external tool?"
- "Why isn't syntax highlighting working?"

---

## Evolution & Roadmap

### Version History

**v1.0 (Current):**
- Unified and side-by-side diff modes
- Syntax highlighting for major languages
- Keyboard navigation
- Multi-file support
- Context lines with expansion
- File operation indicators

---

### Future Enhancements

#### Phase 2: Enhanced Review Tools
- **Word-Level Diffs:** Highlight changed words within lines
- **Diff Export:** Save diffs as .patch or HTML
- **Inline Comments:** Add notes to specific lines
- **Custom Themes:** User-configurable color schemes
- **External Tool Integration:** Open in $DIFFTOOL

**User Value:** More precise review, better collaboration

---

#### Phase 3: Advanced Features
- **Interactive Editing:** Edit changes directly in viewer
- **Split View Options:** Horizontal vs. vertical split
- **Diff History:** Compare multiple versions
- **AI Explanations:** Agent explains each change inline
- **Approval Annotations:** Mark sections approved/questioned

**User Value:** Richer review experience, better learning

---

#### Phase 4: Team & Collaboration
- **Shared Reviews:** Team members review together
- **Review Templates:** Pre-configured review checklists
- **Quality Gates:** Automated checks before approval
- **Review Analytics:** Track review patterns and quality
- **Comment Threading:** Discuss changes inline

**User Value:** Team collaboration, quality standards

---

## Technical References

- **Architecture:** Diff rendering and syntax highlighting system
- **Implementation:** TUI overlay with Bubble Tea framework
- **Related Features:** Tool Approval PRD, TUI Executor PRD
- **Performance:** Virtual scrolling for large files

---

## Changelog

### 2024-12-XX
- Transformed to product-focused PRD format
- Removed technical implementation details (component diagrams, data models)
- Enhanced user experience sections with detailed flows
- Added comprehensive UI mockups and examples
- Expanded competitive analysis
- Added go-to-market considerations
- Improved success metrics with user-focused KPIs

### 2024-12 (Original)
- Initial PRD with technical architecture
- Component structure and data models
- Rendering flow diagrams
