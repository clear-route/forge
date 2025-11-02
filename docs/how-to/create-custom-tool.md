# How to Create a Custom Tool

Step-by-step guide to creating custom tools for Forge agents.

## Overview

Custom tools extend agent capabilities by providing new functions. This guide shows how to create, test, and use custom tools.

**Time to complete:** 15 minutes

**What you'll learn:**
- Implement the Tool interface
- Define tool parameters with JSON Schema
- Handle tool arguments and errors
- Test your tool
- Integrate with agents

---

## Prerequisites

- Go 1.21 or later installed
- Basic understanding of Go interfaces
- Familiarity with JSON Schema (helpful but not required)

---

## Step 1: Define Your Tool

Create a new Go file for your tool:

```go
// tools/weather.go
package tools

import (
    "context"
    "fmt"
)

type WeatherTool struct {
    apiKey string
}

func NewWeatherTool(apiKey string) *WeatherTool {
    return &WeatherTool{
        apiKey: apiKey,
    }
}
```

---

## Step 2: Implement the Tool Interface

The Tool interface requires five methods:

```go
func (w *WeatherTool) Name() string {
    return "get_weather"
}

func (w *WeatherTool) Description() string {
    return "Gets current weather information for a specified location"
}

func (w *WeatherTool) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "location": map[string]interface{}{
                "type":        "string",
                "description": "City name or location (e.g., 'London', 'New York, NY')",
            },
            "units": map[string]interface{}{
                "type":        "string",
                "description": "Temperature units",
                "enum":        []string{"celsius", "fahrenheit"},
                "default":     "celsius",
            },
        },
        "required": []string{"location"},
    }
}

func (w *WeatherTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    // Validate and extract arguments
    location, ok := args["location"].(string)
    if !ok {
        return "", fmt.Errorf("location must be a string")
    }
    
    if location == "" {
        return "", fmt.Errorf("location cannot be empty")
    }
    
    // Get units (with default)
    units := "celsius"
    if u, ok := args["units"].(string); ok {
        units = u
    }
    
    // Call weather API (simplified example)
    weather, err := w.fetchWeather(ctx, location, units)
    if err != nil {
        return "", fmt.Errorf("failed to fetch weather: %w", err)
    }
    
    return weather, nil
}

func (w *WeatherTool) IsLoopBreaking() bool {
    return false // Tool doesn't end the conversation
}
```

---

## Step 3: Implement Tool Logic

Add your tool's actual functionality:

```go
func (w *WeatherTool) fetchWeather(ctx context.Context, location, units string) (string, error) {
    // Example: Call a weather API
    // In practice, you'd use a real API like OpenWeatherMap
    
    // Mock response for demonstration
    temp := "22"
    if units == "fahrenheit" {
        temp = "72"
    }
    
    result := fmt.Sprintf(
        "Weather in %s: Sunny, %s°%s, Humidity: 65%%, Wind: 10 km/h",
        location,
        temp,
        unitSymbol(units),
    )
    
    return result, nil
}

func unitSymbol(units string) string {
    if units == "fahrenheit" {
        return "F"
    }
    return "C"
}
```

**Real implementation example:**

```go
import (
    "encoding/json"
    "net/http"
)

type WeatherResponse struct {
    Main struct {
        Temp     float64 `json:"temp"`
        Humidity int     `json:"humidity"`
    } `json:"main"`
    Weather []struct {
        Description string `json:"description"`
    } `json:"weather"`
    Wind struct {
        Speed float64 `json:"speed"`
    } `json:"wind"`
}

func (w *WeatherTool) fetchWeather(ctx context.Context, location, units string) (string, error) {
    // Build API URL
    url := fmt.Sprintf(
        "https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=%s",
        location,
        w.apiKey,
        units,
    )
    
    // Create request with context
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }
    
    // Make request
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("API returned status %d", resp.StatusCode)
    }
    
    // Parse response
    var weatherResp WeatherResponse
    if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
        return "", err
    }
    
    // Format result
    description := "Unknown"
    if len(weatherResp.Weather) > 0 {
        description = weatherResp.Weather[0].Description
    }
    
    result := fmt.Sprintf(
        "Weather in %s: %s, %.1f°%s, Humidity: %d%%, Wind: %.1f m/s",
        location,
        description,
        weatherResp.Main.Temp,
        unitSymbol(units),
        weatherResp.Main.Humidity,
        weatherResp.Wind.Speed,
    )
    
    return result, nil
}
```

---

## Step 4: Add Error Handling

Handle common error cases:

```go
func (w *WeatherTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    // 1. Validate required arguments
    location, ok := args["location"].(string)
    if !ok {
        return "", fmt.Errorf("location must be a string, got %T", args["location"])
    }
    
    if location == "" {
        return "", fmt.Errorf("location cannot be empty")
    }
    
    // 2. Validate optional arguments
    units := "celsius"
    if u, ok := args["units"].(string); ok {
        validUnits := map[string]bool{"celsius": true, "fahrenheit": true}
        if !validUnits[u] {
            return "", fmt.Errorf("invalid units: %s (must be celsius or fahrenheit)", u)
        }
        units = u
    }
    
    // 3. Check context cancellation
    if err := ctx.Err(); err != nil {
        return "", fmt.Errorf("context cancelled: %w", err)
    }
    
    // 4. Execute with timeout
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    
    // 5. Handle API errors gracefully
    weather, err := w.fetchWeather(ctx, location, units)
    if err != nil {
        // Return user-friendly error message
        return "", fmt.Errorf("unable to get weather for %s: %w", location, err)
    }
    
    return weather, nil
}
```

---

## Step 5: Test Your Tool

Create a test file:

```go
// tools/weather_test.go
package tools

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestWeatherTool_Name(t *testing.T) {
    tool := NewWeatherTool("test-key")
    assert.Equal(t, "get_weather", tool.Name())
}

func TestWeatherTool_Parameters(t *testing.T) {
    tool := NewWeatherTool("test-key")
    params := tool.Parameters()
    
    // Verify structure
    assert.Equal(t, "object", params["type"])
    
    // Verify properties
    properties := params["properties"].(map[string]interface{})
    assert.Contains(t, properties, "location")
    assert.Contains(t, properties, "units")
    
    // Verify required
    required := params["required"].([]string)
    assert.Contains(t, required, "location")
}

func TestWeatherTool_Execute_ValidArgs(t *testing.T) {
    tool := NewWeatherTool("test-key")
    
    args := map[string]interface{}{
        "location": "London",
        "units":    "celsius",
    }
    
    ctx := context.Background()
    result, err := tool.Execute(ctx, args)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, result)
    assert.Contains(t, result, "London")
}

func TestWeatherTool_Execute_MissingLocation(t *testing.T) {
    tool := NewWeatherTool("test-key")
    
    args := map[string]interface{}{
        "units": "celsius",
    }
    
    ctx := context.Background()
    _, err := tool.Execute(ctx, args)
    
    assert.Error(t, err)
}

func TestWeatherTool_Execute_InvalidUnits(t *testing.T) {
    tool := NewWeatherTool("test-key")
    
    args := map[string]interface{}{
        "location": "London",
        "units":    "kelvin", // Invalid
    }
    
    ctx := context.Background()
    _, err := tool.Execute(ctx, args)
    
    assert.Error(t, err)
}

func TestWeatherTool_IsLoopBreaking(t *testing.T) {
    tool := NewWeatherTool("test-key")
    assert.False(t, tool.IsLoopBreaking())
}
```

Run tests:

```bash
go test ./tools -v
```

---

## Step 6: Use Your Tool with an Agent

Integrate your tool into an agent:

```go
package main

import (
    "context"
    "log"
    "os"
    
    "github.com/yourusername/forge/pkg/core"
    "github.com/yourusername/forge/pkg/provider/openai"
    "github.com/yourusername/forge/pkg/memory"
    "github.com/yourusername/forge/pkg/executor/cli"
    "github.com/yourusername/forge/pkg/tool"
    "yourproject/tools"
)

func main() {
    // Get API keys
    openaiKey := os.Getenv("OPENAI_API_KEY")
    weatherKey := os.Getenv("WEATHER_API_KEY")
    
    // Create provider
    provider := openai.NewProvider("gpt-4", openaiKey)
    
    // Create memory
    mem := memory.NewConversationMemory(8000)
    
    // Create tools including your custom tool
    agentTools := []core.Tool{
        tools.NewWeatherTool(weatherKey),
        tool.NewTaskCompletion(),
        tool.NewAskQuestion(executor),
        tool.NewConverse(executor),
    }
    
    // Create agent
    agent, err := core.NewAgent(provider, mem, agentTools)
    if err != nil {
        log.Fatal(err)
    }
    
    // Run agent
    executor := cli.NewExecutor()
    ctx := context.Background()
    
    if err := agent.Run(ctx, executor); err != nil {
        log.Fatal(err)
    }
}
```

---

## Advanced Topics

### Multiple Return Values

Some tools need to return structured data:

```go
type SearchResult struct {
    URL         string
    Title       string
    Description string
}

func (s *SearchTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    results, err := s.search(ctx, args["query"].(string))
    if err != nil {
        return "", err
    }
    
    // Format as readable string for agent
    var output strings.Builder
    output.WriteString(fmt.Sprintf("Found %d results:\n\n", len(results)))
    
    for i, result := range results {
        output.WriteString(fmt.Sprintf("%d. %s\n", i+1, result.Title))
        output.WriteString(fmt.Sprintf("   %s\n", result.URL))
        output.WriteString(fmt.Sprintf("   %s\n\n", result.Description))
    }
    
    return output.String(), nil
}
```

### Tools with State

Tools can maintain state:

```go
type CachingTool struct {
    cache map[string]string
    mu    sync.RWMutex
}

func (c *CachingTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    key := args["key"].(string)
    
    // Check cache
    c.mu.RLock()
    if value, ok := c.cache[key]; ok {
        c.mu.RUnlock()
        return value, nil
    }
    c.mu.RUnlock()
    
    // Fetch and cache
    value, err := c.fetch(ctx, key)
    if err != nil {
        return "", err
    }
    
    c.mu.Lock()
    c.cache[key] = value
    c.mu.Unlock()
    
    return value, nil
}
```

### Loop-Breaking Tools

Tools that end the conversation:

```go
func (t *SaveAndExit) IsLoopBreaking() bool {
    return true // Ends agent loop
}

func (t *SaveAndExit) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    data := args["data"].(string)
    
    if err := t.save(data); err != nil {
        return "", err
    }
    
    return "Data saved successfully. Goodbye!", nil
}
```

---

## Best Practices

### 1. Validate All Arguments

```go
// ✅ Good: Thorough validation
location, ok := args["location"].(string)
if !ok {
    return "", fmt.Errorf("location must be string, got %T", args["location"])
}
if location == "" {
    return "", fmt.Errorf("location cannot be empty")
}

// ❌ Bad: No validation
location := args["location"].(string) // Panics if wrong type
```

### 2. Provide Clear Descriptions

```go
// ✅ Good: Clear, specific description
"description": "City name or location (e.g., 'London', 'New York, NY')"

// ❌ Bad: Vague description
"description": "The location"
```

### 3. Use Enums for Fixed Choices

```go
// ✅ Good: Constrain to valid values
"enum": []string{"celsius", "fahrenheit"}

// ❌ Bad: Allow any string
"type": "string"
```

### 4. Handle Errors Gracefully

```go
// ✅ Good: User-friendly error
return "", fmt.Errorf("unable to get weather for %s: %w", location, err)

// ❌ Bad: Technical error exposed
return "", err
```

### 5. Respect Context

```go
// ✅ Good: Check context cancellation
if err := ctx.Err(); err != nil {
    return "", err
}

// ❌ Bad: Ignore context
// Long operation without checking context
```

### 6. Set Timeouts

```go
// ✅ Good: Limit execution time
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()

// ❌ Bad: No timeout
// Could hang forever
```

---

## Common Patterns

### API Client Tool

```go
type APITool struct {
    client  *http.Client
    baseURL string
}

func (a *APITool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    endpoint := args["endpoint"].(string)
    
    url := a.baseURL + endpoint
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    
    resp, err := a.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    return string(body), nil
}
```

### Database Tool

```go
type DatabaseTool struct {
    db *sql.DB
}

func (d *DatabaseTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    query := args["query"].(string)
    
    rows, err := d.db.QueryContext(ctx, query)
    if err != nil {
        return "", err
    }
    defer rows.Close()
    
    // Format results...
    return results, nil
}
```

### File System Tool

```go
type FileReaderTool struct {
    allowedPaths []string
}

func (f *FileReaderTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    path := args["path"].(string)
    
    // Security: Check path is allowed
    if !f.isAllowed(path) {
        return "", fmt.Errorf("access denied: %s", path)
    }
    
    content, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }
    
    return string(content), nil
}
```

---

## Troubleshooting

### Tool Not Being Called

**Problem:** Agent doesn't use your tool

**Solutions:**
1. Check tool name is clear and descriptive
2. Improve tool description
3. Add examples to description
4. Check parameters are well-defined

### Argument Type Errors

**Problem:** Type assertion panics

**Solution:**
```go
// Instead of:
value := args["key"].(string) // Panics if wrong type

// Use:
value, ok := args["key"].(string)
if !ok {
    return "", fmt.Errorf("expected string, got %T", args["key"])
}
```

### Tool Times Out

**Problem:** Tool execution takes too long

**Solutions:**
1. Add timeout to tool execution
2. Optimize slow operations
3. Use caching
4. Make operations async

---

## Next Steps

- Read [Tool Schema Reference](../reference/tool-schema.md) for advanced schemas
- See [Testing Tools](test-tools.md) for comprehensive testing
- Explore [Example Tools](../examples/calculator-agent.md) for inspiration
- Check [API Reference](../reference/api-reference.md) for full Tool interface

---

## Complete Example

Full working weather tool:

```go
package tools

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type WeatherTool struct {
    apiKey string
    client *http.Client
}

func NewWeatherTool(apiKey string) *WeatherTool {
    return &WeatherTool{
        apiKey: apiKey,
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (w *WeatherTool) Name() string {
    return "get_weather"
}

func (w *WeatherTool) Description() string {
    return "Gets current weather for a location. Use city name like 'London' or 'New York, NY'"
}

func (w *WeatherTool) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "location": map[string]interface{}{
                "type":        "string",
                "description": "City name (e.g., 'London', 'Tokyo')",
            },
            "units": map[string]interface{}{
                "type":        "string",
                "description": "Temperature units",
                "enum":        []string{"celsius", "fahrenheit"},
                "default":     "celsius",
            },
        },
        "required": []string{"location"},
    }
}

func (w *WeatherTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    location, ok := args["location"].(string)
    if !ok || location == "" {
        return "", fmt.Errorf("location must be a non-empty string")
    }
    
    units := "metric"
    if u, ok := args["units"].(string); ok && u == "fahrenheit" {
        units = "imperial"
    }
    
    url := fmt.Sprintf(
        "https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=%s",
        location, w.apiKey, units,
    )
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }
    
    resp, err := w.client.Do(req)
    if err != nil {
        return "", fmt.Errorf("API request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("location not found: %s", location)
    }
    
    var weather struct {
        Main struct {
            Temp     float64 `json:"temp"`
            Humidity int     `json:"humidity"`
        } `json:"main"`
        Weather []struct {
            Description string `json:"description"`
        } `json:"weather"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
        return "", err
    }
    
    desc := "Unknown"
    if len(weather.Weather) > 0 {
        desc = weather.Weather[0].Description
    }
    
    unit := "C"
    if units == "imperial" {
        unit = "F"
    }
    
    return fmt.Sprintf(
        "Weather in %s: %s, %.1f°%s, Humidity: %d%%",
        location, desc, weather.Main.Temp, unit, weather.Main.Humidity,
    ), nil
}

func (w *WeatherTool) IsLoopBreaking() bool {
    return false
}
```

You now have a complete, production-ready custom tool!