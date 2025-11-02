# How to Test Tools

Step-by-step guide to testing custom tools in Forge.

## Overview

Testing ensures your tools work correctly, handle errors gracefully, and integrate properly with agents. This guide covers unit testing, integration testing, and mocking.

**Time to complete:** 15 minutes

**What you'll learn:**
- Write unit tests for tools
- Test tool parameters
- Mock dependencies
- Test error cases
- Integration testing with agents

---

## Prerequisites

- Custom tool implementation
- Go testing basics
- Familiarity with testify/assert (recommended)

---

## Step 1: Set Up Testing

### Install Testing Dependencies

```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
```

### Create Test File

```go
// tools/calculator_test.go
package tools

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

---

## Step 2: Test Tool Interface Methods

### Test Name()

```go
func TestCalculator_Name(t *testing.T) {
    calc := NewCalculator()
    assert.Equal(t, "calculator", calc.Name())
}
```

### Test Description()

```go
func TestCalculator_Description(t *testing.T) {
    calc := NewCalculator()
    desc := calc.Description()
    
    assert.NotEmpty(t, desc)
    assert.Contains(t, desc, "arithmetic")
}
```

### Test IsLoopBreaking()

```go
func TestCalculator_IsLoopBreaking(t *testing.T) {
    calc := NewCalculator()
    assert.False(t, calc.IsLoopBreaking())
}
```

---

## Step 3: Test Parameters Schema

### Validate Schema Structure

```go
func TestCalculator_Parameters(t *testing.T) {
    calc := NewCalculator()
    params := calc.Parameters()
    
    // Check type
    assert.Equal(t, "object", params["type"])
    
    // Check properties exist
    properties, ok := params["properties"].(map[string]interface{})
    require.True(t, ok, "properties must be a map")
    
    assert.Contains(t, properties, "operation")
    assert.Contains(t, properties, "a")
    assert.Contains(t, properties, "b")
}
```

### Validate Required Fields

```go
func TestCalculator_Parameters_Required(t *testing.T) {
    calc := NewCalculator()
    params := calc.Parameters()
    
    required, ok := params["required"].([]string)
    require.True(t, ok, "required must be string slice")
    
    assert.Contains(t, required, "operation")
    assert.Contains(t, required, "a")
    assert.Contains(t, required, "b")
}
```

### Validate Enum Values

```go
func TestCalculator_Parameters_OperationEnum(t *testing.T) {
    calc := NewCalculator()
    params := calc.Parameters()
    
    properties := params["properties"].(map[string]interface{})
    operation := properties["operation"].(map[string]interface{})
    
    enum, ok := operation["enum"].([]string)
    require.True(t, ok, "enum must be string slice")
    
    assert.Contains(t, enum, "add")
    assert.Contains(t, enum, "subtract")
    assert.Contains(t, enum, "multiply")
    assert.Contains(t, enum, "divide")
}
```

---

## Step 4: Test Execute() Method

### Table-Driven Tests

```go
func TestCalculator_Execute(t *testing.T) {
    tests := []struct {
        name      string
        args      map[string]interface{}
        want      string
        wantErr   bool
    }{
        {
            name: "addition",
            args: map[string]interface{}{
                "operation": "add",
                "a":         5.0,
                "b":         3.0,
            },
            want:    "Result: 8",
            wantErr: false,
        },
        {
            name: "subtraction",
            args: map[string]interface{}{
                "operation": "subtract",
                "a":         10.0,
                "b":         4.0,
            },
            want:    "Result: 6",
            wantErr: false,
        },
        {
            name: "multiplication",
            args: map[string]interface{}{
                "operation": "multiply",
                "a":         6.0,
                "b":         7.0,
            },
            want:    "Result: 42",
            wantErr: false,
        },
        {
            name: "division",
            args: map[string]interface{}{
                "operation": "divide",
                "a":         20.0,
                "b":         4.0,
            },
            want:    "Result: 5",
            wantErr: false,
        },
        {
            name: "division by zero",
            args: map[string]interface{}{
                "operation": "divide",
                "a":         10.0,
                "b":         0.0,
            },
            want:    "",
            wantErr: true,
        },
        {
            name: "invalid operation",
            args: map[string]interface{}{
                "operation": "modulo",
                "a":         10.0,
                "b":         3.0,
            },
            want:    "",
            wantErr: true,
        },
    }
    
    calc := NewCalculator()
    ctx := context.Background()
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := calc.Execute(ctx, tt.args)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, result)
            }
        })
    }
}
```

---

## Step 5: Test Error Cases

### Missing Arguments

```go
func TestCalculator_Execute_MissingArguments(t *testing.T) {
    calc := NewCalculator()
    ctx := context.Background()
    
    tests := []struct {
        name string
        args map[string]interface{}
    }{
        {
            name: "missing operation",
            args: map[string]interface{}{
                "a": 5.0,
                "b": 3.0,
            },
        },
        {
            name: "missing a",
            args: map[string]interface{}{
                "operation": "add",
                "b":         3.0,
            },
        },
        {
            name: "missing b",
            args: map[string]interface{}{
                "operation": "add",
                "a":         5.0,
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := calc.Execute(ctx, tt.args)
            assert.Error(t, err)
        })
    }
}
```

### Invalid Argument Types

```go
func TestCalculator_Execute_InvalidTypes(t *testing.T) {
    calc := NewCalculator()
    ctx := context.Background()
    
    tests := []struct {
        name string
        args map[string]interface{}
    }{
        {
            name: "operation not string",
            args: map[string]interface{}{
                "operation": 123,
                "a":         5.0,
                "b":         3.0,
            },
        },
        {
            name: "a not number",
            args: map[string]interface{}{
                "operation": "add",
                "a":         "five",
                "b":         3.0,
            },
        },
        {
            name: "b not number",
            args: map[string]interface{}{
                "operation": "add",
                "a":         5.0,
                "b":         "three",
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := calc.Execute(ctx, tt.args)
            assert.Error(t, err)
        })
    }
}
```

### Context Cancellation

```go
func TestCalculator_Execute_ContextCancellation(t *testing.T) {
    calc := NewCalculator()
    
    // Create canceled context
    ctx, cancel := context.WithCancel(context.Background())
    cancel()
    
    args := map[string]interface{}{
        "operation": "add",
        "a":         5.0,
        "b":         3.0,
    }
    
    _, err := calc.Execute(ctx, args)
    assert.Error(t, err)
    assert.True(t, errors.Is(err, context.Canceled))
}
```

---

## Step 6: Test with Mocks

### Mock HTTP Client

For tools that make HTTP requests:

```go
type MockHTTPClient struct {
    Response *http.Response
    Err      error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.Response, m.Err
}

func TestWeatherTool_Execute_Success(t *testing.T) {
    // Create mock response
    mockResponse := &http.Response{
        StatusCode: 200,
        Body: io.NopCloser(strings.NewReader(`{
            "main": {"temp": 72, "humidity": 65},
            "weather": [{"description": "sunny"}]
        }`)),
    }
    
    // Create tool with mock client
    tool := &WeatherTool{
        client: &MockHTTPClient{Response: mockResponse},
    }
    
    ctx := context.Background()
    args := map[string]interface{}{
        "location": "London",
        "units":    "fahrenheit",
    }
    
    result, err := tool.Execute(ctx, args)
    
    assert.NoError(t, err)
    assert.Contains(t, result, "London")
    assert.Contains(t, result, "72")
}
```

### Mock Database

For tools that query databases:

```go
type MockDB struct {
    QueryFunc func(ctx context.Context, query string) ([]Row, error)
}

func (m *MockDB) Query(ctx context.Context, query string) ([]Row, error) {
    return m.QueryFunc(ctx, query)
}

func TestDatabaseTool_Execute(t *testing.T) {
    mockDB := &MockDB{
        QueryFunc: func(ctx context.Context, query string) ([]Row, error) {
            return []Row{
                {"id": 1, "name": "Alice"},
                {"id": 2, "name": "Bob"},
            }, nil
        },
    }
    
    tool := &DatabaseTool{db: mockDB}
    
    ctx := context.Background()
    args := map[string]interface{}{
        "query": "SELECT * FROM users",
    }
    
    result, err := tool.Execute(ctx, args)
    
    assert.NoError(t, err)
    assert.Contains(t, result, "Alice")
    assert.Contains(t, result, "Bob")
}
```

---

## Step 7: Integration Testing

### Test Tool with Agent

```go
func TestCalculator_WithAgent(t *testing.T) {
    // Create mock provider
    mockProvider := &MockProvider{
        responses: []string{
            `[I'll calculate this]
            
<tool>
{
  "server_name": "local",
  "tool_name": "calculator",
  "arguments": {"operation": "add", "a": 5, "b": 3}
}
</tool>`,
            `The result is 8`,
        },
    }
    
    // Create memory
    mem := memory.NewConversationMemory(1000)
    
    // Create tools
    tools := []core.Tool{
        NewCalculator(),
        tool.NewTaskCompletion(),
    }
    
    // Create agent
    agent, err := core.NewAgent(mockProvider, mem, tools)
    require.NoError(t, err)
    
    // Create mock executor
    mockExecutor := &MockExecutor{
        inputs: []string{"What is 5 + 3?"},
    }
    
    // Run agent
    ctx := context.Background()
    err = agent.Run(ctx, mockExecutor)
    assert.NoError(t, err)
    
    // Verify tool was called
    messages := mem.GetMessages()
    hasToolResult := false
    for _, msg := range messages {
        if msg.Role == "tool" && strings.Contains(msg.Content, "calculator") {
            hasToolResult = true
            assert.Contains(t, msg.Content, "Result: 8")
            break
        }
    }
    assert.True(t, hasToolResult, "Tool should have been called")
}
```

---

## Step 8: Benchmark Tests

### Basic Benchmark

```go
func BenchmarkCalculator_Execute(b *testing.B) {
    calc := NewCalculator()
    ctx := context.Background()
    args := map[string]interface{}{
        "operation": "multiply",
        "a":         12345.0,
        "b":         67890.0,
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        calc.Execute(ctx, args)
    }
}
```

### Benchmark with Different Operations

```go
func BenchmarkCalculator_Operations(b *testing.B) {
    operations := []string{"add", "subtract", "multiply", "divide"}
    calc := NewCalculator()
    ctx := context.Background()
    
    for _, op := range operations {
        b.Run(op, func(b *testing.B) {
            args := map[string]interface{}{
                "operation": op,
                "a":         100.0,
                "b":         50.0,
            }
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                calc.Execute(ctx, args)
            }
        })
    }
}
```

---

## Testing Helpers

### Assert Tool Implementation

```go
func assertImplementsTool(t *testing.T, tool interface{}) {
    t.Helper()
    
    _, ok := tool.(core.Tool)
    assert.True(t, ok, "must implement Tool interface")
}

func TestCalculator_ImplementsTool(t *testing.T) {
    calc := NewCalculator()
    assertImplementsTool(t, calc)
}
```

### Validate JSON Schema

```go
func validateJSONSchema(t *testing.T, schema map[string]interface{}) {
    t.Helper()
    
    // Check required fields
    assert.Contains(t, schema, "type")
    assert.Contains(t, schema, "properties")
    
    // Validate type
    assert.Equal(t, "object", schema["type"])
    
    // Validate properties
    properties, ok := schema["properties"].(map[string]interface{})
    assert.True(t, ok, "properties must be a map")
    assert.NotEmpty(t, properties, "properties cannot be empty")
}

func TestCalculator_SchemaValid(t *testing.T) {
    calc := NewCalculator()
    schema := calc.Parameters()
    validateJSONSchema(t, schema)
}
```

---

## Best Practices

### 1. Test All Public Methods

```go
// ✅ Good: Test all interface methods
func TestMyTool_Name(t *testing.T) { }
func TestMyTool_Description(t *testing.T) { }
func TestMyTool_Parameters(t *testing.T) { }
func TestMyTool_Execute(t *testing.T) { }
func TestMyTool_IsLoopBreaking(t *testing.T) { }

// ❌ Bad: Only test Execute
func TestMyTool_Execute(t *testing.T) { }
```

### 2. Use Table-Driven Tests

```go
// ✅ Good: Many cases in one test
tests := []struct{
    name string
    args map[string]interface{}
    want string
}{
    {"case1", args1, "result1"},
    {"case2", args2, "result2"},
}

// ❌ Bad: Separate test for each case
func TestCase1(t *testing.T) { }
func TestCase2(t *testing.T) { }
```

### 3. Test Error Cases

```go
// ✅ Good: Test errors
func TestMyTool_InvalidInput(t *testing.T) {
    _, err := tool.Execute(ctx, invalidArgs)
    assert.Error(t, err)
}

// ❌ Bad: Only test success
func TestMyTool_Execute(t *testing.T) {
    result, _ := tool.Execute(ctx, args)
    assert.Equal(t, want, result)
}
```

### 4. Use Meaningful Test Names

```go
// ✅ Good: Descriptive name
func TestCalculator_Execute_DivisionByZero_ReturnsError(t *testing.T)

// ❌ Bad: Vague name
func TestCalculator_Error(t *testing.T)
```

### 5. Test Context Handling

```go
// ✅ Good: Test cancellation
ctx, cancel := context.WithCancel(context.Background())
cancel()
_, err := tool.Execute(ctx, args)
assert.Error(t, err)

// ❌ Bad: Never test context
```

---

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Specific Test

```bash
go test -run TestCalculator_Execute ./tools
```

### Run with Coverage

```bash
go test -cover ./tools
go test -coverprofile=coverage.out ./tools
go tool cover -html=coverage.out
```

### Run with Verbose Output

```bash
go test -v ./tools
```

### Run Benchmarks

```bash
go test -bench=. ./tools
```

### Run with Race Detector

```bash
go test -race ./tools
```

---

## Complete Test Suite Example

```go
package tools

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Interface tests
func TestCalculator_ImplementsTool(t *testing.T) {
    calc := NewCalculator()
    var _ core.Tool = calc // Compile-time check
}

func TestCalculator_Name(t *testing.T) {
    calc := NewCalculator()
    assert.Equal(t, "calculator", calc.Name())
}

func TestCalculator_Description(t *testing.T) {
    calc := NewCalculator()
    desc := calc.Description()
    assert.NotEmpty(t, desc)
    assert.Contains(t, desc, "arithmetic")
}

func TestCalculator_IsLoopBreaking(t *testing.T) {
    calc := NewCalculator()
    assert.False(t, calc.IsLoopBreaking())
}

// Schema tests
func TestCalculator_Parameters(t *testing.T) {
    calc := NewCalculator()
    params := calc.Parameters()
    
    assert.Equal(t, "object", params["type"])
    
    properties := params["properties"].(map[string]interface{})
    assert.Contains(t, properties, "operation")
    assert.Contains(t, properties, "a")
    assert.Contains(t, properties, "b")
    
    required := params["required"].([]string)
    assert.Contains(t, required, "operation")
    assert.Contains(t, required, "a")
    assert.Contains(t, required, "b")
}

// Execute tests
func TestCalculator_Execute(t *testing.T) {
    tests := []struct {
        name    string
        args    map[string]interface{}
        want    string
        wantErr bool
    }{
        {
            name: "add",
            args: map[string]interface{}{
                "operation": "add",
                "a": 5.0,
                "b": 3.0,
            },
            want:    "Result: 8",
            wantErr: false,
        },
        {
            name: "divide by zero",
            args: map[string]interface{}{
                "operation": "divide",
                "a": 10.0,
                "b": 0.0,
            },
            want:    "",
            wantErr: true,
        },
    }
    
    calc := NewCalculator()
    ctx := context.Background()
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := calc.Execute(ctx, tt.args)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, result)
            }
        })
    }
}

// Benchmark
func BenchmarkCalculator_Execute(b *testing.B) {
    calc := NewCalculator()
    ctx := context.Background()
    args := map[string]interface{}{
        "operation": "multiply",
        "a": 123.0,
        "b": 456.0,
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        calc.Execute(ctx, args)
    }
}
```

---

## Next Steps

- Read [Testing Reference](../reference/testing.md) for advanced patterns
- See [Create Custom Tool](create-custom-tool.md) for tool development
- Learn [Error Handling](handle-errors.md) for error testing
- Check [API Reference](../reference/api-reference.md) for Tool interface

You're now ready to write comprehensive tests for any tool!