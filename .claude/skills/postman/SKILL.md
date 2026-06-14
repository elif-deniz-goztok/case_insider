# Postman Collection Skill

Generate a Postman collection JSON for all API endpoints.

Read the current routes from `router/router.go` and handlers from `handler/`, then output a valid Postman Collection v2.1 JSON with:

- Collection name: "case_insider API"
- Base URL variable: `{{base_url}}` defaulting to `http://localhost:8080`
- One request per endpoint, organized into two folders: "League" and "Matches"
- Pre-filled request bodies for POST/PUT endpoints
- A logical ordering that demonstrates the full simulation flow:
  1. GET /api/league/table
  2. GET /api/league/weeks
  3. GET /api/league/weeks/1
  4. POST /api/league/next-week (repeat x4)
  5. GET /api/league/predictions
  6. POST /api/league/next-week (x2 more)
  7. GET /api/league/table (final standings)
  8. PUT /api/matches/1 with body {"home_goals": 3, "away_goals": 0}
  9. POST /api/league/reset
  10. POST /api/league/play-all

Save the output to `postman_collection.json` in the project root.
