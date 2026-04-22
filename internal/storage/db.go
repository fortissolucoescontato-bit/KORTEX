package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const (
	dbDir  = ".kortex"
	dbFile = "kortex.db"
)

// DB wraps the sql.DB connection.
type DB struct {
	*sql.DB
}

// Open initializes the SQLite database in the user's home directory.
func Open(homeDir string) (*DB, error) {
	dir := filepath.Join(homeDir, dbDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	path := filepath.Join(dir, dbFile)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	instance := &DB{db}
	if err := instance.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return instance, nil
}

// migrate creates the initial schema for Kortex.
func (db *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS installed_agents (
			agent_id TEXT PRIMARY KEY
		);`,
		`CREATE TABLE IF NOT EXISTS model_assignments (
			agent_id    TEXT,
			phase       TEXT,
			provider_id TEXT,
			model_id    TEXT,
			PRIMARY KEY (agent_id, phase)
		);`,
		`CREATE TABLE IF NOT EXISTS config_kv (
			key   TEXT PRIMARY KEY,
			value TEXT
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
