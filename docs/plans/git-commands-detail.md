# Git Commands Implementation Detail

Detailed technical specifications for `/commit` and `/pr` slash commands.

## `/commit` Command

### Overview
Automatically commits changes made by the agent during the current session. The commit message is optional - if not provided, an LLM generates a conventional commit message based on the actual file changes.

### Usage
```
/commit                              # Auto-generate message
/commit chore: update dependencies   # Custom message
/commit fix navbar rendering issue   # Custom message
```

### File Tracking System

**Approach: Tool-Level Tracking**
- Intercept `write_file` and `apply_diff` tool executions
- Maintain in-memory set of modified file paths
- Store modification metadata (timestamp, operation type)

```go
// pkg/agent/tools/tracking.go (new)
type FileModification struct {
    Path      string
    Operation string    // "write", "diff"
    Timestamp time.Time
}

type ModificationTracker struct {
    mu            sync.RWMutex
    modifications map[string]*FileModification
}
```

### Commit Message Generation

Use LLM to generate conventional commit messages:

```
feat(auth): add user authentication
fix(api): resolve null pointer in handler
refactor(ui): simplify component structure
```

## `/pr` Command

### Overview
Creates pull request from current branch with AI-generated title and description. The title is optional - if not provided, the LLM generates both title and description by analyzing commits and diffs.

### Usage
```
/pr                                  # Fully auto-generated
/pr Add user authentication system   # Custom title, auto description
```

### Base Branch Detection

The command automatically detects where the current branch diverged from:

```go
func DetectBaseBranch(workingDir string) (string, error) {
    // Try common base branches in order
    baseBranches := []string{"main", "master", "develop"}
    
    currentBranch, err := getCurrentBranch(workingDir)
    if err != nil {
        return "", err
    }
    
    for _, base := range baseBranches {
        // Check if base branch exists
        cmd := exec.Command("git", "rev-parse", "--verify", base)
        cmd.Dir = workingDir
        if err := cmd.Run(); err != nil {
            continue
        }
        
        // Find merge base (where branches diverged)
        cmd = exec.Command("git", "merge-base", base, currentBranch)
        cmd.Dir = workingDir
        output, err := cmd.Output()
        if err == nil && len(output) > 0 {
            return base, nil
        }
    }
    
    return "", fmt.Errorf("could not detect base branch")
}
```

### Commit Analysis

Collect all commits between base and current branch:

```go
func GetCommitsSinceBase(workingDir, base, head string) ([]CommitInfo, error) {
    // Get commit list with messages
    cmd := exec.Command("git", "log", "--format=%h|%s", fmt.Sprintf("%s..%s", base, head))
    cmd.Dir = workingDir
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    lines := strings.Split(strings.TrimSpace(string(output)), "\n")
    commits := make([]CommitInfo, 0, len(lines))
    
    for _, line := range lines {
        if line == "" {
            continue
        }
        parts := strings.SplitN(line, "|", 2)
        if len(parts) == 2 {
            commits = append(commits, CommitInfo{
                Hash:    parts[0],
                Message: parts[1],
            })
        }
    }
    
    return commits, nil
}
```

### Diff Analysis

Get material changes (actual code diffs, not just commit messages):

```go
func GetDiffSummary(workingDir, base, head string) (string, error) {
    // Get diff stats
    cmd := exec.Command("git", "diff", "--stat", fmt.Sprintf("%s...%s", base, head))
    cmd.Dir = workingDir
    stats, err := cmd.Output()
    if err != nil {
        return "", err
    }
    
    // Get actual changes (truncated for LLM)
    cmd = exec.Command("git", "diff", fmt.Sprintf("%s...%s", base, head))
    cmd.Dir = workingDir
    diff, err := cmd.Output()
    if err != nil {
        return "", err
    }
    
    // Truncate diff to reasonable size for LLM context
    diffPreview := truncateDiff(string(diff), 5000)
    
    return fmt.Sprintf("Files Changed:\n%s\n\nCode Changes:\n%s", stats, diffPreview), nil
}

func truncateDiff(diff string, maxChars int) string {
    if len(diff) <= maxChars {
        return diff
    }
    return diff[:maxChars] + "\n... (diff truncated)"
}
```

### PR Content Generation

```go
func (g *PRGenerator) Generate(
    commits []CommitInfo,
    diffSummary string,
    baseBranch string,
    headBranch string,
    customTitle string,
) (*PRContent, error) {
    prompt := g.buildPRPrompt(commits, diffSummary, baseBranch, headBranch, customTitle)
    
    response, err := g.llmClient.Generate(context.Background(), prompt)
    if err != nil {
        return nil, fmt.Errorf("failed to generate PR content: %w", err)
    }
    
    content := parsePRContent(response)
    
    // Use custom title if provided
    if customTitle != "" {
        content.Title = customTitle
    }
    
    return content, nil
}

func (g *PRGenerator) buildPRPrompt(
    commits []CommitInfo,
    diffSummary string,
    base, head string,
    customTitle string,
) string {
    var sb strings.Builder
    
    if customTitle != "" {
        sb.WriteString("Generate a pull request description (title already provided).\n\n")
        sb.WriteString(fmt.Sprintf("Title: %s\n\n", customTitle))
    } else {
        sb.WriteString("Generate a pull request title and description.\n\n")
    }
    
    sb.WriteString(fmt.Sprintf("Base: %s -> Head: %s\n\n", base, head))
    
    sb.WriteString("Commits:\n")
    for _, commit := range commits {
        sb.WriteString(fmt.Sprintf("- %s: %s\n", commit.Hash, commit.Message))
    }
    
    sb.WriteString("\nMaterial Changes (from git diff):\n")
    sb.WriteString(diffSummary)
    
    sb.WriteString("\n\nGenerate in this format:\n")
    if customTitle == "" {
        sb.WriteString("TITLE: <concise, actionable summary>\n\n")
    }
    sb.WriteString("DESCRIPTION:\n")
    sb.WriteString("## Summary\n")
    sb.WriteString("<what changed and why>\n\n")
    sb.WriteString("## Changes\n")
    sb.WriteString("- <key changes from actual diffs>\n\n")
    sb.WriteString("## Testing\n")
    sb.WriteString("<how to verify these changes>\n")
    
    return sb.String()
}
```

### Git Provider Detection

```go
func DetectProvider(workingDir string) (GitProvider, error) {
    // Parse git remote URL
    // Detect github.com, gitlab.com, etc.
}
```

### GitHub Integration

- Use GitHub REST API v3
- Require GITHUB_TOKEN environment variable
- Create PR with generated title and description

### Token Management

1. Check GITHUB_TOKEN / GITLAB_TOKEN env vars
2. Fall back to config file
3. Guide user through setup if missing

## Implementation Files

```
pkg/agent/
├── git/
│   ├── tracker.go       # File modification tracking
│   ├── commit.go        # Commit message generation
│   ├── pr.go            # PR creation logic
│   ├── provider.go      # Git provider detection
│   ├── github.go        # GitHub API client
│   └── auth.go          # Token management
└── commands/
    ├── commit.go        # /commit handler
    └── pr.go            # /pr handler
```
