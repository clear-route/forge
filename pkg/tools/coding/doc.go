// Package coding provides reusable tools for file operations and code manipulation.
//
// This package implements the core coding tools used by the Forge TUI agent:
//   - ReadFileTool: Read file contents with optional line ranges
//   - WriteFileTool: Create or overwrite files with validation
//   - ListFilesTool: List directory contents with optional recursion
//   - SearchFilesTool: Search files using regex patterns
//   - ApplyDiffTool: Apply targeted edits using search/replace
//   - ExecuteCommandTool: Execute terminal commands with approval
//
// All tools enforce workspace-level security through the WorkspaceGuard,
// preventing access to files outside the designated workspace directory.
//
// Tools are designed to be reusable across different executors (TUI, CLI, API)
// and integrate with the agent's event system for streaming updates and
// approval workflows.
package coding
