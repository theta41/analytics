package pg

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// NewPG creates new postgres connection.
func NewPG(host, user, password, dbname, sslMode string, port int) (*sql.DB, error) {
	// Connect to Postgres.
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslMode)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	// Ping the database to check if it's alive.
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}
