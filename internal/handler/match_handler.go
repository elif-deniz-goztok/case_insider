package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/elif-deniz-goztok/case_insider/internal/service"
)

// MatchHandler exposes match-level endpoints.
type MatchHandler struct {
	svc MatchEditor
}

// NewMatchHandler creates a MatchHandler with the provided service.
func NewMatchHandler(svc MatchEditor) *MatchHandler {
	return &MatchHandler{svc: svc}
}

type editMatchRequest struct {
	HomeGoals int `json:"home_goals" binding:"min=0"`
	AwayGoals int `json:"away_goals" binding:"min=0"`
}

// EditMatch updates a match result and recalculates standings.
func (h *MatchHandler) EditMatch(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		respondError(c, http.StatusBadRequest, errors.New("invalid match id"))
		return
	}

	var req editMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	match, err := h.svc.EditMatch(c.Request.Context(), id, req.HomeGoals, req.AwayGoals)
	if errors.Is(err, service.ErrMatchNotFound) {
		respondError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondOK(c, match)
}
