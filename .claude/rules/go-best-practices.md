---
name: go-best-practices
paths:
  - "**/*.go"
---

# Go Best Practices

- Use interfaces for all service and repository layers — never depend on concrete types
- Prefer struct composition over inheritance
- Return errors explicitly — never panic in library/service code
- Use meaningful error wrapping: `fmt.Errorf("context: %w", err)`
- Keep functions small and single-purpose
- Exported types, functions, and methods must have doc comments
- Use table-driven tests
- Never use `init()` — initialize explicitly in main or constructors
