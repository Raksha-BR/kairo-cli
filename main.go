// main.go
package main

import (
	"fmt"
	"os"

	"github.com/Raksha-BR/kairo-cli/cmd"
	"github.com/Raksha-BR/kairo-cli/internal/config"
	"github.com/Raksha-BR/kairo-cli/internal/db"
)

func main() {
	// 1. Initialize SQLite Database
	if err := db.InitDB(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// 2. Load Configuration (~/.kairo.yaml)
	if err := config.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 3. Execute CLI Router
	cmd.Execute()
}