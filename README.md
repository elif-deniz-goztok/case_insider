# case_insider

Football league simulation REST API built with Go, Gin, and PostgreSQL.

Four teams play a double round-robin (6 weeks). The API simulates match results week by week, maintains a live league table, and predicts championship probabilities from week 4 onwards.

## Setup

### Prerequisites
- Go 1.21+
- PostgreSQL

### 1. Clone and install dependencies

```bash
git clone https://github.com/elif-deniz-goztok/case_insider.git
cd case_insider
go mod download
```

### 2. Create the database

```bash
createdb case_insider
psql case_insider < db/schema.sql
psql case_insider < db/seed.sql
```

### 3. Configure environment

```bash
cp .env.example .env
# Edit .env with your PostgreSQL credentials
```

### 4. Run

```bash
go run main.go
# Server starts on http://localhost:8080
```

---

## API Reference

All responses follow this envelope:
```json
{ "data": <payload>, "error": "<message or omitted>" }
```

### League

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/league/table` | Current standings |
| GET | `/api/league/weeks` | All weeks with match results |
| GET | `/api/league/weeks/:week` | Matches for a specific week (1–6) |
| POST | `/api/league/next-week` | Simulate the next week |
| POST | `/api/league/play-all` | Simulate all remaining weeks |
| GET | `/api/league/predictions` | Championship % per team (week 4+ only) |
| POST | `/api/league/reset` | Reset league to initial state |

### Matches

| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/api/matches/:id` | Edit a match result |

#### Edit match — request body
```json
{
  "home_goals": 3,
  "away_goals": 1
}
```

---

## Example Postman flow

1. `GET /api/league/table` — confirm all teams start at 0 points
2. `POST /api/league/next-week` × 4 — simulate weeks 1–4
3. `GET /api/league/predictions` — view championship probabilities
4. `POST /api/league/next-week` × 2 — complete the season
5. `GET /api/league/table` — final standings
6. `POST /api/league/reset` — start over
7. `POST /api/league/play-all` — simulate all 6 weeks in one call
8. `PUT /api/matches/1` `{"home_goals":5,"away_goals":0}` — override a result

---

## Architecture

```
main.go          → wires all layers, starts server
config/          → environment variable loading
db/              → PostgreSQL connection, schema.sql, seed.sql
models/          → Team, Match, Standing, Prediction structs
repository/      → interfaces + PostgreSQL implementations (TeamRepository, MatchRepository)
service/         → interfaces + business logic (LeagueService, SimulationService)
handler/         → Gin HTTP handlers (input validation only)
router/          → route registration
```

Dependencies flow strictly downward: handler → service → repository. Every layer boundary is an interface.

## Simulation model

Match outcomes use a **Poisson distribution** parameterised by team strength:

- `homeExpected = 1.5 × (homeStrength / awayStrength) × 1.1`
- `awayExpected = 1.5 × (awayStrength / homeStrength)`

Goals are sampled independently for each side.

Championship predictions use **1000 Monte Carlo simulations** of remaining matches.

## Teams

| Team | Strength |
|------|----------|
| Chelsea | 9 |
| Manchester City | 8 |
| Arsenal | 7 |
| Liverpool | 6 |

## Running tests

```bash
go test ./...
```
