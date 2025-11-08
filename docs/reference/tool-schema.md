# Tool Schema Reference

Complete reference for defining tool parameters using JSON Schema.

## Table of Contents

- [Overview](#overview)
- [Basic Schema Structure](#basic-schema-structure)
- [Data Types](#data-types)
- [Validation Keywords](#validation-keywords)
- [Common Patterns](#common-patterns)
- [Examples](#examples)
- [Best Practices](#best-practices)

---

## Overview

Tools define their parameters using [JSON Schema](https://json-schema.org/), which:

- ✅ Describes expected parameter structure
- ✅ Provides validation rules
- ✅ Gives the LLM parameter documentation
- ✅ Enables type checking

The LLM reads the schema to understand how to call your tool correctly.

---

## Basic Schema Structure

Every tool schema follows this structure:

```go
func (t *MyTool) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            // Parameter definitions
        },
        "required": []string{
            // Required parameter names
        },
        "additionalProperties": false, // Optional: reject unknown parameters
    }
}
```

### Minimal Example

```go
func (t *HelloTool) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "name": map[string]interface{}{
                "type":        "string",
                "description": "Name to greet",
            },
        },
        "required": []string{"name"},
    }
}
```

**LLM sees:**
```json
{
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "description": "Name to greet"
    }
  },
  "required": ["name"]
}
```

**LLM calls:**
```xml
<tool>
{
  "server_name": "local",
  "tool_name": "hello",
  "arguments": {
    "name": "Alice"
  }
}
</tool>
```

---

## Data Types

### String

```go
"param_name": map[string]interface{}{
    "type":        "string",
    "description": "Description for LLM",
}
```

**With constraints:**

```go
"email": map[string]interface{}{
    "type":        "string",
    "description": "User email address",
    "format":      "email",     // Optional: email format
    "minLength":   1,            // Optional: minimum length
    "maxLength":   100,          // Optional: maximum length
    "pattern":     "^[a-z]+$",  // Optional: regex pattern
}
```

**Common formats:**
- `"email"` - Email address
- `"uri"` - URI
- `"date"` - ISO 8601 date
- `"date-time"` - ISO 8601 datetime
- `"uuid"` - UUID

---

### Number

Integer or floating-point number:

```go
"count": map[string]interface{}{
    "type":        "number",
    "description": "Number of items",
    "minimum":     0,      // Optional: minimum value
    "maximum":     100,    // Optional: maximum value
}
```

**Integer only:**

```go
"age": map[string]interface{}{
    "type":        "integer",
    "description": "Age in years",
    "minimum":     0,
    "maximum":     150,
}
```

---

### Boolean

```go
"enabled": map[string]interface{}{
    "type":        "boolean",
    "description": "Whether feature is enabled",
    "default":     false,  // Optional: default value
}
```

---

### Enum (Choice)

Restrict to specific values:

```go
"operation": map[string]interface{}{
    "type":        "string",
    "description": "Arithmetic operation to perform",
    "enum":        []string{"add", "subtract", "multiply", "divide"},
}
```

**With descriptions:**

```go
"priority": map[string]interface{}{
    "type":        "string",
    "description": "Task priority level",
    "enum":        []string{"low", "medium", "high", "urgent"},
    "default":     "medium",
}
```

---

### Array

List of items:

```go
"tags": map[string]interface{}{
    "type":        "array",
    "description": "List of tags",
    "items": map[string]interface{}{
        "type": "string",
    },
    "minItems":    0,     // Optional: minimum items
    "maxItems":    10,    // Optional: maximum items
    "uniqueItems": true,  // Optional: enforce uniqueness
}
```

**Array of numbers:**

```go
"numbers": map[string]interface{}{
    "type":        "array",
    "description": "Numbers to process",
    "items": map[string]interface{}{
        "type": "number",
    },
}
```

**Array of objects:**

```go
"coordinates": map[string]interface{}{
    "type":        "array",
    "description": "List of coordinates",
    "items": map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "x": map[string]interface{}{"type": "number"},
            "y": map[string]interface{}{"type": "number"},
        },
        "required": []string{"x", "y"},
    },
}
```

---

### Object

Nested structure:

```go
"location": map[string]interface{}{
    "type":        "object",
    "description": "Geographic location",
    "properties": map[string]interface{}{
        "lat": map[string]interface{}{
            "type":        "number",
            "description": "Latitude",
            "minimum":     -90,
            "maximum":     90,
        },
        "lng": map[string]interface{}{
            "type":        "number",
            "description": "Longitude",
            "minimum":     -180,
            "maximum":     180,
        },
    },
    "required": []string{"lat", "lng"},
}
```

---

## Validation Keywords

### Required Parameters

```go
return map[string]interface{}{
    "type": "object",
    "properties": map[string]interface{}{
        "required_param": map[string]interface{}{
            "type": "string",
        },
        "optional_param": map[string]interface{}{
            "type": "string",
        },
    },
    "required": []string{"required_param"}, // Only required_param is mandatory
}
```

### Default Values

```go
"timeout": map[string]interface{}{
    "type":        "integer",
    "description": "Timeout in seconds",
    "default":     30,
}
```

### Minimum/Maximum

```go
// Numbers
"score": map[string]interface{}{
    "type":    "number",
    "minimum": 0,
    "maximum": 100,
}

// Strings
"username": map[string]interface{}{
    "type":      "string",
    "minLength": 3,
    "maxLength": 20,
}

// Arrays
"items": map[string]interface{}{
    "type":     "array",
    "minItems": 1,
    "maxItems": 10,
    "items":    map[string]interface{}{"type": "string"},
}
```

### Pattern Matching

```go
"phone": map[string]interface{}{
    "type":        "string",
    "description": "Phone number in format XXX-XXX-XXXX",
    "pattern":     "^\\d{3}-\\d{3}-\\d{4}$",
}

"slug": map[string]interface{}{
    "type":        "string",
    "description": "URL-friendly slug",
    "pattern":     "^[a-z0-9-]+$",
}
```

---

## Common Patterns

### File Path

```go
"file_path": map[string]interface{}{
    "type":        "string",
    "description": "Path to file",
    "pattern":     "^[^<>:\"|?*]+$", // No invalid filename characters
}
```

### URL

```go
"url": map[string]interface{}{
    "type":        "string",
    "description": "Web URL",
    "format":      "uri",
    "pattern":     "^https?://",
}
```

### Email

```go
"email": map[string]interface{}{
    "type":        "string",
    "description": "Email address",
    "format":      "email",
}
```

### Date/Time

```go
"date": map[string]interface{}{
    "type":        "string",
    "description": "Date in YYYY-MM-DD format",
    "format":      "date",
    "pattern":     "^\\d{4}-\\d{2}-\\d{2}$",
}

"datetime": map[string]interface{}{
    "type":        "string",
    "description": "ISO 8601 datetime",
    "format":      "date-time",
}
```

### Key-Value Pairs

```go
"metadata": map[string]interface{}{
    "type":        "object",
    "description": "Arbitrary key-value metadata",
    "additionalProperties": map[string]interface{}{
        "type": "string",
    },
}
```

### OneOf (Union Types)

```go
"value": map[string]interface{}{
    "description": "String or number value",
    "oneOf": []map[string]interface{}{
        {"type": "string"},
        {"type": "number"},
    },
}
```

---

## Examples

### Calculator Tool

```go
func (c *Calculator) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "operation": map[string]interface{}{
                "type":        "string",
                "description": "Mathematical operation",
                "enum":        []string{"add", "subtract", "multiply", "divide"},
            },
            "a": map[string]interface{}{
                "type":        "number",
                "description": "First operand",
            },
            "b": map[string]interface{}{
                "type":        "number",
                "description": "Second operand",
            },
        },
        "required": []string{"operation", "a", "b"},
    }
}
```

### File Search Tool

```go
func (f *FileSearch) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "directory": map[string]interface{}{
                "type":        "string",
                "description": "Directory to search in",
            },
            "pattern": map[string]interface{}{
                "type":        "string",
                "description": "Search pattern (glob or regex)",
            },
            "recursive": map[string]interface{}{
                "type":        "boolean",
                "description": "Whether to search recursively",
                "default":     false,
            },
            "file_type": map[string]interface{}{
                "type":        "string",
                "description": "Filter by file type",
                "enum":        []string{"any", "file", "directory"},
                "default":     "any",
            },
        },
        "required": []string{"directory", "pattern"},
    }
}
```

### API Request Tool

```go
func (a *APIRequest) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "method": map[string]interface{}{
                "type":        "string",
                "description": "HTTP method",
                "enum":        []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
            },
            "url": map[string]interface{}{
                "type":        "string",
                "description": "Request URL",
                "format":      "uri",
            },
            "headers": map[string]interface{}{
                "type":        "object",
                "description": "HTTP headers",
                "additionalProperties": map[string]interface{}{
                    "type": "string",
                },
            },
            "body": map[string]interface{}{
                "type":        "string",
                "description": "Request body (JSON string)",
            },
            "timeout": map[string]interface{}{
                "type":        "integer",
                "description": "Timeout in seconds",
                "minimum":     1,
                "maximum":     300,
                "default":     30,
            },
        },
        "required": []string{"method", "url"},
    }
}
```

### Database Query Tool

```go
func (d *DatabaseQuery) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "query": map[string]interface{}{
                "type":        "string",
                "description": "SQL query to execute",
            },
            "parameters": map[string]interface{}{
                "type":        "array",
                "description": "Query parameters (prevents SQL injection)",
                "items": map[string]interface{}{
                    "oneOf": []map[string]interface{}{
                        {"type": "string"},
                        {"type": "number"},
                        {"type": "boolean"},
                    },
                },
            },
            "max_rows": map[string]interface{}{
                "type":        "integer",
                "description": "Maximum rows to return",
                "minimum":     1,
                "maximum":     1000,
                "default":     100,
            },
        },
        "required": []string{"query"},
    }
}
```

### Email Sender Tool

```go
func (e *EmailSender) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "to": map[string]interface{}{
                "type":        "array",
                "description": "Recipient email addresses",
                "items": map[string]interface{}{
                    "type":   "string",
                    "format": "email",
                },
                "minItems": 1,
            },
            "subject": map[string]interface{}{
                "type":        "string",
                "description": "Email subject",
                "minLength":   1,
                "maxLength":   200,
            },
            "body": map[string]interface{}{
                "type":        "string",
                "description": "Email body (plain text or HTML)",
            },
            "cc": map[string]interface{}{
                "type":        "array",
                "description": "CC recipients",
                "items": map[string]interface{}{
                    "type":   "string",
                    "format": "email",
                },
            },
            "attachments": map[string]interface{}{
                "type":        "array",
                "description": "File paths to attach",
                "items": map[string]interface{}{
                    "type": "string",
                },
            },
        },
        "required": []string{"to", "subject", "body"},
    }
}
```

---

## Best Practices

### 1. Always Include Descriptions

```go
// ✅ Good: Clear description
"query": map[string]interface{}{
    "type":        "string",
    "description": "Search query to find relevant documents",
}

// ❌ Bad: No description
"query": map[string]interface{}{
    "type": "string",
}
```

The LLM reads descriptions to understand how to use parameters.

### 2. Use Enums for Fixed Choices

```go
// ✅ Good: Enum constrains values
"priority": map[string]interface{}{
    "type": "string",
    "enum": []string{"low", "medium", "high"},
}

// ❌ Bad: Free text (LLM might use invalid values)
"priority": map[string]interface{}{
    "type":        "string",
    "description": "Must be low, medium, or high",
}
```

### 3. Set Reasonable Limits

```go
// ✅ Good: Prevents abuse
"count": map[string]interface{}{
    "type":    "integer",
    "minimum": 1,
    "maximum": 100,
}

// ❌ Bad: No limits (could request millions)
"count": map[string]interface{}{
    "type": "integer",
}
```

### 4. Provide Defaults

```go
// ✅ Good: Sensible default
"timeout": map[string]interface{}{
    "type":    "integer",
    "default": 30,
}

// ❌ Bad: No default (LLM must always specify)
"timeout": map[string]interface{}{
    "type": "integer",
}
```

### 5. Use Specific Types

```go
// ✅ Good: Specific integer type
"age": map[string]interface{}{
    "type": "integer",
}

// ❌ Bad: Generic number (allows decimals)
"age": map[string]interface{}{
    "type": "number",
}
```

### 6. Validate in Execute

Schema validation is guidance for LLM, not enforcement:

```go
func (t *MyTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    // Always validate in code
    value, ok := args["param"].(string)
    if !ok {
        return "", fmt.Errorf("param must be string")
    }
    
    if value == "" {
        return "", fmt.Errorf("param cannot be empty")
    }
    
    // Proceed with execution
}
```

### 7. Document Units

```go
// ✅ Good: Clear units
"timeout": map[string]interface{}{
    "type":        "integer",
    "description": "Timeout in seconds",
}

"size": map[string]interface{}{
    "type":        "integer",
    "description": "File size in bytes",
}

// ❌ Bad: Ambiguous
"timeout": map[string]interface{}{
    "type":        "integer",
    "description": "How long to wait",
}
```

### 8. Group Related Parameters

```go
// ✅ Good: Grouped in object
"location": map[string]interface{}{
    "type": "object",
    "properties": map[string]interface{}{
        "lat": map[string]interface{}{"type": "number"},
        "lng": map[string]interface{}{"type": "number"},
    },
}

// ❌ Bad: Separate parameters
"lat": map[string]interface{}{"type": "number"},
"lng": map[string]interface{}{"type": "number"},
```

---

## Testing Your Schema

### Example Test

```go
func TestToolParameters(t *testing.T) {
    tool := NewMyTool()
    params := tool.Parameters()
    
    // Verify structure
    assert.Equal(t, "object", params["type"])
    
    properties := params["properties"].(map[string]interface{})
    assert.Contains(t, properties, "required_param")
    
    required := params["required"].([]string)
    assert.Contains(t, required, "required_param")
}
```

### Manual Testing

Use the schema with real LLM:

```go
agent, _ := core.NewAgent(provider, memory, []core.Tool{
    NewMyTool(),
})

// Try various prompts
// "Use my_tool with param=value"
// "Call my_tool" (should ask for required params)
```

---

## See Also

- [JSON Schema Specification](https://json-schema.org/) - Full specification
- [API Reference](api-reference.md) - Tool interface documentation
- [How-To: Create Custom Tool](../how-to/create-custom-tool.md) - Practical guide
- [Examples](../examples/) - Complete tool examples