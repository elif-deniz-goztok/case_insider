package service

import (
	"testing"

	"github.com/elif-deniz-goztok/case_insider/models"
)

func TestSimulateMatch_ReturnsNonNegativeGoals(t *testing.T) {
	svc := NewSimulationService()
	home := models.Team{ID: 1, Name: "Chelsea", Strength: 9}
	away := models.Team{ID: 2, Name: "Liverpool", Strength: 6}

	for range 100 {
		hg, ag := svc.SimulateMatch(home, away)
		if hg < 0 || ag < 0 {
			t.Errorf("got negative goals: home=%d away=%d", hg, ag)
		}
	}
}

func TestSimulateMatch_StrongerTeamWinsMoreOften(t *testing.T) {
	svc := NewSimulationService()
	strong := models.Team{ID: 1, Name: "Chelsea", Strength: 10}
	weak := models.Team{ID: 2, Name: "Liverpool", Strength: 1}

	strongWins := 0
	const runs = 1000
	for range runs {
		hg, ag := svc.SimulateMatch(strong, weak)
		if hg > ag {
			strongWins++
		}
	}

	// With strength 10 vs 1, strong team should win the vast majority
	if strongWins < runs*60/100 {
		t.Errorf("expected strong team to win >60%% of matches, got %d/%d", strongWins, runs)
	}
}

func TestPredictChampionship_SumTo100(t *testing.T) {
	svc := NewSimulationService()
	teams := []models.Team{
		{ID: 1, Name: "Chelsea", Strength: 9},
		{ID: 2, Name: "Man City", Strength: 8},
		{ID: 3, Name: "Arsenal", Strength: 7},
		{ID: 4, Name: "Liverpool", Strength: 6},
	}
	standings := []models.Standing{
		{Team: teams[0], Points: 9},
		{Team: teams[1], Points: 6},
		{Team: teams[2], Points: 3},
		{Team: teams[3], Points: 0},
	}
	remaining := []models.Match{
		{ID: 1, HomeTeam: teams[2], AwayTeam: teams[3]},
		{ID: 2, HomeTeam: teams[1], AwayTeam: teams[0]},
	}

	predictions := svc.PredictChampionship(teams, standings, remaining, 1000)

	if len(predictions) != len(teams) {
		t.Fatalf("expected %d predictions, got %d", len(teams), len(predictions))
	}

	var total float64
	for _, p := range predictions {
		total += p.ChampionshipPct
	}
	// Allow small floating point tolerance
	if total < 99.0 || total > 101.0 {
		t.Errorf("championship percentages should sum to ~100, got %.2f", total)
	}
}
