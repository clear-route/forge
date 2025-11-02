# Example: Calculator Agent

Complete example of an agent that can perform arithmetic calculations.

## Overview

This example demonstrates:
- Creating a custom tool
- Registering tools with an agent
- Handling tool parameters
- Error handling in tools
- Agent loop with tool usage

**Time to complete:** 10 minutes

---

## Complete Code

### Calculator Tool

```go
// tools/calculator.go
package tools

import (
    "context"
    "fmt"
)

type Calculator struct{}

func NewCalculator() *Calculator {
    return &Calculator{}
}

func (c *Calculator) Name() string {
    return "calculator"
}

func (c *Calculator) Description() string {
    return "Performs arithmetic operations: add, subtract, multiply, divide"
}

func (c *Calculator) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "operation": map[string]interface{}{
                "type":        "string",
                "description": "The arithmetic operation to perform",
                "enum":        []string{"add", "subtract", "multiply", "divide"},
            },
            "a": map[string]interface{}{
                "type":        "number",
                "description": "First number",
            },
            "b": map[string]interface{}{
                "type":        "number",
                "description": "Second number",
            },
        },
        "required": []string{"operation", "a", "b"},
    }
}

func (c *Calculator) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    // Validate operation
    operation, ok := args["operation"].(string)
    if !ok {
        return "", fmt.Errorf("operation must be a string")
    }
    
    // Validate numbers
    a, ok := args["a"].(float64)
    if !ok {
        return "", fmt.Errorf("a must be a number")
    }
    
    b, ok := args["b"].(float64)
    if !ok {
        return "", fmt.Errorf("b must be a number")
    }
    
    // Perform calculation
    var result float64
    switch operation {
    case "add":
        result = a + b
    case "subtract":
        result = a - b
    case "multiply":
        result = a * b
    case "divide":
        if b == 0 {
            return "", fmt.Errorf("division by zero")
        }
        result = a / b
    default:
        return "", fmt.Errorf("unknown operation: %s", operation)
    }
    
    return fmt.Sprintf("Result: %.2f", result), nil
}

func (c *Calculator) IsLoopBreaking() bool {
    return false
}
```

### Main Application

```go
// main.go
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
    // Get API key
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENAI_API_KEY environment variable not set")
    }
    
    // Create provider
    provider := openai.NewProvider("gpt-4", apiKey)
    
    // Create memory
    mem := memory.NewConversationMemory(8000)
    
    // Create tools
    agentTools := []core.Tool{
        tools.NewCalculator(),
        tool.NewTaskCompletion(),
        tool.NewAskQuestion(executor),
        tool.NewConverse(executor),
    }
    
    // Create agent with custom system prompt
    systemPrompt := `You are a helpful math assistant.
When users ask math questions, use the calculator tool.
Show your work and explain the steps.`
    
    agent, err := core.NewAgent(
        provider,
        mem,
        agentTools,
        core.WithSystemPrompt(systemPrompt),
        core.WithMaxIterations(10),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create executor
    executor := cli.NewExecutor()
    
    // Run agent
    ctx := context.Background()
    log.Println("Calculator Agent started. Type 'exit' to quit.")
    
    if err := agent.Run(ctx, executor); err != nil {
        log.Fatal(err)
    }
}
```

---

## Usage Examples

### Example 1: Simple Addition

**User:** What is 15 + 27?

**Agent:**
```
[I'll use the calculator to add these numbers]

<tool>
{
  "server_name": "local",
  "tool_name": "calculator",
  "arguments": {
    "operation": "add",
    "a": 15,
    "b": 27
  }
}
</tool>
```

**Tool:** calculator  
Result: 42.00

**Agent:** 15 + 27 = 42

---

### Example 2: Division with Error

**User:** What is 10 divided by 0?

**Agent:**
```
[Let me calculate that]

<tool>
{
  "server_name": "local",
  "tool_name": "calculator",
  "arguments": {
    "operation": "divide",
    "a": 10,
    "b": 0
  }
}
</tool>
```

**Tool:** calculator  
Error: division by zero

**Agent:**
```
[I see the error - division by zero is not allowed]

I cannot calculate 10 divided by 0 because division by zero is undefined in mathematics.
```

---

### Example 3: Complex Calculation

**User:** Calculate (25 * 4) + (100 / 5)

**Agent:**
```
[I'll break this down into steps]

<tool>
{
  "server_name": "local",
  "tool_name": "calculator",
  "arguments": {
    "operation": "multiply",
    "a": 25,
    "b": 4
  }
}
</tool>
```

**Tool:** calculator  
Result: 100.00

**Agent:**
```
[Now the division part]

<tool>
{
  "server_name": "local",
  "tool_name": "calculator",
  "arguments": {
    "operation": "divide",
    "a": 100,
    "b": 5
  }
}
</tool>
```

**Tool:** calculator  
Result: 20.00

**Agent:**
```
[Finally, add them together]

<tool>
{
  "server_name": "local",
  "tool_name": "calculator",
  "arguments": {
    "operation": "add",
    "a": 100,
    "b": 20
  }
}
</tool>
```

**Tool:** calculator  
Result: 120.00

**Agent:**
```
Breaking down the calculation:
- 25 * 4 = 100
- 100 / 5 = 20
- 100 + 20 = 120

Therefore, (25 * 4) + (100 / 5) = 120
```

---

## Running the Example

### Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/forge.git
cd forge

# Install dependencies
go mod download

# Set API key
export OPENAI_API_KEY="sk-your-key-here"
```

### Run

```bash
# Run the agent
go run examples/calculator-agent/main.go
```

### Test

```bash
# Run tests
go test ./examples/calculator-agent/...
```

---

## Key Concepts Demonstrated

### 1. Tool Implementation

The calculator implements all five required methods:
- `Name()` - Unique identifier
- `Description()` - Tells LLM when to use it
- `Parameters()` - Defines expected arguments
- `Execute()` - Performs the operation
- `IsLoopBreaking()` - Whether it ends the loop

### 2. Parameter Validation

```go
// Type checking
operation, ok := args["operation"].(string)
if !ok {
    return "", fmt.Errorf("operation must be a string")
}

// Value validation
if b == 0 {
    return "", fmt.Errorf("division by zero")
}
```

### 3. Error Handling

The agent gracefully handles errors:
- Division by zero
- Invalid operations
- Type mismatches

### 4. Multi-Step Reasoning

Agent breaks down complex calculations into steps, using the calculator multiple times.

---

## Variations

### Scientific Calculator

Add more operations:

```go
func (c *Calculator) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "operation": map[string]interface{}{
                "type": "string",
                "enum": []string{
                    "add", "subtract", "multiply", "divide",
                    "power", "sqrt", "sin", "cos", "log",
                },
            },
            // ... rest of parameters
        },
    }
}
```

### Unit Conversion

Add unit conversion:

```go
type UnitConverter struct{}

func (u *UnitConverter) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    value := args["value"].(float64)
    from := args["from"].(string)
    to := args["to"].(string)
    
    // Convert between units
    result := convert(value, from, to)
    return fmt.Sprintf("%.2f %s = %.2f %s", value, from, result, to), nil
}
```

### Expression Parser

Parse and evaluate full expressions:

```go
type ExpressionCalculator struct{}

func (e *ExpressionCalculator) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    expression := args["expression"].(string)
    
    // Parse and evaluate
    result, err := parseAndEvaluate(expression)
    if err != nil {
        return "", err
    }
    
    return fmt.Sprintf("%s = %.2f", expression, result), nil
}
```

---

## Testing

### Unit Test

```go
func TestCalculator_Execute_Add(t *testing.T) {
    calc := NewCalculator()
    args := map[string]interface{}{
        "operation": "add",
        "a":         5.0,
        "b":         3.0,
    }
    
    result, err := calc.Execute(context.Background(), args)
    
    assert.NoError(t, err)
    assert.Equal(t, "Result: 8.00", result)
}

func TestCalculator_Execute_DivideByZero(t *testing.T) {
    calc := NewCalculator()
    args := map[string]interface{}{
        "operation": "divide",
        "a":         10.0,
        "b":         0.0,
    }
    
    _, err := calc.Execute(context.Background(), args)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "division by zero")
}
```

---

## Next Steps

- Add more mathematical operations
- Implement expression parsing
- Add unit conversion
- Create statistics functions
- Build graphing capabilities

## Related Examples

- [Weather Agent](weather-agent.md) - External API integration
- [Code Assistant](code-assistant.md) - File operations
- [Database Agent](database-agent.md) - Database queries

---

**Complete code:** [examples/calculator-agent/](../../examples/calculator-agent/)