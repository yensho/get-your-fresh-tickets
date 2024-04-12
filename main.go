package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/yensho/get-your-fresh-tickets/api"
	"golang.org/x/exp/slog"
)

func main() {
	// Load the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
		os.Exit(1)
	}

	// Access the environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	log := slog.New(slog.NewJSONHandler(os.Stdout))
	var db *sqlx.DB
	var connectionError error

	// Implement backoff with a maximum number of attempts
	maxAttempts := 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, connectionError = sqlx.Connect("postgres", "postgres://"+dbUser+":"+dbPassword+"@postgres:5432/postgres?sslmode=disable")
		if connectionError == nil {
			break // Connection successful, exit the loop
		}

		// Connection failed, introduce a delay before the next attempt
		backoffDuration := time.Duration(attempt) * time.Second
		time.Sleep(backoffDuration)
	}

	if connectionError != nil {
		log.Error("Failed to connect to the database after multiple attempts")
		os.Exit(1)
	}

	api.NewApiServer(db, log).ListenAndServe()
	fmt.Println("Shutting down server...")
	os.Exit(1)
}
