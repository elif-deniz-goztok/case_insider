package models

// Standing holds a team's accumulated statistics in the league table.
type Standing struct {
	Team           Team `json:"team"`
	Played         int  `json:"played"`
	Won            int  `json:"won"`
	Drawn          int  `json:"drawn"`
	Lost           int  `json:"lost"`
	GoalsFor       int  `json:"goals_for"`
	GoalsAgainst   int  `json:"goals_against"`
	GoalDifference int  `json:"goal_difference"`
	Points         int  `json:"points"`
}

// Prediction holds the estimated probability of a team winning the championship.
type Prediction struct {
	Team            Team    `json:"team"`
	ChampionshipPct float64 `json:"championship_pct"`
}
