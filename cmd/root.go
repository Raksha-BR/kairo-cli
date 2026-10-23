// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kairo",
	Short: "Kairo — Context-Aware Terminal Task & AI Decomposition Engine",
	Long:  `Kairo is a fast, offline-first terminal task manager built with Go, Cobra, SQLite, Lipgloss, and Bubbletea.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior when running `./kairo` without subcommand is launching TUI
		_ = uiCmd.RunE(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}