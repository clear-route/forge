# Message Format Reference

Complete reference for message formatting in Forge.

## Table of Contents

- [Message Structure](#message-structure)
- [Message Roles](#message-roles)
- [Tool Call Format](#tool-call-format)
- [Thinking Format](#thinking-format)
- [Token Estimation](#token-estimation)
- [Message Examples](#message-examples)

---

## Message Structure

### Core Message Type

```go
type Message struct {
    Role    string  // Message role (user/assistant/system/tool)
    Content string  // Message content
}
```

### Creating Messages

```go
import "github.com/yourusername/forge/pkg/message"

// User message
msg := message.User("Hello, how are you?")

// Assistant message
msg := message.Assistant("I'm doing well, thank you!")

// System message
msg := message.System("You are a helpful assistant")

// Tool result message
msg := message.Tool("calculator", "Result: 42")
```

---

## Message Roles

### System

Defines agent behavior and constraints.

```go
msg := message.System(`You are a Python expert assistant.
- Write clean, idiomatic Python code
- Include type hints
- Follow PEP 8 style guide
- Explain your reasoning`)
```

**Characteristics:**
- **Position:** First message in conversation
- **Count:** Usually one (can be multiple)
- **Preserved:** Never pruned from memory
- **Purpose:** Set agent personality and rules

**Example Usage:**

```go
memory.Add(message.System("You are a helpful coding assistant"))
memory.Add(message.User("Write a hello world function"))
```

---

### User

Messages from the human user.

```go
msg := message.User("What is the capital of France?")
```

**Characteristics:**
- **Source:** Human input
- **Format:** Plain text question or instruction
- **Frequency:** Every turn (user speaks)

**Multi-line:**

```go
msg := message.User(`Please help me with:
1. Understanding the code
2. Fixing the bug
3. Writing tests`)
```

---

### Assistant

Messages from the AI agent.

```go
msg := message.Assistant("The capital of France is Paris.")
```

**Characteristics:**
- **Source:** LLM response
- **Format:** May include thinking + tool calls + text
- **Frequency:** Every turn (agent responds)

**With Thinking:**

```go
msg := message.Assistant(`[Let me search for this information...]

Based on my search, the capital of France is Paris.`)
```

**With Tool Call:**

```go
msg := message.Assistant(`[I'll use the calculator for this]

<tool>
<server_name>local</server_name>
<tool_name>calculator</tool_name>
<arguments>
  <a>5</a>
  <b>3</b>
  <operation>add</operation>
</arguments>
</tool>`)
```

---

### Tool

Results from tool execution.

```go
msg := message.Tool("calculator", "Result: 8")
```

**Characteristics:**
- **Source:** Tool execution
- **Format:** `Tool: {name}\n{result}`
- **Purpose:** Provide tool output to agent

**Full Format:**

```
Tool: calculator
Result: 8
```

**Error Format:**

```
Tool: calculator
Error: division by zero
```

---

## Tool Call Format

### Pure XML Structure

Tools are called using pure XML format with CDATA sections for content:

```xml
<tool>
<server_name>local</server_name>
<tool_name>tool_name</tool_name>
<arguments>
  <param1>value1</param1>
  <param2>value2</param2>
</arguments>
</tool>
```

### Fields

- **server_name:** Always `local` for local tools, or MCP server name for remote tools
- **tool_name:** Exact tool name from `Tool.Name()`
- **arguments:** Nested XML elements for each parameter

### CDATA for Complex Content

Use CDATA sections for content that contains special characters, code, or multi-line text:

```xml
<tool>
<tool_name>write_to_file</tool_name>
<arguments>
  <path>example.go</path>
  <content>CDATA-START[package main

func main() {
	fmt.Println("Hello, World!")
}]CDATA-END</content>
</arguments>
</tool>
```

Note: Replace `CDATA-START[` with `<![CDATA[` and `]CDATA-END` with `]]>` in actual usage.

### Examples

**Simple Tool:**

```xml
<tool>
<server_name>local</server_name>
<tool_name>task_completion</tool_name>
<arguments>
  <result>The calculation result is 42</result>
</arguments>
</tool>
```

**Tool with CDATA Content:**

```xml
<tool>
<tool_name>apply_diff</tool_name>
<arguments>
  <path>app.py</path>
  <diff>CDATA-START[
SEARCH-START SEARCH
def calculate(a, b):
    return a + b
EQUALS
def calculate(a, b):
    """Add two numbers together."""
    return a + b
REPLACE-END REPLACE
]CDATA-END</diff>
</arguments>
</tool>
```

Note: Replace markers with actual diff syntax and CDATA tags.

**Nested Arguments:**

```xml
<tool>
<server_name>mcp-server</server_name>
<tool_name>use_mcp_tool</tool_name>
<arguments>
  <files>
    <file>
      <path>file1.txt</path>
      <content>CDATA-START[content 1]CDATA-END</content>
    </file>
    <file>
      <path>file2.txt</path>
      <content>CDATA-START[content 2]CDATA-END</content>
    </file>
  </files>
</arguments>
</tool>
```

**Multiple Tools:**

Agent can call multiple tools in sequence:

```xml
[I'll need to use multiple tools for this]

<tool>
<server_name>local</server_name>
<tool_name>search</tool_name>
<arguments>
  <query>Python tutorials</query>
</arguments>
</tool>

[Now I'll summarize the results]

<tool>
<tool_name>task_completion</tool_name>
<arguments>
  <result>Here are the top Python tutorials...</result>
</arguments>
</tool>
```

### Migration from JSON Format

The previous JSON-based format is deprecated but still supported:

```xml
<!-- OLD FORMAT (deprecated) -->
<tool>
{
  "server_name": "local",
  "tool_name": "task_completion",
  "arguments": {"result": "Done"}
}
</tool>

<!-- NEW FORMAT (recommended) -->
<tool>
<server_name>local</server_name>
<tool_name>task_completion</tool_name>
<arguments>
  <result>Done</result>
</arguments>
</tool>
```

**Benefits of Pure XML:**
- No JSON escaping issues in complex content
- CDATA sections preserve exact content (code, diffs, etc.)
- Cleaner syntax for LLMs to generate
- More reliable parsing, especially for large files

---

## Thinking Format

[Content continues as before...]