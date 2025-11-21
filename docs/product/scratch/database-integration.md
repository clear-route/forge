# Feature Idea: Database Integration & Query Assistant

**Status:** Draft  
**Priority:** High Impact, Near-Term  
**Last Updated:** November 2025

---

## Overview

Native database connectivity and intelligent query assistance that allows Forge to interact with databases directly—inspect schemas, write and execute queries, manage migrations, and analyze database performance. Transforms Forge into a complete database development companion.

---

## Problem Statement

Developers working with databases face constant friction:
- Switching between code editor and database client breaks flow
- Writing SQL queries requires remembering exact schema
- Testing queries means copy-paste to external tool
- No AI assistance for query optimization
- Schema changes require manual migration writing
- Cannot see query results alongside code
- Debugging database issues requires multiple tools

Current workarounds are painful:
- Use separate database GUIs (TablePlus, DBeaver, pgAdmin)
- Context switch to write/test queries
- Copy-paste results back to code
- Manually inspect schemas
- Write migrations by hand
- No query history or saved queries

This leads to:
- Wasted time context switching
- Syntax errors in queries
- Inefficient query patterns
- Forgotten schema details
- Manual migration errors
- Lost query experimentation

---

## Key Capabilities

### Schema Inspection
- List all databases and tables
- View table structure (columns, types, constraints)
- Inspect indexes and relationships
- Show table statistics (row count, size)
- View foreign key relationships visually
- Export schema as code/diagram

### Query Execution
- Write and run SQL queries interactively
- Streaming results for large queries
- Multiple result set support
- Query parameter binding
- Transaction support
- Prepared statement execution

### Query Builder Assistant
- Generate queries from natural language
- Auto-complete table/column names
- Suggest indexes for slow queries
- Optimize existing queries
- Explain query execution plans
- Convert between SQL dialects

### Migration Management
- Generate migrations from schema changes
- Preview migration effects
- Rollback support
- Migration history
- Seed data management
- Schema versioning

### Query Library
- Save frequently used queries
- Organize queries by project/feature
- Share queries with team
- Query templates with parameters
- Quick query snippets

### Database Support
- PostgreSQL (primary)
- MySQL/MariaDB
- SQLite
- MongoDB (document queries)
- Redis (key-value operations)
- Extensible for more databases

---

## Example Use Cases

### 1. Schema Exploration

```
User: Show me the structure of the users table

Agent: Inspecting users table in database 'production'...

Table: users
Rows: 45,234
Size: 12.3 MB
Indexes: 4

Columns:
┌──────────────┬───────────┬──────────┬─────────┬─────────────┐
│ Column       │ Type      │ Nullable │ Default │ Constraints │
├──────────────┼───────────┼──────────┼─────────┼─────────────┤
│ id           │ bigint    │ NO       │ nextval │ PRIMARY KEY │
│ email        │ varchar   │ NO       │ -       │ UNIQUE      │
│ username     │ varchar   │ NO       │ -       │ UNIQUE      │
│ password_hash│ varchar   │ NO       │ -       │             │
│ created_at   │ timestamp │ NO       │ now()   │             │
│ updated_at   │ timestamp │ NO       │ now()   │             │
│ last_login   │ timestamp │ YES      │ NULL    │             │
│ is_active    │ boolean   │ NO       │ true    │             │
│ role_id      │ bigint    │ YES      │ NULL    │ FOREIGN KEY │
└──────────────┴───────────┴──────────┴─────────┴─────────────┘

Indexes:
  - users_pkey (PRIMARY KEY): id
  - users_email_key (UNIQUE): email
  - users_username_key (UNIQUE): username
  - idx_users_last_login: last_login

Foreign Keys:
  - role_id → roles(id)

Referenced By:
  - orders.user_id
  - sessions.user_id
  - profiles.user_id
```

### 2. Natural Language Query Generation

```
User: Find all users who registered in the last 30 days and have made at least one order

Agent: I'll generate that query for you:

[Generated Query]
SELECT 
    u.id,
    u.email,
    u.username,
    u.created_at,
    COUNT(o.id) as order_count
FROM users u
INNER JOIN orders o ON u.id = o.user_id
WHERE u.created_at > NOW() - INTERVAL '30 days'
GROUP BY u.id, u.email, u.username, u.created_at
HAVING COUNT(o.id) >= 1
ORDER BY u.created_at DESC;

Execute this query? (y/n/edit)
> y

Executing... ████████████████ 100%

Results (47 rows):
┌─────┬────────────────────┬──────────┬─────────────────────┬─────────────┐
│ id  │ email              │ username │ created_at          │ order_count │
├─────┼────────────────────┼──────────┼─────────────────────┼─────────────┤
│ 1001│ user1@example.com  │ user1    │ 2025-11-20 14:32:15 │ 3           │
│ 1002│ user2@example.com  │ user2    │ 2025-11-19 09:15:42 │ 1           │
│ 1003│ user3@example.com  │ user3    │ 2025-11-18 16:45:23 │ 2           │
...
└─────┴────────────────────┴──────────┴─────────────────────┴─────────────┘

Query executed in 143ms
Rows returned: 47

Save this query? (y/n)
```

### 3. Query Optimization

```
User: This query is slow - can you optimize it?

SELECT * FROM orders 
WHERE user_id = 123 
AND created_at > '2025-01-01'
AND status = 'pending'

Agent: Analyzing query performance...

Current Query Plan:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Seq Scan on orders  (cost=0.00..15234.56 rows=234)
  Filter: (user_id = 123 AND created_at > '2025-01-01' AND status = 'pending')
Planning Time: 0.234 ms
Execution Time: 234.567 ms
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Issues Found:
⚠️ Sequential scan (should use index)
⚠️ Selecting all columns (SELECT *)
⚠️ Missing composite index on (user_id, created_at, status)

Optimized Query:
SELECT 
    id, 
    user_id, 
    total_amount, 
    status, 
    created_at
FROM orders 
WHERE user_id = 123 
    AND created_at > '2025-01-01'
    AND status = 'pending';

Recommended Index:
CREATE INDEX idx_orders_user_status_date 
ON orders(user_id, status, created_at);

With this index, estimated execution time: 3.2ms (73x faster)

Create index and use optimized query? (y/n)
```

### 4. Migration Generation

```
User: I need to add a 'deleted_at' column to the users table for soft deletes

Agent: I'll create a migration for soft delete support:

Migration: add_soft_delete_to_users
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

-- Up Migration
ALTER TABLE users 
ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;

CREATE INDEX idx_users_deleted_at 
ON users(deleted_at) 
WHERE deleted_at IS NULL;

COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp';

-- Down Migration
DROP INDEX IF EXISTS idx_users_deleted_at;
ALTER TABLE users DROP COLUMN deleted_at;

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Preview changes:
✓ Adds nullable deleted_at column
✓ Creates partial index (only for non-deleted rows)
✓ Adds column comment for documentation
✓ Fully reversible migration

Write migration files? (y/n)
> y

Created:
  migrations/20251121_add_soft_delete_to_users.up.sql
  migrations/20251121_add_soft_delete_to_users.down.sql

Run migration now? (y/n/test-first)
```

### 5. Interactive Data Exploration

```
User: /db connect postgres://localhost/myapp

Agent: Connected to PostgreSQL database 'myapp'

Database: myapp
Version: PostgreSQL 15.2
Size: 234 MB
Tables: 23

Available commands:
  /db tables          - List all tables
  /db schema [table]  - Show table structure
  /db query           - Enter query mode
  /db explain [query] - Explain query plan
  /db save [name]     - Save current query

User: /db query

Agent: [Query Mode - type SQL, 'help' for assistance, 'exit' to quit]

> SELECT COUNT(*) FROM users WHERE is_active = true;

Executing...

┌───────┐
│ count │
├───────┤
│ 42,145│
└───────┘

Executed in 23ms

> help with joins

I can help you write JOIN queries! Examples:

Inner join:
  SELECT u.*, o.total FROM users u JOIN orders o ON u.id = o.user_id

Left join (include users without orders):
  SELECT u.*, o.total FROM users u LEFT JOIN orders o ON u.id = o.user_id

Multiple joins:
  SELECT u.email, o.total, p.name 
  FROM users u
  JOIN orders o ON u.id = o.user_id
  JOIN products p ON o.product_id = p.id

What would you like to join?
```

---

## Technical Approach

### Database Drivers

**Go SQL Drivers:**
```go
import (
    "database/sql"
    _ "github.com/lib/pq"           // PostgreSQL
    _ "github.com/go-sql-driver/mysql" // MySQL
    _ "github.com/mattn/go-sqlite3"    // SQLite
    "go.mongodb.org/mongo-driver/mongo" // MongoDB
    "github.com/redis/go-redis/v9"      // Redis
)
```

**Connection Management:**
- Connection pooling
- Connection health checks
- Automatic reconnection
- Timeout configuration
- SSL/TLS support
- SSH tunnel support

### Query Execution Engine

**Safe Query Execution:**
```go
type QueryExecutor struct {
    db        *sql.DB
    timeout   time.Duration
    maxRows   int
    streaming bool
}

func (qe *QueryExecutor) Execute(query string, params ...interface{}) (*QueryResult, error) {
    // Validate query (prevent DROP, TRUNCATE without confirmation)
    if err := qe.validateQuery(query); err != nil {
        return nil, err
    }
    
    // Set timeout context
    ctx, cancel := context.WithTimeout(context.Background(), qe.timeout)
    defer cancel()
    
    // Execute with streaming for large results
    if qe.streaming {
        return qe.executeStreaming(ctx, query, params...)
    }
    
    return qe.executeBuffered(ctx, query, params...)
}
```

**Result Formatting:**
- ASCII table rendering
- JSON export
- CSV export
- Markdown tables
- Copy to clipboard
- Save to file

### Schema Introspection

**PostgreSQL Schema Queries:**
```sql
-- List tables
SELECT schemaname, tablename, tableowner 
FROM pg_tables 
WHERE schemaname = 'public';

-- Table structure
SELECT 
    column_name, 
    data_type, 
    is_nullable, 
    column_default
FROM information_schema.columns
WHERE table_name = $1;

-- Indexes
SELECT 
    indexname, 
    indexdef 
FROM pg_indexes 
WHERE tablename = $1;

-- Foreign keys
SELECT
    tc.constraint_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
  ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
  ON ccu.constraint_name = tc.constraint_name
WHERE tc.table_name = $1 AND tc.constraint_type = 'FOREIGN KEY';
```

### Query Builder AI

**Natural Language to SQL:**
```go
func (ai *QueryBuilder) BuildFromNL(prompt string, schema SchemaInfo) (string, error) {
    systemPrompt := fmt.Sprintf(`
You are a SQL query expert. Generate SQL queries based on natural language.

Available schema:
%s

Rules:
- Use proper JOINs
- Include appropriate WHERE clauses
- Order results logically
- Limit results to reasonable numbers
- Use indexes where available
`, schema.String())

    return ai.llm.Generate(systemPrompt, prompt)
}
```

**Query Optimization:**
- Analyze EXPLAIN output
- Suggest missing indexes
- Identify N+1 queries
- Recommend query rewrites
- Check for inefficient patterns

### Migration System

**Migration Files:**
```
migrations/
├── 001_initial_schema.up.sql
├── 001_initial_schema.down.sql
├── 002_add_users.up.sql
├── 002_add_users.down.sql
└── ...
```

**Migration Runner:**
```go
type Migrator struct {
    db          *sql.DB
    migrationsDir string
    versionTable  string
}

func (m *Migrator) Up() error {
    // Get current version
    // Find pending migrations
    // Execute in transaction
    // Update version table
}

func (m *Migrator) Down() error {
    // Rollback last migration
}
```

---

## Value Propositions

### For Backend Developers
- No context switching to database clients
- AI-assisted query writing
- Fast schema exploration
- Query optimization help

### For Full-Stack Developers
- Integrated database workflow
- Quick data inspection
- Test queries alongside code
- Migration generation

### For Data Engineers
- Complex query building
- Performance analysis
- Data exploration
- Schema management

---

## Implementation Phases

### Phase 1: Core Connectivity (3 weeks)
- PostgreSQL connection
- Basic query execution
- Schema inspection
- Result display in TUI

### Phase 2: Query Builder (2 weeks)
- Natural language to SQL
- Query templates
- Auto-completion
- Saved queries

### Phase 3: Optimization (2 weeks)
- EXPLAIN plan analysis
- Index recommendations
- Query rewriting
- Performance tracking

### Phase 4: Migrations (2 weeks)
- Migration generation
- Up/down migrations
- Version tracking
- Seed data support

### Phase 5: Multi-DB (3 weeks)
- MySQL support
- SQLite support
- MongoDB support
- Redis support

---

## Success Metrics

**Adoption:**
- 70%+ developers use database features
- 50%+ use natural language query generation
- 60%+ use schema inspection
- 40%+ use migration generation

**Quality:**
- 85%+ generated queries work first try
- 70%+ optimization suggestions accepted
- 90%+ migrations succeed without manual edits

**Performance:**
- Query execution &lt;100ms overhead
- Schema inspection &lt;1 second
- Support databases with 1000+ tables

**Satisfaction:**
- 4.6+ rating for database features
- "Replaced my database GUI" feedback
- "Query generation is magic" comments
