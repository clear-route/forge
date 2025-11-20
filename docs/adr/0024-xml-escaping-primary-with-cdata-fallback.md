# 24. XML Entity Escaping as Primary Method with CDATA Fallback

**Status:** Accepted
**Date:** 2024-01-20
**Deciders:** Forge Core Team
**Technical Story:** Simplify XML tool call format by making entity escaping the primary method, with CDATA as fallback option
**Supersedes:** ADR-0019 (revises approach based on practical usage)

---

## Context

ADR-0019 established pure XML with CDATA for tool call arguments, mandating CDATA for ALL code/content fields. While this solved JSON parsing issues, it introduced unnecessary verbosity and complexity.

### Background

Modern LLMs are extensively trained on HTML and XML, which use entity escaping (`&amp;`, `&lt;`, `&gt;`, etc.) as the standard way to handle special characters. CDATA sections, while useful, are:
- Less common in LLM training data
- More verbose (adding `<![CDATA[` and `]]>` markers)
- Unnecessary for most content that only contains a few special characters

### Problem Statement

The current mandatory CDATA approach has these issues:

1. **Verbosity**: Every code field requires `<![CDATA[...]]>` even for simple content like `a && b`
2. **Unnecessary Complexity**: LLMs already know XML entity escaping from HTML/XML training
3. **Inconsistent with Web Standards**: HTML/XML typically use entity escaping, not CDATA
4. **Learning Burden**: Forces agents to use CDATA even when escaping would be simpler

**Example of current verbose approach:**
```xml
<search><![CDATA[const x = a && b]]></search>
```

**Simpler alternative using escaping:**
```xml
<search>const x = a &amp;&amp; b</search>
```

### Goals

- Make XML entity escaping the primary/preferred method for handling special characters
- Keep CDATA as a valid fallback option for complex content or when escaping fails
- Reduce verbosity in tool calls
- Leverage LLMs' existing knowledge of HTML/XML entity escaping
- Maintain full backwards compatibility with existing CDATA usage

### Non-Goals

- Deprecate or remove CDATA support
- Force migration of existing tool calls
- Add complex error tracking or retry logic
- Change the XML parser implementation

---

## Decision Drivers

* **LLM Training Data**: XML entity escaping is extensively represented in HTML/XML training data
* **Simplicity**: Most content only needs a few characters escaped (`&`, `<`, `>`)
* **Verbosity**: Entity escaping is less verbose than CDATA for typical use cases
* **Backwards Compatibility**: Must maintain support for existing CDATA-based tool calls
* **Implementation Simplicity**: Should not require complex error tracking or state management
* **Agent Intelligence**: Modern LLMs can recognize errors and try alternatives without explicit tracking

---

## Considered Options

### Option 1: Mandatory XML Entity Escaping (Remove CDATA)

**Description:** Make entity escaping mandatory and remove CDATA support entirely.

**Pros:**
- Simplest approach
- Most consistent
- Smallest prompt

**Cons:**
- Breaking change for existing usage
- Difficult for very large content blocks
- Edge cases with `]]>` in content become harder
- Limits flexibility

### Option 2: Keep Mandatory CDATA (Status Quo)

**Description:** Continue requiring CDATA for all code/content fields as per ADR-0019.

**Pros:**
- No changes needed
- Already documented and understood
- Proven to work

**Cons:**
- Verbose for simple content
- Doesn't leverage LLM training on XML entities
- Forces unnecessary markers on simple content

### Option 3: XML Escaping Primary, CDATA Fallback (CHOSEN)

**Description:** Teach XML entity escaping as the primary method in system prompts, but keep CDATA as a valid fallback option when escaping fails or is complex.

**Pros:**
- Less verbose for typical use cases
- Leverages LLM knowledge of HTML/XML
- Maintains backwards compatibility
- Gives agents flexibility to choose best method
- Simple to implement (just prompt changes)

**Cons:**
- Two valid ways to do the same thing
- Slightly more complex documentation
- Agents might inconsistently switch between methods

### Option 4: Complex Error Tracking with Automatic Fallback

**Description:** Track parsing errors and automatically suggest CDATA after N failures.

**Pros:**
- Intelligent fallback
- Could optimize based on patterns

**Cons:**
- Complex state management
- More code to maintain
- Unnecessary - LLMs can recognize errors themselves
- Over-engineering a simple problem

---

## Decision

**Chosen Option:** Option 3 - XML Escaping Primary, CDATA Fallback

### Rationale

This option strikes the best balance:

1. **Leverages LLM Training**: HTML/XML with entity escaping is abundant in training data
2. **Reduces Verbosity**: Simple content like `a && b` becomes `a &amp;&amp; b` instead of `<![CDATA[a && b]]>`
3. **Maintains Flexibility**: CDATA remains available for complex content or when escaping is problematic
4. **Backwards Compatible**: Existing CDATA-based tool calls continue to work
5. **Simple Implementation**: Only requires updating prompts and documentation, not parser code
6. **Agent Intelligence**: Modern LLMs can recognize parsing errors and try CDATA without explicit tracking

The parser already supports both methods via Go's native `xml.Unmarshal()` (handles CDATA) and `UnmarshalXMLWithFallback()` (handles escaping), so no code changes are needed.

---

## Consequences

### Positive

- **Reduced Verbosity**: Most tool calls will be 10-20% shorter
- **LLM-Friendly**: Aligns with extensive HTML/XML training data
- **Flexible**: Agents can choose the best method for their situation
- **Simple to Implement**: Only prompt and documentation changes needed
- **No Parser Changes**: Existing parser handles both methods perfectly
- **Backwards Compatible**: All existing CDATA tool calls continue to work

### Negative

- **Two Valid Methods**: Documentation must cover both approaches
- **Potential Inconsistency**: Agents might switch between methods unpredictably
- **Learning Curve**: Need to teach when to use each method

### Neutral

- **Error Messages**: Will mention both escaping and CDATA as options
- **Examples**: Will show escaping first, CDATA as alternative
- **Parser Behavior**: Unchanged - already handles both methods

---

## Implementation

### System Prompt Changes (pkg/agent/prompts/static.go)

Replace mandatory CDATA section with:

```
PRIMARY METHOD - XML Entity Escaping:
✅ PREFERRED: Escape special characters using standard XML entities:
  - & becomes &amp;
  - < becomes &lt;
  - > becomes &gt;
  - " becomes &quot;
  - ' becomes &apos;

Examples:
<content>func example() *Config { return &amp;Config{} }</content>
<search>const x = a &amp;&amp; b || c</search>
<pattern>\.go$</pattern>

FALLBACK METHOD - CDATA Sections:
⚠️ Use CDATA if escaping is complex or causes parse errors:
<content><![CDATA[func example() *Config { return &Config{} }]]></content>

IMPORTANT CDATA LIMITATION:
The sequence ]]> terminates a CDATA section and cannot appear within CDATA content.
If your content contains ]]>, you must use entity escaping instead.

WHEN TO USE EACH METHOD:
- Entity Escaping (PREFERRED): Use for most content, especially when you have <5 special characters
- CDATA (FALLBACK): Use for large code blocks or when escaping becomes too complex
- NEVER MIX: Choose ONE method per field - either escape ALL special chars OR wrap in CDATA

Both methods are valid and fully supported by the parser.
```

### Error Recovery Changes (pkg/agent/prompts/error_recovery.go)

Update `buildParseError()` to mention both methods:

```
ERROR: Invalid XML in tool call.

Parse error: {error details}

SOLUTION 1 - Use XML entity escaping (RECOMMENDED FIRST):
Escape special characters: & → &amp;, < → &lt;, > → &gt;
Example: <content>code with &amp; &lt; &gt;</content>

SOLUTION 2 - Use CDATA as fallback for complex content:
Wrap content in CDATA markers (no escaping needed):
Example: <content><![CDATA[code with & < >]]></content>
NOTE: CDATA cannot contain ]]> - use escaping if your content has this sequence.

Both methods are supported. Try entity escaping first, use CDATA as fallback.
```

### XML Example Generator (pkg/agent/prompts/xml_example.go)

Modify `generateStringExample()` to:
- Use entity escaping for fields with typical content
- Use CDATA only for very large content blocks (>1000 characters)
- Add helper function to escape XML entities

### Migration Path

**For Users:**
- No breaking changes - existing CDATA usage continues to work
- New system prompts will teach escaping as primary method
- Agents will gradually adopt escaping for new tool calls
- Both methods remain equally valid

**For Developers:**
- Update system prompts to show escaping first
- Update error messages to mention both options
- Add tests for escaped XML content
- Update ADR-0019 with note referencing this ADR
- No parser code changes needed

### Timeline

Implementation is straightforward:
- System prompt updates: 1-2 hours
- Error recovery updates: 30 minutes  
- XML example generator: 1 hour
- Test additions: 2 hours
- Documentation updates: 1 hour

**Total: ~6 hours of focused work**

---

## Validation

### Success Metrics

- ✅ System prompts teach escaping as primary method
- ✅ Both escaped XML and CDATA parse correctly
- ✅ All existing CDATA-based tests continue passing
- ✅ New tests verify escaping works for all special characters
- ✅ Error messages mention both options clearly
- ✅ Tool calls become 10-20% less verbose on average

### Monitoring

- Monitor parse success rates (should remain >99.9%)
- Track which method agents use (for future optimization)
- Verify no regression in error handling
- Ensure backwards compatibility maintained

---

## Related Decisions

- [ADR-0019](0019-xml-cdata-tool-call-format.md) - Superseded by this ADR
- [ADR-0002](0002-xml-format-for-tool-calls.md) - Original XML format decision
- [ADR-0024](0024-fallback-xml-parsing-with-ampersand-escaping.md) - Parser fallback mechanism

---

## References

- [XML Entity References](https://www.w3.org/TR/xml/#sec-references)
- [HTML Character Entity References](https://html.spec.whatwg.org/multipage/named-characters.html)
- [CDATA Sections in XML](https://www.w3.org/TR/xml/#sec-cdata-sect)
- Go `encoding/xml` package documentation

---

## Notes

**Key Insight:** LLMs are extensively trained on HTML/XML which predominantly uses entity escaping, not CDATA. By teaching the method they already know, we reduce cognitive load and verbosity while maintaining flexibility.

**Parser Already Ready:** The existing parser via `xml.Unmarshal()` natively handles both CDATA and escaped entities, so this is purely a prompt/documentation change with zero parser modifications required.

**Agent Choice:** We trust the agent's intelligence to recognize when escaping fails and try CDATA instead, rather than building complex error tracking logic. Modern LLMs are smart enough to adapt based on error messages.

**Last Updated:** 2024-01-20