package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elif-deniz-goztok/case_insider/internal/models"
)

type teamRepo struct {
	db *sql.DB
}

// NewTeamRepository creates a PostgreSQL-backed team repository.
func NewTeamRepository(db *sql.DB) *teamRepo {
	return &teamRepo{db: db}
}

func (r *teamRepo) GetAll(ctx context.Context) ([]models.Team, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, strength FROM teams ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("teamRepo.GetAll: %w", err)
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Strength); err != nil {
			return nil, fmt.Errorf("teamRepo.GetAll scan: %w", err)
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func (r *teamRepo) GetByID(ctx context.Context, id int) (*models.Team, error) {
	var t models.Team
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, strength FROM teams WHERE id = $1`, id,
	).Scan(&t.ID, &t.Name, &t.Strength)
	if err != nil {
		return nil, fmt.Errorf("teamRepo.GetByID %d: %w", id, err)
	}
	return &t, nil
}
