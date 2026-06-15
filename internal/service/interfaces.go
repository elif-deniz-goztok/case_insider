// Package service contains business logic for the league simulation.
package service

import (
	"context"

	"github.com/elif-deniz-goztok/case_insider/internal/models"
)

// TeamRepository defines the team data access operations that the service requires.
type TeamRepository interface {
	GetAll(ctx context.Context) ([]models.Team, error)
	GetByID(ctx context.Context, id int) (*models.Team, error)
}

// MatchRepository defines the match data access operations that the service requires.
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

// SimulationService handles match outcome generation and championship forecasting.
type SimulationService interface {
	// SimulateMatch returns a randomized score influenced by each team's strength.
	SimulateMatch(home, away models.Team) (homeGoals, awayGoals int)
	// PredictChampionship runs Monte Carlo simulations and returns championship probabilities.
	PredictChampionship(
		teams []models.Team,
		standings []models.Standing,
		remaining []models.Match,
		iterations int,
	) []models.Prediction
}
