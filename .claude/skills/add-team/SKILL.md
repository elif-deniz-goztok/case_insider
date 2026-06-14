# Add Team Skill

Add a new team to the league simulation.

Steps:
1. Ask the user for: team name, strength (1–10)
2. Add an INSERT to `db/seed.sql`:
   ```sql
   INSERT INTO teams (name, strength) VALUES ('<name>', <strength>) ON CONFLICT (name) DO NOTHING;
   ```
3. Add the team's fixtures to `db/seed.sql` — it must play every other team twice (home + away). Update week assignments to keep 2 matches per week.
4. Warn the user: adding a team changes the total number of weeks and match IDs. They must reset the DB:
   ```bash
   psql case_insider < db/schema.sql
   psql case_insider < db/seed.sql
   ```
5. Update the `totalWeeks` constant in `service/league_service.go` to reflect the new schedule length.
6. Run `go build ./...` and `go test ./...` to confirm nothing broke.
