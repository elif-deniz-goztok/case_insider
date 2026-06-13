// Package service contains business logic for the league simulation.
package service

import (
	"context"

	"github.com/elif-deniz-goztok/case_insider/models"
)

// LeagueService orchestrates all league operations.
type LeagueService interface {
	GetStandings(ctx context.Context) ([]models.Standing, error)
	GetWeekResults(ctx context.Context, week int) ([]models.Match, error)
	GetAllWeeks(ctx context.Context) (map[int][]models.Match, error)
	SimulateNextWeek(ctx context.Context) ([]models.Match, error)
	SimulateAll(ctx context.Context) (map[int][]models.Match, error)
	// GetPredictions returns championship probabilities. Returns an error if fewer than 4 weeks have been played.
	GetPredictions(ctx context.Context) ([]models.Prediction, error)
	EditMatch(ctx context.Context, id, homeGoals, awayGoals int) (*models.Match, error)
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
