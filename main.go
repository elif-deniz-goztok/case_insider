package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/elif-deniz-goztok/case_insider/config"
	"github.com/elif-deniz-goztok/case_insider/db"
	"github.com/elif-deniz-goztok/case_insider/handler"
	"github.com/elif-deniz-goztok/case_insider/repository"
	"github.com/elif-deniz-goztok/case_insider/router"
	"github.com/elif-deniz-goztok/case_insider/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading environment variables directly")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	database, err := db.Connect(cfg.DBDSN)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer database.Close()

	// Repository layer
	teamRepo := repository.NewTeamRepository(database)
	matchRepo := repository.NewMatchRepository(database)

	// Service layer
	simSvc := service.NewSimulationService()
	leagueSvc := service.NewLeagueService(teamRepo, matchRepo, simSvc)

	// Handler layer
	leagueHandler := handler.NewLeagueHandler(leagueSvc)
	matchHandler := handler.NewMatchHandler(leagueSvc)

	r := router.New(leagueHandler, matchHandler)

	log.Printf("server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
