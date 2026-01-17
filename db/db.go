package db

import (
	"log/slog"
	"ledgerly/models"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data.db"
	}
	
	slog.Info("Initializing database", "path", dbPath)
	
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect database", "error", err)
		os.Exit(1)
	}

	slog.Info("Running auto-migrations")
	err = DB.AutoMigrate(&models.PettyCashTransaction{}, &models.Expense{}, &models.User{})
	if err != nil {
		slog.Error("Failed to migrate database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database initialized successfully")
}
