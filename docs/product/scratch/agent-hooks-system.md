# Feature Idea: Agent Hooks & Lifecycle System

**Status:** Draft  
**Priority:** CRITICAL - High Impact, Near-Term  
**Last Updated:** November 2025

---

## Overview

A programmable hook system that allows developers to inject custom logic at key points in Forge's agent lifecycle. Define actions that execute before/after tool calls, agent loop iterations, LLM calls, or custom events. Enables workflow customization, validation, quality gates, monitoring, and integration with external systems without modifying Forge's core.

---

## Problem Statement

Forge's agent loop is currently a black box:
- Cannot inject custom validation logic
- No way to enforce team-specific rules
- Cannot integrate with external systems automatically
- No pre/post processing of tool calls
- Cannot implement custom quality checks
- No way to add telemetry/monitoring
- Cannot customize workflow without code changes

Developers need extension points to:
- Validate tool inputs before execution
- Verify tool outputs after execution
- Run tests after code changes
- Format code automatically
- Send notifications to Slack/Discord
- Log to external systems
- Enforce security policies
- Implement approval workflows
- Add custom business logic

Current workarounds are insufficient:
- Manual tool calls - Not automatic
- External scripts - Disconnected from Forge
- System prompts - Limited, no real execution
- Wrapper scripts - Brittle, no integration

This leads to:
- Cannot enforce organizational standards
- No integration with existing tools
- Manual verification required
- Cannot customize behavior per project
- Limited extensibility
- One-size-fits-all approach

---

## Key Capabilities

### Hook Types

**Tool Hooks:**
- `before_tool_call` - Validate/modify tool inputs
- `after_tool_call` - Verify tool outputs, run checks
- `on_tool_error` - Handle tool failures
- `tool_call_filter` - Block certain tool calls conditionally

**Agent Loop Hooks:**
- `before_agent_iteration` - Pre-processing each loop
- `after_agent_iteration` - Post-processing each loop
- `on_task_start` - When task begins
- `on_task_complete` - When task finishes
- `on_task_error` - When task fails

**LLM Hooks:**
- `before_llm_call` - Modify prompts, add context
- `after_llm_call` - Process responses, validate output
- `on_llm_error` - Handle API failures

**File Operation Hooks:**
- `before_file_write` - Validate file changes
- `after_file_write` - Auto-format, lint, test
- `before_file_delete` - Confirm deletion
- `on_file_change` - React to any file modification

**Custom Event Hooks:**
- `on_custom_event` - User-defined events
- `on_state_change` - Agent state transitions

### Hook Configuration

**YAML Configuration (.forge/hooks.yml):**
```yaml
hooks:
  # Run linter after any file write
  after_file_write:
    - name: auto_lint
      command: golangci-lint run --fix {file}
      async: false
      continue_on_error: false
    
    - name: auto_format
      command: gofmt -w {file}
      async: false
    
    - name: notify_slack
      command: ./scripts/notify-slack.sh "{file} modified"
      async: true
      continue_on_error: true

  # Validate tool inputs before execution
  before_tool_call:
    - name: validate_write_permissions
      when: tool == "write_file"
      script: |
        #!/bin/bash
        if [[ "$FILE_PATH" == vendor/* ]]; then
          echo "ERROR: Cannot write to vendor directory"
          exit 1
        fi
      env:
        FILE_PATH: "{args.path}"

  # Run tests after code changes
  after_tool_call:
    - name: run_tests_on_code_change
      when: |
        tool in ["write_file", "apply_diff"] and
        args.path matches "src/**/*.go"
      command: go test ./...
      timeout: 5m
      retry_on_failure: 1
      
  # Quality gates before task completion
  before_task_complete:
    - name: verify_tests_pass
      command: go test ./...
      required: true
      
    - name: verify_lint_passes
      command: golangci-lint run
      required: true
      
    - name: verify_build
      command: go build ./...
      required: true

  # Security scanning on file changes
  after_file_write:
    - name: security_scan
      when: args.path matches "**/*.go"
      command: gosec {file}
      continue_on_error: true
      
  # Custom business logic
  on_custom_event:
    - name: deployment_check
      event: ready_to_deploy
      script: ./scripts/pre-deploy-checks.sh
      required: true

  # Monitoring and observability
  after_agent_iteration:
    - name: log_metrics
      command: |
        curl -X POST https://metrics.company.com/forge \
          -d "iteration=$ITERATION" \
          -d "tools_used=$TOOLS_USED"
      async: true
      continue_on_error: true
```

**Programmatic Hooks (Go API):**
```go
forge.RegisterHook("after_file_write", func(ctx HookContext) error {
    filePath := ctx.Args["path"].(string)
    
    // Auto-format
    if strings.HasSuffix(filePath, ".go") {
        cmd := exec.Command("gofmt", "-w", filePath)
        if err := cmd.Run(); err != nil {
            return err
        }
    }
    
    // Run tests
    if strings.Contains(filePath, "src/") {
        cmd := exec.Command("go", "test", "./...")
        return cmd.Run()
    }
    
    return nil
})
```

### Hook Context

**Available Data in Hooks:**
```yaml
# Context variables available in all hooks
context:
  workspace_root: /path/to/workspace
  task_id: task-123
  iteration: 5
  tool_name: write_file
  tool_args:
    path: src/main.go
    content: "package main..."
  tool_result: success
  llm_model: claude-3-5-sonnet
  tokens_used: 1234
  files_modified: ["src/main.go", "src/util.go"]
  
# Environment variables automatically set
env:
  FORGE_WORKSPACE: /path/to/workspace
  FORGE_TASK_ID: task-123
  FORGE_ITERATION: 5
  FORGE_TOOL: write_file
  FORGE_FILE: src/main.go  # for file hooks
```

### Hook Execution

**Execution Models:**
- **Synchronous** - Block agent until hook completes
- **Asynchronous** - Run in background, don't block
- **Required** - Must succeed or task fails
- **Optional** - Failure logged but doesn't block

**Error Handling:**
```yaml
hooks:
  after_file_write:
    - name: lint
      command: eslint {file}
      continue_on_error: false  # Stop on failure
      retry_on_failure: 2       # Retry twice
      timeout: 30s              # Kill after 30s
      
    - name: notify
      command: ./notify.sh
      continue_on_error: true   # Continue even if fails
      async: true               # Don't wait
```

### Conditional Execution

**When Conditions:**
```yaml
hooks:
  before_tool_call:
    # Only for write operations
    - name: backup
      when: tool == "write_file"
      command: cp {args.path} {args.path}.bak
    
    # Only for Go files
    - name: go_lint
      when: |
        tool == "write_file" and
        args.path ends_with ".go"
      command: golangci-lint run {args.path}
    
    # Only in CI mode
    - name: strict_checks
      when: env.CI == "true"
      command: ./strict-validation.sh
    
    # Complex conditions
    - name: security_review
      when: |
        tool == "write_file" and
        (args.path matches "auth/**" or args.path matches "security/**") and
        iteration > 10
      command: ./security-scan.sh {args.path}
```

**Pattern Matching:**
- `==`, `!=` - Equality
- `matches` - Glob/regex matching
- `in` - List membership
- `>`, `<`, `>=`, `<=` - Numeric comparison
- `and`, `or`, `not` - Boolean logic
- `starts_with`, `ends_with`, `contains` - String matching

---

## Example Use Cases

### 1. Automatic Code Quality Enforcement

```yaml
# .forge/hooks.yml
hooks:
  # Format on every file write
  after_file_write:
    - name: auto_format
      when: args.path ends_with ".go"
      command: gofmt -w {file}
      async: false
      
    - name: organize_imports
      when: args.path ends_with ".go"
      command: goimports -w {file}
      async: false
      
    - name: lint_fix
      when: args.path ends_with ".go"
      command: golangci-lint run --fix {file}
      continue_on_error: true

  # Run affected tests after code changes
  after_tool_call:
    - name: run_tests
      when: |
        tool in ["write_file", "apply_diff"] and
        args.path matches "src/**/*.go"
      command: go test ./... -run TestsFor{basename}
      timeout: 5m
      
  # Quality gates before completing task
  before_task_complete:
    - name: full_test_suite
      command: go test ./...
      required: true
      
    - name: lint_check
      command: golangci-lint run
      required: true
      
    - name: security_scan
      command: gosec ./...
      required: true
```

**Execution:**
```
User: Add a new user handler

[Agent loop starts]

[Agent calls write_file on src/handlers/user.go]

Hook: after_file_write/auto_format
  Running: gofmt -w src/handlers/user.go
  ✓ Completed (45ms)

Hook: after_file_write/organize_imports
  Running: goimports -w src/handlers/user.go
  ✓ Completed (123ms)

Hook: after_file_write/lint_fix
  Running: golangci-lint run --fix src/handlers/user.go
  ⚠ Found 2 issues, fixed automatically
  ✓ Completed (1.2s)

Hook: after_tool_call/run_tests
  Running: go test ./handlers -run TestsForUser
  ✓ All tests passed (3.4s)

[Agent continues...]

[Agent calls task_completion]

Hook: before_task_complete/full_test_suite
  Running: go test ./...
  ✓ 245 tests passed (12.3s)

Hook: before_task_complete/lint_check
  Running: golangci-lint run
  ✓ No issues found (2.1s)

Hook: before_task_complete/security_scan
  Running: gosec ./...
  ✓ No vulnerabilities found (1.8s)

✅ All quality gates passed!

Task completed successfully.
```

### 2. Integration with External Systems

```yaml
hooks:
  # Notify team on code changes
  after_file_write:
    - name: slack_notification
      command: |
        curl -X POST $SLACK_WEBHOOK \
          -d "{\"text\": \"Forge modified: {file}\"}"
      async: true
      continue_on_error: true
      env:
        SLACK_WEBHOOK: ${SLACK_WEBHOOK_URL}
  
  # Log to monitoring system
  after_agent_iteration:
    - name: log_metrics
      script: |
        #!/bin/bash
        curl -X POST https://metrics.company.com/api/events \
          -H "Authorization: Bearer $API_TOKEN" \
          -d "{
            \"event\": \"forge_iteration\",
            \"iteration\": $FORGE_ITERATION,
            \"tool\": \"$FORGE_TOOL\",
            \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
          }"
      async: true
      continue_on_error: true
  
  # Create Jira ticket on errors
  on_task_error:
    - name: create_jira_ticket
      script: ./scripts/create-jira-incident.sh
      env:
        ERROR_MESSAGE: "{error.message}"
        TASK_ID: "{task.id}"
```

### 3. Security Policy Enforcement

```yaml
hooks:
  # Prevent writing to protected directories
  before_tool_call:
    - name: protect_vendor
      when: tool == "write_file" and args.path starts_with "vendor/"
      script: |
        #!/bin/bash
        echo "ERROR: Cannot modify vendor directory"
        exit 1
      required: true
    
    - name: protect_secrets
      when: tool == "write_file"
      script: |
        #!/bin/bash
        if grep -q "API_KEY\|PASSWORD\|SECRET" "$FORGE_FILE_CONTENT"; then
          echo "ERROR: Possible secret detected in file"
          exit 1
        fi
      env:
        FORGE_FILE_CONTENT: "{args.content}"
      required: true
  
  # Scan for vulnerabilities after changes
  after_file_write:
    - name: security_scan
      when: args.path matches "**/*.{go,js,py}"
      command: trivy fs {file}
      continue_on_error: true
      
  # Audit log all operations
  after_tool_call:
    - name: audit_log
      script: |
        #!/bin/bash
        echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) - $FORGE_TOOL - $FORGE_USER - $FORGE_FILE" >> audit.log
      async: true
```

### 4. Development Workflow Automation

```yaml
hooks:
  # Auto-generate mocks when interfaces change
  after_file_write:
    - name: generate_mocks
      when: |
        args.path matches "**/*_interface.go" or
        content contains "type.*interface"
      command: mockgen -source={file} -destination={file}_mock.go
      
  # Update API documentation
  after_file_write:
    - name: update_swagger
      when: args.path matches "api/**/*.go"
      command: swag init
      async: true
      
  # Regenerate GraphQL schema
  after_file_write:
    - name: graphql_codegen
      when: args.path ends_with ".graphql"
      command: go run github.com/99designs/gqlgen generate
      
  # Database migrations
  after_file_write:
    - name: create_migration
      when: args.path matches "models/**/*.go"
      script: |
        #!/bin/bash
        if git diff --name-only | grep -q models/; then
          ./scripts/generate-migration.sh
        fi
```

### 5. Custom Approval Workflows

```yaml
hooks:
  # Require approval for production changes
  before_tool_call:
    - name: require_approval_prod
      when: |
        env.ENVIRONMENT == "production" and
        tool in ["write_file", "apply_diff", "execute_command"]
      script: |
        #!/bin/bash
        echo "Production change requires approval"
        echo "Waiting for approval..."
        
        # Call approval API
        APPROVAL_ID=$(curl -X POST https://approvals.company.com/api/requests \
          -d "{\"tool\": \"$FORGE_TOOL\", \"file\": \"$FORGE_FILE\"}" \
          | jq -r .id)
        
        # Poll for approval (timeout 5 min)
        for i in {1..60}; do
          STATUS=$(curl https://approvals.company.com/api/requests/$APPROVAL_ID \
            | jq -r .status)
          
          if [ "$STATUS" == "approved" ]; then
            echo "✓ Approved!"
            exit 0
          elif [ "$STATUS" == "rejected" ]; then
            echo "✗ Rejected"
            exit 1
          fi
          
          sleep 5
        done
        
        echo "✗ Timeout waiting for approval"
        exit 1
      required: true
      timeout: 6m
```

### 6. Testing & Validation Pipeline

```yaml
hooks:
  # Progressive test running
  after_file_write:
    # Fast: Unit tests for changed file only
    - name: unit_tests_fast
      when: args.path matches "src/**/*.go"
      command: go test -run TestsFor{basename} ./...
      timeout: 1m
      
  after_tool_call:
    # Medium: Integration tests for related modules
    - name: integration_tests
      when: |
        tool in ["write_file", "apply_diff"] and
        iteration % 5 == 0
      command: go test -tags=integration ./...
      timeout: 5m
      
  before_task_complete:
    # Slow: Full test suite before completion
    - name: full_test_suite
      command: go test ./...
      required: true
      timeout: 10m
      
    # End-to-end tests
    - name: e2e_tests
      command: npm run test:e2e
      required: true
      timeout: 15m
      
    # Performance tests
    - name: benchmark
      command: go test -bench=. -benchmem ./...
      continue_on_error: true
```

---

## Technical Approach

### Hook Registry

**Hook Manager:**
```go
type HookManager struct {
    hooks map[HookType][]Hook
    mu    sync.RWMutex
}

type Hook struct {
    Name              string
    Type              HookType
    Command           string
    Script            string
    When              *Condition
    Async             bool
    ContinueOnError   bool
    Required          bool
    Timeout           time.Duration
    RetryOnFailure    int
    Env               map[string]string
}

type HookContext struct {
    HookType    HookType
    ToolName    string
    ToolArgs    map[string]interface{}
    ToolResult  interface{}
    WorkspaceRoot string
    TaskID      string
    Iteration   int
    FilesModified []string
    LLMModel    string
    TokensUsed  int
}

func (hm *HookManager) ExecuteHooks(hookType HookType, ctx *HookContext) error {
    hm.mu.RLock()
    hooks := hm.hooks[hookType]
    hm.mu.RUnlock()
    
    for _, hook := range hooks {
        // Check condition
        if hook.When != nil &amp;&amp; !hook.When.Evaluate(ctx) {
            continue
        }
        
        // Execute hook
        if err := hm.executeHook(hook, ctx); err != nil {
            if hook.Required || !hook.ContinueOnError {
                return err
            }
            log.Warn("Hook failed but continuing: %s", hook.Name)
        }
    }
    
    return nil
}
```

### Hook Execution

**Script Runner:**
```go
func (hm *HookManager) executeHook(hook Hook, ctx *HookContext) error {
    // Prepare environment
    env := hm.buildEnv(hook, ctx)
    
    // Prepare command
    cmd := hm.interpolateCommand(hook.Command, ctx)
    
    // Create exec.Command
    execCmd := exec.Command("sh", "-c", cmd)
    execCmd.Env = env
    execCmd.Dir = ctx.WorkspaceRoot
    
    // Set timeout
    cmdCtx, cancel := context.WithTimeout(context.Background(), hook.Timeout)
    defer cancel()
    
    // Execute
    if hook.Async {
        go hm.runAsync(execCmd, hook)
        return nil
    }
    
    // Synchronous execution with retry
    var err error
    for attempt := 0; attempt &lt;= hook.RetryOnFailure; attempt++ {
        err = execCmd.Run()
        if err == nil {
            return nil
        }
        
        if attempt &lt; hook.RetryOnFailure {
            time.Sleep(time.Second * time.Duration(attempt+1))
        }
    }
    
    return err
}
```

### Condition Evaluator

**Expression Parser:**
```go
type Condition struct {
    expression string
    ast        *ASTNode
}

func (c *Condition) Evaluate(ctx *HookContext) bool {
    evaluator := &Evaluator{ctx: ctx}
    result := evaluator.Eval(c.ast)
    return result.(bool)
}

// Example expressions:
// - tool == "write_file"
// - args.path matches "src/**/*.go"
// - iteration > 10 and env.CI == "true"
// - args.path in ["main.go", "util.go"]
```

### Configuration Loader

**YAML Parser:**
```go
type HooksConfig struct {
    Hooks map[string][]HookDefinition `yaml:"hooks"`
}

type HookDefinition struct {
    Name            string            `yaml:"name"`
    Command         string            `yaml:"command"`
    Script          string            `yaml:"script"`
    When            string            `yaml:"when"`
    Async           bool              `yaml:"async"`
    ContinueOnError bool              `yaml:"continue_on_error"`
    Required        bool              `yaml:"required"`
    Timeout         string            `yaml:"timeout"`
    RetryOnFailure  int               `yaml:"retry_on_failure"`
    Env             map[string]string `yaml:"env"`
}

func LoadHooksConfig(path string) (*HooksConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var config HooksConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

### Integration Points

**Agent Loop Integration:**
```go
// Before each tool call
if err := hookManager.ExecuteHooks(BeforeToolCall, &HookContext{
    ToolName: toolName,
    ToolArgs: args,
}); err != nil {
    return err
}

// Execute tool
result, err := executor.ExecuteTool(toolName, args)

// After tool call
if err := hookManager.ExecuteHooks(AfterToolCall, &HookContext{
    ToolName:   toolName,
    ToolArgs:   args,
    ToolResult: result,
}); err != nil {
    return err
}
```

---

## Value Propositions

### For All Teams
- Customize Forge to team needs
- Enforce organizational standards
- Integrate with existing tools
- No code changes to Forge

### For Platform Teams
- Build internal developer platforms
- Standardize workflows
- Policy enforcement
- Compliance automation

### For Security Teams
- Security policy enforcement
- Audit logging
- Vulnerability scanning
- Approval workflows

---

## Implementation Phases

### Phase 1: Core Hook System (2 weeks)
- Hook registry and manager
- Basic hook types (tool, file)
- YAML configuration
- Synchronous execution

### Phase 2: Advanced Features (2 weeks)
- Condition system
- Async execution
- Retry logic
- Environment variables

### Phase 3: Integration (1 week)
- Agent loop integration
- LLM call hooks
- Custom events
- Error handling

### Phase 4: Polish (1 week)
- Hook templates/presets
- Documentation
- Testing
- Performance optimization

---

## Success Metrics

**Adoption:**
- 70%+ teams define custom hooks
- 50%+ use quality gate hooks
- 60%+ integrate external systems
- 80%+ report improved workflow

**Impact:**
- 90% reduction in manual quality checks
- 80% faster integration with tools
- 70% better code quality consistency
- 100% policy compliance

**Satisfaction:**
- 4.9+ rating for hooks system
- "Game changer for customization"
- "Finally can enforce our standards"
