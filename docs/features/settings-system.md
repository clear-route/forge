# Settings System

## Overview

The Forge settings system provides a flexible, extensible architecture for managing application configuration. It supports both auto-approval of tools and command whitelisting, with a dedicated `/settings` command in the TUI for viewing current settings.

## Architecture

### Core Components

#### 1. Configuration Package (`pkg/config/`)

The configuration system is built around a section-based architecture that allows for easy extension:

**Section Interface** (`section.go`)
```go
type Section interface {
    ID() string
    Title() string
    Description() string
    Data() map[string]interface{}
    SetData(data map[string]interface{}) error
    Validate() error
    Reset()
}
```

**Manager** (`manager.go`)
- Coordinates all configuration sections
- Handles loading/saving of configuration
- Provides thread-safe access to settings

**Store** (`store.go`)
- Persists configuration to `~/.forge/config.json`
- Atomic file writes for data safety
- Thread-safe with mutex protection

#### 2. Available Sections

**Auto-Approval Section** (`auto_approval.go`)
- Dynamic tool discovery - tools are added as encountered
- Per-tool approval toggles
- Default: all tools require approval (safety-first)
- Helper: `config.IsToolAutoApproved(toolName)`

**Command Whitelist Section** (`whitelist.go`)
- Pattern-based command matching
- Supports wildcards (e.g., "npm" matches all npm commands)
- Specific matches (e.g., "npm install" matches only "npm install")
- Helper: `config.IsCommandWhitelisted(command)`

### 3. TUI Integration

**Settings Overlay** (`pkg/executor/tui/settings.go`)
- Read-only view of current settings
- Shows all registered configuration sections
- Displays summary statistics
- Accessible via `/settings` command

**Slash Command** (`pkg/executor/tui/slash_commands.go`)
- `/settings` - Opens the settings overlay
- Registered in the command registry

## Usage

### Accessing Settings in Code

```go
import "github.com/entrhq/forge/pkg/config"

// Check if a tool is auto-approved
if config.IsToolAutoApproved("read_file") {
    // Auto-approve
}

// Check if a command is whitelisted
if config.IsCommandWhitelisted("npm install") {
    // Auto-approve
}
```

### Viewing Settings in TUI

Type `/settings` in the Forge TUI to view the current configuration:

```
⚙️  Settings

Press [Esc] or [q] to close • Settings are saved to ~/.forge/config.json

▸ Tool Auto-Approval
  Automatically approve specified tools without prompting
  2/5 tools auto-approved
    ✓ read_file
    ✓ list_files
    ... and 0 more

▸ Command Whitelist
  Whitelisted command patterns for execute_command tool
  3 command pattern(s) whitelisted
    ✓ npm install - Install npm dependencies
    ✓ go build - Build Go projects
    ✓ git status - Check git status
```

## Configuration File Format

Settings are stored in `~/.forge/config.json`:

```json
{
  "auto_approval": {
    "read_file": true,
    "write_to_file": false,
    "execute_command": false,
    "apply_diff": true,
    "list_files": true
  },
  "command_whitelist": {
    "patterns": [
      {
        "pattern": "npm install",
        "description": "Install npm dependencies"
      },
      {
        "pattern": "go build",
        "description": "Build Go projects"
      },
      {
        "pattern": "git status",
        "description": "Check git status"
      }
    ]
  }
}
```

## Integration with Agent Flow

The auto-approval system is integrated into the agent's approval flow (`pkg/agent/default.go`):

1. When a tool call requires approval, the agent first checks the configuration
2. For `execute_command`, it checks the command whitelist
3. For other tools, it checks the auto-approval settings
4. If auto-approved, the tool executes immediately without user prompt
5. Otherwise, the normal approval flow is followed

```go
// For execute_command, check command whitelist
if toolCall.ToolName == "execute_command" {
    if cmd, ok := argsMap["command"].(string); ok {
        if config.IsCommandWhitelisted(cmd) {
            return true, false // Auto-approve
        }
    }
} else if config.IsToolAutoApproved(toolCall.ToolName) {
    return true, false // Auto-approve
}
```

## Extensibility

### Adding New Settings Sections

To add a new settings section:

1. Create a new file in `pkg/config/` (e.g., `my_section.go`)
2. Implement the `Section` interface
3. Register it with the manager in `config.go`:

```go
func Initialize(configPath string) error {
    // ... existing code ...
    
    mySection := &MySection{}
    if err := globalManager.RegisterSection(mySection); err != nil {
        return err
    }
    
    // ... existing code ...
}
```

4. The new section will automatically appear in the `/settings` overlay

### Pattern Matching Rules

Command whitelist patterns follow these rules:

1. **Exact match**: Pattern "npm install" matches only "npm install"
2. **Prefix match**: Pattern "npm" matches "npm install", "npm run dev", etc.
3. **Case-sensitive**: Patterns are matched exactly as specified
4. **First match wins**: First matching pattern determines approval

## Future Enhancements

The current implementation provides a read-only view. Future enhancements include:

- [ ] Interactive settings editor with keyboard navigation
- [ ] Toggle switches for boolean settings
- [ ] List management for command patterns
- [ ] Real-time validation and feedback
- [ ] Import/export settings profiles
- [ ] Setting presets for common workflows

## Security Considerations

- Auto-approval should be used carefully for security-sensitive tools
- Command whitelist patterns should be as specific as possible
- Default behavior is safe: all tools require approval
- Configuration file permissions should be restricted (user-only access)
- Review auto-approved tools periodically

## Files Modified/Created

### New Files
- `pkg/config/section.go` - Section interface
- `pkg/config/manager.go` - Configuration manager
- `pkg/config/store.go` - File persistence
- `pkg/config/auto_approval.go` - Auto-approval section
- `pkg/config/whitelist.go` - Command whitelist section
- `pkg/config/config.go` - Global initialization and helpers
- `pkg/executor/tui/settings.go` - Settings overlay UI

### Modified Files
- `cmd/forge/main.go` - Initialize config system
- `pkg/agent/default.go` - Integrate auto-approval checks
- `pkg/executor/tui/slash_commands.go` - Register `/settings` command
- `pkg/executor/tui/overlay.go` - Add `OverlayModeSettings` enum

## Testing

To test the settings system:

1. Run Forge: `go run ./cmd/forge`
2. Type `/settings` to view current configuration
3. Exit and edit `~/.forge/config.json` manually
4. Restart Forge and verify changes appear
5. Test auto-approval by adding tools to the config
6. Test command whitelist by adding patterns

## Notes

- Settings are loaded once at startup
- Changes to the config file require restart to take effect
- Invalid JSON in config file will cause initialization to fail
- Missing config file is created automatically with defaults