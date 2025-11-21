# Feature Idea: Enhanced File System Operations

**Status:** Draft  
**Priority:** High Impact, Near-Term  
**Last Updated:** November 2025

---

## Overview

Expand file system capabilities beyond basic read/write/list operations to enable sophisticated file manipulation, batch operations, file watching, and intelligent file management. This transforms Forge from a simple file editor into a powerful file system automation tool.

---

## Problem Statement

Current file operations are limited:
- No bulk file operations (rename, move, copy multiple files)
- Cannot watch files for changes
- No advanced search capabilities (content search, regex in filenames)
- Missing file metadata operations (permissions, timestamps)
- No temporary file/directory management
- Cannot handle binary files effectively
- No file compression/decompression
- Missing symbolic link handling

This leads to:
- Manual repetitive file operations
- Inability to react to file system changes
- Limited refactoring capabilities
- Poor support for complex file reorganization
- Manual handling of generated files

---

## Key Capabilities

### Batch File Operations
- Move/copy/rename multiple files at once
- Pattern-based file operations (*.test.js → test/*.js)
- Recursive operations with filters
- Dry-run mode to preview changes
- Atomic operations (all or nothing)

### File Watching
- Monitor files/directories for changes
- Trigger actions on file modifications
- Watch patterns (*.go, test/*)
- Debounced change detection
- Ignore patterns (.git/, node_modules/)

### Advanced Search
- Full-text content search across files
- Multi-pattern search (AND/OR logic)
- Search with file metadata filters
- Search result ranking by relevance
- Export search results

### File Metadata Operations
- Read/modify file permissions
- Update timestamps
- Get file size, type, ownership
- Check file existence efficiently
- Compare file hashes (detect duplicates)

### Smart File Management
- Create temporary files/directories
- Automatic cleanup of temp files
- Safe file operations (backup before modify)
- File versioning (keep N backups)
- Atomic file replacement

### Archive Operations
- Create/extract zip, tar, tar.gz archives
- Archive entire directories
- Selective extraction
- Archive inspection without extraction
- Progress tracking for large archives

---

## Example Use Cases

### 1. Bulk File Reorganization

```
User: Reorganize test files - move all *_test.go files to test/ directory

Agent: I'll reorganize your test files:

Found 47 test files to move:
  auth_test.go → test/auth_test.go
  user_test.go → test/user_test.go
  payment_test.go → test/payment_test.go
  ... (44 more)

[Shows full list in expandable section]

This will:
✓ Create test/ directory if needed
✓ Move 47 files
✓ Update import paths in 12 files
✓ Run tests to verify everything works

Dry run completed - no errors detected
Proceed with move?
```

### 2. File Watching for Auto-Regeneration

```
User: Watch the schema.prisma file and regenerate client when it changes

Agent: I'll set up a file watcher:

Watching: prisma/schema.prisma
On change: Run `npx prisma generate`
Debounce: 500ms (avoid rapid successive runs)

Watcher started ✓

[5 minutes later]
Detected: schema.prisma changed
Running: npx prisma generate
[Shows streaming output]
✓ Client regenerated successfully
```

### 3. Content Search and Replace

```
User: Find all files that import the old logging library and show me where

Agent: Searching for imports of old logger...

Found 23 files with matches:

📄 src/auth/service.go:3
  import "github.com/old/logger"

📄 src/user/repository.go:5
  import "github.com/old/logger"

📄 handlers/api.go:7
  import (
    "github.com/old/logger"
  )

... and 20 more files

Would you like me to:
1. Replace with new logger import
2. Show all matches in detail
3. Export list to file
```

### 4. Duplicate File Detection

```
User: Check for duplicate files in the assets directory

Agent: Scanning assets/ for duplicates...

Found 3 sets of duplicates (12 files, 4.5 MB wasted):

Set 1 (1.2 MB each):
  assets/images/logo.png
  assets/backup/logo.png
  public/logo.png
  [Same content hash: a3f2c9...]

Set 2 (890 KB each):
  assets/icons/arrow.svg
  components/icons/arrow.svg
  [Same content hash: b8e1d4...]

Set 3 (350 KB each):
  styles/fonts/roboto.woff2
  public/fonts/roboto.woff2
  dist/fonts/roboto.woff2
  [Same content hash: c7d2e8...]

Recommended actions:
- Keep one copy of each file
- Update references to use single copy
- Delete duplicates

Apply recommendations?
```

### 5. Archive Creation with Exclusions

```
User: Create a backup archive of the entire project, excluding build artifacts and dependencies

Agent: Creating project backup...

Archive: project-backup-2025-11-21.tar.gz
Including: All project files
Excluding:
  - node_modules/ (14,523 files)
  - dist/ (234 files)
  - .git/ (1,892 files)
  - *.log files (45 files)
  - .DS_Store files (12 files)

Files to archive: 1,234 files (45 MB)
Estimated archive size: 8 MB (82% compression)

[Progress bar]
Creating archive... ████████████████ 100%

✓ Archive created: project-backup-2025-11-21.tar.gz (7.8 MB)
Location: ../backups/
Verify integrity? y
✓ Archive integrity verified (checksum: d9e3f1...)
```

---

## Technical Approach

### Batch Operations Engine

**Operation Planning:**
- Analyze file tree before operations
- Detect conflicts (overwrites, permission issues)
- Calculate disk space requirements
- Create execution plan
- Support dry-run mode

**Safe Execution:**
- Atomic operations where possible
- Create backups before destructive operations
- Rollback on failure
- Track all changes for undo
- Verify operations post-execution

### File Watching System

**Watcher Implementation:**
- Use OS-native file watching APIs (inotify, FSEvents, ReadDirectoryChangesW)
- Efficient monitoring (don't poll)
- Pattern-based filtering
- Debouncing to avoid event floods
- Multiple watchers support

**Event Handling:**
- Detect create, modify, delete, rename
- Capture old/new paths for renames
- Filter noise (temp files, editor backups)
- Queue events for processing
- Handle watcher lifecycle

### Search Engine

**Indexing:**
- Build in-memory index for fast search
- Update index incrementally
- Support incremental updates
- Respect gitignore patterns
- Configurable index scope

**Search Capabilities:**
- Full-text content search
- Regex pattern matching
- File metadata filtering
- Combined queries (content + metadata)
- Result ranking by relevance

### Metadata Operations

**Cross-Platform:**
- Abstract OS-specific details
- Handle permission differences (Unix vs Windows)
- Graceful degradation on unsupported operations
- Clear error messages for platform limitations

**Performance:**
- Batch metadata reads
- Cache results when appropriate
- Minimize syscalls
- Parallel operations where safe

---

## Value Propositions

### For All Developers
- Automate repetitive file operations
- React to file changes automatically
- Find files and content quickly
- Manage project files efficiently

### For Refactoring Heavy Work
- Bulk rename/move operations
- Update imports automatically
- Reorganize directory structures safely
- Maintain consistency across moves

### For Large Codebases
- Fast content search
- Duplicate detection
- Efficient file management
- Archive/backup automation

---

## Implementation Phases

### Phase 1: Batch Operations (2 weeks)
- Copy, move, rename multiple files
- Pattern matching
- Dry-run mode
- Atomic operations

### Phase 2: File Watching (2 weeks)
- File/directory watching
- Pattern-based filtering
- Action triggers
- Debouncing

### Phase 3: Advanced Search (2 weeks)
- Content search
- Metadata filtering
- Result ranking
- Search indexing

### Phase 4: Metadata & Archives (2 weeks)
- Permission operations
- Archive creation/extraction
- Duplicate detection
- Safe file operations

---

## Open Questions

1. **Performance:** How to handle large file trees efficiently?
2. **Watching:** Should watchers persist across sessions?
3. **Search Index:** Store index on disk or rebuild each time?
4. **Platforms:** Which OS-specific features to support?
5. **Limits:** Max files for batch operations?

---

## Related Features

**Synergies with:**
- **Code Intelligence** - Update imports during file moves
- **Git Integration** - Coordinate with git operations
- **Testing** - Auto-run tests on file changes

---

## Success Metrics

**Adoption:**
- 60%+ use batch operations
- 40%+ use file watching
- 50%+ use content search
- 30%+ use archives

**Impact:**
- 70% reduction in manual file operations
- 50% faster project reorganization
- 80% faster finding files/content
- 90% reduction in duplicate files

**Satisfaction:**
- 4.5+ rating for file operations
- "Saves hours of manual work" feedback
