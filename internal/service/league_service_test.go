package service

import (
	"testing"

	"github.com/elif-deniz-goztok/case_insider/internal/models"
)

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
