# Product Requirements: Settings System

**Feature:** Settings and Configuration Management  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Core Team  
**Last Updated:** December 2024

---

## Overview

The Settings System provides a comprehensive, user-friendly way to configure Forge's behavior, preferences, and integrations. It includes both a visual settings interface (TUI overlay) and persistent storage, allowing users to customize their experience without editing configuration files manually.

---

## Problem Statement

Users need to configure various aspects of Forge but face challenges:

1. **Complexity:** Many configuration options across different domains (LLM, UI, security)
2. **Discovery:** Users don't know what settings are available
3. **Validation:** Manual file editing leads to syntax errors and invalid values
4. **Persistence:** Settings need to survive across sessions
5. **Context:** Users need help understanding what each setting does
6. **Accessibility:** Editing config files requires leaving the TUI

Traditional configuration approaches (CLI flags, env vars, config files) are error-prone and have poor discoverability.

---

## Goals

### Primary Goals

1. **User-Friendly Interface:** Provide visual settings UI within TUI
2. **Comprehensive Coverage:** Support all configurable aspects of Forge
3. **Safe Defaults:** Ship with sensible defaults that work for most users
4. **Easy Customization:** Make common adjustments quick and intuitive
5. **Persistence:** Automatically save and restore settings
6. **Validation:** Prevent invalid configurations

### Non-Goals

1. **Programmatic API:** This is NOT for runtime configuration changes via code
2. **Multi-User:** Does NOT support team/shared settings (per-user only)
3. **Remote Config:** Does NOT sync settings across machines
4. **Version Control:** Settings files are NOT meant to be committed to git

---

## User Personas

### Primary: Customization-Focused Developer
- **Background:** Experienced developer who optimizes tools for workflow
- **Workflow:** Tweaks settings frequently to match preferences
- **Pain Points:** Wants control over every aspect without complexity
- **Goals:** Perfect configuration for personal workflow

### Secondary: Team Lead
- **Background:** Sets standards for team's tool usage
- **Workflow:** Configures tools once, uses consistently
- **Pain Points:** Needs reliable, documented settings
- **Goals:** Predictable behavior, easy to replicate setup

### Tertiary: New User
- **Background:** First-time Forge user
- **Workflow:** Exploring tool capabilities
- **Pain Points:** Overwhelmed by options, needs guidance
- **Goals:** Get started quickly with good defaults

---

## Requirements

### Functional Requirements

#### FR1: Settings Categories
- **R1.1:** General settings (workspace, iterations, notifications)
- **R1.2:** LLM settings (provider, model, API keys, parameters)
- **R1.3:** Auto-approval settings (rules, patterns, tool configuration)
- **R1.4:** Display settings (theme, colors, syntax highlighting, diff style)
- **R1.5:** Advanced settings (debug, logging, performance)

#### FR2: Settings UI (TUI Overlay)
- **R2.1:** Multi-tab interface for different categories
- **R2.2:** Navigate tabs with Tab/Shift+Tab or numbers
- **R2.3:** Arrow keys to navigate within tab
- **R2.4:** Edit values inline
- **R2.5:** Dropdown/picker for enumerated values
- **R2.6:** Text input for string/numeric values
- **R2.7:** Toggle for boolean values
- **R2.8:** Help text for each setting
- **R2.9:** Validation feedback on change
- **R2.10:** Save/Cancel/Reset options

#### FR3: General Settings
- **R3.1:** Workspace path (directory for agent operations)
- **R3.2:** Max agent iterations (safety limit)
- **R3.3:** Enable/disable toast notifications
- **R3.4:** Default editor (for file editing)
- **R3.5:** Language preference (future: i18n)
- **R3.6:** Session timeout (optional)

#### FR4: LLM Settings
- **R4.1:** Provider selection (OpenAI, Anthropic, Local, etc.)
- **R4.2:** Model selection (provider-specific)
- **R4.3:** API key management (per provider)
- **R4.4:** API base URL (for custom endpoints)
- **R4.5:** Temperature (creativity parameter)
- **R4.6:** Max tokens (context window limit)
- **R4.7:** Top-p (nucleus sampling)
- **R4.8:** Streaming enabled/disabled

#### FR5: Auto-Approval Settings
- **R5.1:** Enable/disable auto-approval globally
- **R5.2:** Auto-approve read operations (default: true)
- **R5.3:** Path-based whitelist rules
- **R5.4:** Path-based blacklist rules
- **R5.5:** Command pattern whitelist
- **R5.6:** Tool-specific approval settings
- **R5.7:** Rule priority/ordering
- **R5.8:** Rule enable/disable toggles
- **R5.9:** Add/edit/delete rules

#### FR6: Display Settings
- **R6.1:** Color theme (dark/light/custom)
- **R6.2:** Syntax highlighting enabled/disabled
- **R6.3:** Highlighting color scheme
- **R6.4:** Diff display style (unified vs side-by-side)
- **R6.5:** Show line numbers (in diffs, code blocks)
- **R6.6:** Font preferences (if terminal supports)
- **R6.7:** Compact vs spacious layout
- **R6.8:** Show/hide status bar elements

#### FR7: Persistence
- **R7.1:** Save settings to `~/.config/forge/settings.json`
- **R7.2:** Auto-save on change (no explicit save button)
- **R7.3:** Load settings on startup
- **R7.4:** Atomic writes (prevent corruption)
- **R7.5:** Backup previous settings version
- **R7.6:** Migration for settings schema changes

#### FR8: Validation
- **R8.1:** Validate settings before applying
- **R8.2:** Show error messages for invalid values
- **R8.3:** Prevent saving invalid configurations
- **R8.4:** Type checking (string, number, boolean)
- **R8.5:** Range validation (min/max for numbers)
- **R8.6:** Path validation (workspace exists)
- **R8.7:** API key format validation

#### FR9: Environment Variables
- **R9.1:** Override settings via environment variables
- **R9.2:** `FORGE_WORKSPACE` for workspace path
- **R9.3:** `OPENAI_API_KEY` for API key
- **R9.4:** `FORGE_MODEL` for model selection
- **R9.5:** Env vars take precedence over file settings
- **R9.6:** Show which settings are overridden in UI

#### FR10: Import/Export
- **R10.1:** Export settings to file (for backup/sharing)
- **R10.2:** Import settings from file
- **R10.3:** Export as JSON or YAML
- **R10.4:** Validate imported settings
- **R10.5:** Merge vs replace import options

### Non-Functional Requirements

#### NFR1: Performance
- **N1.1:** Settings overlay opens within 100ms
- **N1.2:** Setting changes apply within 50ms
- **N1.3:** File save completes within 200ms
- **N1.4:** No lag when navigating settings
- **N1.5:** Efficient memory usage (<5MB for settings)

#### NFR2: Reliability
- **N2.1:** Never corrupt settings file
- **N2.2:** Graceful handling of malformed settings
- **N2.3:** Fallback to defaults if settings invalid
- **N2.4:** Transaction-safe file writes
- **N2.5:** Automatic recovery from errors

#### NFR3: Usability
- **N3.1:** Intuitive tab navigation
- **N3.2:** Clear labels and descriptions
- **N3.3:** Immediate visual feedback on changes
- **N3.4:** Consistent interaction patterns
- **N3.5:** Help text always visible

#### NFR4: Security
- **N4.1:** API keys stored securely (not plain text)
- **N4.2:** Settings file has restrictive permissions (600)
- **N4.3:** No sensitive data in logs
- **N4.4:** Validate all external input
- **N4.5:** Sanitize paths to prevent traversal

#### NFR5: Compatibility
- **N5.1:** Settings format backwards compatible
- **N5.2:** Migration path for schema changes
- **N5.3:** Work across different OS (Linux, macOS, Windows)
- **N5.4:** Handle different config directories gracefully

---

## User Experience

### Core Workflows

#### Workflow 1: First-Time Setup
1. User launches Forge for first time
2. No settings file exists
3. Defaults are used
4. User opens `/settings`
5. Sees all default values
6. Changes API key
7. Settings auto-save
8. User closes overlay
9. Agent uses new API key

**Success Criteria:** User can configure API key in under 1 minute

#### Workflow 2: Changing Display Preferences
1. User wants side-by-side diffs
2. Opens settings with `/settings`
3. Navigates to "Display" tab
4. Finds "Diff Style" setting
5. Changes from "unified" to "side-by-side"
6. Setting saves automatically
7. Next diff shows side-by-side
8. User confirms change works

**Success Criteria:** Visual preferences apply immediately

#### Workflow 3: Configuring Auto-Approval Rules
1. User tired of approving read operations
2. Opens settings
3. Goes to "Auto-Approval" tab
4. Sees "Auto-approve reads" setting
5. Toggles to enabled
6. Adds path whitelist rule: `docs/*`
7. Rules save automatically
8. Future reads auto-approve

**Success Criteria:** Rules configured without editing files

#### Workflow 4: Switching LLM Providers
1. User wants to try different model
2. Opens settings
3. Goes to "LLM" tab
4. Changes provider from OpenAI to Anthropic
5. Enters Anthropic API key
6. Selects Claude model
7. Settings validate and save
8. Next agent call uses Claude

**Success Criteria:** Provider switch works seamlessly

#### Workflow 5: Resetting to Defaults
1. User's settings are misconfigured
2. Opens settings
3. Clicks "Reset to Defaults" button
4. Confirmation dialog appears
5. User confirms
6. All settings revert to defaults
7. Settings save
8. User reconfigures as needed

**Success Criteria:** Easy recovery from bad configuration

---

## Technical Architecture

### Component Structure

```
Settings System
├── Settings Manager
│   ├── Configuration Store
│   ├── Validation Engine
│   ├── Migration Manager
│   └── Environment Handler
├── Settings Overlay (TUI)
│   ├── Tab Container
│   ├── General Tab
│   ├── LLM Tab
│   ├── Auto-Approval Tab
│   ├── Display Tab
│   └── Input Handlers
├── Persistence Layer
│   ├── File Writer
│   ├── File Reader
│   ├── JSON Parser
│   └── Backup Manager
└── Settings Schema
    ├── Type Definitions
    ├── Validators
    ├── Defaults
    └── Migrations
```

### Data Model

```go
type Settings struct {
    General     GeneralSettings     `json:"general"`
    LLM         LLMSettings         `json:"llm"`
    AutoApproval AutoApprovalSettings `json:"auto_approval"`
    Display     DisplaySettings     `json:"display"`
    Advanced    AdvancedSettings    `json:"advanced"`
    Version     string              `json:"version"`
}

type GeneralSettings struct {
    Workspace        string `json:"workspace"`
    MaxIterations    int    `json:"max_iterations"`
    ToastEnabled     bool   `json:"toast_enabled"`
    DefaultEditor    string `json:"default_editor"`
}

type LLMSettings struct {
    Provider    string  `json:"provider"`
    Model       string  `json:"model"`
    APIKey      string  `json:"api_key"`
    BaseURL     string  `json:"base_url"`
    Temperature float64 `json:"temperature"`
    MaxTokens   int     `json:"max_tokens"`
}

type AutoApprovalSettings struct {
    Enabled         bool              `json:"enabled"`
    AutoApproveReads bool             `json:"auto_approve_reads"`
    PathWhitelist   []string          `json:"path_whitelist"`
    PathBlacklist   []string          `json:"path_blacklist"`
    CommandPatterns []string          `json:"command_patterns"`
    ToolRules       map[string]bool   `json:"tool_rules"`
}

type DisplaySettings struct {
    Theme              string `json:"theme"`
    SyntaxHighlighting bool   `json:"syntax_highlighting"`
    ColorScheme        string `json:"color_scheme"`
    DiffStyle          string `json:"diff_style"`
    ShowLineNumbers    bool   `json:"show_line_numbers"`
}
```

### Settings Lifecycle

```
Startup
  ↓
Load from ~/.config/forge/settings.json
  ↓
Merge with environment variables
  ↓
Validate & apply defaults for missing values
  ↓
Settings available to application
  ↓
User opens settings overlay
  ↓
User modifies value
  ↓
Validate change
  ↓
Apply to runtime
  ↓
Save to disk (atomic write)
  ↓
Settings active for session
```

---

## Design Decisions

### Why Auto-Save Instead of Explicit Save?
- **User expectation:** Modern apps auto-save (Google Docs, Notion)
- **Safety:** No risk of losing changes
- **Simplicity:** Fewer buttons, less cognitive load
- **Immediate feedback:** Changes take effect right away
- **Atomic writes:** Safe to save on every change

### Why Multi-Tab Interface?
- **Organization:** Groups related settings logically
- **Discoverability:** Easier to find specific settings
- **Screen space:** Avoids overwhelming single-page list
- **Progressive disclosure:** Advanced settings in separate tab
- **Familiar pattern:** Standard in desktop/web applications

### Why JSON Instead of YAML/TOML?
- **Ubiquity:** JSON parsers everywhere
- **Go native:** Standard library support
- **Simplicity:** No indentation sensitivity
- **Validation:** Easy to validate with JSON Schema
- **Performance:** Fast parsing
- **Web-friendly:** Easy to export/import via web interfaces (future)

### Why Store API Keys in Settings File?
**Alternatives considered:**
1. **System keychain:** OS-specific, complex integration
2. **Environment variables only:** Harder for users to manage
3. **Separate secrets file:** More files to manage

**Why settings file won:**
- Simple, works everywhere
- File permissions provide security (600)
- Environment variable override available
- Future: Encryption layer if needed

---

## Settings Reference

### General Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| workspace | string | `$PWD` | Working directory for agent operations |
| max_iterations | int | 25 | Maximum agent loop iterations |
| toast_enabled | bool | true | Show toast notifications |
| default_editor | string | `$EDITOR` | Text editor for file editing |

### LLM Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| provider | string | "openai" | LLM provider (openai, anthropic, local) |
| model | string | "gpt-4" | Model name |
| api_key | string | "" | API key for provider |
| base_url | string | "" | Custom API endpoint |
| temperature | float | 0.7 | Creativity (0.0-2.0) |
| max_tokens | int | 4096 | Context window size |

### Auto-Approval Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| enabled | bool | true | Enable auto-approval system |
| auto_approve_reads | bool | true | Auto-approve read operations |
| path_whitelist | []string | [] | Paths to auto-approve |
| path_blacklist | []string | [] | Paths to never approve |
| command_patterns | []string | [] | Shell command patterns to approve |

### Display Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| theme | string | "dark" | Color theme |
| syntax_highlighting | bool | true | Enable syntax highlighting |
| color_scheme | string | "monokai" | Chroma color scheme |
| diff_style | string | "unified" | Diff display (unified/side-by-side) |
| show_line_numbers | bool | true | Show line numbers in code |

---

## Success Metrics

### Adoption Metrics
- **Settings usage:** >70% of users access settings at least once
- **Customization rate:** >50% of users modify at least one setting
- **API key setup:** >90% successfully configure API keys
- **Auto-approval:** >40% enable at least one auto-approval rule

### Usability Metrics
- **Time to configure:** p95 under 2 minutes for common changes
- **Error rate:** <5% of setting changes result in validation errors
- **Help usage:** >30% read help text before changing settings
- **Discovery:** >60% find desired setting without external help

### Quality Metrics
- **Corruption rate:** 0% settings file corruption
- **Migration success:** 100% of schema upgrades succeed
- **Default satisfaction:** >80% of users keep most default settings
- **Recovery rate:** 100% recovery from invalid settings

---

## Dependencies

### External Dependencies
- File system access (read/write to ~/.config/)
- JSON standard library
- Environment variable access

### Internal Dependencies
- TUI framework (for overlay rendering)
- Agent core (applies settings to behavior)
- Tool approval system (uses auto-approval rules)
- LLM providers (use API keys and parameters)

### Platform Requirements
- Writable home directory
- Standard config directory (~/.config/ on Unix)
- File permissions support (chmod 600)

---

## Risks & Mitigations

### Risk 1: Settings File Corruption
**Impact:** High  
**Probability:** Low  
**Mitigation:**
- Atomic file writes
- Backup previous version
- Validate before writing
- Fallback to defaults on corruption
- Automatic recovery

### Risk 2: API Key Security
**Impact:** High  
**Probability:** Medium  
**Mitigation:**
- File permissions (600)
- Warn users about security
- Support environment variable override
- Future: Encryption at rest
- Never log API keys

### Risk 3: Breaking Changes in Settings Schema
**Impact:** Medium  
**Probability:** High (as features evolve)  
**Mitigation:**
- Version settings schema
- Automatic migration on upgrade
- Backwards compatibility
- Clear upgrade notes
- Validate migrated settings

### Risk 4: Too Many Settings (Overwhelming)
**Impact:** Medium  
**Probability:** Medium  
**Mitigation:**
- Good defaults (most users don't need to change)
- Progressive disclosure (basic vs advanced)
- Search/filter in settings (future)
- Help text for each setting
- Presets for common configurations (future)

---

## Future Enhancements

### Phase 2 Ideas
- **Encryption:** Encrypt sensitive settings (API keys)
- **Presets:** Common configuration templates
- **Profiles:** Switch between different setting sets
- **Search:** Find settings by keyword
- **Validation UI:** Show which settings are invalid with fixes

### Phase 3 Ideas
- **Cloud Sync:** Sync settings across machines
- **Team Settings:** Share approved settings with team
- **Settings History:** Track changes over time
- **Import from Other Tools:** Migrate from Cursor, Copilot, etc.
- **Settings API:** Programmatic configuration

---

## Open Questions

1. **Should we encrypt API keys by default?**
   - Pro: Better security
   - Con: Adds complexity, key management
   - Decision: Phase 2 feature, file permissions sufficient for now

2. **Should we support YAML/TOML formats?**
   - Pro: More human-friendly for manual editing
   - Con: More dependencies, parsing complexity
   - Decision: JSON only for now, consider YAML export in Phase 2

3. **Should we have per-workspace settings?**
   - Use case: Different settings per project
   - Complexity: Merging global + workspace settings
   - Decision: Phase 3 feature if requested

4. **Should settings overlay have search?**
   - Pro: Easier to find specific settings
   - Con: More UI complexity
   - Decision: Add when >50 total settings

---

## Related Documentation

- [ADR-0017: Auto-Approval and Settings System](../adr/0017-auto-approval-and-settings-system.md)
- [How-to: Use TUI Interface - Settings](../how-to/use-tui-interface.md#settings)
- [Configuration Guide](../reference/configuration.md) (if exists)
- [Security Best Practices](../../SECURITY.md)

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2024-12 | 1.0 | Initial PRD creation |
