package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kerem-ozt/GoodBlast_API/routes"
	"github.com/kerem-ozt/GoodBlast_API/services"
)

// @title GoLang Rest API Starter Doc
// @version 1.0
// @description GoLang - Gin - RESTful - MongoDB - Redis
// @termsOfService https://swagger.io/terms/

// @contact.name Muhammet Kerem Ozturk
// @contact.url https://github.com/kerem-ozt
// @contact.email mkeremozt@gmail.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3002
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Bearer-Token

func main() {
	// Load configuration and initialize MongoDB
	services.LoadConfig()
	services.InitMongoDB()

	// Init leaderboard
	_, _ = services.EnsureLeaderboardInitialized("global")
	// services.EnsureLeaderboardInitialized("global")

	// Check Redis connection if configured
	if services.Config.UseRedis {
		services.CheckRedisConnection()
	}

	// Initialize Gin router
	routes.InitGin()
	router := routes.New()

	// Create HTTP server
	server := &http.Server{
		Addr:         services.Config.ServerAddr + ":" + services.Config.ServerPort,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Handler:      router,
	}

	// Start HTTP server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Server ListenAndServe error: %s\n", err)
		}
	}()

	// Create the first tournament
	_, _ = services.CreateTournament()

	// Schedule routine to create a new tournament every 24 hours
	go func() {
		now := time.Now().UTC()

		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		durationUntilMidnight := nextMidnight.Sub(now)
		time.Sleep(durationUntilMidnight)

		_, _ = services.CreateTournament()

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			_, _ = services.CreateTournament()
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
