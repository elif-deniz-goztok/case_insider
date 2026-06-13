# case_insider

Football league simulation backend — Insider Development Intern Hiring Day Task.

**Priority:** Code quality is the primary evaluation criterion. This project is reviewed by AI code review. Interface-based design, idiomatic Go, and clean separation of concerns must be maintained throughout.

## Commands
- Build: `go build ./...`
- Run: `go run main.go`
- Test: `go test ./...`
- Format: `gofmt -w .`
- Lint: `golangci-lint run`

## Architecture
GoLang REST API (no frontend). Testable via Postman.

```
.
├── main.go              # Entry point, router setup, dependency wiring
├── models/              # Pure data structs (Team, Match, League, Week)
├── repository/          # DB access — always behind an interface
├── service/             # Business logic — always behind an interface
├── handler/             # HTTP handlers, input validation only
└── db/                  # SQL schema and seed queries
```

### Layer Rules
- `handler` → calls `service` interface only, no DB access
- `service` → calls `repository` interface only, contains all business logic
- `repository` → SQL queries only, no business logic
- Every layer dependency must be an interface, never a concrete type

## Key Endpoints
- `GET  /league/table` — current standings
- `GET  /league/weeks/:week` — match results for a given week
- `POST /league/next-week` — simulate the next week's matches
- `POST /league/play-all` — simulate all remaining weeks at once (extra)
- `PUT  /matches/:id` — edit a match result, recalculate standings (extra)
- `GET  /league/predictions` — championship probability per team (shown after week 4)

## Conventions
- Language: Go (required — no Java, .NET, Ruby, etc.)
- Design: interface-based with struct composition (required)
- Naming: PascalCase exported, camelCase unexported, snake_case SQL
- Indentation: tabs (enforced by gofmt)
- Error handling: always wrap with `fmt.Errorf("context: %w", err)`
- Doc comments: required on all exported types and functions

## League Rules (Premier League style)
- 4 teams, each plays every other team twice (home + away) = 6 weeks, 2 matches/week
- Win: 3 pts | Draw: 1 pt | Loss: 0 pts
- Standings: Points → Goal Difference → Goals Scored
- Championship predictions shown starting from week 4

## Team Strengths
Teams have a numeric strength attribute that influences simulated match outcomes (affects goal probability distribution). Strength values are seeded in the database.

## .claude/ Config
- `rules/go-best-practices.md` — Go idioms and interface rules (auto-loaded for .go files)
- `rules/api-design.md` — endpoint structure and response format (auto-loaded for handlers)
- `rules/testing.md` — test conventions (auto-loaded for _test.go files)
- `skills/code-review` — `/code-review` skill for AI-review-style self-review
- `skills/add-endpoint` — `/add-endpoint` skill for safely scaffolding new endpoints

## Important Context
- No frontend — all interaction via REST API (Postman)
- SQL schema and queries are a required deliverable
- Deployment is a plus (share access link)
- Extras (strong plus): play-all simulation with per-week breakdown, editable match results
