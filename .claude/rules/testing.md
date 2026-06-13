---
name: testing
paths:
  - "**/*_test.go"
---

# Testing Rules

- Use table-driven tests for all business logic
- Mock dependencies via interfaces — never hit a real DB in unit tests
- Test file naming: `<file>_test.go` in the same package
- Test function naming: `Test<FunctionName>_<scenario>`
- Service layer must have unit tests; handler layer integration tests are a plus
