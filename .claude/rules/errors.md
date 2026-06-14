---
name: errors
paths:
  - "**/*.go"
---

# Error Handling Rules

## Wrapping
- Always wrap errors with context: `fmt.Errorf("ServiceName.MethodName: %w", err)`
- The context prefix must identify where the error originated (package + function)
- Never discard errors with `_`

## Sentinel errors
- Define sentinel errors for conditions callers need to handle differently:
  ```go
  var ErrNotFound = errors.New("not found")
  var ErrLeagueFinished = errors.New("all weeks have been played")
  ```
- Use `errors.Is()` to check for sentinel errors — never compare error strings
- Define sentinels in the service layer, not the handler layer

## HTTP mapping (handler layer)
- `ErrLeagueFinished` → 409 Conflict
- `ErrPredictionTooEarly` → 400 Bad Request
- `sql.ErrNoRows` → 404 Not Found
- All other errors → 500 Internal Server Error
- Never expose raw internal error messages to API consumers in production

## Panics
- Never use `panic` in library or service code
- Only `log.Fatalf` in `main()` for unrecoverable startup errors
