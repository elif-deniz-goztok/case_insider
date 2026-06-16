package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/elif-deniz-goztok/case_insider/internal/models"
)

// --- mock repository implementations ---

type mockTeamRepo struct {
	teams []models.Team
}

func (r *mockTeamRepo) GetAll(_ context.Context) ([]models.Team, error) { return r.teams, nil }
func (r *mockTeamRepo) GetByID(_ context.Context, _ int) (*models.Team, error) {
	return nil, nil
}

type mockMatchRepo struct {
	currentWeek  int
	allMatches   []models.Match
	weekMatches  []models.Match
	updateErr    error
	updateResult *models.Match
}

func (m *mockMatchRepo) GetAll(_ context.Context) ([]models.Match, error) {
	return m.allMatches, nil
}
func (m *mockMatchRepo) GetByWeek(_ context.Context, _ int) ([]models.Match, error) {
	return m.weekMatches, nil
}
func (m *mockMatchRepo) GetCurrentWeek(_ context.Context) (int, error) {
	return m.currentWeek, nil
}
func (m *mockMatchRepo) MarkPlayed(_ context.Context, _, _, _ int) error { return nil }
func (m *mockMatchRepo) UpdateResult(_ context.Context, _, _, _ int) (*models.Match, error) {
	return m.updateResult, m.updateErr
}
func (m *mockMatchRepo) Reset(_ context.Context) error { return nil }

func TestComputeStandings(t *testing.T) {
	ptr := func(n int) *int { return &n }

	chelsea := models.Team{ID: 1, Name: "Chelsea"}
	arsenal := models.Team{ID: 2, Name: "Arsenal"}
	liverpool := models.Team{ID: 3, Name: "Liverpool"}

	tests := []struct {
		name       string
		teams      []models.Team
		matches    []models.Match
		wantFirst  string
		wantPoints map[string]int
		wantDrawn  map[string]int
		wantPlayed map[string]int
	}{
		{
			name:  "win gives 3 points to winner and 0 to loser",
			teams: []models.Team{chelsea, arsenal},
			matches: []models.Match{
				{HomeTeam: chelsea, AwayTeam: arsenal, HomeGoals: ptr(2), AwayGoals: ptr(1), Played: true},
			},
			wantFirst:  "Chelsea",
			wantPoints: map[string]int{"Chelsea": 3, "Arsenal": 0},
		},
		{
			name:  "draw gives both teams 1 point and 1 drawn",
			teams: []models.Team{chelsea, arsenal},
			matches: []models.Match{
				{HomeTeam: chelsea, AwayTeam: arsenal, HomeGoals: ptr(1), AwayGoals: ptr(1), Played: true},
			},
			wantPoints: map[string]int{"Chelsea": 1, "Arsenal": 1},
			wantDrawn:  map[string]int{"Chelsea": 1, "Arsenal": 1},
		},
		{
			name:  "goal difference breaks points tie",
			teams: []models.Team{chelsea, arsenal, liverpool},
			matches: []models.Match{
				{HomeTeam: chelsea, AwayTeam: liverpool, HomeGoals: ptr(1), AwayGoals: ptr(0), Played: true},
				{HomeTeam: arsenal, AwayTeam: liverpool, HomeGoals: ptr(3), AwayGoals: ptr(0), Played: true},
			},
			wantFirst: "Arsenal",
		},
		{
			name:  "unplayed matches do not affect standings",
			teams: []models.Team{chelsea, arsenal},
			matches: []models.Match{
				{HomeTeam: chelsea, AwayTeam: arsenal, Played: false},
			},
			wantPoints: map[string]int{"Chelsea": 0, "Arsenal": 0},
			wantPlayed: map[string]int{"Chelsea": 0, "Arsenal": 0},
		},
		{
			name:  "away goals break tie when points, GD and goals scored are equal",
			teams: []models.Team{chelsea, arsenal},
			matches: []models.Match{
				// Chelsea wins 1-0 at home → Chelsea: 3pts, GF=1, GA=0, awayGF=0
				{HomeTeam: chelsea, AwayTeam: arsenal, HomeGoals: ptr(1), AwayGoals: ptr(0), Played: true},
				// Arsenal wins 2-1 at home → Arsenal: 3pts, GF=2+0=2, GA=1+1=2; Chelsea scores 1 away goal
				{HomeTeam: arsenal, AwayTeam: chelsea, HomeGoals: ptr(2), AwayGoals: ptr(1), Played: true},
			},
			// Both: 3pts, GD=0, GF=2 — Chelsea wins via away goals (1 vs 0)
			wantFirst: "Chelsea",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeStandings(tt.teams, tt.matches)

			if tt.wantFirst != "" && got[0].Team.Name != tt.wantFirst {
				t.Errorf("rank 1: got %s, want %s", got[0].Team.Name, tt.wantFirst)
			}
			for _, s := range got {
				if want, ok := tt.wantPoints[s.Team.Name]; ok && s.Points != want {
					t.Errorf("%s points: got %d, want %d", s.Team.Name, s.Points, want)
				}
				if want, ok := tt.wantDrawn[s.Team.Name]; ok && s.Drawn != want {
					t.Errorf("%s drawn: got %d, want %d", s.Team.Name, s.Drawn, want)
				}
				if want, ok := tt.wantPlayed[s.Team.Name]; ok && s.Played != want {
					t.Errorf("%s played: got %d, want %d", s.Team.Name, s.Played, want)
				}
			}
		})
	}
}

func TestSimulateNextWeek_ErrLeagueFinished(t *testing.T) {
	svc := NewLeagueService(&mockTeamRepo{}, &mockMatchRepo{currentWeek: TotalWeeks}, NewSimulationService())
	_, err := svc.SimulateNextWeek(context.Background())
	if !errors.Is(err, ErrLeagueFinished) {
		t.Errorf("got %v, want ErrLeagueFinished", err)
	}
}

func TestSimulateAll_ErrLeagueFinished(t *testing.T) {
	svc := NewLeagueService(&mockTeamRepo{}, &mockMatchRepo{currentWeek: TotalWeeks}, NewSimulationService())
	_, err := svc.SimulateAll(context.Background())
	if !errors.Is(err, ErrLeagueFinished) {
		t.Errorf("got %v, want ErrLeagueFinished", err)
	}
}

func TestGetPredictions_ErrPredictionTooEarly(t *testing.T) {
	svc := NewLeagueService(&mockTeamRepo{}, &mockMatchRepo{currentWeek: predictionMinWeek - 1}, NewSimulationService())
	_, err := svc.GetPredictions(context.Background())
	if !errors.Is(err, ErrPredictionTooEarly) {
		t.Errorf("got %v, want ErrPredictionTooEarly", err)
	}
}

func TestEditMatch_ErrMatchNotFound(t *testing.T) {
	svc := NewLeagueService(&mockTeamRepo{}, &mockMatchRepo{updateErr: sql.ErrNoRows}, NewSimulationService())
	_, err := svc.EditMatch(context.Background(), 99, 0, 0)
	if !errors.Is(err, ErrMatchNotFound) {
		t.Errorf("got %v, want ErrMatchNotFound", err)
	}
}
