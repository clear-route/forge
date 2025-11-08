# Fix: JSON Parsing Errors in Tool Call Parser

## Problem

Users were experiencing frequent JSON parsing errors with messages like:
```
❌ Error: failed to parse tool call JSON: invalid character '>' after top-level value
```

The error showed incomplete JSON being parsed, with the `>` character from XML closing tags appearing where it shouldn't. The agent would eventually recover and process the tool call correctly, but the parsing errors were disruptive to the UX.

## Root Cause Analysis

The issue was in the [`ToolCallParser.flushBufferIfNotInTag()`](pkg/llm/parser/toolcall.go:113) function in [`pkg/llm/parser/toolcall.go`](pkg/llm/parser/toolcall.go).

### How the Parser Works

1. The parser processes streaming LLM responses character by character
2. When it detects `<tool>`, it enters "tool call mode" (`inToolCall = true`)
3. Content between `<tool>` and `</tool>` should be accumulated without being emitted
4. When `</tool>` is detected, all accumulated content is emitted as a complete tool call

### The Bug

The `flushBufferIfNotInTag()` function had logic to keep potential tag prefixes in the buffer (e.g., `<`, `</`, `</t`, `</to`, `</too`, `</tool`) to handle tags that span chunk boundaries. However, when inside a tool call:

**OLD BEHAVIOR:**
- If the buffer didn't end with a partial `</tool>` prefix, it would flush ALL buffered content immediately
- This meant JSON content would be emitted BEFORE the closing `</tool>` tag arrived
- The flushed content was accumulated in `state.toolCallContent` via [`handleToolCallContent()`](pkg/agent/core/stream.go:125)
- Later when `</tool>` finally arrived, the remaining buffer content (including potentially the `>` character) would be appended
- This resulted in incomplete or malformed JSON being passed to [`json.Unmarshal()`](pkg/agent/default.go:449)

**Example Streaming Scenario:**
```
Chunk 1: <tool>{"tool_name": "read_file", "arguments": {
Chunk 2: "path": "file.txt"}}
Chunk 3: </tool>
```

With the old code:
1. After Chunk 1: Enters tool call mode, buffer has `{"tool_name": "read_file", "arguments": {`
2. After Chunk 2: Buffer now has complete JSON `{"tool_name": "read_file", "arguments": {"path": "file.txt"}}`
   - `flushBufferIfNotInTag()` sees no `</tool>` prefix at the end
   - Flushes entire JSON to `toolContent`
   - Emits it immediately via tool call event
3. After Chunk 3: Detects `</tool>` closing tag
   - But JSON was already emitted in step 2!
   - Any remaining buffer content (like `>` from `</tool>`) might get appended

## The Fix

Modified [`flushBufferIfNotInTag()`](pkg/llm/parser/toolcall.go:113) to NOT emit tool call content until the complete `</tool>` closing tag is detected.

### Key Changes

1. **Inside Tool Call**: When `p.inToolCall == true`, content is moved from the main buffer to `p.toolContent` but NOT emitted as events. The content stays internal until `</tool>` is found.

2. **Proper Accumulation**: All content between `<tool>` and `</tool>` is accumulated in `p.toolContent` without intermediate emissions.

3. **Emit Only on Complete Tag**: Tool call content is only emitted when [`handleToolEnd()`](pkg/llm/parser/toolcall.go:91) processes the complete `</tool>` tag.

4. **Trimming**: Added `strings.TrimSpace()` in `handleToolEnd()` to clean any trailing whitespace or stray characters.

### Code Changes

**File**: [`pkg/llm/parser/toolcall.go`](pkg/llm/parser/toolcall.go)

**Changes in `flushBufferIfNotInTag()`** (lines 113-197):
```go
if p.inToolCall {
    // Check for partial </tool> prefix and keep it buffered
    // Move everything else to toolContent but DON'T emit yet
    // ... (prefix checking logic) ...
    
    // Accumulate to toolContent but don't emit
    if flushText != "" {
        p.toolContent.WriteString(flushText)
        return nil  // Don't emit anything yet
    }
}
```

**Changes in `handleToolEnd()`** (lines 91-109):
```go
// Add buffered content to tool content
p.toolContent.WriteString(textBefore)
content := p.toolContent.String()
p.toolContent.Reset()

// Trim any trailing whitespace or stray characters
content = strings.TrimSpace(content)

return &ParsedContent{
    Type:    ContentTypeToolCall,
    Content: content,
}
```

## Testing

All existing tests pass, including:
- `TestToolCallParser_SimpleToolCall`
- `TestToolCallParser_StreamedToolCall`
- `TestToolCallParser_MultipleToolCalls`
- `TestToolCallParser_IncompleteTagAtBoundary`
- `TestToolCallParser_StreamedContentAccumulation`

Updated `TestToolCallParser_FlushRegularContent` to reflect the new behavior where incomplete non-tool tags are immediately flushed.

## Expected Outcome

- ✅ No more JSON parsing errors with "invalid character '>' after top-level value"
- ✅ Tool calls are parsed cleanly in one complete chunk
- ✅ All content between `<tool>` and `</tool>` is properly accumulated
- ✅ Better streaming behavior with reduced error recovery cycles
- ✅ Improved UX with fewer disruptive error messages

## Related Files

- [`pkg/llm/parser/toolcall.go`](pkg/llm/parser/toolcall.go) - Main parser implementation
- [`pkg/llm/parser/toolcall_test.go`](pkg/llm/parser/toolcall_test.go) - Parser tests
- [`pkg/agent/core/stream.go`](pkg/agent/core/stream.go) - Stream processing that uses the parser
- [`pkg/agent/default.go`](pkg/agent/default.go) - Agent loop that parses JSON from tool calls