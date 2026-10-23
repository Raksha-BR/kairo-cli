// internal/db/db.go
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %v", err)
	}

	dbPath := filepath.Join(homeDir, ".kairo.db")

	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("could not open database: %v", err)
	}

	if err := database.Ping(); err != nil {
		return fmt.Errorf("could not ping database: %v", err)
	}

	DB = database
	if err := createTables(); err != nil {
		return err
	}

	return migrateSchema()
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		status TEXT DEFAULT 'TODO',
		priority TEXT DEFAULT 'MEDIUM',
		due_date DATETIME,
		time_spent_seconds INTEGER DEFAULT 0,
		repo_path TEXT,
		branch_name TEXT,
		parent_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(query)
	return err
}

func migrateSchema() error {
	migrations := []string{
		"ALTER TABLE tasks ADD COLUMN priority TEXT DEFAULT 'MEDIUM';",
		"ALTER TABLE tasks ADD COLUMN due_date DATETIME;",
		"ALTER TABLE tasks ADD COLUMN time_spent_seconds INTEGER DEFAULT 0;",
	}

	for _, m := range migrations {
		_, _ = DB.Exec(m)
	}
	return nil
}