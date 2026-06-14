---
name: database
paths:
  - "repository/**/*.go"
  - "db/**/*.sql"
---

# Database Rules

## Go (repository layer)
- Always use parameterized queries (`$1`, `$2`, ...) — never interpolate values into SQL strings
- Always pass `context.Context` as the first argument to every DB call
- Always call `rows.Close()` with `defer` immediately after checking the query error
- Always check `rows.Err()` after iterating rows
- Wrap all errors with `fmt.Errorf("repoName.MethodName: %w", err)`
- Never put business logic in repository methods — only SQL and scanning

## SQL files
- Use `IF NOT EXISTS` on all `CREATE TABLE` statements
- Use `ON CONFLICT DO NOTHING` on seed inserts to make them idempotent
- Foreign key columns must have explicit `REFERENCES` constraints
- Use `SERIAL` for auto-increment primary keys
