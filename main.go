// main.go
package main

import (
	"fmt"
	"os"

	"github.com/Raksha-BR/kairo-cli/cmd"
	"github.com/Raksha-BR/kairo-cli/internal/db"
)

func main() {
	// Initialize the database before running any CLI commands
	if err := db.InitDB(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Execute the CLI router
	cmd.Execute()
}