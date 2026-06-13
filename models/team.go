// Package models contains the core domain types for the league simulation.
package models

// Team represents a football club with a strength rating that influences match simulation.
type Team struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Strength int    `json:"strength"`
}
