// Package router wires HTTP routes to handlers.
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/elif-deniz-goztok/case_insider/handler"
)

// New returns a configured Gin engine with all API routes registered.
func New(league *handler.LeagueHandler, match *handler.MatchHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		l := api.Group("/league")
		{
			l.GET("/table", league.GetTable)
			l.GET("/weeks", league.GetAllWeeks)
			l.GET("/weeks/:week", league.GetWeek)
			l.POST("/next-week", league.NextWeek)
			l.POST("/play-all", league.PlayAll)
			l.GET("/predictions", league.GetPredictions)
			l.POST("/reset", league.Reset)
		}

		api.PUT("/matches/:id", match.EditMatch)
	}

	return r
}
