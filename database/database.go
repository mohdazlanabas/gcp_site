package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

// InitDB initializes the PostgreSQL database connection
func InitDB() (*pgx.Conn, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// For App Engine Standard, DB_HOST is the Cloud SQL connection name
		// For local development, it's a standard host:port
		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")

		if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
			return nil, fmt.Errorf("missing database environment variables (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME or DATABASE_URL)")
		}

		// Check if DB_HOST is a Cloud SQL connection name (starts with /cloudsql/)
		if len(dbHost) > 9 && dbHost[:9] == "/cloudsql/" {
			// App Engine Standard connection string
			connStr = fmt.Sprintf("user=%s password=%s database=%s host=%s",
				dbUser, dbPassword, dbName, dbHost)
		} else {
			// Local or other environment connection string
			connStr = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
				dbUser, dbPassword, dbHost, dbName)
		}
	}

	log.Printf("Attempting to connect to database: %s", connStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL!")
	return conn, nil
}

// MigrateDB creates the messages table if it doesn't exist
func MigrateDB(conn *pgx.Conn) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := conn.Exec(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create messages table: %w", err)
	}

	log.Println("Messages table migration complete.")
	return nil
}
