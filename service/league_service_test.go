package service

import (
	"testing"

	"github.com/elif-deniz-goztok/case_insider/models"
)

func TestComputeStandings_CorrectPoints(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Name: "Chelsea"},
		{ID: 2, Name: "Arsenal"},
	}
	h, a := 2, 1
	matches := []models.Match{
		{HomeTeam: teams[0], AwayTeam: teams[1], HomeGoals: &h, AwayGoals: &a, Played: true},
	}

	standings := computeStandings(teams, matches)

	if standings[0].Team.Name != "Chelsea" {
		t.Errorf("expected Chelsea top, got %s", standings[0].Team.Name)
	}
	if standings[0].Points != 3 {
		t.Errorf("winner should have 3 points, got %d", standings[0].Points)
	}
	if standings[1].Points != 0 {
		t.Errorf("loser should have 0 points, got %d", standings[1].Points)
	}
}

func TestComputeStandings_DrawPoints(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Name: "Chelsea"},
		{ID: 2, Name: "Arsenal"},
	}
	g := 1
	matches := []models.Match{
		{HomeTeam: teams[0], AwayTeam: teams[1], HomeGoals: &g, AwayGoals: &g, Played: true},
	}

	standings := computeStandings(teams, matches)

	for _, s := range standings {
		if s.Points != 1 {
			t.Errorf("%s should have 1 point after draw, got %d", s.Team.Name, s.Points)
		}
		if s.Drawn != 1 {
			t.Errorf("%s should have 1 draw, got %d", s.Team.Name, s.Drawn)
		}
	}
}

func TestComputeStandings_GoalDifferenceTiebreak(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Name: "Chelsea"},
		{ID: 2, Name: "Arsenal"},
		{ID: 3, Name: "Liverpool"},
	}
	// Chelsea wins 1-0, Arsenal wins 3-0 — both 3 pts, Arsenal has better GD
	h1, a1 := 1, 0
	h2, a2 := 3, 0
	matches := []models.Match{
		{HomeTeam: teams[0], AwayTeam: teams[2], HomeGoals: &h1, AwayGoals: &a1, Played: true},
		{HomeTeam: teams[1], AwayTeam: teams[2], HomeGoals: &h2, AwayGoals: &a2, Played: true},
	}

	standings := computeStandings(teams, matches)

	if standings[0].Team.Name != "Arsenal" {
		t.Errorf("Arsenal should rank 1st by GD tiebreak, got %s", standings[0].Team.Name)
	}
}

func TestComputeStandings_UnplayedMatchesIgnored(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Name: "Chelsea"},
		{ID: 2, Name: "Arsenal"},
	}
	matches := []models.Match{
		{HomeTeam: teams[0], AwayTeam: teams[1], Played: false},
	}

	standings := computeStandings(teams, matches)

	for _, s := range standings {
		if s.Points != 0 || s.Played != 0 {
			t.Errorf("unplayed match should not affect standings, got %+v", s)
		}
	}
}
