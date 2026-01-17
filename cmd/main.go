package main

import (
	"log/slog"
	"ledgerly/db"
	"ledgerly/routes"
	"os"

	_ "ledgerly/docs" // Swagger docs

	"github.com/joho/godotenv"
)

// @title Ledgerly API
// @version 1.0
// @description Petty cash and expense management API
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using defaults")
	} else {
		slog.Info("Loaded configuration from .env")
	}

	db.InitDB()

	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Ensure port starts with :
	if port[0] != ':' {
		port = ":" + port
	}

	slog.Info("Starting server", "port", port)
	if err := r.Run(port); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
