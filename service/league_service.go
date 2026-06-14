package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"github.com/elif-deniz-goztok/case_insider/models"
	"github.com/elif-deniz-goztok/case_insider/repository"
)

const totalWeeks = 6
const predictionMinWeek = 4

// ErrPredictionTooEarly is returned when predictions are requested before week 4.
var ErrPredictionTooEarly = errors.New("predictions are available from week 4 onwards")

// ErrLeagueFinished is returned when trying to simulate a week after all matches are played.
var ErrLeagueFinished = errors.New("all weeks have been played")

// ErrMatchNotFound is returned when a match ID does not exist.
var ErrMatchNotFound = errors.New("match not found")

type leagueService struct {
	teams   repository.TeamRepository
	matches repository.MatchRepository
	sim     SimulationService
}

// NewLeagueService wires the league service with its dependencies.
func NewLeagueService(
	teams repository.TeamRepository,
	matches repository.MatchRepository,
	sim SimulationService,
) LeagueService {
	return &leagueService{teams: teams, matches: matches, sim: sim}
}

func (s *leagueService) GetStandings(ctx context.Context) ([]models.Standing, error) {
	teams, err := s.teams.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetStandings: %w", err)
	}
	matches, err := s.matches.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetStandings: %w", err)
	}
	return computeStandings(teams, matches), nil
}

func (s *leagueService) GetWeekResults(ctx context.Context, week int) ([]models.Match, error) {
	matches, err := s.matches.GetByWeek(ctx, week)
	if err != nil {
		return nil, fmt.Errorf("GetWeekResults week %d: %w", week, err)
	}
	return matches, nil
}

func (s *leagueService) GetAllWeeks(ctx context.Context) (map[int][]models.Match, error) {
	all, err := s.matches.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllWeeks: %w", err)
	}
	byWeek := make(map[int][]models.Match)
	for _, m := range all {
		byWeek[m.Week] = append(byWeek[m.Week], m)
	}
	return byWeek, nil
}

func (s *leagueService) SimulateNextWeek(ctx context.Context) ([]models.Match, error) {
	current, err := s.matches.GetCurrentWeek(ctx)
	if err != nil {
		return nil, fmt.Errorf("SimulateNextWeek: %w", err)
	}
	next := current + 1
	if next > totalWeeks {
		return nil, ErrLeagueFinished
	}

	weekMatches, err := s.matches.GetByWeek(ctx, next)
	if err != nil {
		return nil, fmt.Errorf("SimulateNextWeek fetch week %d: %w", next, err)
	}

	for i, m := range weekMatches {
		hg, ag := s.sim.SimulateMatch(m.HomeTeam, m.AwayTeam)
		if err := s.matches.MarkPlayed(ctx, m.ID, hg, ag); err != nil {
			return nil, fmt.Errorf("SimulateNextWeek mark played: %w", err)
		}
		weekMatches[i].HomeGoals = &hg
		weekMatches[i].AwayGoals = &ag
		weekMatches[i].Played = true
	}
	return weekMatches, nil
}

func (s *leagueService) SimulateAll(ctx context.Context) (map[int][]models.Match, error) {
	current, err := s.matches.GetCurrentWeek(ctx)
	if err != nil {
		return nil, fmt.Errorf("SimulateAll: %w", err)
	}
	if current >= totalWeeks {
		return nil, ErrLeagueFinished
	}

	results := make(map[int][]models.Match)
	for week := current + 1; week <= totalWeeks; week++ {
		played, err := s.simulateWeek(ctx, week)
		if err != nil {
			return nil, err
		}
		results[week] = played
	}
	return results, nil
}

func (s *leagueService) simulateWeek(ctx context.Context, week int) ([]models.Match, error) {
	weekMatches, err := s.matches.GetByWeek(ctx, week)
	if err != nil {
		return nil, fmt.Errorf("simulateWeek %d: %w", week, err)
	}
	for i, m := range weekMatches {
		if m.Played {
			continue
		}
		hg, ag := s.sim.SimulateMatch(m.HomeTeam, m.AwayTeam)
		if err := s.matches.MarkPlayed(ctx, m.ID, hg, ag); err != nil {
			return nil, fmt.Errorf("simulateWeek %d mark played: %w", week, err)
		}
		weekMatches[i].HomeGoals = &hg
		weekMatches[i].AwayGoals = &ag
		weekMatches[i].Played = true
	}
	return weekMatches, nil
}

func (s *leagueService) GetPredictions(ctx context.Context) ([]models.Prediction, error) {
	current, err := s.matches.GetCurrentWeek(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetPredictions: %w", err)
	}
	if current < predictionMinWeek {
		return nil, ErrPredictionTooEarly
	}

	teams, err := s.teams.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetPredictions: %w", err)
	}
	all, err := s.matches.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetPredictions: %w", err)
	}

	standings := computeStandings(teams, all)

	var remaining []models.Match
	for _, m := range all {
		if !m.Played {
			remaining = append(remaining, m)
		}
	}

	return s.sim.PredictChampionship(teams, standings, remaining, mcIterations), nil
}

func (s *leagueService) EditMatch(ctx context.Context, id, homeGoals, awayGoals int) (*models.Match, error) {
	m, err := s.matches.UpdateResult(ctx, id, homeGoals, awayGoals)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMatchNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("EditMatch %d: %w", id, err)
	}
	return m, nil
}

func (s *leagueService) Reset(ctx context.Context) error {
	if err := s.matches.Reset(ctx); err != nil {
		return fmt.Errorf("Reset: %w", err)
	}
	return nil
}

// computeStandings derives the league table from match results in memory.
// Ordered by: Points desc → GoalDifference desc → GoalsFor desc.
func computeStandings(teams []models.Team, matches []models.Match) []models.Standing {
	table := make(map[int]*models.Standing, len(teams))
	for i := range teams {
		table[teams[i].ID] = &models.Standing{Team: teams[i]}
	}

	for _, m := range matches {
		if !m.Played || m.HomeGoals == nil || m.AwayGoals == nil {
			continue
		}
		h := table[m.HomeTeam.ID]
		a := table[m.AwayTeam.ID]
		hg, ag := *m.HomeGoals, *m.AwayGoals

		h.Played++
		a.Played++
		h.GoalsFor += hg
		h.GoalsAgainst += ag
		a.GoalsFor += ag
		a.GoalsAgainst += hg

		switch {
		case hg > ag:
			h.Won++
			h.Points += 3
			a.Lost++
		case hg < ag:
			a.Won++
			a.Points += 3
			h.Lost++
		default:
			h.Drawn++
			a.Drawn++
			h.Points++
			a.Points++
		}
	}

	standings := make([]models.Standing, 0, len(teams))
	for _, st := range table {
		st.GoalDifference = st.GoalsFor - st.GoalsAgainst
		standings = append(standings, *st)
	}

	sort.Slice(standings, func(i, j int) bool {
		a, b := standings[i], standings[j]
		if a.Points != b.Points {
			return a.Points > b.Points
		}
		if a.GoalDifference != b.GoalDifference {
			return a.GoalDifference > b.GoalDifference
		}
		return a.GoalsFor > b.GoalsFor
	})

	return standings
}
