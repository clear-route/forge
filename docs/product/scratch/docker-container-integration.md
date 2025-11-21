# Feature Idea: Docker & Container Integration

**Status:** Draft  
**Priority:** High Impact, Near-Term  
**Last Updated:** November 2025

---

## Overview

Native Docker and container workflow integration that makes Forge a first-class container development tool. Build, run, debug, and manage containers directly from Forge without leaving the development environment. Brings container development into the same seamless workflow as code editing.

---

## Problem Statement

Container development requires constant context switching:
- Jump to terminal to run `docker build`
- Check another terminal for container logs
- Use Docker Desktop GUI to manage containers
- Manually inspect container filesystem
- SSH into containers for debugging
- Copy files in/out manually
- Restart containers repeatedly during development
- No integration between code changes and containers

Current workflow is fragmented:
- Write code in editor
- Switch to terminal for docker commands
- Switch to browser for Docker Desktop
- Switch to another terminal for logs
- Repeat 20+ times per day

This leads to:
- Wasted time context switching
- Forgotten running containers
- Difficulty debugging container issues
- Slow iteration cycles
- Port conflicts from orphaned containers
- Manual cleanup needed constantly
- Cannot easily inspect container state
- Multi-container orchestration is painful

---

## Key Capabilities

### Container Lifecycle Management

**Build Operations:**
- Build images from Dockerfiles
- Multi-stage build optimization
- Build cache management
- Layer inspection
- Build argument injection
- Tag management
- Push to registries
- Buildkit support

**Run & Control:**
- Start containers with smart defaults
- Stop/restart containers
- Pause/unpause execution
- Kill hung containers
- Remove containers and volumes
- Attach to running containers
- Execute commands in containers

**Status Monitoring:**
- List all containers (running, stopped, exited)
- Show container resource usage
- View port mappings
- Check health status
- Monitor restart counts
- Track container events

### Live Development Workflow

**Hot Reload Integration:**
- Detect code changes
- Auto-rebuild images
- Smart incremental rebuilds
- Volume mounting for development
- Preserve container state
- Fast iteration cycles

**Interactive Debugging:**
- Attach debugger to containerized apps
- Set breakpoints in container code
- Inspect container filesystem
- View environment variables
- Check network connectivity
- Analyze running processes

**Log Management:**
- Aggregate logs from multiple containers
- Color-coded output per container
- Filter and search logs
- Follow logs in real-time
- Export log history
- Timestamp correlation

### Container Inspection

**Filesystem Explorer:**
- Browse container filesystem
- View file contents
- Copy files in/out
- Edit files in running containers
- Compare with image filesystem
- Diff between containers

**Network Analysis:**
- View network interfaces
- Check DNS resolution
- Test connectivity
- Inspect iptables rules
- Monitor network traffic
- Debug network issues

**Resource Monitoring:**
- CPU usage tracking
- Memory consumption
- Network I/O
- Disk I/O
- Process list
- Resource limits

### Multi-Container Orchestration

**Docker Compose Integration:**
- Parse docker-compose.yml
- Start entire stacks
- Manage service dependencies
- Scale services
- View service logs
- Health check monitoring

**Service Discovery:**
- List running services
- Show service endpoints
- Test service connectivity
- View service dependencies
- Monitor inter-service traffic

**Orchestration Control:**
- Start/stop entire stack
- Restart individual services
- Scale services up/down
- Update service configs
- Rolling updates
- Rollback support

### Registry Operations

**Image Management:**
- List local images
- Search Docker Hub
- Pull images
- Push to registries
- Tag images
- Delete unused images
- Prune dangling images

**Registry Integration:**
- Docker Hub
- GitHub Container Registry
- Google Container Registry
- Amazon ECR
- Azure Container Registry
- Private registries

**Security Scanning:**
- Scan images for vulnerabilities
- Check for outdated base images
- Detect exposed secrets
- Verify image signatures
- License compliance

---

## Example Use Cases

### 1. Rapid Development Workflow

```
User: Start my development environment

Agent: Starting containers for local development...

[Reading docker-compose.yml]

Found services: postgres, redis, backend, frontend

Starting in dependency order:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ postgres (postgres:15)
  Port: 5432 → 5432
  Status: Healthy (accepting connections)
  
✓ redis (redis:7-alpine)
  Port: 6379 → 6379
  Status: Ready
  
⠋ backend (building...)
  Building image: myapp-backend:dev
  [1/4] FROM golang:1.21
  [2/4] COPY go.mod go.sum ./
  [3/4] RUN go mod download
  [4/4] COPY . .
  Built in 23s (cache: 89%)
  
✓ backend (myapp-backend:dev)
  Port: 8080 → 8080
  Mounted: ./backend → /app (for hot reload)
  Status: Healthy (/health returned 200)
  
⠋ frontend (building...)
  Building image: myapp-frontend:dev
  [1/3] FROM node:20-alpine
  [2/3] COPY package*.json ./
  [3/3] RUN npm install
  Built in 45s (cache: 95%)
  
✓ frontend (myapp-frontend:dev)
  Port: 3000 → 3000
  Mounted: ./frontend → /app
  Status: Webpack dev server started

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

All services running! 🚀

Services:
┌───────────┬──────────────┬──────┬────────┬──────────┐
│ Service   │ Container ID │ Port │ Status │ Health   │
├───────────┼──────────────┼──────┼────────┼──────────┤
│ postgres  │ a1b2c3d4     │ 5432 │ Up     │ Healthy  │
│ redis     │ e5f6g7h8     │ 6379 │ Up     │ Healthy  │
│ backend   │ i9j0k1l2     │ 8080 │ Up     │ Healthy  │
│ frontend  │ m3n4o5p6     │ 3000 │ Up     │ Healthy  │
└───────────┴──────────────┴──────┴────────┴──────────┘

Your app: http://localhost:3000
API docs: http://localhost:8080/docs

Hot reload enabled:
  ✓ Backend: code changes auto-rebuild
  ✓ Frontend: webpack watching

[User edits backend/handlers/user.go]

Agent: Code change detected in backend...

Rebuilding backend container:
⠋ Syncing files to container...
✓ Files synced (234ms)
⠋ Recompiling Go binary...
✓ Compiled (1.2s)
⠋ Restarting application...
✓ Backend restarted

Backend updated in 2.1 seconds! 🎉

Health check: ✓ /health returned 200

[Showing combined logs]
[postgres]  2025-11-21 15:30:15 | Connection from backend
[backend]   2025-11-21 15:30:15 | Server started on :8080
[backend]   2025-11-21 15:30:16 | Hot reload: user.go modified
[backend]   2025-11-21 15:30:17 | Reloaded successfully
[frontend]  2025-11-21 15:30:18 | Webpack compiled
```

### 2. Container Debugging

```
User: My backend container keeps crashing, help debug it

Agent: Let's investigate the backend container...

[Checking container status]

Container Status:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Name: backend (myapp-backend:dev)
Status: Restarting
Restart Count: 5 in last 10 minutes
Last Exit Code: 1
Last Exit Time: 30 seconds ago

This container is crash-looping!

Recent logs:
[backend] 2025-11-21 15:25:10 | Starting application...
[backend] 2025-11-21 15:25:11 | Connecting to database...
[backend] 2025-11-21 15:25:11 | ERROR: dial tcp 127.0.0.1:5432: connect: connection refused
[backend] 2025-11-21 15:25:11 | Failed to start: database connection failed
[backend] 2025-11-21 15:25:11 | exit status 1

Issue found: Backend trying to connect to 127.0.0.1:5432
But postgres is in a different container!

[Inspecting network configuration]

Network Analysis:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

backend network: myapp_default
postgres network: myapp_default

DNS resolution inside backend:
$ docker exec backend nslookup postgres
Name: postgres
Address: 172.18.0.2

The issue: Your connection string uses "localhost"
Should use: "postgres" (service name)

Current config:
DATABASE_URL=postgresql://localhost:5432/myapp

Should be:
DATABASE_URL=postgresql://postgres:5432/myapp

Would you like me to:
1. Fix the environment variable
2. Show how to debug this yourself
3. Restart with correct config

User: 1

Agent: Fixing environment variable...

[Updating docker-compose.yml]

✓ Updated backend service environment:
  DATABASE_URL: postgresql://postgres:5432/myapp

✓ Restarting backend with new config...
⠋ Stopping backend...
✓ Stopped
⠋ Starting backend...
✓ Started

Monitoring startup...
[backend] 2025-11-21 15:26:30 | Starting application...
[backend] 2025-11-21 15:26:31 | Connecting to database...
[backend] 2025-11-21 15:26:31 | Connected to postgres successfully
[backend] 2025-11-21 15:26:32 | Running migrations...
[backend] 2025-11-21 15:26:33 | Server started on :8080

✓ Backend is now healthy! Problem solved.

Health check: ✓ /health returned 200
Uptime: 45 seconds (no restarts)
```

### 3. Multi-Container Log Analysis

```
User: /docker logs

Agent: Streaming logs from all containers...

[Live combined log view]

[15:30:15] [postgres]  ✓ Database system ready to accept connections
[15:30:16] [redis]     ✓ Ready to accept connections
[15:30:17] [backend]   ✓ Server started on :8080
[15:30:18] [frontend]  ✓ Webpack dev server listening on 3000
[15:30:25] [backend]   → GET /api/users - 200 (45ms)
[15:30:25] [postgres]  → Query: SELECT * FROM users (12ms)
[15:30:26] [frontend]  → Rendering UserList component
[15:30:30] [backend]   → POST /api/orders - 201 (234ms)
[15:30:30] [postgres]  → BEGIN transaction
[15:30:30] [backend]   → Publishing event to Redis
[15:30:30] [redis]     → PUBLISH orders:created
[15:30:31] [postgres]  → COMMIT
[15:30:31] [backend]   ✓ Order created successfully
[15:30:35] [backend]   ⚠️ Slow query detected (1.2s)
[15:30:35] [postgres]  → Query: SELECT * FROM orders WHERE... (1200ms)
[15:30:40] [backend]   ✗ ERROR: Connection timeout to payment service
[15:30:40] [backend]   → Stack trace: payment.go:67

[Commands]
  /filter backend     - Show only backend logs
  /filter error       - Show only errors
  /since 5m           - Last 5 minutes
  /search "payment"   - Search for term
  /export logs.txt    - Save to file

User: /filter error

Agent: Filtering for errors...

[15:30:40] [backend]   ✗ ERROR: Connection timeout to payment service
[15:30:40] [backend]   → Stack trace:
                         payment.go:67: callPaymentAPI()
                         order.go:234: CreateOrder()
                         handler.go:45: handleCreateOrder()
[15:32:15] [backend]   ✗ ERROR: Redis connection lost
[15:32:15] [redis]     ✗ Client disconnected unexpectedly
[15:32:16] [backend]   → Retrying Redis connection...
[15:32:17] [redis]     ✓ New client connected
[15:32:17] [backend]   ✓ Redis connection restored

Found 2 error patterns in last hour:
1. Payment service timeouts (3 occurrences)
2. Redis connection drops (1 occurrence, recovered)

Would you like me to:
1. Investigate payment service connectivity
2. Check Redis container health
3. View full context for any error
```

### 4. Container Resource Optimization

```
User: /docker stats

Agent: Monitoring container resource usage...

Container Resource Usage
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

┌───────────┬────────┬──────────┬──────────┬──────────┐
│ Container │ CPU    │ Memory   │ Net I/O  │ Block I/O│
├───────────┼────────┼──────────┼──────────┼──────────┤
│ frontend  │ 45.2%  │ 523 MB   │ 2.3 MB   │ 45 KB    │
│ backend   │ 12.5%  │ 245 MB   │ 5.7 MB   │ 234 KB   │
│ postgres  │ 2.1%   │ 156 MB   │ 1.2 MB   │ 12 MB    │
│ redis     │ 0.3%   │ 12 MB    │ 234 KB   │ 0 B      │
└───────────┴────────┴──────────┴──────────┴──────────┘

Total: 60.1% CPU, 936 MB Memory

⚠️ Optimization Opportunities:

1. Frontend Container - High Memory Usage
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Memory: 523 MB (50% higher than typical)
Cause: Webpack dev server + large bundle

Suggestions:
✓ Enable webpack cache to disk (reduce memory)
✓ Use production build for testing (160 MB)
✓ Limit webpack workers (--max-workers=2)
✓ Add memory limit to container

Potential savings: ~200 MB

2. Backend Container - CPU Spikes
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CPU usage spiking to 45% during requests
Baseline: 2-3%, Spikes: 45%

Analysis:
[Checking container processes]

Top processes:
  PID 1: /app/server (45% CPU during spike)
  PID 23: gopls (language server, 8% CPU)

Recommendations:
✓ Profile application to find hotspots
✓ Add CPU limits to prevent resource hogging
✓ Consider disabling gopls in container

3. Image Size Optimization
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Current image sizes:
- frontend: 1.2 GB
- backend: 856 MB

Potential optimizations:

Frontend (1.2 GB → 400 MB):
✓ Use multi-stage build
✓ Alpine base image (node:20-alpine)
✓ Remove dev dependencies in final image
✓ .dockerignore node_modules

Backend (856 MB → 25 MB):
✓ Use scratch or alpine base
✓ Build static binary
✓ Multi-stage build
✓ Strip debug symbols

Total savings: 1.6 GB (60% reduction)

Would you like me to:
1. Generate optimized Dockerfiles
2. Add resource limits to docker-compose.yml
3. Build and test optimized images
```

### 5. Container Filesystem Inspection

```
User: /docker inspect backend filesystem

Agent: Inspecting backend container filesystem...

Container Filesystem: backend (i9j0k1l2)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Mount Points:
/app              → ./backend (read-write, host volume)
/go/pkg/mod       → go-mod-cache (volume)
/tmp              → tmpfs (memory)

Filesystem Structure:
/app/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── models/
│   └── services/
├── go.mod
├── go.sum
└── tmp/
    └── main (compiled binary)

Modified files (since container start):
- /app/tmp/main (2 min ago) - recompiled binary
- /app/internal/handlers/user.go (5 min ago) - your edit
- /tmp/go-build-cache/ (ongoing) - build cache

Unexpected files:
⚠️ /app/.env found (should be in .dockerignore)
⚠️ Contains: DATABASE_PASSWORD, JWT_SECRET

Security issue: Secrets in container image!

Recommendations:
✓ Add .env to .dockerignore
✓ Use environment variables instead
✓ Rebuild image without secrets

[Commands]
  /docker cp backend:/app/.env .    - Copy file out
  /docker exec backend cat /app/.env - View file content
  /docker files diff                 - Show changed files

User: Show me the .env file

Agent: Reading /app/.env from container...

⚠️ WARNING: This file contains sensitive data!

.env contents:
DATABASE_URL=postgresql://user:password@postgres:5432/myapp
JWT_SECRET=super-secret-key-123
API_KEY=sk_live_abc123def456
REDIS_URL=redis://redis:6379

Issues found:
🔴 CRITICAL: Production API key in development container
🔴 CRITICAL: JWT secret in plain text
⚠️ WARNING: Database password exposed

This file should NOT be in the container image!

Recommended fixes:
1. Add .env to .dockerignore
2. Pass secrets via environment variables
3. Use Docker secrets or env files
4. Rotate exposed credentials immediately

Updated docker-compose.yml:
```yaml
services:
  backend:
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - JWT_SECRET=${JWT_SECRET}
      - API_KEY=${API_KEY}
    env_file:
      - .env.local  # Git-ignored file
```

Shall I:
1. Fix the Dockerfile and rebuild
2. Rotate the exposed credentials
3. Update docker-compose.yml
```

---

## Technical Approach

### Docker API Integration

**Docker Client:**
```go
import (
    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types"
)

type DockerManager struct {
    client *client.Client
    ctx    context.Context
}

func NewDockerManager() (*DockerManager, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        return nil, err
    }
    
    return &DockerManager{
        client: cli,
        ctx:    context.Background(),
    }, nil
}

func (dm *DockerManager) BuildImage(dockerfile, tag string) error {
    buildContext, _ := archive.TarWithOptions(".", &archive.TarOptions{})
    
    opts := types.ImageBuildOptions{
        Tags:       []string{tag},
        Dockerfile: dockerfile,
        Remove:     true,
    }
    
    resp, err := dm.client.ImageBuild(dm.ctx, buildContext, opts)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Stream build output
    return dm.streamBuildOutput(resp.Body)
}
```

### Container Monitoring

**Resource Tracking:**
```go
func (dm *DockerManager) MonitorContainer(containerID string) (*ContainerStats, error) {
    stats, err := dm.client.ContainerStats(dm.ctx, containerID, false)
    if err != nil {
        return nil, err
    }
    defer stats.Body.Close()
    
    var v *types.StatsJSON
    if err := json.NewDecoder(stats.Body).Decode(&v); err != nil {
        return nil, err
    }
    
    // Calculate percentages
    cpuPercent := calculateCPUPercent(v)
    memPercent := float64(v.MemoryStats.Usage) / float64(v.MemoryStats.Limit) * 100
    
    return &ContainerStats{
        CPUPercent:    cpuPercent,
        MemoryUsage:   v.MemoryStats.Usage,
        MemoryLimit:   v.MemoryStats.Limit,
        MemoryPercent: memPercent,
        NetworkRx:     v.Networks["eth0"].RxBytes,
        NetworkTx:     v.Networks["eth0"].TxBytes,
    }, nil
}
```

### Log Aggregation

**Multi-Container Logs:**
```go
func (dm *DockerManager) StreamLogs(containers []string) (chan LogEntry, error) {
    logChan := make(chan LogEntry, 100)
    
    for _, containerID := range containers {
        go func(id string) {
            opts := types.ContainerLogsOptions{
                ShowStdout: true,
                ShowStderr: true,
                Follow:     true,
                Timestamps: true,
            }
            
            reader, err := dm.client.ContainerLogs(dm.ctx, id, opts)
            if err != nil {
                return
            }
            defer reader.Close()
            
            scanner := bufio.NewScanner(reader)
            for scanner.Scan() {
                logChan <- LogEntry{
                    ContainerID: id,
                    Timestamp:   time.Now(),
                    Message:     scanner.Text(),
                }
            }
        }(containerID)
    }
    
    return logChan, nil
}
```

### Docker Compose Integration

**Compose Parser:**
```go
import "gopkg.in/yaml.v3"

type ComposeFile struct {
    Version  string
    Services map[string]Service
}

type Service struct {
    Image       string
    Build       BuildConfig
    Ports       []string
    Environment map[string]string
    Volumes     []string
    DependsOn   []string
}

func (dm *DockerManager) StartCompose(composeFile string) error {
    var compose ComposeFile
    data, _ := os.ReadFile(composeFile)
    yaml.Unmarshal(data, &compose)
    
    // Start services in dependency order
    return dm.startServicesInOrder(compose.Services)
}
```

---

## Value Propositions

### For All Developers
- Unified development workflow
- No context switching
- Faster iteration cycles
- Integrated debugging

### For Backend Developers
- Microservices development simplified
- Easy multi-container orchestration
- Database container management
- Service mesh testing

### For DevOps Engineers
- Container optimization tools
- Resource monitoring
- Security scanning
- Image size analysis

---

## Implementation Phases

### Phase 1: Core Docker Operations (2 weeks)
- Docker API integration
- Build/run/stop containers
- Log streaming
- Basic monitoring

### Phase 2: Development Workflow (2 weeks)
- Hot reload support
- Volume mounting
- Multi-container coordination
- Docker Compose integration

### Phase 3: Debugging Tools (2 weeks)
- Filesystem inspection
- Process monitoring
- Network debugging
- Interactive shell access

### Phase 4: Optimization (1 week)
- Resource monitoring
- Image size analysis
- Security scanning
- Performance profiling

---

## Success Metrics

**Adoption:**
- 75%+ use Docker features
- 60%+ replace Docker Desktop
- 80%+ use hot reload
- 50%+ debug in containers

**Impact:**
- 50% faster iteration cycles
- 70% reduction in manual docker commands
- 60% fewer "works on my machine" issues
- 40% faster container debugging

**Satisfaction:**
- 4.8+ rating
- "Never leave Forge for Docker" feedback
- "Development workflow transformed" comments
