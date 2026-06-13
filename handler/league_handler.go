// Package handler contains Gin HTTP handlers. Handlers validate input and delegate to the service layer.
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/elif-deniz-goztok/case_insider/service"
)

// LeagueHandler exposes league simulation endpoints.
type LeagueHandler struct {
	svc service.LeagueService
}

// NewLeagueHandler creates a LeagueHandler with the provided service.
func NewLeagueHandler(svc service.LeagueService) *LeagueHandler {
	return &LeagueHandler{svc: svc}
}

// GetTable returns the current league standings.
func (h *LeagueHandler) GetTable(c *gin.Context) {
	standings, err := h.svc.GetStandings(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondOK(c, standings)
}

// GetAllWeeks returns all weeks with their match results.
func (h *LeagueHandler) GetAllWeeks(c *gin.Context) {
	weeks, err := h.svc.GetAllWeeks(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondOK(c, weeks)
}

// GetWeek returns match results for a specific week.
func (h *LeagueHandler) GetWeek(c *gin.Context) {
	week, err := strconv.Atoi(c.Param("week"))
	if err != nil || week < 1 || week > 6 {
		respondError(c, http.StatusBadRequest, errors.New("week must be an integer between 1 and 6"))
		return
	}
	matches, err := h.svc.GetWeekResults(c.Request.Context(), week)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondOK(c, matches)
}

// NextWeek simulates the next unplayed week.
func (h *LeagueHandler) NextWeek(c *gin.Context) {
	matches, err := h.svc.SimulateNextWeek(c.Request.Context())
	if errors.Is(err, service.ErrLeagueFinished) {
		respondError(c, http.StatusConflict, err)
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondOK(c, matches)
}

// PlayAll simulates all remaining weeks.
func (h *LeagueHandler) PlayAll(c *gin.Context) {
	results, err := h.svc.SimulateAll(c.Request.Context())
	if errors.Is(err, service.ErrLeagueFinished) {
		respondError(c, http.StatusConflict, err)
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondOK(c, results)
}

// GetPredictions returns championship probabilities (available from week 4 onwards).
func (h *LeagueHandler) GetPredictions(c *gin.Context) {
	predictions, err := h.svc.GetPredictions(c.Request.Context())
	if errors.Is(err, service.ErrPredictionTooEarly) {
		respondError(c, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondOK(c, predictions)
}

// Reset clears all match results and restores the league to its initial state.
func (h *LeagueHandler) Reset(c *gin.Context) {
	if err := h.svc.Reset(c.Request.Context()); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondOK(c, gin.H{"message": "league reset successfully"})
}
