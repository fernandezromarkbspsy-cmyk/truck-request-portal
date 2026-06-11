package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse database url: %v", err)
	}

	// Connection pooling settings for scalability
	config.MaxConns = 10
	config.MinConns = 2

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %v", err)
	}

	// Test connection
	err = DB.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("unable to ping database: %v", err)
	}

	fmt.Println("Successfully connected to Supabase PostgreSQL")
	return nil
}
