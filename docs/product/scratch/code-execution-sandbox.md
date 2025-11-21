# Feature Idea: Code Execution Sandbox

**Status:** Draft  
**Priority:** High Impact, Medium-Term  
**Last Updated:** November 2025

---

## Overview

Safe, isolated environment for executing code snippets, running scripts, and testing code changes without affecting the host system. Enables interactive REPL-like workflows, quick prototyping, and safe experimentation with generated code before committing to files.

---

## Problem Statement

Developers need to test code before writing it to files:
- Generated code may have runtime errors
- Need to verify logic works before committing
- Want to experiment with APIs interactively
- Prototyping requires safe environment
- Running untrusted code is dangerous

Current workarounds are painful:
- Manually create temp files → slow and cluttered
- Risk breaking existing code → anxiety
- No isolation → security concerns
- Can't easily rollback → manual cleanup
- Context switching to external REPLs → lost flow

This leads to:
- Accepting code without testing
- Manual verification overhead
- Production bugs from untested code
- Hesitation to experiment
- Slower iteration cycles

---

## Key Capabilities

### Interactive Code Execution
- Run code snippets instantly
- See output in real-time
- Multi-statement execution
- Persistent REPL session
- Import project code

### Language Support
- Go: playground-style execution
- Python: IPython-like REPL
- JavaScript/TypeScript: Node.js REPL
- Shell: safe bash/zsh execution
- SQL: query execution
- More languages extensible

### Safe Sandboxing
- Isolated environment per language
- Resource limits (CPU, memory, time)
- Network access control
- File system restrictions
- Process isolation

### Result Visualization
- Syntax-highlighted output
- Pretty-printed data structures
- Table formatting for data
- Image/graph rendering
- Error traces with context

### Integration with Agent
- Execute before writing to files
- Test generated code automatically
- Quick prototyping workflow
- Verify API responses
- Debug complex logic

---

## Example Use Cases

### 1. Test Before Writing

```
User: Create a function to calculate compound interest

Agent: I'll create and test the function first:

[Sandbox execution]
def compound_interest(principal, rate, time, n=12):
    """Calculate compound interest"""
    return principal * (1 + rate/(n*100))**(n*time)

# Test cases
print(compound_interest(1000, 5, 10))  # $1,647.01
print(compound_interest(5000, 3.5, 5)) # $5,941.29

Output:
1647.0094888843174
5941.290166907242

✓ Tests passed! Function works correctly.

Write this to file? (y/n)
```

### 2. Interactive API Exploration

```
User: Help me figure out how to use the Stripe API for payments

Agent: Let's explore the Stripe API interactively:

[Sandbox - Python REPL]
>>> import stripe
>>> stripe.api_key = "sk_test_..."
>>> 
>>> # List available payment methods
>>> methods = stripe.PaymentMethod.list(type="card", limit=3)
>>> for method in methods:
...     print(f"{method.id}: {method.card.brand} ending in {method.card.last4}")

Output:
pm_1AB2CD: visa ending in 4242
pm_2EF3GH: mastercard ending in 5555
pm_3IJ4KL: amex ending in 1234

>>> # Create a test payment intent
>>> intent = stripe.PaymentIntent.create(
...     amount=1000,
...     currency="usd",
...     payment_method="pm_1AB2CD"
... )
>>> print(f"Created intent: {intent.id}, status: {intent.status}")

Output:
Created intent: pi_1MN5OP6QR, status: requires_confirmation

✓ API working! Ready to implement?
```

### 3. Quick Prototyping

```
User: I need to parse JSON logs and extract errors

Agent: Let me prototype a solution:

[Sandbox - JavaScript]
const logs = `
{"level":"info","msg":"Server started"}
{"level":"error","msg":"Database connection failed","error":"ECONNREFUSED"}
{"level":"info","msg":"Request processed"}
{"level":"error","msg":"API timeout","error":"ETIMEDOUT"}
`;

const errors = logs
  .split('\n')
  .filter(line => line.trim())
  .map(line => JSON.parse(line))
  .filter(log => log.level === 'error')
  .map(log => ({
    message: log.msg,
    error: log.error
  }));

console.log(errors);

Output:
[
  { message: 'Database connection failed', error: 'ECONNREFUSED' },
  { message: 'API timeout', error: 'ETIMEDOUT' }
]

✓ Prototype working! Here's the full implementation:
[Shows complete script with file I/O]
```

### 4. Safe Execution of Generated Code

```
User: Generate a script to clean up old log files

Agent: I'll generate and test the script in sandbox first:

[Sandbox - Bash (restricted)]
#!/bin/bash
# Cleanup script - RUNNING IN SANDBOX

# Find log files older than 30 days
find /tmp/sandbox/logs -name "*.log" -mtime +30 -type f

# (In sandbox, this shows test files)
Output:
/tmp/sandbox/logs/app-2025-10-15.log
/tmp/sandbox/logs/error-2025-10-10.log
/tmp/sandbox/logs/access-2025-10-08.log

# Test delete (dry-run)
find /tmp/sandbox/logs -name "*.log" -mtime +30 -type f -print -delete

Output:
Would delete 3 files (45 MB)

⚠️ Safety check: This will DELETE 3 files from logs/
Run in real environment? (y/n)
```

### 5. Database Query Testing

```
User: Help me write a SQL query to find inactive users

Agent: Let's build the query step by step:

[Sandbox - PostgreSQL (test database)]
-- First, let's see the users table structure
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'users';

Output:
 column_name |     data_type
-------------+------------------
 id          | integer
 email       | varchar
 last_login  | timestamp
 created_at  | timestamp
 is_active   | boolean

-- Now find inactive users (no login in 90 days)
SELECT id, email, last_login,
       NOW() - last_login as days_inactive
FROM users
WHERE last_login < NOW() - INTERVAL '90 days'
  AND is_active = true
ORDER BY last_login ASC
LIMIT 10;

Output:
 id |        email        |     last_login      | days_inactive
----+---------------------+---------------------+---------------
  5 | user5@example.com   | 2025-06-15 10:23:45 | 159 days
  12| user12@example.com  | 2025-07-01 14:32:11 | 143 days
  23| user23@example.com  | 2025-07-20 09:15:33 | 124 days
(10 rows)

✓ Query working! Found 47 total inactive users.
Apply to production database? (y/n)
```

---

## Technical Approach

### Sandbox Isolation

**Container-Based (Primary):**
- Docker containers for each language
- Pre-built images with common libraries
- Network isolation (optional internet access)
- File system limits (tmpfs, quotas)
- CPU/memory resource limits
- Auto-cleanup after execution

**Process-Based (Fallback):**
- Separate process per execution
- Restricted permissions
- Resource limits (ulimit, cgroups)
- Timeout enforcement
- Clean environment variables

### Language Runtimes

**Go Playground:**
- Compile and run in isolated container
- Capture stdout/stderr
- Timeout after 30 seconds
- Memory limit 100MB
- No network access (or restricted)

**Python REPL:**
- IPython kernel in container
- Persistent session
- Rich output (tables, plots)
- Import safety checks
- Standard library available

**JavaScript/TypeScript:**
- Node.js REPL in container
- npm packages available
- Console capture
- Error stack traces
- Module resolution

**Shell Execution:**
- Restricted bash/zsh
- Limited commands (whitelist)
- No system modification
- Read-only system files
- Safe path handling

### Security Model

**Defense in Depth:**
1. Container isolation (primary boundary)
2. Non-root user execution
3. Read-only file systems (except /tmp)
4. Network policies (deny by default)
5. Resource limits (prevent DoS)
6. Execution timeouts
7. Output size limits

**Safe Defaults:**
- No network access unless explicitly allowed
- Cannot access host file system
- Cannot execute privileged operations
- Automatic cleanup of all artifacts
- Audit logging of all executions

### Result Handling

**Output Processing:**
- Stream output in real-time
- Syntax highlight by type
- Pretty-print JSON, tables
- Truncate large outputs (with expand option)
- Capture stdout, stderr separately

**Error Handling:**
- Parse stack traces
- Highlight error lines
- Show context around errors
- Suggest fixes for common errors
- Link to documentation

---

## Value Propositions

### For All Developers
- Test before committing
- Experiment safely
- Quick prototyping
- Interactive learning
- Faster feedback loops

### For AI-Generated Code
- Verify correctness automatically
- Catch errors before writing files
- Build confidence in AI output
- Iterate on logic interactively
- Reduce bugs in production

### For Learning & Exploration
- Try new APIs safely
- Experiment with libraries
- Learn by doing
- No setup required
- Immediate feedback

---

## Implementation Phases

### Phase 1: Basic Execution (3 weeks)
- Docker-based sandboxing
- Go, Python, JavaScript support
- Basic output capture
- Timeout and resource limits
- Simple TUI integration

### Phase 2: REPL Sessions (2 weeks)
- Persistent REPL for Python/JS
- Multi-statement execution
- Session state management
- Import project code
- Rich output formatting

### Phase 3: Advanced Features (3 weeks)
- Network access control
- Database query execution
- File system operations (sandboxed)
- Result visualization
- Error analysis and suggestions

### Phase 4: Integration (2 weeks)
- Auto-test generated code
- Pre-commit verification
- Batch execution
- Result comparison
- Performance profiling

---

## Open Questions

1. **Resource Management:** How many concurrent sandboxes?
2. **Persistence:** Should REPL sessions persist across agent restarts?
3. **Network:** Allow external API calls? If so, how to control?
4. **Images:** Pre-build language images or pull on-demand?
5. **Cleanup:** Aggressive cleanup or keep sandboxes warm for performance?
6. **Cost:** Docker overhead acceptable? Need lighter alternatives?

---

## Related Features

**Enhanced by:**
- **LSP Integration** - Type checking in sandbox
- **Testing Automation** - Run tests in sandbox
- **Code Intelligence** - Analyze sandbox execution

**Enables:**
- **Interactive Debugging** - Step through code in sandbox
- **Performance Testing** - Benchmark in isolated environment

---

## Success Metrics

**Adoption:**
- 50%+ test code in sandbox before writing
- 40%+ use interactive REPL
- 30%+ use for prototyping
- 60%+ auto-test generated code

**Quality:**
- 80% reduction in runtime errors from AI-generated code
- 90% of sandbox executions complete successfully
- 95% of tested code written to files without issues

**Performance:**
- Sandbox startup &lt;2 seconds
- Execution results in &lt;5 seconds
- Support 10+ concurrent sandboxes
- Memory usage &lt;500MB per sandbox

**Satisfaction:**
- 4.5+ rating for sandbox feature
- "Game changer for confidence in AI code" feedback
- "Love the interactive prototyping" comments

---

## Risks and Mitigations

### Risk: Security Vulnerabilities
**Impact:** Sandbox escape, host system compromise  
**Mitigation:**
- Multiple isolation layers
- Regular security audits
- Keep container images updated
- Principle of least privilege
- Comprehensive testing

### Risk: Resource Exhaustion
**Impact:** System slowdown, DoS  
**Mitigation:**
- Strict resource limits
- Concurrent sandbox limits
- Automatic cleanup
- Monitoring and alerts
- Graceful degradation

### Risk: Complex Setup
**Impact:** Users can't use feature  
**Mitigation:**
- Auto-install Docker if missing
- Pre-built images
- Fallback to process-based
- Clear error messages
- Setup troubleshooting guide

---

## Next Steps

1. **Prototype** - Build basic Go sandbox (1 week)
2. **Security Review** - Validate isolation approach
3. **User Testing** - Test with 5-10 developers
4. **Iterate** - Refine based on feedback
5. **Expand** - Add more languages
6. **Document** - Write usage guides
