# Feature Idea: Process Monitoring & Management

**Status:** Draft  
**Priority:** High Impact, Near-Term  
**Last Updated:** November 2025

---

## Overview

Built-in process monitoring and management that lets Forge track, control, and interact with running processes—dev servers, background jobs, databases, build processes, and test suites. Provides a unified interface for managing the entire development environment without juggling terminal windows.

---

## Problem Statement

Developers constantly juggle multiple running processes:
- Dev server on port 3000
- API backend on port 8080
- Database server
- Redis cache
- File watcher/bundler
- Test runner in watch mode
- Background job worker
- Proxy/tunnel services

Current situation is chaotic:
- 5-10 terminal tabs/windows open
- Lost track of what's running where
- Processes orphaned when terminal closes
- No unified view of system state
- Hard to see logs across processes
- Kill processes manually by finding PIDs
- Restart requires finding right terminal
- Port conflicts from forgotten processes

This leads to:
- "It works on my machine" issues
- Forgotten background processes eating resources
- Port conflicts requiring `lsof -i :3000 | grep LISTEN | awk '{print $2}' | xargs kill`
- Lost productivity switching between terminals
- Unclear system state
- Difficulty debugging multi-service issues

---

## Key Capabilities

### Process Lifecycle Management

**Start & Stop:**
- Launch processes from Forge
- Stop individual processes or groups
- Restart with preserved state
- Graceful shutdown with configurable timeout
- Force kill for hung processes
- Auto-restart on crash

**Process Groups:**
- Group related processes (frontend + backend + db)
- Start/stop entire groups at once
- Dependency management (start DB before API)
- Named configurations ("dev", "test", "prod-like")
- Save/load process configurations

**Smart Process Detection:**
- Auto-detect common dev processes
- Suggest processes to start
- Detect port conflicts
- Find orphaned processes
- Identify resource hogs

### Real-Time Monitoring

**Live Status Dashboard:**
- Process list with status (running, stopped, crashed)
- CPU and memory usage per process
- Port bindings
- Uptime
- Restart count
- Health checks

**Log Aggregation:**
- Combined log view across all processes
- Color-coded by process
- Filtering by process, level, search term
- Tail live logs
- Log history with timestamps
- Export logs

**Resource Tracking:**
- CPU usage over time
- Memory consumption
- Network activity
- File descriptor usage
- Thread counts
- Disk I/O

### Interactive Control

**Log Streaming:**
- Live tail of any process
- Search/filter logs
- Follow mode (auto-scroll)
- Pause/resume streaming
- Save log snapshots

**Process Communication:**
- Send signals (SIGTERM, SIGKILL, SIGUSR1, etc.)
- Send input to stdin
- Interactive REPL access
- Attach to running process
- Debug mode toggle

**Port Management:**
- List all listening ports
- Identify process by port
- Kill process by port
- Detect port conflicts
- Suggest alternative ports

### Configuration Profiles

**Named Environments:**
```yaml
# .forge/processes.yml
dev:
  - name: frontend
    command: npm run dev
    cwd: ./frontend
    port: 3000
    env:
      NODE_ENV: development
    health_check: http://localhost:3000/health
    
  - name: backend
    command: go run cmd/server/main.go
    cwd: ./backend
    port: 8080
    depends_on: [postgres, redis]
    
  - name: postgres
    command: docker run -p 5432:5432 postgres:15
    health_check: pg_isready
    
  - name: redis
    command: redis-server
    port: 6379

test:
  - name: test-db
    command: docker run -p 5433:5432 postgres:15
  - name: test-runner
    command: npm test -- --watch
    depends_on: [test-db]
```

**Auto-Save State:**
- Remember running processes
- Restore on Forge restart
- Per-workspace configurations
- Export/import configs

---

## Example Use Cases

### 1. Starting Development Environment

```
User: /dev start

Agent: Starting development environment...

[Reading .forge/processes.yml]

Starting process group: dev
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ postgres (PID 12345)
  Port: 5432
  Status: Healthy
  
✓ redis (PID 12346)
  Port: 6379
  Status: Ready
  
⠋ backend (PID 12347)
  Port: 8080
  Status: Starting...
  Waiting for dependencies...
  
⠋ frontend (PID 12348)
  Port: 3000
  Status: Building...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[2 minutes later]

✓ All processes started successfully!

Services:
┌───────────┬──────┬──────┬────────┬──────────┐
│ Name      │ PID  │ Port │ CPU    │ Memory   │
├───────────┼──────┼──────┼────────┼──────────┤
│ postgres  │ 12345│ 5432 │ 0.1%   │ 45 MB    │
│ redis     │ 12346│ 6379 │ 0.0%   │ 12 MB    │
│ backend   │ 12347│ 8080 │ 2.3%   │ 156 MB   │
│ frontend  │ 12348│ 3000 │ 15.2%  │ 423 MB   │
└───────────┴──────┴──────┴────────┴──────────┘

Your app is ready at: http://localhost:3000

Commands:
  /dev stop          - Stop all processes
  /dev restart       - Restart all processes
  /dev logs          - View combined logs
  /dev logs backend  - View backend logs only
  /dev status        - Show status
```

### 2. Monitoring Processes

```
User: /dev status

Agent: 

Development Environment Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Active Processes: 4
Total Memory: 636 MB
Total CPU: 17.6%
Uptime: 2h 34m

┌───────────┬─────────┬──────┬────────┬──────────┬──────────┐
│ Name      │ Status  │ Port │ CPU    │ Memory   │ Uptime   │
├───────────┼─────────┼──────┼────────┼──────────┼──────────┤
│ postgres  │ Running │ 5432 │ 0.1%   │ 45 MB    │ 2h 34m   │
│ redis     │ Running │ 6379 │ 0.0%   │ 12 MB    │ 2h 34m   │
│ backend   │ Running │ 8080 │ 2.5%   │ 158 MB   │ 2h 34m   │
│ frontend  │ Running │ 3000 │ 14.8%  │ 421 MB   │ 2h 33m   │
└───────────┴─────────┴──────┴────────┴──────────┴──────────┘

Health Checks:
✓ postgres - Accepting connections
✓ redis - PONG received
✓ backend - HTTP 200 from /health
✓ frontend - Webpack compiled successfully

Recent Events:
  [14:23:45] backend - Hot reload: user.go modified
  [14:22:10] frontend - Compiled successfully
  [14:15:33] backend - Database query: 145ms (slow)

Warnings:
⚠️ backend - 3 slow queries in last hour
⚠️ frontend - Bundle size increased by 150KB

View detailed logs? (y/n)
```

### 3. Debugging with Logs

```
User: /dev logs

Agent: Streaming combined logs from all processes...

[Live log view with color coding by process]

[postgres] 2025-11-21 14:30:15 | Connection from 127.0.0.1:54321
[backend]  2025-11-21 14:30:15 | GET /api/users - 200 (45ms)
[backend]  2025-11-21 14:30:16 | SELECT * FROM users WHERE id=123
[postgres] 2025-11-21 14:30:16 | Query execution: 12ms
[frontend] 2025-11-21 14:30:17 | Rendering UserProfile component
[backend]  2025-11-21 14:30:18 | ERROR: Connection timeout
[redis]    2025-11-21 14:30:18 | Client disconnected
[backend]  2025-11-21 14:30:18 | Retrying Redis connection...
[redis]    2025-11-21 14:30:19 | New client connected
[backend]  2025-11-21 14:30:19 | Redis connection restored

[Commands available]
  /filter backend     - Show only backend logs
  /search error       - Search for 'error'
  /pause              - Pause streaming
  /save logs.txt      - Save to file
  /clear              - Clear screen
  q                   - Quit log view

User: /filter backend

Agent: Filtering logs for: backend

[backend]  2025-11-21 14:30:15 | GET /api/users - 200 (45ms)
[backend]  2025-11-21 14:30:16 | SELECT * FROM users WHERE id=123
[backend]  2025-11-21 14:30:18 | ERROR: Connection timeout
[backend]  2025-11-21 14:30:18 | Retrying Redis connection...
[backend]  2025-11-21 14:30:19 | Redis connection restored
[backend]  2025-11-21 14:30:25 | POST /api/orders - 201 (234ms)
[backend]  2025-11-21 14:30:26 | Sending confirmation email

User: /search error

Agent: Found 3 error entries in last hour:

[backend]  2025-11-21 14:30:18 | ERROR: Connection timeout
[backend]  2025-11-21 13:45:12 | ERROR: Rate limit exceeded
[backend]  2025-11-21 13:22:08 | ERROR: Invalid JSON in request

Would you like to:
1. See full context for any error
2. Restart backend to clear errors
3. View error trends
```

### 4. Port Conflict Resolution

```
User: /dev start

Agent: Starting development environment...

⚠️ Port Conflict Detected!

Port 3000 is already in use:
┌──────────────┬──────┬──────────────────┬──────────┐
│ Process      │ PID  │ Command          │ Started  │
├──────────────┼──────┼──────────────────┼──────────┤
│ node         │ 8432 │ npm run dev      │ 2h ago   │
└──────────────┴──────┴──────────────────┴──────────┘

This appears to be an old development server.

Options:
1. Kill old process and start fresh
2. Use different port (suggest: 3001)
3. Attach to existing process
4. Cancel

> 1

Killing process 8432...
✓ Process terminated

Starting frontend on port 3000...
✓ Frontend started (PID 15234)
```

### 5. Crash Recovery

```
[Automatic detection]

Agent: 

⚠️ Process Crashed!

backend (PID 12347) exited unexpectedly
Exit code: 1
Last output:
  panic: runtime error: invalid memory address
  goroutine 1 [running]:
  main.handleRequest(...)
      /app/handlers/user.go:45

This is crash #3 in the last 10 minutes.

Actions:
1. Auto-restart (default in 5s)
2. View full crash log
3. Start in debug mode
4. Keep stopped for investigation

Countdown: 3s... (press any key to cancel)

User: 2

Agent: Full crash log:

[Shows complete stack trace and context]

Crash appears to be caused by nil pointer in user.go:45

Recent changes to this file:
  - Modified 5 minutes ago
  - Added GetUserOrders() function

Suggested actions:
✓ Review recent changes to user.go
✓ Add nil checks in GetUserOrders()
✓ Run tests before restarting

Would you like me to:
1. Restart with debug logging enabled
2. Revert recent changes to user.go
3. Open file at crash location
```

### 6. Resource Monitoring

```
User: /dev monitor

Agent: Real-time resource monitoring...

System Overview
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CPU:    28.5% [████████░░░░░░░░░░░░] 
Memory: 4.2GB / 16GB [█████░░░░░░░░░░░]

Process Breakdown:

frontend (PID 12348)
CPU:    [▁▂▃▅▇█▇▅▃▂▁] 15.2% avg (spike to 45% during compile)
Memory: [▁▁▂▃▃▃▃▃▃▃▃] 421 MB (growing slowly)
Threads: 12
Ports: 3000, 35729 (websocket)

backend (PID 12347)  
CPU:    [▁▁▁▂▁▁▁▁▁▁▁] 2.5% avg (steady)
Memory: [▁▂▂▂▂▂▂▂▂▂▂] 158 MB (stable)
Threads: 8
Ports: 8080
DB Connections: 5/10 pool

postgres (PID 12345)
CPU:    [▁▁▁▁▁▂▁▁▁▁▁] 0.1% avg (occasional spikes)
Memory: [▁▁▁▁▁▁▁▁▁▁▁] 45 MB (stable)
Connections: 3
Cache hit rate: 94.2%

⚠️ Alerts:
  - frontend memory increased 50MB in last 10 min (possible leak?)
  - backend had 3 requests >1s in last hour

Refreshing every 2s... (press q to quit)
```

---

## Technical Approach

### Process Management Core

**Process Controller:**
```go
type ProcessManager struct {
    processes map[string]*ManagedProcess
    groups    map[string]*ProcessGroup
    mu        sync.RWMutex
}

type ManagedProcess struct {
    Name      string
    PID       int
    Cmd       *exec.Cmd
    StartTime time.Time
    Restarts  int
    
    // I/O
    Stdout    io.ReadCloser
    Stderr    io.ReadCloser
    Stdin     io.WriteCloser
    
    // Monitoring
    CPUUsage  float64
    MemUsage  uint64
    
    // Control
    StopChan  chan struct{}
    DoneChan  chan error
}

func (pm *ProcessManager) Start(name, command string, opts ProcessOpts) error {
    cmd := exec.Command("sh", "-c", command)
    cmd.Dir = opts.WorkingDir
    cmd.Env = append(os.Environ(), opts.Env...)
    
    // Capture output
    stdout, _ := cmd.StdoutPipe()
    stderr, _ := cmd.StderrPipe()
    stdin, _ := cmd.StdinPipe()
    
    // Start process
    if err := cmd.Start(); err != nil {
        return err
    }
    
    proc := &ManagedProcess{
        Name:      name,
        PID:       cmd.Process.Pid,
        Cmd:       cmd,
        Stdout:    stdout,
        Stderr:    stderr,
        Stdin:     stdin,
        StartTime: time.Now(),
    }
    
    pm.mu.Lock()
    pm.processes[name] = proc
    pm.mu.Unlock()
    
    // Monitor in background
    go pm.monitorProcess(proc)
    
    return nil
}
```

**Resource Monitoring:**
```go
import "github.com/shirou/gopsutil/v3/process"

func (pm *ProcessManager) monitorProcess(proc *ManagedProcess) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    p, err := process.NewProcess(int32(proc.PID))
    if err != nil {
        return
    }
    
    for {
        select {
        case <-ticker.C:
            // CPU usage
            cpu, _ := p.CPUPercent()
            proc.CPUUsage = cpu
            
            // Memory usage
            mem, _ := p.MemoryInfo()
            proc.MemUsage = mem.RSS
            
            // Check health
            if proc.HealthCheck != nil {
                proc.HealthCheck()
            }
            
        case <-proc.StopChan:
            return
        }
    }
}
```

### Log Aggregation

**Log Multiplexer:**
```go
type LogAggregator struct {
    streams map[string]*LogStream
    output  chan LogEntry
    mu      sync.RWMutex
}

type LogEntry struct {
    Timestamp time.Time
    Process   string
    Level     LogLevel
    Message   string
}

func (la *LogAggregator) AddStream(name string, reader io.Reader) {
    stream := &LogStream{
        Name:   name,
        Reader: bufio.NewReader(reader),
    }
    
    la.mu.Lock()
    la.streams[name] = stream
    la.mu.Unlock()
    
    go la.readStream(stream)
}

func (la *LogAggregator) readStream(stream *LogStream) {
    scanner := bufio.NewScanner(stream.Reader)
    for scanner.Scan() {
        line := scanner.Text()
        
        entry := LogEntry{
            Timestamp: time.Now(),
            Process:   stream.Name,
            Message:   line,
            Level:     detectLogLevel(line),
        }
        
        la.output <- entry
    }
}
```

### Port Detection

**Port Scanner:**
```go
import "net"

func FindProcessByPort(port int) (*ProcessInfo, error) {
    // Use lsof on Unix, netstat on Windows
    cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-t")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    pid, _ := strconv.Atoi(strings.TrimSpace(string(output)))
    
    proc, _ := process.NewProcess(int32(pid))
    name, _ := proc.Name()
    cmdline, _ := proc.Cmdline()
    
    return &ProcessInfo{
        PID:     pid,
        Name:    name,
        Command: cmdline,
    }, nil
}

func IsPortAvailable(port int) bool {
    ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        return false
    }
    ln.Close()
    return true
}
```

### Configuration Management

**Process Config Parser:**
```go
type ProcessConfig struct {
    Name         string            `yaml:"name"`
    Command      string            `yaml:"command"`
    WorkingDir   string            `yaml:"cwd"`
    Env          map[string]string `yaml:"env"`
    Port         int               `yaml:"port"`
    DependsOn    []string          `yaml:"depends_on"`
    HealthCheck  string            `yaml:"health_check"`
    AutoRestart  bool              `yaml:"auto_restart"`
}

func LoadConfig(path string) (map[string][]ProcessConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var config map[string][]ProcessConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return config, nil
}
```

---

## Value Propositions

### For All Developers
- Single pane of glass for all processes
- No more lost terminal windows
- Easy start/stop entire environments
- Unified log viewing
- Resource monitoring

### For Backend Developers
- Manage microservices easily
- Database/cache control
- Background job monitoring
- Multi-service debugging

### For Full-Stack Developers
- Frontend + backend + DB together
- Coordinated restarts
- Dependency management
- Clear system state

---

## Implementation Phases

### Phase 1: Core Process Control (2 weeks)
- Start/stop processes
- Basic monitoring (CPU, memory)
- Log capture
- TUI interface

### Phase 2: Advanced Monitoring (2 weeks)
- Resource graphs
- Health checks
- Crash detection
- Auto-restart

### Phase 3: Configuration & Groups (1 week)
- YAML config files
- Process groups
- Dependency management
- Named environments

### Phase 4: Advanced Features (2 weeks)
- Port management
- Interactive process control
- Log search/filtering
- Process communication

---

## Success Metrics

**Adoption:**
- 70%+ use process management daily
- 60%+ save process configurations
- 80%+ use log aggregation
- 50%+ monitor resources

**Impact:**
- 50% reduction in "forgot to start X" issues
- 60% faster environment setup
- 80% reduction in orphaned processes
- 40% faster debugging multi-service issues

**Satisfaction:**
- 4.7+ rating
- "No more terminal juggling" feedback
- "Saved my sanity" comments
