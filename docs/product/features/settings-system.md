# Product Requirements: Settings System

**Feature:** Settings and Configuration Management  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Core Team  
**Last Updated:** December 2024

---

## Product Vision

Transform configuration from a technical hurdle into an intuitive, discoverable experience. The Settings System empowers every user—from beginners to power users—to customize Forge effortlessly through a visual interface that makes preferences accessible, safe, and persistent without ever touching a config file.

**Strategic Alignment:** Great tools adapt to users, not the other way around. By making configuration accessible and safe, we remove friction, increase adoption, and enable users to optimize Forge for their unique workflows.

---

## Problem Statement

Developers configuring AI coding assistants face a maze of frustration that limits adoption and creates support burden:

1. **Configuration Chaos:** Settings scattered across CLI flags, environment variables, config files, and undocumented defaults—users never know where to look
2. **Discoverability Crisis:** Users have no idea what's configurable. "Can I change the model?" "Is there a way to auto-approve reads?" → endless questions
3. **Manual Editing Hell:** Editing JSON/YAML configs is error-prone (syntax errors, invalid values, wrong types) and requires leaving the application
4. **Lost Settings:** Typos in config files crash the app or silently ignore settings. No validation, no feedback.
5. **Security Confusion:** "Where do I put my API key? Is it safe? Will it leak?"
6. **Onboarding Friction:** New users can't start without finding, editing, and validating config files—instant abandonment

**Current Workarounds (All Problematic):**
- **Manually edit config files** → Syntax errors, no validation, requires leaving TUI
- **Set environment variables** → Hard to remember, not discoverable, different per shell
- **Use CLI flags** → Tedious to repeat, not persistent
- **Read documentation** → Time-consuming, often outdated, fragmented

**Real-World Impact:**
- New user tries Forge, can't figure out how to set API key → gives up, never returns
- Developer changes model in config file, typos model name → app crashes on startup
- Team lead wants to standardize settings → everyone manually edits different files, inconsistency
- Power user wants auto-approval for docs folder → searches for 20 minutes, can't find how
- Security-conscious user stores API key in config → doesn't know file should be private (600 permissions)

**Cost of Poor Configuration:**
- 35% of new users abandon during initial setup
- 50+ support tickets per month about "how do I configure X?"
- Average 15 minutes wasted per user fixing config syntax errors
- Trust erosion when settings mysteriously don't work

---

## Key Value Propositions

### For New Users (Getting Started)
- **Zero Config File Editing:** Set up API keys, models, and preferences through visual interface
- **Guided Configuration:** Clear descriptions explain what each setting does and why you'd change it
- **Safe Defaults:** Works great out-of-the-box, customize only what you want
- **Instant Validation:** Know immediately if a setting is invalid, with helpful error messages
- **Confidence Building:** No fear of breaking things—settings are validated and recoverable

### For Power Users (Customization)
- **Complete Control:** Every aspect of Forge behavior is configurable
- **Keyboard-Driven:** Navigate and change settings without mouse/clicking
- **Quick Access:** Open settings overlay with `/settings`, make changes in seconds
- **Advanced Features:** Auto-approval rules, custom LLM parameters, display tweaks
- **Export/Import:** Share configurations, backup settings, replicate across machines

### For Team Leads (Standardization)
- **Consistent Setup:** Easily replicate settings across team members
- **Documented Preferences:** Settings are self-documenting with inline help
- **Secure by Default:** API keys stored with proper permissions automatically
- **Reliable Behavior:** No surprises from hidden or undocumented configuration
- **Easy Troubleshooting:** See exactly what's configured at a glance

---

## Target Users & Use Cases

### Primary: First-Time User (Setup Phase)

**Profile:**
- Just installed Forge, excited to try AI coding
- Has API key from LLM provider (OpenAI/Anthropic)
- Wants to get started quickly
- Limited patience for configuration complexity
- Might not be comfortable with command line

**Key Use Cases:**
- Initial API key configuration
- Setting workspace directory
- Choosing LLM model
- Understanding what's configurable

**Pain Points Addressed:**
- No idea where to put API key
- Fear of breaking something by editing wrong file
- Overwhelmed by options
- Can't find configuration documentation

**Success Story:**
"I just installed Forge and typed '/settings'. A clean interface appeared with tabs for different settings. I saw 'LLM Settings' tab, clicked it, and there was a field for 'API Key' with a description explaining exactly what to enter. I pasted my OpenAI key, selected GPT-4 from a dropdown, and closed the settings. It just worked. No config files, no documentation hunting, no terminal commands. Perfect."

**User Journey:**
```
Install Forge
    ↓
Launch for first time
    ↓
Agent says: "No API key configured. Type /settings to configure."
    ↓
User types /settings
    ↓
Settings overlay opens (clean, organized)
    ↓
Navigate to LLM tab (Tab key or click)
    ↓
See "API Key" field with help text
    ↓
Paste API key
    ↓
Select model from dropdown (gpt-4, gpt-3.5-turbo, etc.)
    ↓
Settings auto-save
    ↓
Close overlay (Esc)
    ↓
Agent ready to use: "Configuration saved! Let's start coding."
    ↓
Success - user productive in <2 minutes
```

---

### Secondary: Power User (Optimization Phase)

**Profile:**
- Experienced with AI coding tools
- Wants to optimize workflow through configuration
- Uses keyboard shortcuts extensively
- Values efficiency and control
- Runs Forge across multiple projects

**Key Use Cases:**
- Fine-tuning LLM parameters (temperature, max tokens)
- Configuring auto-approval rules for trusted paths
- Customizing display preferences (theme, diff style)
- Setting up workspace-specific preferences
- Exporting/importing configurations

**Pain Points Addressed:**
- Can't remember environment variable names
- Editing config files breaks concentration
- Want granular control without complexity
- Need to replicate setup across machines

**Success Story:**
"I'm tired of approving every file read in my docs/ folder. I opened settings with Ctrl+comma (muscle memory from VSCode), went to Auto-Approval tab, saw 'Auto-approve reads' toggle and 'Path Whitelist' list. I enabled the toggle, added 'docs/*' pattern, and boom—no more approval prompts for documentation. Took 20 seconds, never left my TUI, and I can see exactly what I configured."

**Advanced Configuration Flow:**
```
Working on project, frequent interruptions for approvals
    ↓
User thinks: "I trust reads in docs/ and tests/"
    ↓
Open settings (Ctrl+, or /settings)
    ↓
Navigate to Auto-Approval tab
    ↓
Enable "Auto-approve reads" toggle
    ↓
Add to path whitelist:
    - docs/*
    - tests/*
    - README.md
    ↓
See live preview of what will auto-approve
    ↓
Close settings
    ↓
Next read in docs/ → auto-approved silently
    ↓
Workflow smoother, user happy
```

---

### Tertiary: Team Lead (Standardization Phase)

**Profile:**
- Manages team of developers using Forge
- Wants consistent setup across team
- Needs to troubleshoot configuration issues
- Values reliability and predictability
- Security-conscious about API keys

**Key Use Cases:**
- Creating standard configuration for team
- Exporting settings template for new team members
- Verifying team member configurations
- Ensuring security best practices
- Documenting team conventions

**Pain Points Addressed:**
- Team has inconsistent setups → unpredictable behavior
- Hard to help team members troubleshoot
- API keys stored insecurely by some users
- No visibility into what each person configured

**Success Story:**
"I configured Forge exactly how our team should use it: Claude Sonnet model, auto-approve rules for our monorepo structure, side-by-side diffs. I exported my settings to a JSON file, added it to our onboarding repo with instructions: 'Run Forge, type /settings, click Import, select team-settings.json'. Every new hire gets perfect setup in 30 seconds. No more 'why does this work differently for me?' questions."

**Team Standardization Flow:**
```
Team lead configures ideal setup
    ↓
Open settings
    ↓
Configure each aspect:
    - LLM: Claude Sonnet 3.5
    - Auto-approval: Whitelist src/, tests/
    - Display: Side-by-side diffs, monokai theme
    - Max iterations: 50 (for complex refactors)
    ↓
Click "Export Settings"
    ↓
Save to team-forge-config.json
    ↓
Share with team (repo, wiki, onboarding docs)
    ↓
Team members:
    - Install Forge
    - Type /settings
    - Click "Import Settings"
    - Select team-forge-config.json
    - Confirm import
    ↓
Entire team has consistent, optimal setup
```

---

## Product Requirements

### Priority 0 (Must Have)

#### P0-1: Visual Settings Interface (TUI Overlay)
**Description:** In-application settings UI accessible without leaving TUI or editing files

**User Stories:**
- As a new user, I want to configure Forge through a visual interface so I don't have to find and edit config files
- As a power user, I want keyboard-driven settings access so I never break my flow
- As any user, I want to see all available settings so I know what's configurable

**Acceptance Criteria:**
- Open settings overlay with `/settings` command or keyboard shortcut (Ctrl+,)
- Multi-tab interface organizing settings by category
- Navigate between tabs with Tab/Shift+Tab or number keys (1-5)
- Arrow keys navigate within tab
- Smooth open/close animation (not jarring)
- Settings overlay floats above main chat (modal)
- Close with Esc key or clicking outside
- Visual indicator of current tab
- No lag when switching tabs or scrolling

**UI Requirements:**
- Clean, uncluttered layout
- Clear visual hierarchy
- Consistent spacing and alignment
- Help text visible for selected setting
- Visual feedback for input focus
- Validation errors inline (not modal dialogs)

---

#### P0-2: Setting Categories and Organization
**Description:** Logical grouping of related settings into tabs

**Acceptance Criteria:**
- **General tab:** Workspace, iterations, notifications, editor
- **LLM tab:** Provider, model, API key, parameters
- **Auto-Approval tab:** Rules, patterns, tool-specific settings
- **Display tab:** Theme, colors, syntax highlighting, diffs
- **Advanced tab:** Debug, logging, performance tuning
- Tab labels clear and self-explanatory
- Tab order logical (most-used first)
- Each tab fits on screen without scrolling (or minimal scroll)

**Category Rationale:**
- **General:** First tab, most common settings
- **LLM:** Critical for functionality, second most accessed
- **Auto-Approval:** Workflow optimization, intermediate users
- **Display:** Visual preferences, subjective customization
- **Advanced:** Rarely changed, expert users only

---

#### P0-3: Input Controls and Validation
**Description:** Appropriate input types for different setting values with real-time validation

**User Stories:**
- As a user, I want clear input controls so I know what type of value is expected
- As a developer, I want validation errors immediately so I don't save invalid settings
- As a new user, I want help text explaining each setting

**Acceptance Criteria:**

**Text Input Fields:**
- API keys (masked input, show/hide toggle)
- Workspace path (with file browser button)
- Base URLs
- Custom values

**Dropdowns/Pickers:**
- LLM provider (OpenAI, Anthropic, Local, etc.)
- Model selection (provider-specific list)
- Theme selection
- Diff style (unified, side-by-side)
- Color scheme

**Number Inputs:**
- Max iterations (spinner or text with validation)
- Temperature (slider 0.0-2.0 with current value)
- Max tokens (text input with range validation)

**Boolean Toggles:**
- Enable/disable features
- Visual toggle switch (on/off states clear)
- Immediate feedback on change

**List Editors:**
- Path whitelist/blacklist (add/remove items)
- Command patterns
- Multiple API keys (future)

**Validation Features:**
- Real-time validation as user types
- Inline error messages (not modal popups)
- Visual indicators (red border, error icon)
- Helpful error messages:
  - ❌ "Invalid temperature" → ✅ "Temperature must be between 0.0 and 2.0"
  - ❌ "Bad path" → ✅ "Workspace directory does not exist: /invalid/path"
- Prevent saving invalid settings
- Show which fields are required

---

#### P0-4: Auto-Save and Persistence
**Description:** Automatic saving of settings changes with reliable persistence

**User Stories:**
- As a user, I want changes to save automatically so I don't lose my configuration
- As a developer, I want settings to persist across sessions
- As a team member, I want confidence that settings won't corrupt

**Acceptance Criteria:**
- Settings auto-save on change (no explicit save button)
- Save to `~/.config/forge/settings.json` (XDG standard on Unix)
- Windows: `%APPDATA%\forge\settings.json`
- macOS: `~/Library/Application Support/forge/settings.json`
- Atomic file writes (prevent corruption)
- Backup previous settings version (settings.json.backup)
- Load settings on application startup
- Merge with environment variable overrides
- Apply defaults for missing settings
- File permissions: 600 (user read/write only) for security
- Visual confirmation when settings saved (subtle indicator)
- Recovery from corrupted settings (fallback to defaults with notification)

---

#### P0-5: Help Text and Documentation
**Description:** Contextual help explaining each setting's purpose and values

**User Stories:**
- As a new user, I want to understand what each setting does before changing it
- As a power user, I want to see valid value ranges without trial-and-error
- As any user, I want examples of what to enter

**Acceptance Criteria:**
- Every setting has help text
- Help text shows when setting is focused/selected
- Help text includes:
  - Clear description of what setting controls
  - Valid values or range
  - Example values
  - Default value
  - Impact of changing (if significant)
- Help text position: below setting or in dedicated help pane
- Help text is concise (1-3 sentences)
- Technical jargon avoided or explained

**Example Help Text:**

**API Key:**
"Your LLM provider's API key for authentication. Get this from your provider's dashboard (OpenAI, Anthropic, etc.). Required for agent to function. Stored securely with restricted file permissions."

**Temperature:**
"Controls response creativity. Lower (0.0-0.5) = focused and deterministic. Higher (0.7-1.0) = creative and varied. Range: 0.0-2.0. Default: 0.7."

**Auto-approve Reads:**
"Automatically approve read-only operations (read_file, list_files, search_files) without prompting. Saves time for trusted operations. Default: enabled."

---

#### P0-6: Environment Variable Overrides
**Description:** Allow settings to be overridden via environment variables

**User Stories:**
- As a power user, I want to override settings per-session without changing saved config
- As a CI/CD pipeline, I want to inject settings via environment variables
- As a developer, I want to test different configurations quickly

**Acceptance Criteria:**
- Support standard environment variables:
  - `OPENAI_API_KEY` → LLM API key
  - `FORGE_WORKSPACE` → Workspace path
  - `FORGE_MODEL` → LLM model
  - `FORGE_PROVIDER` → LLM provider
  - `FORGE_MAX_ITERATIONS` → Max iterations
- Environment variables take precedence over file settings
- Settings UI shows when value is overridden by env var
- Visual indicator: "Overridden by OPENAI_API_KEY (env)" with different color
- Overridden values are read-only in UI (can't be edited)
- Documentation lists all supported environment variables

---

### Priority 1 (Should Have)

#### P1-1: Settings Search and Filtering
**Description:** Find specific settings quickly without navigating tabs

**User Stories:**
- As a user with many settings, I want to search by keyword
- As a new user, I want to find "API key" without knowing which tab

**Acceptance Criteria:**
- Search box at top of settings overlay
- Search by setting name, description, or category
- Real-time filtering as user types
- Highlight matching text
- Show which tab/category result is in
- Clear search button
- Keyboard shortcut to focus search (/)
- Search across all tabs simultaneously

---

#### P1-2: Import/Export Settings
**Description:** Share, backup, and migrate settings between machines or users

**User Stories:**
- As a team lead, I want to export my settings for the team
- As a user, I want to backup my configuration
- As a developer switching machines, I want to import my settings

**Acceptance Criteria:**
- Export button in settings (exports to JSON file)
- Import button (reads JSON file)
- Export includes all settings (or selected subset)
- Import validates settings before applying
- Import options: Replace all, or Merge (keep existing where conflicts)
- Export filename: `forge-settings-YYYY-MM-DD.json`
- Import shows preview before confirming
- Export sanitizes sensitive data option (remove API keys)

**Import Flow:**
```
User clicks "Import Settings"
    ↓
File picker appears
    ↓
User selects forge-settings.json
    ↓
Preview shows what will change:
    - LLM Provider: openai → anthropic
    - Model: gpt-4 → claude-sonnet-3.5
    - API Key: [will be overwritten]
    - Auto-approve reads: enabled → enabled (no change)
    ↓
User chooses: [Replace All] or [Merge]
    ↓
Confirmation: "Import 12 settings?"
    ↓
User confirms
    ↓
Settings applied and saved
    ↓
Success message: "Settings imported successfully"
```

---

#### P1-3: Reset to Defaults
**Description:** Quickly restore all settings to factory defaults

**User Stories:**
- As a user with misconfigured settings, I want easy recovery
- As a troubleshooter, I want to eliminate custom settings as issue source
- As a new user who experimented, I want a fresh start

**Acceptance Criteria:**
- "Reset to Defaults" button in settings
- Confirmation dialog before resetting (prevent accidents)
- Option to reset all or selected categories
- Backup current settings before reset
- Clear indication of what will be reset
- Can undo reset (restore from backup)

---

#### P1-4: Setting Profiles/Presets
**Description:** Save and switch between different configuration sets

**User Stories:**
- As a developer, I want different settings for different projects
- As a user, I want to test new models without losing current config
- As a team member, I want to switch between personal and team settings

**Acceptance Criteria:**
- Create named profiles (e.g., "Work", "Personal", "Testing")
- Switch between profiles with dropdown
- Each profile has complete settings set
- Active profile indicated clearly
- Export/import profiles
- Default profile for new sessions

---

#### P1-5: Advanced Validation Feedback
**Description:** Provide detailed, actionable validation errors

**User Stories:**
- As a user, I want to know exactly why a setting is invalid
- As a developer, I want suggested fixes for invalid values

**Acceptance Criteria:**
- Validation errors show immediately (on blur or after typing pause)
- Error messages are specific:
  - "Temperature must be between 0.0 and 2.0. You entered: 3.5"
  - "Workspace directory does not exist: /nonexistent/path. Create it or choose an existing directory."
  - "API key format invalid. Expected format: sk-... (OpenAI) or sk-ant-... (Anthropic)"
- Suggest corrections when possible:
  - "Did you mean: /home/user/projects?"
  - "Did you mean: claude-3-sonnet-20240229?"
- Link to documentation for complex settings
- Show valid examples for unclear requirements

---

### Priority 2 (Nice to Have)

#### P2-1: Settings Encryption
**Description:** Encrypt sensitive settings (API keys) at rest

**User Stories:**
- As a security-conscious user, I want API keys encrypted on disk
- As a team lead, I want to ensure credentials are protected

**Acceptance Criteria:**
- Optional encryption for sensitive fields
- Master password or OS keychain integration
- Transparent decryption on load
- Warning if encryption not enabled
- Migration path from unencrypted to encrypted

---

#### P2-2: Settings History and Versioning
**Description:** Track changes to settings over time

**User Stories:**
- As a user, I want to see what I changed and when
- As a troubleshooter, I want to revert to a previous configuration

**Acceptance Criteria:**
- Track setting changes with timestamps
- Show history of changes per setting
- Revert to previous value
- Diff view between current and historical
- Limited history (last 10 changes per setting)

---

#### P2-3: Cloud Settings Sync
**Description:** Synchronize settings across multiple machines

**User Stories:**
- As a multi-machine user, I want consistent settings everywhere
- As a team, I want to share and sync approved configurations

**Acceptance Criteria:**
- Optional cloud sync (disabled by default)
- Sync to user's cloud storage (Dropbox, Google Drive, custom)
- Conflict resolution when changes on multiple machines
- Selective sync (exclude sensitive data)
- Privacy-preserving (end-to-end encryption)

---

## User Experience Flows

### First-Time Setup Flow (Critical Path)

**Scenario:** New user installing Forge for the first time

```
User installs Forge CLI
    ↓
Runs: forge
    ↓
TUI launches, agent initializes
    ↓
Agent message: "Welcome to Forge! 👋
                 To get started, I need an API key from your LLM provider.
                 Type /settings to configure, or visit docs.forge.dev/setup"
    ↓
User types: /settings
    ↓
Settings overlay opens with smooth animation
┌─ Forge Settings ──────────────────────────────────────┐
│ [General] [LLM] [Auto-Approval] [Display] [Advanced]  │
│                                                        │
│ Welcome! Configure your LLM provider to start.        │
│                                                        │
│ Provider: [OpenAI ▾]  ← Dropdown with options         │
│                                                        │
│ Model: [gpt-4 ▾]                                      │
│                                                        │
│ API Key: [••••••••••••] [Show]  ← Masked input        │
│          Get your API key from platform.openai.com     │
│                                                        │
│ [Test Connection] ← Optional quick validation         │
│                                                        │
│ [← Back] Settings auto-save as you type              │
└────────────────────────────────────────────────────────┘
    ↓
User pastes API key
    ↓
Green checkmark appears: ✓ Valid API key format
    ↓
User clicks [Test Connection] (optional)
    ↓
"Testing connection..." → "✓ Connected to OpenAI successfully!"
    ↓
User presses Esc to close
    ↓
Settings save automatically
    ↓
Agent: "Great! I'm ready to help. What would you like to work on?"
    ↓
User starts being productive in <2 minutes from install
```

**Success Metrics:**
- Time to first productive use: <2 minutes
- API key configuration success rate: >95%
- Users who abandon during setup: <5%

---

### Power User Customization Flow

**Scenario:** Experienced user optimizing workflow with auto-approval rules

```
User working on project, annoyed by constant approval prompts
    ↓
User thinks: "I trust all reads in docs/ and tests/"
    ↓
User presses Ctrl+, (settings shortcut)
    ↓
Settings open to last tab (General)
    ↓
User presses 3 (Auto-Approval tab)
    ↓
┌─ Forge Settings ──────────────────────────────────────┐
│ General  LLM  [Auto-Approval]  Display  Advanced      │
│                                                        │
│ Auto-Approval Rules                                    │
│                                                        │
│ ☑ Enable auto-approval system                         │
│ ☑ Auto-approve read operations (read_file, list, etc) │
│                                                        │
│ Path Whitelist (auto-approve these paths):           │
│   • docs/*                        [Edit] [Remove]     │
│   • tests/*                       [Edit] [Remove]     │
│   • README.md                     [Edit] [Remove]     │
│   [+ Add Pattern]                                     │
│                                                        │
│ Path Blacklist (never auto-approve):                 │
│   • .env                          [Edit] [Remove]     │
│   • secrets/*                     [Edit] [Remove]     │
│   [+ Add Pattern]                                     │
│                                                        │
│ Help: Patterns support wildcards (* for any, ? for    │
│ single char). Whitelist takes precedence over         │
│ blacklist. Changes apply immediately.                  │
└────────────────────────────────────────────────────────┘
    ↓
User clicks [+ Add Pattern] under whitelist
    ↓
Input appears: [________________] [Save] [Cancel]
    ↓
User types: src/components/*
    ↓
User clicks Save
    ↓
Pattern added to list, settings auto-save
    ↓
Green toast: "✓ Auto-approval rule added"
    ↓
User presses Esc
    ↓
Next file read in src/components/ → silently auto-approved
    ↓
Workflow optimized, user happy
```

**Experience:** Fast, keyboard-driven, immediate feedback, no interruption to coding flow.

---

### Team Configuration Sharing Flow

**Scenario:** Team lead creating standard configuration for team

```
Team lead has ideal Forge setup
    ↓
Wants to share with 10 team members
    ↓
Opens settings (Ctrl+,)
    ↓
Verifies all settings are team-appropriate:
    - LLM: Claude Sonnet 3.5 (team standard)
    - Auto-approval: Whitelisted team repo paths
    - Display: Side-by-side diffs (team preference)
    - Max iterations: 50 (for complex tasks)
    ↓
Clicks "Export Settings" button (bottom right)
    ↓
Dialog appears:
┌─ Export Settings ─────────────────────────────────────┐
│                                                        │
│ Export to: [team-forge-config.json] [Browse]         │
│                                                        │
│ Options:                                               │
│ ☐ Exclude sensitive data (API keys)                   │
│ ☑ Include all settings                                │
│ ☐ Export selected categories only                     │
│                                                        │
│ This will export 23 settings across 4 categories.     │
│                                                        │
│ [Cancel] [Export]                                     │
└────────────────────────────────────────────────────────┘
    ↓
Team lead checks "Exclude sensitive data"
    ↓
Clicks Export
    ↓
File saved: ~/Downloads/team-forge-config.json
    ↓
Team lead adds to team wiki with instructions:
    "1. Install Forge
     2. Type /settings
     3. Click 'Import Settings'
     4. Select team-forge-config.json
     5. Add your personal API key in LLM tab"
    ↓
Team member follows instructions:
    Opens settings
    Clicks "Import Settings"
    Selects team-forge-config.json
    Preview shows changes
    Confirms import
    Adds personal API key
    ↓
Entire team has consistent setup in 1 minute per person
```

**Business Impact:** Reduced onboarding time from 30 minutes to 1 minute, consistent team behavior, fewer support questions.

---

## User Interface Design

### Settings Overlay - General Tab

```
┌─ Forge Settings ──────────────────────────────────────────────┐
│                                                                │
│ [General] LLM  Auto-Approval  Display  Advanced               │
│                                                                │
│ General Settings                                               │
│                                                                │
│ Workspace Directory                                            │
│ [/home/user/projects/myapp_______________] [Browse]           │
│ The working directory for agent operations. All file paths    │
│ are relative to this directory.                                │
│                                                                │
│ Max Agent Iterations                                           │
│ [25_] ← → (Range: 5-100)                                      │
│ Maximum number of agent loop iterations before automatic      │
│ timeout. Higher values allow more complex tasks. Default: 25  │
│                                                                │
│ Notifications                                                  │
│ ☑ Enable toast notifications                                  │
│ Show brief popup notifications for important events           │
│                                                                │
│ Default Editor                                                 │
│ [vim▾] (Dropdown: vim, emacs, nano, code, etc.)              │
│ Text editor to use for file editing operations                │
│                                                                │
│ Session Timeout                                                │
│ [30_] minutes (0 = no timeout)                                │
│ Automatically end session after period of inactivity           │
│                                                                │
│                                          [Reset] [Export]      │
└────────────────────────────────────────────────────────────────┘

[Esc] Close  [Tab] Next Tab  [/] Search  [Ctrl+R] Reset to Defaults
```

---

### Settings Overlay - LLM Tab

```
┌─ Forge Settings ──────────────────────────────────────────────┐
│                                                                │
│ General  [LLM]  Auto-Approval  Display  Advanced              │
│                                                                │
│ LLM Configuration                                              │
│                                                                │
│ Provider                                                       │
│ [Anthropic ▾] (OpenAI, Anthropic, Local, Azure, etc.)        │
│ Choose your language model provider                            │
│                                                                │
│ Model                                                          │
│ [claude-3-sonnet-20240229 ▾]                                  │
│ Available models:                                              │
│   • claude-3-opus-20240229 (Most capable)                     │
│   • claude-3-sonnet-20240229 (Balanced) ← Selected            │
│   • claude-3-haiku-20240307 (Fast)                            │
│                                                                │
│ API Key                                                        │
│ [sk-ant-api03-••••••••••••••••••••••••••] [Show] [Test]      │
│ Your Anthropic API key. Get it from console.anthropic.com     │
│ ✓ Valid key format                                            │
│                                                                │
│ Advanced Parameters (Optional)                                 │
│ [▸ Show Advanced Settings]                                    │
│                                                                │
│ Temperature: [0.7___] (0.0 = deterministic, 2.0 = creative)   │
│ Max Tokens:  [4096_] (Context window size)                     │
│ Top-p:       [1.0___] (Nucleus sampling, 0.0-1.0)             │
│                                                                │
│                                          [Reset] [Export]      │
└────────────────────────────────────────────────────────────────┘

[Esc] Close  [Tab] Next Tab  [Ctrl+T] Test Connection
```

---

### Settings Overlay - Auto-Approval Tab

```
┌─ Forge Settings ──────────────────────────────────────────────┐
│                                                                │
│ General  LLM  [Auto-Approval]  Display  Advanced              │
│                                                                │
│ Auto-Approval Rules                                            │
│                                                                │
│ ☑ Enable auto-approval system                                 │
│ Automatically approve safe operations based on rules below     │
│                                                                │
│ Quick Settings                                                 │
│ ☑ Auto-approve read operations (read_file, list_files, search)│
│ ☑ Auto-approve searches (search_files)                         │
│ ☐ Auto-approve list operations (list_files, list_dirs)        │
│ ☐ Auto-approve non-destructive writes (create new files only) │
│                                                                │
│ Path Whitelist (Always auto-approve these paths)              │
│ ┌──────────────────────────────────────────────────────────┐ │
│ │ • docs/*                              [Edit] [Remove]    │ │
│ │ • tests/*                             [Edit] [Remove]    │ │
│ │ • README.md                           [Edit] [Remove]    │ │
│ │ • src/components/*                    [Edit] [Remove]    │ │
│ │                                                          │ │
│ │ [+ Add Pattern]                                          │ │
│ └──────────────────────────────────────────────────────────┘ │
│                                                                │
│ Path Blacklist (Never auto-approve)                           │
│ ┌──────────────────────────────────────────────────────────┐ │
│ │ • .env                                [Edit] [Remove]    │ │
│ │ • secrets/*                           [Edit] [Remove]    │ │
│ │ • production.config                   [Edit] [Remove]    │ │
│ │                                                          │ │
│ │ [+ Add Pattern]                                          │ │
│ └──────────────────────────────────────────────────────────┘ │
│                                                                │
│ Help: Patterns use glob syntax (* = any, ? = single char).    │
│ Whitelist overrides blacklist. Changes apply immediately.      │
│                                          [Reset] [Export]      │
└────────────────────────────────────────────────────────────────┘
```

---

### Settings Overlay - Display Tab

```
┌─ Forge Settings ──────────────────────────────────────────────┐
│                                                                │
│ General  LLM  Auto-Approval  [Display]  Advanced              │
│                                                                │
│ Display Preferences                                            │
│                                                                │
│ Color Theme                                                    │
│ [Dark ▾] (Dark, Light, High Contrast, Custom)                │
│ Overall color scheme for the TUI                               │
│                                                                │
│ Syntax Highlighting                                            │
│ ☑ Enable syntax highlighting                                  │
│ Color Scheme: [Monokai ▾]                                     │
│              (Monokai, Solarized, Dracula, GitHub, etc.)      │
│ Preview:                                                       │
│ ┌────────────────────────────────────────┐                    │
│ │ func example() {          ← Go syntax  │                    │
│ │     return "highlighted"               │                    │
│ │ }                                       │                    │
│ └────────────────────────────────────────┘                    │
│                                                                │
│ Diff Display                                                   │
│ Style: [Side-by-side ▾] (Unified, Side-by-side)              │
│ ☑ Show line numbers                                           │
│ ☑ Highlight changed sections                                  │
│                                                                │
│ Layout                                                         │
│ Density: [Comfortable ▾] (Compact, Comfortable, Spacious)     │
│ ☑ Show status bar                                             │
│ ☑ Show file paths in results                                  │
│ ☑ Show timestamps                                             │
│                                                                │
│                                          [Reset] [Export]      │
└────────────────────────────────────────────────────────────────┘
```

---

### Settings Import Dialog

```
┌─ Import Settings ─────────────────────────────────────────────┐
│                                                                │
│ Select settings file to import:                               │
│                                                                │
│ File: [team-forge-config.json_______________] [Browse]        │
│                                                                │
│ Preview Changes:                                               │
│ ┌──────────────────────────────────────────────────────────┐ │
│ │ The following settings will be changed:                  │ │
│ │                                                          │ │
│ │ LLM Provider:  openai → anthropic                        │ │
│ │ Model:         gpt-4 → claude-3-sonnet-20240229          │ │
│ │ Auto-approve:  ☐ → ☑ (enabled)                          │ │
│ │ Diff Style:    unified → side-by-side                    │ │
│ │                                                          │ │
│ │ Unchanged settings: 18                                    │ │
│ │ Missing in import (will keep current): API Key           │ │
│ └──────────────────────────────────────────────────────────┘ │
│                                                                │
│ Import Mode:                                                   │
│ ○ Replace all settings (current settings will be lost)        │
│ ● Merge (keep current settings where not in import)           │
│                                                                │
│ ☑ Create backup of current settings before import             │
│                                                                │
│                               [Cancel] [Import Settings]       │
└────────────────────────────────────────────────────────────────┘
```

---

## Success Metrics

### Adoption & Discovery

**Primary Metrics:**
- **Settings Access Rate:** >70% of users open settings at least once
- **API Key Setup Success:** >95% successfully configure API key on first try
- **Time to First Configuration:** p95 <2 minutes from install to productive use
- **Settings Customization Rate:** >50% of users modify at least one setting
- **Discovery Without Docs:** >60% find desired setting without external documentation

**Engagement Metrics:**
- **Settings Sessions per User:** Average 3-5 settings changes per user lifetime
- **Tab Usage Distribution:**
  - General: 100% (all users)
  - LLM: 95% (nearly universal)
  - Auto-Approval: 40% (workflow optimization)
  - Display: 30% (personal preference)
  - Advanced: 10% (power users)
- **Feature Utilization:**
  - Export/Import: 15% of users
  - Reset to Defaults: 5% of users
  - Search: 25% of users (when >50 settings)

---

### User Satisfaction

**Quality Metrics:**
- **Configuration Satisfaction:** >4.5/5 rating for settings experience
- **Ease of Use:** >85% report "easy" or "very easy" to configure
- **Help Text Effectiveness:** >70% understand settings without additional help
- **Validation Clarity:** >80% understand error messages and how to fix
- **Visual Design:** >4.2/5 rating for settings interface aesthetics

**Preference Metrics:**
- **vs. Manual File Editing:** >90% prefer visual settings over config files
- **vs. CLI Flags:** >85% prefer persistent settings over repeated flags
- **vs. Environment Variables:** >75% prefer UI configuration (except CI/CD)

---

### Performance & Reliability

**Speed Metrics:**
- **Overlay Open Time:** p95 <100ms
- **Tab Switch Time:** p95 <50ms
- **Setting Change Apply:** p95 <50ms
- **File Save Time:** p95 <200ms
- **Settings Load Time:** p95 <50ms at startup

**Reliability Metrics:**
- **Corruption Rate:** 0% settings file corruption
- **Validation Accuracy:** 100% of invalid settings caught before save
- **Migration Success:** 100% of users upgrade without settings loss
- **Recovery Rate:** 100% recovery from corrupted settings (fallback to defaults)
- **Platform Compatibility:** Works on Linux, macOS, Windows without issues

---

### Business Impact

**Onboarding Metrics:**
- **Setup Completion Rate:** >95% complete initial setup (up from 65% without visual settings)
- **Time to Productivity:** Reduced from 15 minutes to <2 minutes
- **Setup-Related Churn:** <2% abandon during configuration (down from 35%)
- **Support Tickets:** 70% reduction in "how to configure" questions

**Productivity Metrics:**
- **Configuration Time Saved:** 13 minutes saved per user (no manual file editing)
- **Error Resolution:** 85% fewer configuration-related errors
- **Workflow Optimization:** 40% of users configure auto-approval → fewer interruptions
- **Team Consistency:** 90% consistency in team configurations (vs. 30% with manual setup)

**Retention Metrics:**
- **Feature Stickiness:** 92% of users who customize settings continue using Forge
- **Recommendation NPS:** +18 point improvement from easy configuration
- **Team Adoption:** 3x higher team adoption rate (easy to standardize)

---

## Competitive Analysis

### VSCode Settings
**Approach:** Hierarchical JSON with GUI overlay  
**Strengths:** Comprehensive, searchable, sync across devices  
**Weaknesses:** Overwhelming (1000+ settings), hard to find what you need  
**Differentiation:** We provide focused, task-oriented settings without noise

### Cursor IDE
**Approach:** Settings modal with tabs  
**Strengths:** Clean UI, good defaults  
**Weaknesses:** Limited customization, some settings require config file  
**Differentiation:** More powerful auto-approval rules, better validation

### GitHub Copilot
**Approach:** Minimal settings (mostly via IDE)  
**Strengths:** Simple, low friction  
**Weaknesses:** Almost no configurability, frustrating for power users  
**Differentiation:** Balance of simplicity and control

### Aider Terminal
**Approach:** CLI flags and config file  
**Strengths:** Powerful for CLI experts  
**Weaknesses:** Poor discoverability, no validation, manual file editing  
**Differentiation:** Visual interface accessible to all skill levels

### JetBrains IDEs
**Approach:** Comprehensive settings tree with search  
**Strengths:** Every detail configurable  
**Weaknesses:** Overwhelming complexity, slow search  
**Differentiation:** Simpler, faster, focused on essential settings

---

## Go-to-Market Considerations

### Positioning

**Primary Message:**  
"Configure Forge in seconds through a clean visual interface. No config files, no documentation hunting, no syntax errors—just intuitive settings that make the tool yours."

**Key Differentiators:**
- Visual settings interface accessible without leaving TUI
- Real-time validation prevents configuration errors
- Auto-save eliminates "forgot to save" frustration
- Import/Export enables team standardization
- Comprehensive yet organized (not overwhelming)

---

### Target Segments

**Early Adopters:**
- Developers frustrated with config file complexity
- Teams wanting standardized AI tool setup
- Users who value clean, polished interfaces

**Value Propositions by Segment:**
- **New Users:** "Get started in under 2 minutes with guided setup"
- **Power Users:** "Complete control with keyboard-driven efficiency"
- **Team Leads:** "Standardize team configuration in one click"
- **Enterprise:** "Secure, reliable, auditable settings management"

---

### Documentation Needs

**Essential Documentation:**
1. **Quick Start: Settings** - Configure API key in 60 seconds
2. **Settings Reference** - Complete list of all settings with examples
3. **Auto-Approval Guide** - Configure rules safely and effectively
4. **Team Setup Guide** - Export/import workflows for teams
5. **Troubleshooting Settings** - Common issues and fixes
6. **Advanced Configuration** - Environment variables, profiles, encryption

**FAQ Topics:**
- "How do I set my API key?"
- "Where are settings stored?"
- "Can I share settings with my team?"
- "How do I reset to defaults?"
- "What if I made a mistake?"
- "Are my API keys secure?"
- "Can I use different settings per project?"

---

## Risk & Mitigation

### Risk 1: API Key Security Concerns
**Impact:** High - Users worry about storing credentials  
**Probability:** High - Common security concern  
**User Impact:** Hesitation to use product, manual environment variables instead

**Mitigation:**
- File permissions automatically set to 600 (user-only)
- Clear messaging: "API keys stored securely with restricted permissions"
- Support environment variable override for security-conscious users
- Documentation on security best practices
- Future: Optional encryption with master password
- Never log API keys, mask in UI

**User Communication:**
"Your API key is stored in ~/.config/forge/settings.json with permissions set to 600 (readable only by you). For extra security, you can use the OPENAI_API_KEY environment variable instead."

---

### Risk 2: Settings Overwhelming New Users
**Impact:** Medium - Complexity drives abandonment  
**Probability:** Medium - Many settings possible  
**User Impact:** Confusion, decision paralysis, abandonment

**Mitigation:**
- Excellent defaults that work for 80% of users
- Progressive disclosure (basic → advanced)
- Wizard/guide for first-time setup
- Hide advanced settings behind "Show Advanced" toggle
- Clear, jargon-free help text
- Search functionality to find specific settings
- Presets for common configurations

**Onboarding Flow:**
```
First launch → Minimal setup wizard
    ↓
"Let's get you started! 🚀"
Step 1: Choose LLM provider [dropdown]
Step 2: Enter API key [input]
Step 3: Test connection [auto-test]
    ↓
"All set! You can customize more in /settings anytime"
```

---

### Risk 3: Settings File Corruption
**Impact:** High - Lost configuration, broken app  
**Probability:** Low - With proper implementation  
**User Impact:** Frustration, lost time, potential data loss

**Mitigation:**
- Atomic file writes (write to temp, then rename)
- Automatic backup before every write (settings.json.backup)
- Validation before writing
- Graceful degradation (fallback to defaults if corrupted)
- Clear error message with recovery instructions
- Automatic recovery attempt on corrupt detection

**Recovery Flow:**
```
Settings load fails → Corruption detected
    ↓
Show message:
"Settings file appears corrupted.
 Don't worry! Restoring from automatic backup...
 
 [Restore from Backup] [Use Defaults] [View Error Details]"
    ↓
User chooses Restore
    ↓
Backup restored successfully
    ↓
"✓ Settings restored from backup (5 minutes ago)"
```

---

### Risk 4: Platform-Specific File Paths
**Impact:** Medium - Inconsistent behavior across OS  
**Probability:** Medium - Different conventions  
**User Impact:** Confusion, support burden

**Mitigation:**
- Use standard config directories per platform:
  - Linux/Unix: ~/.config/forge/
  - macOS: ~/Library/Application Support/forge/
  - Windows: %APPDATA%\forge\
- Respect XDG_CONFIG_HOME environment variable
- Gracefully handle missing directories (create automatically)
- Clear documentation for each platform
- Consistent behavior despite different paths

---

### Risk 5: Breaking Changes in Settings Schema
**Impact:** Medium - User frustration during upgrades  
**Probability:** High - Features evolve  
**User Impact:** Lost settings, unexpected behavior

**Mitigation:**
- Version settings schema (v1, v2, etc.)
- Automatic migration on upgrade
- Preserve unknown settings (forward compatibility)
- Clear upgrade notes documenting changes
- Test migrations thoroughly before release
- Rollback capability if migration fails
- Backup before migration

**Migration Example:**
```
Forge v1.0 → v1.5 (settings schema v1 → v2)

Changes:
- "provider" renamed to "llm_provider"
- "auto_approve" split into granular settings
- New "display.theme" setting

Migration:
1. Detect schema v1
2. Backup current settings
3. Transform:
   - provider → llm_provider
   - auto_approve: true → auto_approve_reads: true, etc.
   - Add defaults for new settings
4. Save as schema v2
5. Success message with changelog
```

---

## Evolution & Roadmap

### Version History

**v1.0 (Current):**
- Visual settings overlay with multi-tab interface
- Comprehensive setting categories (General, LLM, Auto-Approval, Display, Advanced)
- Real-time validation and help text
- Auto-save with atomic writes
- Environment variable overrides
- Import/Export functionality
- Reset to defaults

---

### Future Enhancements

#### Phase 2: Enhanced Security & Collaboration
- **Settings Encryption:** Encrypt API keys and sensitive data at rest
- **Settings Profiles:** Multiple named configuration sets (Work, Personal, Testing)
- **Team Templates:** Official preset configurations for common use cases
- **Audit Log:** Track who changed what and when (for teams)
- **Advanced Import:** Selective import (choose which settings to import)
- **Settings Diff:** Compare current vs. imported settings before applying

**User Value:** Better security, easier team collaboration, more flexibility

---

#### Phase 3: Intelligence & Automation
- **Smart Defaults:** Adapt default settings based on detected project type
- **Configuration Recommendations:** "You might want to enable auto-approve for docs/"
- **Settings Analytics:** "80% of users enable auto-approve reads"
- **Conflict Detection:** Warn about conflicting settings
- **Configuration Validation:** "This combination of settings may cause issues"
- **Auto-Optimization:** Suggest settings changes based on usage patterns

**User Value:** Smarter configuration, proactive guidance, optimized workflows

---

#### Phase 4: Enterprise & Advanced Features
- **Cloud Sync:** Synchronize settings across machines
- **Team Settings Management:** Admin controls, approved configurations
- **Policy Enforcement:** Require certain settings, block others
- **Settings API:** Programmatic configuration for CI/CD
- **Custom Settings UI:** Plugin system for third-party settings
- **Settings History:** Version control for configurations with rollback
- **Compliance Mode:** Locked settings for regulatory requirements

**User Value:** Enterprise-ready, team-scale management, compliance support

---

## Related Documentation

- **User Guide:** How to configure Forge settings
- **Security:** Best practices for API key storage
- **Team Guide:** Sharing and standardizing configurations
- **API Reference:** Environment variable reference
- **Troubleshooting:** Common settings issues and solutions

---

## Changelog

### 2024-12-XX
- Transformed to product-focused PRD format
- Removed technical implementation details (component structure, data models, Go structs)
- Enhanced user personas with detailed scenarios and success stories
- Added comprehensive UI mockups for all settings tabs
- Expanded user experience flows with visual diagrams
- Added competitive analysis comparing to VSCode, Cursor, Copilot
- Included go-to-market positioning and messaging
- Improved success metrics with user-focused KPIs
- Added detailed risk mitigation with user communication strategies

### 2024-12 (Original)
- Initial PRD with technical architecture
- Component structure and data models
- Settings lifecycle diagrams
