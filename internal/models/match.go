package models

// Match represents a scheduled or completed fixture between two teams.
// HomeGoals and AwayGoals are nil when the match has not yet been played.
type Match struct {
	ID        int  `json:"id"`
	Week      int  `json:"week"`
	HomeTeam  Team `json:"home_team"`
	AwayTeam  Team `json:"away_team"`
	HomeGoals *int `json:"home_goals"`
	AwayGoals *int `json:"away_goals"`
	Played    bool `json:"played"`
}
