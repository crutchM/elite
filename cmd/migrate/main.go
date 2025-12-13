package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/crutchm/elite/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	var (
		command = flag.String("command", "up", "migration command: up, down, status")
	)
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	migrationsDir := "./migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Fatalf("Migrations directory does not exist: %s", migrationsDir)
	}

	switch *command {
	case "up":
		if err := goose.Up(db, migrationsDir); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		if err := goose.Down(db, migrationsDir); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Println("Migrations rolled back successfully")
	case "status":
		if err := goose.Status(db, migrationsDir); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}
