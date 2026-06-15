package service

import (
	"testing"

	"github.com/elif-deniz-goztok/case_insider/internal/models"
)

func TestSimulateMatch(t *testing.T) {
	svc := &simulationService{}

	tests := []struct {
		name       string
		home, away models.Team
		runs       int
		minWinRate float64
	}{
		{
			name: "goals are never negative",
			home: models.Team{ID: 1, Name: "Chelsea", Strength: 9},
			away: models.Team{ID: 2, Name: "Liverpool", Strength: 6},
			runs: 100,
		},
		{
			name:       "stronger home team wins more than 60% of the time",
			home:       models.Team{ID: 1, Strength: 10},
			away:       models.Team{ID: 2, Strength: 1},
			runs:       1000,
			minWinRate: 0.60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeWins := 0
			for range tt.runs {
				hg, ag := svc.SimulateMatch(tt.home, tt.away)
				if hg < 0 || ag < 0 {
					t.Errorf("negative goals: home=%d away=%d", hg, ag)
				}
				if hg > ag {
					homeWins++
				}
			}
			if tt.minWinRate > 0 {
				got := float64(homeWins) / float64(tt.runs)
				if got < tt.minWinRate {
					t.Errorf("home win rate %.2f < minimum %.2f", got, tt.minWinRate)
				}
			}
		})
	}
}

func TestPredictChampionship(t *testing.T) {
	svc := &simulationService{}

	teams := []models.Team{
		{ID: 1, Name: "Chelsea", Strength: 9},
		{ID: 2, Name: "Man City", Strength: 8},
		{ID: 3, Name: "Arsenal", Strength: 7},
		{ID: 4, Name: "Liverpool", Strength: 6},
	}

	tests := []struct {
		name       string
		standings  []models.Standing
		remaining  []models.Match
		iterations int
		wantLen    int
		sumNear100 bool
	}{
		{
			name: "championship percentages sum to ~100",
			standings: []models.Standing{
				{Team: teams[0], Points: 9},
				{Team: teams[1], Points: 6},
				{Team: teams[2], Points: 3},
				{Team: teams[3], Points: 0},
			},
			remaining: []models.Match{
				{ID: 1, HomeTeam: teams[2], AwayTeam: teams[3]},
				{ID: 2, HomeTeam: teams[1], AwayTeam: teams[0]},
			},
			iterations: 1000,
			wantLen:    len(teams),
			sumNear100: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.PredictChampionship(teams, tt.standings, tt.remaining, tt.iterations)

			if len(got) != tt.wantLen {
				t.Fatalf("prediction count: got %d, want %d", len(got), tt.wantLen)
			}
			if tt.sumNear100 {
				var total float64
				for _, p := range got {
					total += p.ChampionshipPct
				}
				if total < 99.0 || total > 101.0 {
					t.Errorf("percentages sum to %.2f, want ~100", total)
				}
			}
		})
	}
}
