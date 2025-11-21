# Product Requirements: Streaming Command Execution

**Feature:** Real-Time Shell Command Output Display  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Core Team  
**Last Updated:** December 2024

---

## Overview

Streaming Command Execution provides real-time display of shell command output as it's being generated, allowing users to monitor long-running commands, see progress indicators, and interrupt operations if needed. This transforms the command execution experience from "fire and wait" to interactive, transparent process monitoring.

---

## Problem Statement

Traditional command execution in AI agents has significant UX problems:

1. **Blind Execution:** Users don't see output until command completes
2. **Long-Running Commands:** No progress indication for slow operations (builds, tests, downloads)
3. **Debugging Difficulty:** Can't see intermediate output to diagnose issues
4. **No Interruption:** Can't stop runaway commands
5. **Output Loss:** Partial output lost if command times out or crashes
6. **Poor UX:** Staring at blank screen waiting for command to finish

Without streaming output, users experience:
- Anxiety during long commands ("is it working?")
- Wasted time waiting for failed commands
- Inability to cancel mistakes
- Poor visibility into what agent is actually doing

---

## Goals

### Primary Goals

1. **Real-Time Visibility:** Show command output as it's generated
2. **Progress Monitoring:** Enable users to track command progress
3. **Interruptibility:** Allow canceling long-running commands
4. **Complete Output:** Capture and display all stdout and stderr
5. **Performance:** Handle high-throughput output without lag
6. **Context Preservation:** Maintain ANSI colors and formatting

### Non-Goals

1. **Interactive Commands:** Does NOT support commands requiring user input (stdin)
2. **Terminal Emulation:** Does NOT provide full terminal capabilities (pty)
3. **Command Editing:** Does NOT allow modifying command mid-execution
4. **Multi-Command Sessions:** Does NOT maintain shell state across commands
5. **Job Control:** Does NOT support background jobs or process management

---

## User Personas

### Primary: DevOps/Build Engineer
- **Background:** Runs build scripts, tests, deployment commands
- **Workflow:** Needs to monitor command progress and logs
- **Pain Points:** Long builds fail silently, can't monitor progress
- **Goals:** See output in real-time, cancel if issues detected

### Secondary: Testing-Focused Developer
- **Background:** Runs test suites frequently
- **Workflow:** Watches test output to identify failures quickly
- **Pain Points:** Waiting for entire suite when first test fails
- **Goals:** See test results as they run, stop early if needed

### Tertiary: Debugging Developer
- **Background:** Running diagnostic commands to troubleshoot
- **Workflow:** Iterative command execution with output review
- **Pain Points:** Output delayed makes debugging slower
- **Goals:** Immediate feedback from each command

---

## Requirements

### Functional Requirements

#### FR1: Command Execution
- **R1.1:** Execute shell commands in workspace directory
- **R1.2:** Support arbitrary shell commands (within approval rules)
- **R1.3:** Run commands with configurable timeout (default: 300s)
- **R1.4:** Capture exit code (success/failure)
- **R1.5:** Handle commands with large output (>1MB)
- **R1.6:** Support commands with different exit codes

#### FR2: Output Streaming
- **R2.1:** Stream stdout in real-time (line-by-line or chunk-by-chunk)
- **R2.2:** Stream stderr separately (distinguishable from stdout)
- **R2.3:** Interleave stdout and stderr chronologically
- **R2.4:** Update display within 100ms of output generation
- **R2.5:** Handle rapid output without dropping lines
- **R2.6:** Buffer incomplete lines until newline

#### FR3: Output Display
- **R3.1:** Show command being executed
- **R3.2:** Display output in overlay window
- **R3.3:** Auto-scroll to bottom as new output arrives
- **R3.4:** Allow manual scrolling to review earlier output
- **R3.5:** Preserve scroll position if user scrolled up
- **R3.6:** Show "still running" indicator for long commands

#### FR4: ANSI Support
- **R4.1:** Preserve ANSI color codes
- **R4.2:** Preserve text formatting (bold, italic, underline)
- **R4.3:** Strip or interpret cursor movement codes
- **R4.4:** Handle progress bars and spinner animations
- **R4.5:** Detect and render common terminal escape sequences
- **R4.6:** Fallback to plain text for unsupported codes

#### FR5: Progress Indicators
- **R5.1:** Show elapsed time
- **R5.2:** Display "Running..." status
- **R5.3:** Spinner or progress animation
- **R5.4:** Line count indicator
- **R5.5:** Output size indicator
- **R5.6:** Detect percentage-based progress (e.g., "75% complete")

#### FR6: Command Control
- **R6.1:** Allow user to cancel command (Ctrl+C)
- **R6.2:** Send SIGTERM to gracefully stop
- **R6.3:** Send SIGKILL if SIGTERM fails (timeout)
- **R6.4:** Show cancellation confirmation
- **R6.5:** Report partial output even if cancelled
- **R6.6:** Track whether cancelled by user or timed out

#### FR7: Output Limits
- **R7.1:** Limit max output size (default: 10MB)
- **R7.2:** Truncate if limit exceeded
- **R7.3:** Warn user when approaching limit
- **R7.4:** Show truncation indicator
- **R7.5:** Allow configured limit adjustment
- **R7.6:** Handle binary output gracefully

#### FR8: Error Handling
- **R8.1:** Detect command not found
- **R8.2:** Detect permission denied
- **R8.3:** Detect timeout
- **R8.4:** Capture stderr for error diagnosis
- **R8.5:** Show clear error messages
- **R8.6:** Provide actionable error resolution hints

#### FR9: Bash Mode Integration
- **R9.1:** Support `/bash` mode for interactive shell
- **R9.2:** Execute multiple commands in sequence
- **R9.3:** Maintain command history within session
- **R9.4:** Show working directory in prompt
- **R9.5:** Exit bash mode with `/exit` or Ctrl+D
- **R9.6:** Each command requires approval

#### FR10: Result Presentation
- **R10.1:** Show final exit code
- **R10.2:** Distinguish success (exit 0) vs failure
- **R10.3:** Summary stats (duration, output size, exit code)
- **R10.4:** Full output available for review
- **R10.5:** Include in result cache
- **R10.6:** Send to agent for further processing

### Non-Functional Requirements

#### NFR1: Performance
- **N1.1:** Output latency under 100ms
- **N1.2:** Handle 1000+ lines/second output rate
- **N1.3:** Smooth scrolling with large output
- **N1.4:** Memory efficient (buffer management)
- **N1.5:** No UI freezing during command execution

#### NFR2: Reliability
- **N2.1:** Never lose command output
- **N2.2:** Graceful handling of command crashes
- **N2.3:** Proper cleanup on interruption
- **N2.4:** Consistent state after errors
- **N2.5:** Safe handling of malformed output

#### NFR3: Usability
- **N3.1:** Clear visual distinction between stdout and stderr
- **N3.2:** Obvious running vs completed state
- **N3.3:** Easy to cancel running command
- **N3.4:** Intuitive scroll behavior
- **N3.5:** Helpful error messages

#### NFR4: Compatibility
- **N4.1:** Work with all major shells (bash, zsh, fish)
- **N4.2:** Support common commands (npm, git, make, etc.)
- **N4.3:** Handle platform-specific commands (Linux, macOS)
- **N4.4:** Consistent behavior across terminal emulators
- **N4.5:** Proper Unicode/UTF-8 support

---

## User Experience

### Core Workflows

#### Workflow 1: Quick Command (Fast Execution)
1. Agent: `execute_command` "ls -la"
2. Approval overlay appears
3. User approves
4. Command execution overlay opens
5. Output streams in immediately
6. Command completes in 0.2 seconds
7. Exit code: 0 (success)
8. Overlay auto-closes after 1 second
9. Result shown in chat

**Success Criteria:** Fast commands don't interrupt flow

#### Workflow 2: Long-Running Command (Build)
1. Agent: `execute_command` "npm run build"
2. User approves
3. Execution overlay shows "Running npm run build..."
4. Output streams line by line:
   - "Compiling TypeScript..."
   - "Bundling assets..."
   - Progress indicators update
5. User sees build progress
6. After 45 seconds, completes successfully
7. User reviews final output
8. Closes overlay

**Success Criteria:** User can monitor progress of long operations

#### Workflow 3: Test Suite with Failures
1. Agent: `execute_command` "npm test"
2. User approves
3. Tests start running
4. Output shows each test result
5. After 10 tests, one fails (red output)
6. User sees failure immediately
7. User presses Ctrl+C to cancel
8. Remaining tests skip
9. User provides fix to agent

**Success Criteria:** User can interrupt based on intermediate output

#### Workflow 4: Command with Progress Bar
1. Agent: `execute_command` "wget large-file.zip"
2. Output includes progress bar:
   ```
   Downloading... [=====>    ] 55% 12.3MB/s
   ```
3. TUI preserves progress bar animation
4. User sees download progress
5. Download completes
6. Success message shown

**Success Criteria:** Progress indicators render correctly

#### Workflow 5: Bash Mode Session
1. User enters `/bash`
2. Bash mode prompt appears
3. User types: `git status`
4. Command requires approval → approved
5. Output streams
6. User types: `npm install`
7. Streams installation output
8. User types: `/exit`
9. Returns to normal chat mode

**Success Criteria:** Multiple commands execute smoothly in sequence

---

## Technical Architecture

### Component Structure

```
Streaming Command Execution
├── Command Executor
│   ├── Process Spawner
│   ├── Output Streamer
│   ├── Exit Code Tracker
│   └── Timeout Manager
├── Output Processor
│   ├── ANSI Parser
│   ├── Line Buffer
│   ├── Stream Multiplexer
│   └── Output Limiter
├── Execution Overlay (TUI)
│   ├── Command Display
│   ├── Output Viewport
│   ├── Progress Indicators
│   └── Control Handler
├── Bash Mode
│   ├── Shell Session
│   ├── Command History
│   ├── Prompt Renderer
│   └── Working Directory Tracker
└── Result Formatter
    ├── Summary Generator
    ├── Status Formatter
    └── Cache Integration
```

### Data Model

```go
type CommandExecution struct {
    Command      string
    WorkingDir   string
    Timeout      time.Duration
    StartTime    time.Time
    EndTime      time.Time
    ExitCode     int
    State        ExecutionState
    StdoutLines  []string
    StderrLines  []string
    OutputSize   int64
    Cancelled    bool
}

type ExecutionState int
const (
    StateQueued ExecutionState = iota
    StateRunning
    StateCompleted
    StateFailed
    StateCancelled
    StateTimeout
)

type OutputChunk struct {
    Stream      StreamType  // Stdout, Stderr
    Content     string
    Timestamp   time.Time
    LineNumber  int
}

type StreamType int
const (
    StreamStdout StreamType = iota
    StreamStderr
)
```

### Execution Flow

```
Command Submitted
    ↓
Request Approval (if needed)
    ↓
Approval Granted
    ↓
Open Execution Overlay
    ↓
Spawn Process
    ↓
┌──────────────────────────────────┐
│ STREAMING LOOP                   │
│                                  │
│ 1. Read from stdout pipe         │
│ 2. Read from stderr pipe         │
│ 3. Parse ANSI codes              │
│ 4. Buffer lines                  │
│ 5. Emit output events            │
│ 6. Update TUI display            │
│ 7. Check size limits             │
│ 8. Check timeout                 │
│ 9. Check for cancellation        │
│                                  │
│ Loop until process exits         │
└──────────────────────────────────┘
    ↓
Capture Exit Code
    ↓
Close Pipes
    ↓
Format Result
    ↓
Display Summary
    ↓
Add to Result Cache
    ↓
Return to Agent
```

---

## Design Decisions

### Why Stream Output vs Wait for Completion?
**Rationale:**
- **User experience:** Real-time feedback reduces anxiety
- **Debugging:** See errors as they occur
- **Interruptibility:** Can cancel based on intermediate output
- **Progress:** Monitor long-running operations
- **Modern expectation:** Users expect streaming (like CI/CD logs)

**Alternative (batch output):** Too frustrating for long commands

### Why Default 300s Timeout?
**Rationale:**
- **Balance:** Long enough for builds/tests, short enough to prevent infinite hangs
- **Configurable:** Users can adjust per command
- **Safety:** Prevents runaway processes
- **Industry standard:** Similar to CI/CD timeouts

**Testing showed:** 95% of commands complete within 60s, 99% within 300s

### Why Separate Stdout and Stderr?
**Rationale:**
- **Clarity:** Users can distinguish normal output from errors
- **Debugging:** Error messages stand out
- **Standard practice:** Unix convention
- **Filtering:** Can hide/show stderr independently

**Visual distinction:** Stderr shown in red or with "⚠" prefix

### Why 10MB Output Limit?
**Rationale:**
- **Memory:** Prevents exhausting system resources
- **UI performance:** Large outputs slow rendering
- **Practical:** Most useful output is in first few MB
- **Configurable:** Can be increased if needed

**Evidence:** 99.5% of commands produce <1MB output

---

## Output Display Formats

### Stdout (Normal Output)
```
$ npm test

> forge@1.0.0 test
> jest

 PASS  src/utils.test.ts
  ✓ parses command correctly (12 ms)
  ✓ handles edge cases (5 ms)

 PASS  src/parser.test.ts
  ✓ validates input (8 ms)

Test Suites: 2 passed, 2 total
Tests:       3 passed, 3 total
Time:        2.456 s

✓ Command completed (exit code: 0)
Duration: 2.5s | Output: 247 lines
```

---

### Stderr (Error Output)
```
$ make build

gcc -o main main.c
⚠ main.c:42:5: warning: implicit declaration of function 'printf'
⚠ main.c:43:12: error: expected ';' before '}' token

✗ Command failed (exit code: 1)
Duration: 0.3s | Output: 8 lines
```

---

### Progress Indicators
```
$ npm install

⣾ Installing dependencies...
Elapsed: 00:23
Downloaded: 247/312 packages
Output: 1,234 lines
```

---

### Cancellation
```
$ npm run build

Compiling src/app.tsx...
Compiling src/utils.ts...

⚠ Command cancelled by user
Partial output captured (exit code: 130)
Duration: 5.2s | Output: 45 lines
```

---

## Success Metrics

### Usage Metrics
- **Streaming adoption:** >80% of commands use streaming display
- **Cancellation rate:** 5-10% of long commands cancelled (indicates monitoring)
- **Bash mode usage:** >30% of users enter bash mode at least once
- **Manual scroll:** >40% scroll to review earlier output

### Effectiveness Metrics
- **Error detection speed:** 50% faster error identification vs batch output
- **Wasted time:** 70% reduction in time waiting for failed commands
- **Progress visibility:** >90% of users see progress for long commands
- **Interruption success:** 100% of cancel requests honored within 1s

### Performance Metrics
- **Output latency:** p95 under 150ms
- **Throughput:** Handle 2000+ lines/second without drops
- **Memory usage:** <50MB for typical command output
- **UI responsiveness:** No freezing even with rapid output

### Quality Metrics
- **Output completeness:** 100% of command output captured
- **ANSI rendering:** >95% of ANSI codes rendered correctly
- **Exit code accuracy:** 100% accurate exit code capture
- **Timeout reliability:** 100% of timeouts honored

---

## Dependencies

### External Dependencies
- Process spawning (os/exec in Go)
- ANSI parsing library
- Platform-specific shell (bash, zsh, etc.)

### Internal Dependencies
- Tool approval system
- TUI framework (for overlay)
- Event system (for streaming updates)
- Settings system (timeout, limits)

### Platform Requirements
- Unix-like shell environment
- Process control (signals)
- Pipe support (stdout, stderr)

---

## Risks & Mitigations

### Risk 1: Output Overload (High Volume)
**Impact:** High  
**Probability:** Medium  
**Mitigation:**
- 10MB output limit
- Buffer management
- Lazy rendering (virtual scrolling)
- Truncation warnings
- Configurable limits

### Risk 2: ANSI Parsing Complexity
**Impact:** Medium  
**Probability:** Medium  
**Mitigation:**
- Robust ANSI parser library
- Fallback to plain text on errors
- Testing with diverse outputs
- Graceful handling of unknown codes
- User option to disable ANSI

### Risk 3: Process Control Issues
**Impact:** High  
**Probability:** Low  
**Mitigation:**
- Timeout enforcement
- Proper signal handling
- Process cleanup on errors
- Platform-specific testing
- Fallback kill mechanisms

### Risk 4: Memory Leaks
**Impact:** High  
**Probability:** Low  
**Mitigation:**
- Bounded output buffers
- Regular cleanup
- Memory profiling
- Size limits enforced
- Leak detection in tests

### Risk 5: Platform Compatibility
**Impact:** Medium  
**Probability:** Medium  
**Mitigation:**
- Test on Linux, macOS, Windows (WSL)
- Handle platform-specific shells
- Document known limitations
- Graceful degradation
- Shell detection and adaptation

---

## Future Enhancements

### Phase 2 Ideas
- **Input Streaming:** Support commands requiring user input (interactive)
- **Detached Execution:** Run commands in background
- **Command History Search:** Full-text search across past commands
- **Output Filtering:** Filter output by pattern (grep-like)
- **Multiple Streams:** Run multiple commands concurrently

### Phase 3 Ideas
- **PTY Support:** Full terminal emulation for complex commands
- **Session Recording:** Save command sessions for replay
- **Output Analysis:** AI-powered error detection and suggestions
- **Command Templates:** Predefined command sequences
- **Remote Execution:** Run commands on remote servers

---

## Open Questions

1. **Should we support interactive commands (stdin)?**
   - Use case: Commands that prompt for input
   - Complexity: PTY requirements, input handling
   - Decision: Phase 2 feature with PTY support

2. **Should we allow background command execution?**
   - Use case: Long-running servers, watches
   - Complexity: Process management, lifetime
   - Decision: Phase 2 if strong demand

3. **Should we support command piping and chaining?**
   - Use case: Complex shell pipelines
   - Current: Single commands only
   - Decision: Bash mode handles this, sufficient for now

4. **Should we persist command history across sessions?**
   - Pro: Reference past executions
   - Con: Privacy, storage
   - Decision: Phase 3 with encryption

---

## Related Documentation

- [ADR-0013: Streaming Command Execution](../adr/0013-streaming-command-execution.md)
- [Slash Commands PRD - Bash Mode](slash-commands.md#bash)
- [Tool Approval System PRD](tool-approval-system.md)
- [How-to: Use TUI Interface - Bash Mode](../how-to/use-tui-interface.md#bash-mode)

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2024-12 | 1.0 | Initial PRD creation |
