// Package repository defines data access interfaces and their PostgreSQL implementations.
package repository

import (
	"context"

	"github.com/elif-deniz-goztok/case_insider/models"
)

// TeamRepository defines read operations for team data.
type TeamRepository interface {
	GetAll(ctx context.Context) ([]models.Team, error)
	GetByID(ctx context.Context, id int) (*models.Team, error)
}

// MatchRepository defines all data access operations for matches.
type MatchRepository interface {
	GetAll(ctx context.Context) ([]models.Match, error)
	GetByWeek(ctx context.Context, week int) ([]models.Match, error)
	// GetCurrentWeek returns the highest week number that has been fully played.
	// Returns 0 if no matches have been played yet.
	GetCurrentWeek(ctx context.Context) (int, error)
	// MarkPlayed sets a match as played with the given score.
	MarkPlayed(ctx context.Context, id, homeGoals, awayGoals int) error
	// UpdateResult changes the score of an already-played match.
	UpdateResult(ctx context.Context, id, homeGoals, awayGoals int) (*models.Match, error)
	// Reset clears all match results, reverting the league to its initial state.
	Reset(ctx context.Context) error
}
