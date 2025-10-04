package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func InitDB() (*pgxpool.Pool, error) {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not loaded")
	}

	// Buat connection string dari .env
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Test koneksi
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	DB = pool
	fmt.Println("✅ Connected to PostgreSQL")

	return pool, nil
}
