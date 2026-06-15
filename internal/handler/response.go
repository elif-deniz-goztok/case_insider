package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type envelope struct {
	Data  any    `json:"data"`
	Error string `json:"error,omitempty"`
}

func respondOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, envelope{Data: data})
}

func respondError(c *gin.Context, status int, err error) {
	c.JSON(status, envelope{Error: err.Error()})
}
