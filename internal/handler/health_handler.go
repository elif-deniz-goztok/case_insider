package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// pinger is satisfied by *sql.DB and any value with a PingContext method.
type pinger interface {
	PingContext(ctx context.Context) error
}

// HealthHandler exposes a liveness/readiness endpoint.
type HealthHandler struct {
	db pinger
}

// NewHealthHandler creates a HealthHandler that pings the given database.
func NewHealthHandler(db pinger) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check returns 200 if the server and database are reachable, 503 otherwise.
func (h *HealthHandler) Check(c *gin.Context) {
	if err := h.db.PingContext(c.Request.Context()); err != nil {
		respondError(c, http.StatusServiceUnavailable, err)
		return
	}
	respondOK(c, gin.H{"status": "ok"})
}
