package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elif-deniz-goztok/case_insider/models"
)

type matchRepo struct {
	db *sql.DB
}

// NewMatchRepository creates a PostgreSQL-backed MatchRepository.
func NewMatchRepository(db *sql.DB) MatchRepository {
	return &matchRepo{db: db}
}

const matchSelectQuery = `
	SELECT
		m.id, m.week, m.played, m.home_goals, m.away_goals,
		ht.id, ht.name, ht.strength,
		at.id, at.name, at.strength
	FROM matches m
	JOIN teams ht ON m.home_team_id = ht.id
	JOIN teams at ON m.away_team_id = at.id`

func scanMatch(row interface{ Scan(...any) error }) (models.Match, error) {
	var m models.Match
	err := row.Scan(
		&m.ID, &m.Week, &m.Played, &m.HomeGoals, &m.AwayGoals,
		&m.HomeTeam.ID, &m.HomeTeam.Name, &m.HomeTeam.Strength,
		&m.AwayTeam.ID, &m.AwayTeam.Name, &m.AwayTeam.Strength,
	)
	return m, err
}

func (r *matchRepo) GetAll(ctx context.Context) ([]models.Match, error) {
	rows, err := r.db.QueryContext(ctx, matchSelectQuery+` ORDER BY m.week, m.id`)
	if err != nil {
		return nil, fmt.Errorf("matchRepo.GetAll: %w", err)
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		m, err := scanMatch(rows)
		if err != nil {
			return nil, fmt.Errorf("matchRepo.GetAll scan: %w", err)
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

func (r *matchRepo) GetByWeek(ctx context.Context, week int) ([]models.Match, error) {
	rows, err := r.db.QueryContext(ctx,
		matchSelectQuery+` WHERE m.week = $1 ORDER BY m.id`, week,
	)
	if err != nil {
		return nil, fmt.Errorf("matchRepo.GetByWeek %d: %w", week, err)
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		m, err := scanMatch(rows)
		if err != nil {
			return nil, fmt.Errorf("matchRepo.GetByWeek scan: %w", err)
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

func (r *matchRepo) GetCurrentWeek(ctx context.Context) (int, error) {
	var week int
	err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(MAX(week), 0) FROM matches WHERE played = TRUE`,
	).Scan(&week)
	if err != nil {
		return 0, fmt.Errorf("matchRepo.GetCurrentWeek: %w", err)
	}
	return week, nil
}

func (r *matchRepo) MarkPlayed(ctx context.Context, id, homeGoals, awayGoals int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE matches SET played = TRUE, home_goals = $1, away_goals = $2 WHERE id = $3`,
		homeGoals, awayGoals, id,
	)
	if err != nil {
		return fmt.Errorf("matchRepo.MarkPlayed %d: %w", id, err)
	}
	return nil
}

func (r *matchRepo) UpdateResult(ctx context.Context, id, homeGoals, awayGoals int) (*models.Match, error) {
	_, err := r.db.ExecContext(ctx,
		`UPDATE matches SET home_goals = $1, away_goals = $2, played = TRUE WHERE id = $3`,
		homeGoals, awayGoals, id,
	)
	if err != nil {
		return nil, fmt.Errorf("matchRepo.UpdateResult %d: %w", id, err)
	}

	m, err := scanMatch(r.db.QueryRowContext(ctx,
		matchSelectQuery+` WHERE m.id = $1`, id,
	))
	if err != nil {
		return nil, fmt.Errorf("matchRepo.UpdateResult fetch: %w", err)
	}
	return &m, nil
}

func (r *matchRepo) Reset(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE matches SET played = FALSE, home_goals = NULL, away_goals = NULL`,
	)
	if err != nil {
		return fmt.Errorf("matchRepo.Reset: %w", err)
	}
	return nil
}
