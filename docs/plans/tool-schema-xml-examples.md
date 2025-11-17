# Diagnosis: Tool Schema XML Examples Issue

## Problem Summary

The current tool schema generation system shows JSON schemas but provides only generic XML examples like:
```xml
<tool>
<server_name>local</server_name>
<tool_name>tool_name</tool_name>
<arguments>
  ...  <!-- ❌ NOT HELPFUL -->
</arguments>
</tool>
```

This causes LLMs to:
1. **Guess at XML structure** for complex nested types (arrays, objects)
2. **Misuse CDATA** - wrapping structures instead of content
3. **Generate invalid XML** - unescaped `&` characters, wrong entity refs

## Evidence from User's Screenshots

**Error 1**: XML parsing failure
```
XML syntax error on line 7: invalid character entity &SlashCommand (no semicolon)
```
→ LLM generated `&SlashCommand` without escaping or CDATA

**Error 2**: Type mismatch
```
json: cannot unmarshal string into Go struct field .edits of type []struct { Search string; Replace string }
```
→ LLM generated `<edits><![CDATA[...]]></edits>` instead of nested `<edit>` elements

**LLM's Own Analysis** (from screenshot):
- "I need to use CDATA sections for the apply_diff content since it contains special characters like '&'"
- "The edits parameter should be an array of objects with 'search' and 'replace' keys"  
- "I'm passing the edits as a CDATA string when it should be structured XML elements"

## Root Cause Analysis

### Current Schema System

1. **Tool Interface** ([`pkg/agent/tools/tool.go`](../../pkg/agent/tools/tool.go)):
   ```go
   type Tool interface {
       Name() string
       Description() string
       Execute(args json.RawMessage) (string, error)
       Schema() map[string]interface{}  // ← JSON Schema only
       IsLoopBreaking() bool
   }
   ```

2. **Schema Formatter** ([`pkg/agent/prompts/formatter.go:72`](../../pkg/agent/prompts/formatter.go#L72)):
   ```go
   builder.WriteString(fmt.Sprintf("**Example:**\n```xml\n<tool>\n<server_name>local</server_name>\n<tool_name>%s</tool_name>\n<arguments>\n  ...\n</arguments>\n</tool>\n```\n\n",
       tool.Name()))
   ```
   ↑ **Generic, unhelpful example**

3. **No XML-Specific Metadata**: Tools don't provide:
   - Concrete XML examples
   - When to use CDATA
   - How to structure nested arrays/objects

## Proposed Solution Architecture

### Solution 1: Add XMLExample() Method to Tool Interface

**Pros:**
- Most explicit and clear
- Each tool provides its own concrete example
- Full control over example quality

**Cons:**
- Requires updating all tools
- More code to maintain
- Breaking change to Tool interface

```go
type Tool interface {
    Name() string
    Description() string  
    Execute(args json.RawMessage) (string, error)
    Schema() map[string]interface{}
    XMLExample() string  // ← NEW
    IsLoopBreaking() bool
}
```

### Solution 2: Auto-Generate XML Examples from JSON Schema

**Pros:**
- No interface changes needed
- Works for all tools automatically
- Single source of truth (JSON Schema)

**Cons:**
- Complex generation logic
- Harder to handle edge cases
- May not always generate optimal examples

```go
func GenerateXMLExample(schema map[string]interface{}, toolName string) string {
    // Parse JSON Schema
    // Generate appropriate XML structure
    // Add CDATA where needed
    // Handle nested arrays/objects
}
```

### Solution 3: Hybrid - Optional XMLExample() with Auto-Generation Fallback

**Pros:**
- Best of both worlds
- Tools can override auto-generation
- Graceful degradation

**Cons:**
- Most complex implementation
- Two code paths to maintain

```go
type Tool interface {
    Name() string
    Description() string
    Execute(args json.RawMessage) (string, error)
    Schema() map[string]interface{}
    IsLoopBreaking() bool
}

// Optional interface
type XMLExampleProvider interface {
    XMLExample() string
}

func FormatToolSchema(tool Tool) string {
    var example string
    if provider, ok := tool.(XMLExampleProvider); ok {
        example = provider.XMLExample()
    } else {
        example = GenerateXMLExample(tool.Schema(), tool.Name())
    }
    // ... use example
}
```

## Recommended Approach

**Solution 3 (Hybrid)** is best because:

1. **Immediate fix**: Can add `XMLExample()` to problematic tools (`apply_diff`) right now
2. **Future-proof**: Auto-generation handles simple tools
3. **Flexible**: Complex tools can provide custom examples
4. **Non-breaking**: Existing tools keep working

## Implementation Plan

### Phase 1: Core Infrastructure

1. **Create XML example generator** (`pkg/agent/prompts/xml_example.go`):
   ```go
   func GenerateXMLExample(schema map[string]interface{}, toolName string) string
   ```

2. **Update formatter** to use examples:
   ```go
   func FormatToolSchema(tool Tool) string {
       // Check for XMLExampleProvider
       // Fallback to auto-generation
       // Include in schema output
   }
   ```

3. **Add CDATA usage guidelines** to static prompts

### Phase 2: Tool-Specific Examples

4. **Create XMLExampleProvider interface**
5. **Implement for `apply_diff` tool** with correct nested structure
6. **Implement for other complex tools** as needed

### Phase 3: Validation & Testing

7. **Test auto-generation** with all existing tools
8. **Compare auto vs manual** examples for quality
9. **Add tests** for edge cases

## XML Example Generation Rules

The generator must handle:

1. **Simple types** → Direct XML elements
   ```xml
   <path>./file.go</path>
   <count>42</count>
   <enabled>true</enabled>
   ```

2. **Strings with special chars** → CDATA
   ```xml
   <content><![CDATA[code with & < > chars]]></content>
   ```

3. **Arrays** → Repeated elements or nested structure
   ```xml
   <exclude>vendor</exclude>
   <exclude>node_modules</exclude>
   ```
   OR
   ```xml
   <edits>
     <edit>
       <search><![CDATA[old]]></search>
       <replace><![CDATA[new]]></replace>
     </edit>
   </edits>
   ```

4. **Objects** → Nested XML elements
   ```xml
   <config>
     <name>prod</name>
     <port>8080</port>
   </config>
   ```

## CDATA Usage Rules (for prompts)

Add to system prompts:

```
CDATA Usage Rules:
1. Use CDATA for CONTENT/TEXT that contains special characters (&, <, >, quotes)
2. Use CDATA for code, diffs, file content, JSON, HTML, etc.
3. DO NOT use CDATA for STRUCTURE (arrays, objects)
4. DO NOT wrap entire argument sections in CDATA

✅ CORRECT - CDATA for content:
<search><![CDATA[func &doThing() { }]]></search>

❌ WRONG - CDATA for structure:
<edits><![CDATA[{ search: "...", replace: "..." }]]></edits>

✅ CORRECT - Nested XML for structure:
<edits>
  <edit>
    <search><![CDATA[old code]]></search>
    <replace><![CDATA[new code]]></replace>
  </edit>
</edits>
```

## Success Metrics

- ✅ No more XML entity errors (`&SlashCommand`)
- ✅ No more type mismatch errors (string→array)
- ✅ LLMs generate correct XML on first try >95% of time
- ✅ All tools have clear, unambiguous examples
- ✅ CDATA used correctly (content, not structure)

## Next Steps

1. **Confirm approach** with user
2. **Implement** hybrid XML example system
3. **Add** `apply_diff` concrete example
4. **Update** static prompts with CDATA rules
5. **Test** with actual LLM usage
6. **Monitor** parse success rates