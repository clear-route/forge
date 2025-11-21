# Feature Idea: Language Server Protocol (LSP) Integration

**Status:** Draft  
**Priority:** High Impact, Near-Term  
**Last Updated:** November 2025

---

## Overview

Integrate Language Server Protocol to give Forge real-time feedback about code correctness, type information, compilation errors, and IDE-level intelligence. This would transform Forge from a text-manipulation tool to a code-aware assistant that understands semantics, types, and project structure.

---

## Problem Statement

Currently, Forge operates on code as text files without understanding:
- Whether code will compile
- Type correctness
- Symbol definitions and references
- Import requirements
- Code structure and semantics

This leads to:
- Generated code that doesn't compile
- Type errors discovered only after writing
- Manual verification of code correctness
- Inability to do safe refactoring
- Missing imports or incorrect signatures

---

## Key Capabilities

### Code Validation
- Check syntax errors before writing files
- Validate type correctness in real-time
- Detect compilation errors immediately
- Get warnings and diagnostics
- Understand code semantics beyond text matching

### Intelligent Navigation
- Go to definition across files
- Find all references to symbols
- Show type information on hover
- Display function signatures and documentation
- Navigate type hierarchies

### Smart Code Operations
- Rename symbols project-wide with LSP guarantee
- Auto-import missing dependencies
- Organize imports intelligently
- Format code using language-specific rules
- Apply quick fixes and code actions

### Contextual Understanding
- Know what types are available in scope
- Understand function signatures when calling
- Access documentation inline
- See method completions for objects
- Understand type constraints and generics

---

## Example Use Cases

### 1. Type-Aware Code Generation

```
User: Add a method to UserService that fetches users by email

Agent: I'll add the method with proper types:
[Uses LSP to check UserService structure]
[Verifies User type definition]
[Checks database interface methods]
[Generates type-correct implementation]

Generated code guaranteed to compile!
```

### 2. Error Prevention

```
User: Create a function to process payment data

Agent: Creating function...
[LSP detects: PaymentData type doesn't exist]

Agent: I notice PaymentData isn't defined. Should I:
1. Create a new type
2. Use existing Payment type
3. Import from payment package
```

### 3. Refactoring with Confidence

```
User: Rename calculateTotal to computeTotal

Agent: Using LSP to find all references...
Found 47 usages across 12 files
[Shows impact analysis]
[Performs safe rename across entire codebase]

All references updated, no compilation errors
```

### 4. Import Management

```
Agent: I'm adding code that uses the uuid library
[LSP detects missing import]

Agent: Adding import for github.com/google/uuid
[Organizes imports automatically]

Code ready with proper imports!
```

---

## Technical Approach

### LSP Server Management

**Auto-detect and Start Servers:**
- Detect language from file extensions
- Start appropriate LSP servers (gopls, typescript-language-server, etc.)
- Manage server lifecycle (start, stop, restart)
- Handle multiple servers for polyglot projects
- Cache server instances for performance

**Server Support (Initial):**
- Go → gopls
- TypeScript/JavaScript → typescript-language-server
- Python → pylsp, pyright
- Rust → rust-analyzer
- Java → jdtls

### Integration Points

**Before Writing Files:**
- Query LSP for diagnostics
- Check for compilation errors
- Validate types and imports
- Show warnings to user

**During Code Generation:**
- Get type information for context
- Query available symbols in scope
- Check function signatures
- Validate method calls

**For Refactoring:**
- Use LSP rename for safety
- Find all references
- Apply code actions
- Organize imports

**For Navigation:**
- Go to definition
- Find implementations
- Show call hierarchy
- Navigate symbols

### User Experience

**Automatic Setup:**
- No manual LSP configuration required
- Auto-detect and install servers
- Workspace-specific initialization
- Fallback when LSP unavailable

**Visual Feedback:**
- Show diagnostics in diff previews
- Display type info in agent reasoning
- Warning prompts for errors before approval
- Quick fix suggestions in chat
- Error highlighting in TUI

**Graceful Degradation:**
- Work without LSP if unavailable
- Clear messaging when LSP fails
- Fallback to text-based operations
- Optional LSP usage (can disable)

---

## Value Propositions

### For All Users
- Generate correct code the first time
- Catch errors before they're written
- Understand existing code better
- Navigate large codebases easily
- Build confidence in AI-generated code

### For Type-Safe Language Users (Go, TypeScript, Rust)
- Full type checking before writes
- Type-aware refactoring
- Automatic import management
- Signature help during generation
- Zero compilation errors from agent

### For Quality-Focused Developers
- Lint-compliant code generation
- Proper formatting automatically
- Best practices enforcement
- Validation before approval
- Professional-grade output

---

## Implementation Phases

### Phase 1: Basic Integration (2-3 weeks)
- Start/stop LSP servers for Go
- Get diagnostics after edits
- Show errors in TUI
- Basic error prevention

**Deliverables:**
- LSP server lifecycle management
- Diagnostic display in diffs
- Error warnings before approval
- Go language support

### Phase 2: Code Intelligence (3-4 weeks)
- Type-aware code generation
- Auto-import management
- Symbol completion
- Hover information
- Go-to-definition

**Deliverables:**
- Type information in agent context
- Automatic import insertion
- Symbol navigation tools
- Enhanced code generation

### Phase 3: Advanced Features (4-5 weeks)
- Project-wide refactoring
- Code actions and quick fixes
- Workspace symbols search
- Call hierarchy
- Format on save

**Deliverables:**
- Safe rename operations
- Quick fix application
- Workspace-wide search
- Code action support

### Phase 4: Multi-Language (Ongoing)
- TypeScript/JavaScript support
- Python support
- Rust support
- Language-specific optimizations
- Custom LSP server configs

**Deliverables:**
- 3-5 language support
- Polyglot project handling
- Performance optimizations
- Configuration options

---

## Open Questions

1. **Performance:** How to handle LSP overhead for large projects?
   - Cache responses aggressively?
   - Lazy-load servers per file type?
   - Set timeouts for LSP queries?

2. **Configuration:** Auto-detect vs manual LSP server config?
   - Start with auto-detect, allow overrides?
   - Provide server installation helpers?
   - Support custom server binaries?

3. **Multi-Language:** How to prioritize server support?
   - Start with Go (our implementation language)?
   - Add TypeScript next (web projects)?
   - Community vote on priorities?

4. **Errors:** How aggressive about blocking on LSP errors?
   - Block on compilation errors?
   - Warn on linting issues?
   - Allow override for all errors?

5. **Offline:** Graceful degradation when LSP unavailable?
   - Clear error messages?
   - Fallback to text operations?
   - Cache previous LSP data?

6. **Resource Usage:** How to manage memory/CPU?
   - Limit concurrent LSP servers?
   - Stop idle servers?
   - Monitor resource consumption?

---

## Related Features

**Synergies with:**
- **Code Intelligence & Navigation** - LSP provides foundation
- **Testing & Quality Automation** - LSP validates test code
- **Documentation & Knowledge Base** - LSP extracts signatures
- **Performance & Optimization** - Cache LSP responses
- **Smart Refactoring** - LSP ensures safety

**Dependencies:**
- None - can be standalone feature
- Enhanced by caching system
- Better with multi-workspace support

---

## Success Metrics

**Adoption:**
- 80%+ of users have LSP enabled
- Top 3 languages have LSP support
- <5% LSP failure rate

**Quality:**
- 95%+ of generated code compiles first try
- 90%+ reduction in type errors
- 50%+ faster refactoring operations

**Performance:**
- LSP queries complete in <500ms
- No noticeable latency impact
- Memory usage under 100MB per server

**User Satisfaction:**
- 4.5+ rating for LSP feature
- "Game changer" feedback
- Cited as key differentiator

---

## Risks and Mitigations

### Risk: LSP Server Reliability
**Impact:** Servers crash or become unresponsive  
**Mitigation:**
- Auto-restart on failure
- Timeout all queries
- Graceful fallback to text mode
- Clear error messages

### Risk: Performance Impact
**Impact:** LSP queries slow down agent  
**Mitigation:**
- Aggressive caching
- Async queries where possible
- Timeout limits
- Resource monitoring

### Risk: Configuration Complexity
**Impact:** Users struggle with setup  
**Mitigation:**
- Auto-detect and install servers
- Zero-config for common languages
- Clear setup guides
- Helpful error messages

### Risk: Multi-Language Complexity
**Impact:** Supporting many languages is difficult  
**Mitigation:**
- Start with 1-2 languages
- Incremental rollout
- Community contributions
- Standardized interface

---

## Next Steps

1. **Prototype** - Build basic Go LSP integration (1 week)
2. **User Testing** - Test with 5-10 users, gather feedback
3. **Refinement** - Iterate based on feedback
4. **Full Implementation** - Phase 1 development
5. **Documentation** - Write setup guides and best practices
6. **Launch** - Release as beta feature

---

## References

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [gopls Documentation](https://github.com/golang/tools/tree/master/gopls)
- [typescript-language-server](https://github.com/typescript-language-server/typescript-language-server)
- [LSP Implementations List](https://langserver.org/)
