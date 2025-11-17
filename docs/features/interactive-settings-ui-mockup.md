# Interactive Settings UI Mockup

This document provides a visual representation of the interactive settings editor accessible via `/settings` command.

## Color Palette

Using Forge's official color scheme:
- **Salmon Pink** (#FFB3BA) - Primary accent, selections, active elements
- **Mint Green** (#A8E6CF) - Success states, enabled toggles
- **Muted Gray** (#6B7280) - Disabled states, descriptions
- **Dark Background** (#111827) - Main background
- **Bright White** (#F9FAFB) - Primary text

---

## Layout Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          ⚙️  Settings                                    │
│                                                                          │
│  ↑↓/jk: Navigate • Tab/←→/hl: Switch section • Space/Enter: Toggle      │
│  Ctrl+S: Save • Esc/q: Close                                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ▸ Auto-Approval Settings                    [Section 1 of 2]          │
│    Configure which tools can execute without manual approval            │
│                                                                          │
│    ▸ [ ] read_file                                                      │
│      [✓] list_files                                                     │
│      [ ] search_files                                                   │
│      [ ] write_file                                                     │
│      [ ] apply_diff                                                     │
│      [ ] execute_command                                                │
│      [✓] task_completion                                                │
│      [✓] ask_question                                                   │
│      [✓] converse                                                       │
│                                                                          │
│  ▸ Command Whitelist                        [Section 2 of 2]          │
│    Patterns for commands that auto-approve                              │
│                                                                          │
│      ✓ npm - All npm commands                                           │
│      ✓ git status - Git status check                                    │
│      ✓ ls - List directory contents                                     │
│                                                                          │
│                                                                          │
├─────────────────────────────────────────────────────────────────────────┤
│  ● Unsaved changes - Press Ctrl+S to save                              │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## State 1: Initial View (No Selection)

When user first opens `/settings`:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          ⚙️  Settings                                    │
│                                                                          │
│  ↑↓/jk: Navigate • Tab/←→/hl: Switch section • Space/Enter: Toggle      │
│  Ctrl+S: Save • Esc/q: Close                                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ▸ Auto-Approval Settings                                              │
│    Configure which tools can execute without manual approval            │
│                                                                          │
│    [ ] read_file                                                        │
│    [ ] list_files                                                       │
│    [ ] search_files                                                     │
│    [ ] write_file                                                       │
│    [ ] apply_diff                                                       │
│    [ ] execute_command                                                  │
│    [✓] task_completion                                                  │
│    [✓] ask_question                                                     │
│    [✓] converse                                                         │
│                                                                          │
│  ▸ Command Whitelist                                                   │
│    Patterns for commands that auto-approve                              │
│                                                                          │
│    (no patterns configured)                                             │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Colors:**
- Title "⚙️ Settings": Salmon Pink (#FFB3BA)
- Help text: Muted Gray (#6B7280)
- Section titles: Mint Green (#A8E6CF)
- Descriptions: Muted Gray (#6B7280)
- Enabled checkboxes [✓]: Mint Green (#A8E6CF)
- Disabled checkboxes [ ]: Muted Gray (#6B7280)
- Tool names: Bright White (#F9FAFB)

---

## State 2: Item Selected (Salmon Pink Highlight)

When user navigates to "list_files":

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          ⚙️  Settings                                    │
│                                                                          │
│  ↑↓/jk: Navigate • ←→/hl: Switch section • Space/Enter: Toggle          │
│  Ctrl+S: Save • Esc/q: Close                                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ▸ Auto-Approval Settings                                              │
│    Configure which tools can execute without manual approval            │
│                                                                          │
│    [ ] read_file                                                        │
│  ▸ [ ] list_files                    ← [SALMON PINK HIGHLIGHT]         │
│    [ ] search_files                                                     │
│    [ ] write_file                                                       │
│    [ ] apply_diff                                                       │
│    [ ] execute_command                                                  │
│    [✓] task_completion                                                  │
│    [✓] ask_question                                                     │
│    [✓] converse                                                         │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Selection Indicator:**
- "▸" prefix in Salmon Pink
- Entire line has Salmon Pink foreground
- Bold text weight for selected item

---

## State 3: After Toggle (Modified State)

After pressing Space on "list_files":

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          ⚙️  Settings                                    │
│                                                                          │
│  ↑↓/jk: Navigate • ←→/hl: Switch section • Space/Enter: Toggle          │
│  Ctrl+S: Save • Esc/q: Close                                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ▸ Auto-Approval Settings                                              │
│    Configure which tools can execute without manual approval            │
│                                                                          │
│    [ ] read_file                                                        │
│  ▸ [✓] list_files *                  ← [SELECTED + MODIFIED]           │
│    [ ] search_files                                                     │
│    [ ] write_file                                                       │
│    [ ] apply_diff                                                       │
│    [ ] execute_command                                                  │
│    [✓] task_completion                                                  │
│    [✓] ask_question                                                     │
│    [✓] converse                                                         │
│                                                                          │
├─────────────────────────────────────────────────────────────────────────┤
│  ● Unsaved changes - Press Ctrl+S to save                              │
└─────────────────────────────────────────────────────────────────────────┘
```

**Modified Indicator:**
- Asterisk (*) in Salmon Pink after item name
- Status bar appears at bottom showing unsaved changes
- Checkbox changed from [ ] to [✓] in Mint Green

---

## State 4: Section Switching (Command Whitelist)

When user presses Tab or → to switch sections:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          ⚙️  Settings                                    │
│                                                                          │
│  ↑↓/jk: Navigate • ←→/hl: Switch section • Space/Enter: Toggle          │
│  Ctrl+S: Save • Esc/q: Close                                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ▸ Auto-Approval Settings                [Section 1 of 2]              │
│    Configure which tools can execute without manual approval            │
│                                                                          │
│    [ ] read_file                                                        │
│    [✓] list_files *                                                     │
│    [ ] search_files                                                     │
│    [ ] write_file                                                       │
│                                                                          │
│  ▸ Command Whitelist                    [Section 2 of 2] ← [ACTIVE]    │
│    Patterns for commands that auto-approve                              │
│                                                                          │
│  ▸ ✓ npm - All npm commands              ← [SELECTED]                  │
│    ✓ git status - Git status check                                      │
│    ✓ ls - List directory contents                                       │
│                                                                          │
├─────────────────────────────────────────────────────────────────────────┤
│  ● Unsaved changes - Press Ctrl+S to save                              │
└─────────────────────────────────────────────────────────────────────────┘
```

**Section Switching:**
- Active section title in Salmon Pink
- Inactive section title in Mint Green
- First item in new section auto-selected
- Section indicator shows current/total

---

## State 5: Scrolling (Long Lists)

When list exceeds viewport height:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          ⚙️  Settings                                    │
│                                                                          │
│  ↑↓/jk: Navigate • ←→/hl: Switch section • Space/Enter: Toggle          │
│  Ctrl+S: Save • Esc/q: Close                                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ▸ Auto-Approval Settings                                              │
│    Configure which tools can execute without manual approval            │
│                                                                          │
│    [ ] apply_diff                                                       │
│    [ ] execute_command                                                  │
│  ▸ [✓] task_completion                   ← [SELECTED]                  │
│    [✓] ask_question                                                     │
│    [✓] converse                                                         │
│    [ ] custom_tool_1                                                    │
│    [ ] custom_tool_2                                                    │
│                       ⋮                                                 │
│                      [8/15]                  ← [SCROLL INDICATOR]       │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Scrolling Behavior:**
- Auto-scroll to keep selected item visible
- Scroll indicator shows position (current/total)
- Ellipsis (⋮) indicates more content above/below

---

## Keyboard Shortcuts Reference

### Navigation
- **↑ / k**: Move selection up
- **↓ / j**: Move selection down  
- **← / h**: Previous section
- **→ / l**: Next section
- **Tab**: Next section
- **Shift+Tab**: Previous section
- **Home**: Jump to first item
- **End**: Jump to last item

### Actions
- **Space**: Toggle current checkbox
- **Enter**: Toggle current checkbox / Confirm input
- **a**: Add new item (in command whitelist section)
- **d**: Delete selected item (in command whitelist section)
- **e**: Edit selected item (in command whitelist section)
- **Ctrl+S**: Save all changes
- **Esc**: Cancel current action / Close settings (prompts if unsaved)
- **q**: Close settings (prompts if unsaved changes)

### Future (Not Yet Implemented)
- **r**: Reset section to defaults
- **/**: Search/filter settings
---

## Adding/Editing Command Whitelist Items

### State 6: Adding a New Command Pattern

When user is in Command Whitelist section and presses **a**:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          ⚙️  Settings                                    │
│                                                                          │
│  Enter: Confirm • Esc: Cancel                                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ▸ Command Whitelist                    [Section 2 of 2]              │
│    Patterns for commands that auto-approve                              │
│                                                                          │
│    ✓ npm - All npm commands                                             │
│    ✓ git status - Git status check                                      │
│    ✓ ls - List directory contents                                       │
│                                                                          │
│  ┌─ Add New Command Pattern ──────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  Pattern (command or prefix):                                   │   │
│  │  ▸ docker_                              ← [TEXT INPUT CURSOR]   │   │
│  │                                                                  │   │
│  │  Description:                                                    │   │
│  │    Docker commands                      ← [TEXT INPUT]          │   │
│  │                                                                  │   │
│  │  Pattern type:                                                   │   │
│  │    ○ Prefix match    ● Exact match      ← [RADIO SELECTION]    │   │
│  │                                                                  │   │
│  │                 [Enter to Add] [Esc to Cancel]                  │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Input Dialog Features:**
- Modal overlay with salmon pink border
- Text input fields with cursor indicator (▸)
- Radio button selection for pattern type
- Tab to move between fields
- Enter to confirm, Esc to cancel
- Real-time validation feedback

---

### State 7: Editing an Existing Pattern

When user selects a command pattern and presses **e**:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          ⚙️  Settings                                    │
│                                                                          │
│  Enter: Save changes • Esc: Cancel                                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ▸ Command Whitelist                    [Section 2 of 2]              │
│    Patterns for commands that auto-approve                              │
│                                                                          │
│  ▸ ✓ npm - All npm commands              ← [SELECTED FOR EDIT]         │
│    ✓ git status - Git status check                                      │
│    ✓ ls - List directory contents                                       │
│                                                                          │
│  ┌─ Edit Command Pattern ───────────────────────────────────────────┐  │
│  │                                                                   │  │
│  │  Pattern (command or prefix):                                    │  │
│  │  ▸ npm                                  ← [TEXT INPUT CURSOR]    │  │
│  │                                                                   │  │
│  │  Description:                                                     │  │
│  │    All npm commands_                    ← [TEXT INPUT]           │  │
│  │                                                                   │  │
│  │  Pattern type:                                                    │  │
│  │    ● Prefix match    ○ Exact match      ← [RADIO SELECTION]     │  │
│  │                                                                   │  │
│  │                [Enter to Save] [Esc to Cancel]                   │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Edit Dialog:**
- Pre-populated with existing values
- Same interface as add dialog
- Changes marked with asterisk when modified

---

### State 8: Delete Confirmation

When user selects a pattern and presses **d**:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          ⚙️  Settings                                    │
│                                                                          │
│  y: Confirm delete • n/Esc: Cancel                                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ▸ Command Whitelist                    [Section 2 of 2]              │
│    Patterns for commands that auto-approve                              │
│                                                                          │
│    ✓ npm - All npm commands                                             │
│  ▸ ✓ git status - Git status check       ← [SELECTED FOR DELETE]       │
│    ✓ ls - List directory contents                                       │
│                                                                          │
│  ┌─ Confirm Delete ──────────────────────────────────────────────────┐  │
│  │                                                                    │  │
│  │  ⚠️  Are you sure you want to delete this pattern?                │  │
│  │                                                                    │  │
│  │  Pattern: git status                                              │  │
│  │  Description: Git status check                                    │  │
│  │                                                                    │  │
│  │  This command will require manual approval after deletion.        │  │
│  │                                                                    │  │
│  │              [y] Yes, delete    [n] No, cancel                    │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Delete Confirmation:**
- Warning icon (⚠️) in salmon pink
- Shows what will be deleted
- Explains consequences
- y/n keyboard shortcuts

---

### Adding Command Pattern: Full Workflow

**Step-by-Step Process:**

1. **Navigate to Command Whitelist section**
   - Use Tab or arrow keys to switch to Command Whitelist

2. **Press 'a' to add new pattern**
   - Add dialog overlay appears
   - Focus on Pattern field

3. **Enter pattern**
   - Type command or prefix (e.g., `docker`, `python`, `npm install`)
   - Tab to move to next field

4. **Enter description** (optional but recommended)
   - Type human-readable description
   - Tab to pattern type

5. **Select pattern type**
   - Use Space or arrow keys to toggle:
     - **Prefix match**: Pattern matches commands starting with this text
     - **Exact match**: Pattern must match command exactly

6. **Confirm or cancel**
   - Press Enter to add the pattern
   - Press Esc to cancel without adding

7. **Pattern appears in list**
   - New pattern added with asterisk (modified)
   - Status bar shows unsaved changes
   - Press Ctrl+S to save

**Pattern Type Examples:**

**Prefix Match** (recommended for most cases):
- Pattern: `npm` → Matches: `npm install`, `npm run dev`, `npm test`
- Pattern: `git` → Matches: `git status`, `git commit`, `git push`
- Pattern: `docker` → Matches: `docker ps`, `docker build`, `docker run`

**Exact Match** (for specific commands):
- Pattern: `ls -la` → Matches: `ls -la` only (not `ls`)
- Pattern: `git status` → Matches: `git status` only
- Pattern: `make test` → Matches: `make test` only

---

### Validation Rules

**Pattern Field:**
- Cannot be empty
- Must not contain only whitespace
- Shows error: "Pattern cannot be empty" in red

**Duplicate Detection:**
- Shows warning: "Pattern already exists" in yellow
- Allows saving but highlights the duplicate

**Description Field:**
- Optional
- Max 100 characters
- Shows character count when typing

**Real-time Feedback:**
```
┌─ Add New Command Pattern ──────────────────────────────────────┐
│                                                                 │
│  Pattern (command or prefix):                                  │
│  ▸ npm                                  ✓ Valid                │
│                                                                 │
│  Description:                                           [25/100]│
│    Node package manager_                                        │
│                                                                 │
│  Pattern type:                                                  │
│    ● Prefix match    ○ Exact match                             │
│                                                                 │
│                 [Enter to Add] [Esc to Cancel]                 │
└─────────────────────────────────────────────────────────────────┘
```


---

## Visual States Summary

### Item States
1. **Unselected + Disabled**: Gray checkbox [ ], white text
2. **Unselected + Enabled**: Green checkbox [✓], white text
3. **Selected + Disabled**: Salmon pink ▸ prefix, bold salmon pink text, gray checkbox
4. **Selected + Enabled**: Salmon pink ▸ prefix, bold salmon pink text, green checkbox
5. **Modified**: Salmon pink asterisk (*) suffix

### Section States
1. **Inactive Section**: Mint green title
2. **Active Section**: Salmon pink title
3. **Collapsed Section** (future): ▸ arrow, only title visible
4. **Expanded Section**: ▾ arrow, all items visible

### Overlay States
1. **No Changes**: No status bar
2. **Unsaved Changes**: Salmon pink status bar with save prompt
3. **Saving**: "Saving..." message in status bar
4. **Saved**: Brief "Saved!" toast, then disappears
5. **Error**: Red error message in status bar

---

## Example: Complete Interaction Flow

### Step 1: Open Settings
```
User types: /settings
Result: Overlay appears, first item selected
```

### Step 2: Navigate to Tool
```
User presses: ↓ ↓ ↓ (3 times)
Result: Selection on "write_file"
```

### Step 3: Enable Auto-Approval
```
User presses: Space
Result: [✓] appears, * indicator added, status bar shows unsaved
```

### Step 4: Switch to Whitelist
```
User presses: Tab
Result: Command Whitelist section selected
```

### Step 5: Save Changes
```
User presses: Ctrl+S
Result: "Saving..." → "Saved!" toast → status bar disappears
```

### Step 6: Close Settings
```
User presses: Esc
Result: Overlay closes, returns to conversation
```

---

## Implementation Notes

The UI is rendered using:
- **Bubble Tea**: TUI framework for state management and rendering
- **Lipgloss**: Terminal styling for colors, borders, and layout
- **Custom rendering**: Line-by-line composition to avoid padding issues

Key rendering functions:
- `View()`: Main render method, composes entire overlay
- `renderSection()`: Renders a section with its items
- `renderItem()`: Renders individual toggles/list items
- `renderToggle()`: Formats checkbox and label
- `buildHelpText()`: Dynamically shows relevant shortcuts
