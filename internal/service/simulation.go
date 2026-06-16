package service

import (
	"math"
	"math/rand"

	"github.com/elif-deniz-goztok/case_insider/internal/models"
)

const (
	baseGoals     = 1.5
	homeAdvantage = 1.1
	mcIterations  = 1000
)

type simulationService struct{}

// NewSimulationService creates a simulation service using a Poisson-based match model.
func NewSimulationService() *simulationService {
	return &simulationService{}
}

// SimulateMatch generates a match score using Poisson-distributed expected goals.
// Home team receives a slight advantage multiplier.
func (s *simulationService) SimulateMatch(home, away models.Team) (int, int) {
	homeExp := baseGoals * (float64(home.Strength) / float64(away.Strength)) * homeAdvantage
	awayExp := baseGoals * (float64(away.Strength) / float64(home.Strength))
	return poissonSample(homeExp), poissonSample(awayExp)
}

// PredictChampionship runs Monte Carlo simulations over remaining fixtures and
// returns each team's estimated probability of finishing top of the table.
func (s *simulationService) PredictChampionship(
	teams []models.Team,
	standings []models.Standing,
	remaining []models.Match,
	iterations int,
) []models.Prediction {
	wins := make(map[int]int, len(teams))

	for range iterations {
		pts := make(map[int]int, len(standings))
		for _, st := range standings {
			pts[st.Team.ID] = st.Points
		}

		for _, m := range remaining {
			hg, ag := s.SimulateMatch(m.HomeTeam, m.AwayTeam)
			switch {
			case hg > ag:
				pts[m.HomeTeam.ID] += 3
			case hg < ag:
				pts[m.AwayTeam.ID] += 3
			default:
				pts[m.HomeTeam.ID]++
				pts[m.AwayTeam.ID]++
			}
		}

		winner := topTeamID(pts)
		wins[winner]++
	}

	predictions := make([]models.Prediction, 0, len(teams))
	for _, t := range teams {
		predictions = append(predictions, models.Prediction{
			Team:            t,
			ChampionshipPct: math.Round(float64(wins[t.ID])/float64(iterations)*10000) / 100,
		})
	}
	return predictions
}

// poissonSample draws a sample from a Poisson distribution with the given lambda.
func poissonSample(lambda float64) int {
	l := math.Exp(-lambda)
	k := 0
	p := 1.0
	for p > l {
		k++
		p *= rand.Float64()
	}
	return k - 1
}

// topTeamID returns the team ID with the highest points in the given map.
// Ties are broken by lower team ID to ensure deterministic results across iterations.
func topTeamID(pts map[int]int) int {
	topID, topPts := 0, -1
	for id, p := range pts {
		if p > topPts || (p == topPts && id < topID) {
			topPts = p
			topID = id
		}
	}
	return topID
}
