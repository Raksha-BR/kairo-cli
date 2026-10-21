// cmd/ui.go
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/Raksha-BR/kairo-cli/internal/ui"
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Launch the interactive Terminal UI dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(ui.InitialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error launching UI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}