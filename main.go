package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	healthHandler := handler.NewHealthHandler(database)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router.New(leagueHandler, matchHandler, healthHandler),
	}

	// Start server in background so we can listen for shutdown signals.
	go func() {
		log.Printf("server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server stopped")
}
