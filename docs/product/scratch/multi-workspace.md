# Feature Idea: Multi-Workspace & Project Management

**Status:** Draft  
**Priority:** High Impact, Medium-Term  
**Last Updated:** November 2025

---

## Overview

Handle multiple projects and complex workspace structures seamlessly. Enable developers to work across multiple repositories, switch between project contexts, share configurations, bootstrap new projects from templates, and maintain separate conversation histories for different work streams.

---

## Problem Statement

Developers work on multiple projects:
- Switching between projects requires restarting Forge
- Each project loses conversation context
- Configurations don't transfer between projects
- Starting new projects is manual and repetitive
- No way to organize work by project

This leads to:
- Lost productivity from context switching
- Repeated setup for each project
- Mixed conversation histories
- Manual project initialization
- Difficulty managing multiple codebases

---

## Key Capabilities

### Workspace Profiles
- Save workspace configurations
- Switch contexts instantly
- Preserve conversation history per workspace
- Project-specific settings
- Quick workspace switching

### Multi-Repo Operations
- Work across multiple repositories
- Coordinate changes across repos
- Sync related repositories
- Monorepo support
- Cross-repo refactoring

### Project Templates
- Bootstrap new projects from templates
- Community template marketplace
- Custom template creation
- Project scaffolding with AI assistance
- Best practices enforcement

### Task Management
- Track coding tasks within Forge
- Link commits to task IDs
- Task prioritization
- Progress tracking
- Sprint planning integration

### Session Bookmarks
- Save conversation states
- Resume from specific points
- Branch conversation threads
- Share session snapshots
- Session versioning

---

## Example Use Cases

### 1. Switch Between Projects

```
User: /workspace switch api-backend

Agent: Switching to workspace: api-backend
📂 /Users/dev/projects/api-backend
🔧 Loading project settings...
💬 Restored conversation (last active: 2 hours ago)

Previous context:
We were implementing the payment processing endpoint.
Last action: Created test suite with 85% coverage.

Ready to continue!
```

### 2. Multi-Repo Changes

```
User: Update the API version to 2.0 across all repos

Agent: Found 3 related repositories:
- api-backend (current)
- frontend-app  
- mobile-app

Changes needed:
api-backend:
  - Update openapi.yaml version
  - Bump package.json version
  - Update README

frontend-app:
  - Update API_VERSION constant
  - Update fetch URLs
  
mobile-app:
  - Update base URL config
  - Update API models

Apply changes to all repos?
```

### 3. Create Project from Template

```
User: Create a new Go microservice project

Agent: I'll bootstrap a new project for you.

Template: go-microservice-template
Includes:
✓ Go modules setup
✓ Docker configuration
✓ CI/CD pipeline (GitHub Actions)
✓ Health check endpoints
✓ Logging & metrics
✓ Test structure

Project name? > user-service
Database? > PostgreSQL
Auth type? > JWT

[Creates project structure]
[Initializes git repository]
[Configures CI/CD]

Project ready! Run 'go run cmd/server/main.go' to start.
```

### 4. Task Tracking

```
User: /tasks list

Agent: Current tasks for api-backend:

🟢 In Progress:
- #1: Implement payment processing (75% complete)
  Last: Created test suite
  Next: Add error handling

🔴 Blocked:
- #2: Database migration (waiting on DBA)

⚪ Todo:
- #3: Add rate limiting
- #4: Implement caching
- #5: Update documentation

/tasks show #1 for details
```

### 5. Session Management

```
User: /session bookmark payment-feature

Agent: Session bookmarked: payment-feature

Captured state:
- 47 messages in conversation
- 3 files modified
- 2 branches created
- Current task: Payment processing

Resume anytime with: /session restore payment-feature

Available bookmarks:
- payment-feature (just now)
- auth-refactor (2 days ago)
- api-v2-migration (1 week ago)
```

---

## Technical Approach

### Workspace Management

**Workspace Configuration:**
```json
{
  "name": "api-backend",
  "path": "/Users/dev/projects/api-backend",
  "settings": {
    "autoApproval": {...},
    "llmProvider": "openai",
    "defaultBranch": "main"
  },
  "relatedRepos": [
    "/Users/dev/projects/frontend-app",
    "/Users/dev/projects/mobile-app"
  ],
  "lastActive": "2025-11-21T14:30:00Z"
}
```

**Switching:**
- Save current workspace state
- Load target workspace config
- Restore conversation history
- Update file watchers
- Switch working directory

### Multi-Repo Support

**Repository Detection:**
- Scan for .git directories
- Detect monorepo structure
- Identify related projects
- Track dependencies

**Coordinated Operations:**
- Batch file operations across repos
- Synchronized git operations
- Cross-repo refactoring
- Dependency updates

### Template System

**Template Structure:**
```
template/
├── template.json (metadata)
├── files/ (template files)
├── scripts/ (setup scripts)
└── README.md (instructions)
```

**Template Processing:**
- Variable substitution
- Conditional file inclusion
- Script execution
- Git initialization
- Dependency installation

### Session Persistence

**Session State:**
- Conversation history
- File modifications
- Tool executions
- User preferences
- Bookmark metadata

**Storage:**
- SQLite for session data
- JSON for configurations
- File snapshots
- Incremental updates

---

## Value Propositions

### For Multi-Project Developers
- Seamless project switching
- Organized workflows
- Preserved context
- Efficient multitasking

### For Teams
- Shared workspace configs
- Standard project templates
- Consistent setups
- Knowledge transfer

### For Newcomers
- Quick project setup
- Best practices built-in
- Guided initialization
- Example projects

---

## Implementation Phases

### Phase 1: Basic Workspaces (2-3 weeks)
- Workspace configuration
- Switch workspaces
- Preserve settings
- Quick switching UI

### Phase 2: Multi-Repo (3 weeks)
- Detect related repos
- Coordinate operations
- Cross-repo changes
- Monorepo support

### Phase 3: Templates (2-3 weeks)
- Template structure
- Variable substitution
- Script execution
- Template marketplace

### Phase 4: Session Management (2-3 weeks)
- Save sessions
- Restore sessions
- Session bookmarks
- Session search

### Phase 5: Task Integration (2 weeks)
- Task tracking
- Link to commits
- Progress visualization
- External tool integration

---

## Open Questions

1. **Storage:** Where to store workspace configs?
2. **Sync:** Sync workspaces across machines?
3. **Templates:** How to distribute/discover templates?
4. **Sessions:** How much history to preserve?
5. **Multi-Repo:** Auto-detect related repositories?

---

## Related Features

**Synergies with:**
- **Settings System** - Per-workspace settings
- **Memory System** - Per-workspace conversation history
- **Git Integration** - Multi-repo operations

---

## Success Metrics

**Adoption:**
- 50%+ users create multiple workspaces
- 40%+ use templates
- 30%+ work with multi-repo
- 60%+ use session bookmarks

**Impact:**
- 40% faster project switching
- 60% faster project initialization
- 50% reduction in setup errors
- 70% better context preservation

**Satisfaction:**
- 4.5+ rating for workspace features
- "Essential for my workflow" feedback
