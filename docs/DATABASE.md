# Database Schema Documentation

This document describes the database schema, relationships, and migration strategy for the Go Lang Template project.

## Table of Contents

- [Overview](#overview)
- [Database Schema](#database-schema)
- [Tables](#tables)
- [Indexes](#indexes)
- [Relationships](#relationships)
- [Migration Strategy](#migration-strategy)
- [Data Types](#data-types)
- [Best Practices](#best-practices)
- [Examples](#examples)

## Overview

The application uses PostgreSQL as its primary database backend. The database connection is optional and can be enabled via configuration. The schema is managed through versioned migrations using [golang-migrate](https://github.com/golang-migrate/migrate).

**Database Support:**
- Primary: PostgreSQL 12+
- Optional: Can run without database when `Database.Enabled=false`

## Database Schema

### Entity Relationship Diagram

```
┌─────────────────────────────────────┐
│              users                  │
├─────────────────────────────────────┤
│ id         SERIAL PRIMARY KEY       │
│ username   VARCHAR(255) NOT NULL    │
│ email      VARCHAR(255) NOT NULL    │
│ created_at TIMESTAMP NOT NULL       │
│ updated_at TIMESTAMP NOT NULL       │
└─────────────────────────────────────┘
         │
         │ Indexes:
         │ - idx_users_username (username)
         │ - idx_users_email (email)
         │
         │ Constraints:
         │ - UNIQUE (username)
         │ - UNIQUE (email)
         └─ Auto-increment ID (SERIAL)
```

## Tables

### users

The `users` table stores user account information.

**Columns:**

| Column      | Type          | Nullable | Default              | Description                    |
|-------------|---------------|----------|----------------------|--------------------------------|
| id          | SERIAL        | NO       | AUTO_INCREMENT       | Unique user identifier         |
| username    | VARCHAR(255)  | NO       | -                    | Unique username for login      |
| email       | VARCHAR(255)  | NO       | -                    | Unique email address           |
| created_at  | TIMESTAMP     | NO       | CURRENT_TIMESTAMP    | Record creation timestamp      |
| updated_at  | TIMESTAMP     | NO       | CURRENT_TIMESTAMP    | Record last update timestamp   |

**Constraints:**

- **Primary Key:** `id`
- **Unique Constraints:**
  - `username` - Ensures usernames are unique across the system
  - `email` - Ensures email addresses are unique across the system

**Indexes:**

- `idx_users_username` - B-tree index on `username` for fast lookups
- `idx_users_email` - B-tree index on `email` for fast email searches

**Migration Files:**
- Up: `migrations/000001_create_users_table.up.sql`
- Down: `migrations/000001_create_users_table.down.sql`

## Indexes

Indexes are created to optimize query performance for common access patterns.

### Performance Considerations

| Index               | Columns   | Type   | Purpose                              |
|---------------------|-----------|--------|--------------------------------------|
| idx_users_username  | username  | B-tree | Fast username lookups for login      |
| idx_users_email     | email     | B-tree | Fast email searches and validation   |

**Query Optimization:**

```sql
-- These queries will use indexes:
SELECT * FROM users WHERE username = 'john_doe';  -- Uses idx_users_username
SELECT * FROM users WHERE email = 'john@example.com';  -- Uses idx_users_email

-- These queries will be fast due to primary key:
SELECT * FROM users WHERE id = 1;  -- Uses primary key index
```

## Relationships

Currently, the schema contains a single table with no foreign key relationships. This is a starter template, and relationships can be added as needed.

### Future Relationship Patterns

When extending the schema, consider these common patterns:

**One-to-Many:**
```sql
-- Example: Users have many posts
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

**Many-to-Many:**
```sql
-- Example: Users can have many roles, roles can be assigned to many users
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE user_roles (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);
```

**One-to-One:**
```sql
-- Example: User has one profile
CREATE TABLE user_profiles (
    user_id INTEGER PRIMARY KEY,
    bio TEXT,
    avatar_url VARCHAR(500),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## Migration Strategy

The project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema versioning.

### Migration File Naming

Migrations follow the pattern: `{version}_{description}.{direction}.sql`

- **version:** Sequential number (e.g., 000001, 000002)
- **description:** Snake_case description of the change
- **direction:** `up` for applying, `down` for reverting

**Examples:**
- `000001_create_users_table.up.sql`
- `000001_create_users_table.down.sql`
- `000002_add_user_roles.up.sql`
- `000002_add_user_roles.down.sql`

### Migration Commands

**Apply all pending migrations:**
```bash
./bin/go-lang-template migrate up
```

**Rollback last migration:**
```bash
./bin/go-lang-template migrate down
```

**Apply specific number of migrations:**
```bash
./bin/go-lang-template migrate steps 2  # Apply next 2 migrations
./bin/go-lang-template migrate steps -1  # Rollback last migration
```

**Check current migration version:**
```bash
./bin/go-lang-template migrate version
```

**Force migration version (use with caution):**
```bash
./bin/go-lang-template migrate force 1
```

### Migration Best Practices

1. **Always create reversible migrations:** Every `.up.sql` must have a corresponding `.down.sql`
2. **Test migrations thoroughly:** Test both up and down migrations in development
3. **Make migrations atomic:** Each migration should be a single logical change
4. **Avoid data loss:** Use `IF EXISTS` and `IF NOT EXISTS` clauses
5. **Backup before migrating:** Always backup production data before running migrations
6. **Use transactions:** Migrations run in transactions by default (for supported operations)

### Migration Workflow

```bash
# 1. Create new migration files
touch migrations/000002_add_new_feature.up.sql
touch migrations/000002_add_new_feature.down.sql

# 2. Write SQL for both directions
# Edit the .up.sql file with your schema changes
# Edit the .down.sql file to reverse those changes

# 3. Test in development
./bin/go-lang-template migrate up
# Test your application
./bin/go-lang-template migrate down
# Verify rollback works correctly

# 4. Apply to production (with backup)
pg_dump production_db > backup.sql
./bin/go-lang-template migrate up
```

## Data Types

### PostgreSQL Type Guidelines

| PostgreSQL Type | Use Case                           | Example                    |
|-----------------|------------------------------------|-----------------------------|
| SERIAL          | Auto-incrementing integers         | User IDs, primary keys      |
| VARCHAR(n)      | Variable-length strings            | Usernames, emails           |
| TEXT            | Long text content                  | Descriptions, blog posts    |
| INTEGER         | Whole numbers                      | Counts, foreign keys        |
| BIGINT          | Large whole numbers                | Large counters              |
| BOOLEAN         | True/false values                  | Active flags                |
| TIMESTAMP       | Date and time                      | Created/updated timestamps  |
| TIMESTAMPTZ     | Date and time with timezone        | Event times across zones    |
| NUMERIC(p,s)    | Precise decimal numbers            | Prices, financial data      |
| JSON/JSONB      | Semi-structured data               | Configuration, metadata     |
| UUID            | Globally unique identifiers        | Distributed system IDs      |

### Type Selection Guidelines

- **Use SERIAL for auto-incrementing IDs** in single-database setups
- **Use UUID for distributed systems** where IDs must be globally unique
- **Use VARCHAR with reasonable limits** instead of TEXT when possible for better indexing
- **Use TIMESTAMP for timestamps** (application handles timezone conversion)
- **Use TIMESTAMPTZ for user-facing times** that need timezone awareness
- **Use JSONB instead of JSON** for better performance and indexing

## Best Practices

### Schema Design

1. **Normalization:**
   - Normalize to 3NF (Third Normal Form) for transactional tables
   - Consider denormalization for read-heavy tables with proper justification

2. **Primary Keys:**
   - Always use a primary key on every table
   - Use SERIAL for simple cases, UUID for distributed systems

3. **Foreign Keys:**
   - Always define foreign key constraints for data integrity
   - Use appropriate `ON DELETE` actions (CASCADE, RESTRICT, SET NULL)

4. **Indexes:**
   - Index foreign key columns
   - Index columns used frequently in WHERE clauses
   - Index columns used in ORDER BY and JOIN operations
   - Don't over-index - each index has overhead on writes

5. **Timestamps:**
   - Include `created_at` on all tables
   - Include `updated_at` on tables with mutable data
   - Consider `deleted_at` for soft deletes

6. **Constraints:**
   - Use NOT NULL for required fields
   - Use UNIQUE constraints where appropriate
   - Use CHECK constraints for domain validation

### Naming Conventions

- **Tables:** Plural lowercase with underscores (e.g., `users`, `user_roles`)
- **Columns:** Lowercase with underscores (e.g., `created_at`, `user_id`)
- **Indexes:** `idx_{table}_{column(s)}` (e.g., `idx_users_email`)
- **Foreign Keys:** `fk_{table}_{referenced_table}` (e.g., `fk_posts_users`)
- **Unique Constraints:** `uq_{table}_{column(s)}` (e.g., `uq_users_email`)

### Security Considerations

1. **Never store plain text passwords** - use bcrypt or argon2
2. **Use prepared statements** to prevent SQL injection (handled by Go's database/sql)
3. **Limit user permissions** - application user should not have DROP/CREATE permissions
4. **Encrypt sensitive data** at rest when required
5. **Use SSL/TLS connections** in production (sslmode=require or verify-full)

## Examples

### Creating a Migration

**000002_add_user_status.up.sql:**
```sql
-- Add status column to users table
ALTER TABLE users ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'active';

-- Add check constraint
ALTER TABLE users ADD CONSTRAINT check_user_status
    CHECK (status IN ('active', 'inactive', 'suspended'));

-- Create index for status queries
CREATE INDEX idx_users_status ON users(status);
```

**000002_add_user_status.down.sql:**
```sql
-- Remove index
DROP INDEX IF EXISTS idx_users_status;

-- Remove check constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_status;

-- Remove status column
ALTER TABLE users DROP COLUMN IF EXISTS status;
```

**Example: Creating and dropping a table**

000001_create_users_table.up.sql:
```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

000001_create_users_table.down.sql:
```sql
DROP TABLE IF EXISTS users;
```

### Querying with Go

```go
// Example: Query user by username
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
    var user User
    query := `
        SELECT id, username, email, created_at, updated_at
        FROM users
        WHERE username = $1
    `
    err := db.QueryRow(query, username).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Example: Insert new user
func CreateUser(db *sql.DB, username, email string) (*User, error) {
    var user User
    query := `
        INSERT INTO users (username, email)
        VALUES ($1, $2)
        RETURNING id, username, email, created_at, updated_at
    `
    err := db.QueryRow(query, username, email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

### Common Queries

```sql
-- Find user by ID
SELECT * FROM users WHERE id = 1;

-- Find user by username (uses index)
SELECT * FROM users WHERE username = 'john_doe';

-- Find user by email (uses index)
SELECT * FROM users WHERE email = 'john@example.com';

-- Get recently created users
SELECT * FROM users ORDER BY created_at DESC LIMIT 10;

-- Count total users
SELECT COUNT(*) FROM users;

-- Update user email
UPDATE users SET email = 'new@example.com', updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- Delete user
DELETE FROM users WHERE id = 1;
```

## Maintenance

### Regular Maintenance Tasks

1. **Vacuum:** Run VACUUM ANALYZE regularly to reclaim space and update statistics
   ```sql
   VACUUM ANALYZE users;
   ```

2. **Reindex:** Rebuild indexes if they become fragmented
   ```sql
   REINDEX TABLE users;
   ```

3. **Monitor slow queries:** Enable `log_min_duration_statement` in PostgreSQL
4. **Review query plans:** Use EXPLAIN ANALYZE for slow queries
5. **Monitor index usage:** Check `pg_stat_user_indexes` for unused indexes
6. **Backup regularly:** Automated backups with point-in-time recovery

### Monitoring Queries

```sql
-- Check table size
SELECT pg_size_pretty(pg_total_relation_size('users'));

-- Check index usage
SELECT indexrelname, idx_scan, idx_tup_read, idx_tup_fetch
FROM pg_stat_user_indexes
WHERE schemaname = 'public';

-- Find slow queries (requires pg_stat_statements extension)
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

## Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [Go database/sql Tutorial](https://go.dev/doc/database/sql-injection)
- [PostgreSQL Performance Tips](https://wiki.postgresql.org/wiki/Performance_Optimization)

---

**Note:** This schema is a starting point for the template. Extend it based on your application's specific requirements while following the best practices outlined in this document.
