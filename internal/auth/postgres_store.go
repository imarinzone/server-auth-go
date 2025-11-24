package auth

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgresStore implements Store using PostgreSQL
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new Postgres store
func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	// Ensure table exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS clients (
			client_id VARCHAR(255) PRIMARY KEY,
			client_secret VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

// VerifyCredentials checks if the clientID and clientSecret match in the DB
func (s *PostgresStore) VerifyCredentials(clientID, clientSecret string) (bool, error) {
	var storedSecret string
	err := s.db.QueryRow("SELECT client_secret FROM clients WHERE client_id = $1", clientID).Scan(&storedSecret)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("query failed: %w", err)
	}

	// In a real app, compare hashes. Here we assume plain text or pre-hashed match.
	return storedSecret == clientSecret, nil
}

// AddClient adds a new client (helper for seeding)
func (s *PostgresStore) AddClient(clientID, clientSecret string) error {
	_, err := s.db.Exec("INSERT INTO clients (client_id, client_secret) VALUES ($1, $2) ON CONFLICT (client_id) DO UPDATE SET client_secret = $2", clientID, clientSecret)
	return err
}
