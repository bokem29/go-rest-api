package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigration(pool *pgxpool.Pool) error {
	ctx := context.Background()

	query := `
	CREATE TABLE IF NOT EXISTS characters (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		role TEXT NOT NULL,
		game TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW(),
		deleted_at TIMESTAMPTZ NULL
	);`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("gagal menjalankan migration: %w", err)
	}

	fmt.Println("✅ Migration berhasil dijalankan")
	return nil
}
