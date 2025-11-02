# How to Deploy to Production

Step-by-step guide to deploying Forge agents in production environments.

## Overview

Production deployment requires careful configuration for reliability, security, monitoring, and performance. This guide covers best practices and common patterns.

**Time to complete:** 30 minutes

**What you'll learn:**
- Configure for production
- Handle secrets securely
- Implement logging and monitoring
- Set up error handling
- Deploy to various platforms
- Scale your agent

---

## Prerequisites

- Working Forge agent
- Basic understanding of deployment concepts
- Access to deployment platform

---

## Step 1: Production Configuration

### Environment-Based Configuration

```go
package main

import (
    "log"
    "os"
    "strconv"
    "time"
)

type Config struct {
    Environment string
    OpenAIKey   string
    Model       string
    Temperature float64
    MaxTokens   int
    Timeout     time.Duration
    MemoryLimit int
    MaxIter     int
    LogLevel    string
}

func LoadConfig() (*Config, error) {
    env := getEnv("ENVIRONMENT", "development")
    
    cfg := &Config{
        Environment: env,
        OpenAIKey:   os.Getenv("OPENAI_API_KEY"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
    }
    
    // Validate required fields
    if cfg.OpenAIKey == "" {
        return nil, fmt.Errorf("OPENAI_API_KEY is required")
    }
    
    // Environment-specific settings
    switch env {
    case "production":
        cfg.Model = getEnv("MODEL", "gpt-4-turbo-preview")
        cfg.Temperature = getEnvFloat("TEMPERATURE", 0.5)
        cfg.MaxTokens = getEnvInt("MAX_TOKENS", 2000)
        cfg.Timeout = getEnvDuration("TIMEOUT", 60*time.Second)
        cfg.MemoryLimit = getEnvInt("MEMORY_LIMIT", 8000)
        cfg.MaxIter = getEnvInt("MAX_ITERATIONS", 15)
        
    case "development":
        cfg.Model = getEnv("MODEL", "gpt-3.5-turbo")
        cfg.Temperature = getEnvFloat("TEMPERATURE", 0.7)
        cfg.MaxTokens = getEnvInt("MAX_TOKENS", 1000)
        cfg.Timeout = getEnvDuration("TIMEOUT", 30*time.Second)
        cfg.MemoryLimit = getEnvInt("MEMORY_LIMIT", 4000)
        cfg.MaxIter = getEnvInt("MAX_ITERATIONS", 10)
        
    default:
        return nil, fmt.Errorf("unknown environment: %s", env)
    }
    
    return cfg, nil
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

func getEnvInt(key string, fallback int) int {
    if value := os.Getenv(key); value != "" {
        if i, err := strconv.Atoi(value); err == nil {
            return i
        }
    }
    return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
    if value := os.Getenv(key); value != "" {
        if f, err := strconv.ParseFloat(value, 64); err == nil {
            return f
        }
    }
    return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if d, err := time.ParseDuration(value); err == nil {
            return d
        }
    }
    return fallback
}
```

### Environment File (.env)

```bash
# .env.production
ENVIRONMENT=production
OPENAI_API_KEY=sk-your-production-key
MODEL=gpt-4-turbo-preview
TEMPERATURE=0.5
MAX_TOKENS=2000
TIMEOUT=60s
MEMORY_LIMIT=8000
MAX_ITERATIONS=15
LOG_LEVEL=info
```

---

## Step 2: Secure Secrets Management

### Never Commit Secrets

```bash
# .gitignore
.env
.env.*
!.env.example
*.key
secrets/
```

### Use Environment Variables

```go
// ✅ Good: Read from environment
apiKey := os.Getenv("OPENAI_API_KEY")

// ❌ Bad: Hardcoded
apiKey := "sk-hardcoded-key"
```

### Use Secret Management Services

#### AWS Secrets Manager

```go
import (
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
)

func getSecret(secretName string) (string, error) {
    sess := session.Must(session.NewSession())
    svc := secretsmanager.New(sess)
    
    input := &secretsmanager.GetSecretValueInput{
        SecretId: aws.String(secretName),
    }
    
    result, err := svc.GetSecretValue(input)
    if err != nil {
        return "", err
    }
    
    return *result.SecretString, nil
}

// Usage
apiKey, err := getSecret("prod/openai/api-key")
```

#### HashiCorp Vault

```go
import "github.com/hashicorp/vault/api"

func getVaultSecret(path string) (string, error) {
    client, err := api.NewClient(api.DefaultConfig())
    if err != nil {
        return "", err
    }
    
    secret, err := client.Logical().Read(path)
    if err != nil {
        return "", err
    }
    
    return secret.Data["value"].(string), nil
}

// Usage
apiKey, err := getVaultSecret("secret/data/openai/api-key")
```

---

## Step 3: Implement Logging

### Structured Logging

```go
import (
    "log/slog"
    "os"
)

func setupLogging(level string) *slog.Logger {
    var logLevel slog.Level
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }
    
    opts := &slog.HandlerOptions{
        Level: logLevel,
    }
    
    handler := slog.NewJSONHandler(os.Stdout, opts)
    return slog.New(handler)
}

// Usage
logger := setupLogging(cfg.LogLevel)

logger.Info("agent started",
    "environment", cfg.Environment,
    "model", cfg.Model,
)

logger.Error("provider error",
    "error", err,
    "model", cfg.Model,
    "tokens", tokenCount,
)
```

### Request Logging

```go
type RequestLogger struct {
    logger *slog.Logger
}

func (r *RequestLogger) LogRequest(ctx context.Context, req Request) {
    r.logger.Info("request started",
        "request_id", req.ID,
        "user_id", req.UserID,
        "model", req.Model,
    )
}

func (r *RequestLogger) LogResponse(ctx context.Context, req Request, resp Response, err error) {
    if err != nil {
        r.logger.Error("request failed",
            "request_id", req.ID,
            "error", err,
            "duration", req.Duration,
        )
        return
    }
    
    r.logger.Info("request completed",
        "request_id", req.ID,
        "tokens", resp.Tokens,
        "duration", req.Duration,
        "cost", resp.Cost,
    )
}
```

---

## Step 4: Monitoring and Metrics

### Prometheus Metrics

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

var (
    requestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "agent_requests_total",
            Help: "Total number of requests",
        },
        []string{"status"},
    )
    
    requestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "agent_request_duration_seconds",
            Help:    "Request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"model"},
    )
    
    tokensUsed = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "agent_tokens_used_total",
            Help: "Total tokens used",
        },
        []string{"model", "type"},
    )
    
    costTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "agent_cost_total_usd",
            Help: "Total cost in USD",
        },
    )
)

func recordMetrics(model string, duration time.Duration, tokens int, cost float64, err error) {
    status := "success"
    if err != nil {
        status = "error"
    }
    
    requestsTotal.WithLabelValues(status).Inc()
    requestDuration.WithLabelValues(model).Observe(duration.Seconds())
    tokensUsed.WithLabelValues(model, "total").Add(float64(tokens))
    costTotal.Add(cost)
}

// Start metrics server
func startMetricsServer(port string) {
    http.Handle("/metrics", promhttp.Handler())
    go http.ListenAndServe(":"+port, nil)
}
```

### Health Checks

```go
type HealthChecker struct {
    provider core.Provider
    memory   core.Memory
}

func (h *HealthChecker) Check() map[string]string {
    status := make(map[string]string)
    
    // Check provider
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    _, err := h.provider.Complete(ctx, []message.Message{
        message.User("health check"),
    })
    
    if err != nil {
        status["provider"] = "unhealthy: " + err.Error()
    } else {
        status["provider"] = "healthy"
    }
    
    // Check memory
    if h.memory.Size() >= 0 {
        status["memory"] = "healthy"
    } else {
        status["memory"] = "unhealthy"
    }
    
    return status
}

// HTTP handler
func healthHandler(checker *HealthChecker) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        status := checker.Check()
        
        // Determine overall health
        healthy := true
        for _, s := range status {
            if !strings.HasPrefix(s, "healthy") {
                healthy = false
                break
            }
        }
        
        if healthy {
            w.WriteHeader(http.StatusOK)
        } else {
            w.WriteHeader(http.StatusServiceUnavailable)
        }
        
        json.NewEncoder(w).Encode(status)
    }
}
```

---

## Step 5: Error Handling and Recovery

### Graceful Shutdown

```go
func main() {
    // Setup
    cfg, _ := LoadConfig()
    agent, _ := setupAgent(cfg)
    
    // Graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Handle signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        log.Println("Shutdown signal received, gracefully stopping...")
        cancel()
    }()
    
    // Run agent
    if err := agent.Run(ctx, executor); err != nil {
        if errors.Is(err, context.Canceled) {
            log.Println("Agent stopped gracefully")
        } else {
            log.Fatalf("Agent error: %v", err)
        }
    }
}
```

### Circuit Breaker with Monitoring

```go
type MonitoredCircuitBreaker struct {
    breaker *CircuitBreaker
    logger  *slog.Logger
    metrics prometheus.Counter
}

func (m *MonitoredCircuitBreaker) Call(fn func() error) error {
    err := m.breaker.Call(fn)
    
    if err != nil && err.Error() == "circuit breaker is open" {
        m.logger.Warn("circuit breaker opened",
            "consecutive_failures", m.breaker.consecutiveFails,
        )
        m.metrics.Inc()
    }
    
    return err
}
```

---

## Step 6: Deployment Platforms

### Docker

**Dockerfile:**

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o agent ./cmd/agent

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/agent .

# Non-root user
RUN addgroup -g 1000 agent && \
    adduser -D -u 1000 -G agent agent
USER agent

EXPOSE 8080

CMD ["./agent"]
```

**docker-compose.yml:**

```yaml
version: '3.8'

services:
  agent:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENVIRONMENT=production
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - MODEL=gpt-4-turbo-preview
      - LOG_LEVEL=info
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### Kubernetes

**deployment.yaml:**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: forge-agent
spec:
  replicas: 3
  selector:
    matchLabels:
      app: forge-agent
  template:
    metadata:
      labels:
        app: forge-agent
    spec:
      containers:
      - name: agent
        image: your-registry/forge-agent:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: openai-secret
              key: api-key
        - name: MODEL
          value: "gpt-4-turbo-preview"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: forge-agent
spec:
  selector:
    app: forge-agent
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

### AWS Lambda

```go
package main

import (
    "context"
    "github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
    UserMessage string `json:"user_message"`
}

type Response struct {
    AgentResponse string `json:"agent_response"`
    Tokens        int    `json:"tokens"`
    Cost          float64 `json:"cost"`
}

func handleRequest(ctx context.Context, req Request) (*Response, error) {
    // Setup agent (reuse across invocations)
    agent, err := setupAgent()
    if err != nil {
        return nil, err
    }
    
    // Process request
    result, tokens, cost, err := processWithAgent(ctx, agent, req.UserMessage)
    if err != nil {
        return nil, err
    }
    
    return &Response{
        AgentResponse: result,
        Tokens:        tokens,
        Cost:          cost,
    }, nil
}

func main() {
    lambda.Start(handleRequest)
}
```

---

## Step 7: Scaling Strategies

### Horizontal Scaling

```go
type LoadBalancer struct {
    agents []*Agent
    mu     sync.Mutex
    index  int
}

func (lb *LoadBalancer) GetAgent() *Agent {
    lb.mu.Lock()
    defer lb.mu.Unlock()
    
    agent := lb.agents[lb.index]
    lb.index = (lb.index + 1) % len(lb.agents)
    
    return agent
}

// Create multiple agent instances
func setupLoadBalancer(count int) *LoadBalancer {
    lb := &LoadBalancer{
        agents: make([]*Agent, count),
    }
    
    for i := 0; i < count; i++ {
        agent, _ := setupAgent()
        lb.agents[i] = agent
    }
    
    return lb
}
```

### Rate Limiting

```go
import "golang.org/x/time/rate"

type RateLimiter struct {
    limiter *rate.Limiter
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
    }
}

func (rl *RateLimiter) Allow(ctx context.Context) error {
    if err := rl.limiter.Wait(ctx); err != nil {
        return fmt.Errorf("rate limit exceeded: %w", err)
    }
    return nil
}

// Usage
limiter := NewRateLimiter(10) // 10 requests per second

func handleRequest(ctx context.Context) error {
    if err := limiter.Allow(ctx); err != nil {
        return err
    }
    
    // Process request
    return nil
}
```

---

## Production Checklist

### Before Deployment

- [ ] Environment-specific configuration
- [ ] Secrets in secure storage (not code)
- [ ] Structured logging implemented
- [ ] Metrics and monitoring setup
- [ ] Health checks configured
- [ ] Error handling and retries
- [ ] Circuit breakers in place
- [ ] Graceful shutdown handling
- [ ] Resource limits set
- [ ] Rate limiting configured

### After Deployment

- [ ] Monitor error rates
- [ ] Track response times
- [ ] Watch token usage and costs
- [ ] Check memory usage
- [ ] Review logs regularly
- [ ] Test health endpoints
- [ ] Verify metrics collection
- [ ] Test failover scenarios
- [ ] Monitor API rate limits
- [ ] Set up alerts

---

## Monitoring Dashboards

### Key Metrics to Track

1. **Request Metrics:**
   - Requests per second
   - Success rate
   - Error rate
   - Average latency

2. **Resource Metrics:**
   - CPU usage
   - Memory usage
   - Token usage
   - Cost per request

3. **Business Metrics:**
   - Daily active users
   - Total cost
   - Cost per user
   - User satisfaction

---

## Next Steps

- Read [Configuration Reference](../reference/configuration.md) for all options
- See [Error Handling Guide](handle-errors.md) for production errors
- Learn [Performance Optimization](optimize-performance.md) for scaling
- Check [Security Best Practices](#) for hardening

Your agent is now production-ready!