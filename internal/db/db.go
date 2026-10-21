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

// InitDB initializes the SQLite database connection and creates tables
func InitDB() error {
	// 1. Find the user's home directory to store the database
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %v", err)
	}

	// 2. Define the path to the SQLite file
	dbPath := filepath.Join(homeDir, ".kairo.db")

	// 3. Open the connection
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("could not open database: %v", err)
	}

	// 4. Test the connection
	if err := database.Ping(); err != nil {
		return fmt.Errorf("could not ping database: %v", err)
	}

	DB = database

	// 5. Create the schema (based on our PRD)
	return createTables()
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		status TEXT DEFAULT 'TODO',
		repo_path TEXT,
		branch_name TEXT,
		parent_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(query)
	return err
}