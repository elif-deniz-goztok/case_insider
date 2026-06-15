package handler

import (
	"context"

	"github.com/elif-deniz-goztok/case_insider/internal/models"
)

// LeagueService is the subset of service operations that LeagueHandler requires.
type LeagueService interface {
	GetStandings(ctx context.Context) ([]models.Standing, error)
	GetWeekResults(ctx context.Context, week int) ([]models.Match, error)
	GetAllWeeks(ctx context.Context) (map[int][]models.Match, error)
	SimulateNextWeek(ctx context.Context) ([]models.Match, error)
	SimulateAll(ctx context.Context) (map[int][]models.Match, error)
	GetPredictions(ctx context.Context) ([]models.Prediction, error)
	Reset(ctx context.Context) error
}

// MatchEditor is the subset of service operations that MatchHandler requires.
type MatchEditor interface {
	EditMatch(ctx context.Context, id, homeGoals, awayGoals int) (*models.Match, error)
}
