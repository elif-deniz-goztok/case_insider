# Deploy Skill

Deploy the project to Railway.

Steps:
1. Run `go build ./...` — confirm zero errors before deploying
2. Run `go test ./...` — confirm all tests pass
3. Run `gofmt -l .` — confirm no unformatted files
4. Run `git status` — check for uncommitted changes; commit them if needed
5. Run `git push origin main` — Railway auto-deploys on push
6. Check Railway logs for startup errors: look for "server starting on" in the logs
7. Hit `GET <RAILWAY_URL>/api/league/table` to confirm the deployed API responds

If the deployment fails:
- Check that all Railway env vars are set: PORT, DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
- Check Railway build logs for Go compilation errors
- Confirm the PostgreSQL plugin is provisioned and the DATABASE_URL is linked
