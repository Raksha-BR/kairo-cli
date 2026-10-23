// cmd/stats.go
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Raksha-BR/kairo-cli/internal/db"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

var statsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"summary", "metrics"},
	Short:   "Display workspace task metrics and summary statistics",
	Run: func(cmd *cobra.Command, args []string) {
		gitCtx := git.GetCurrentContext()
		stats, err := db.GetStats(gitCtx.RepoPath, gitCtx.BranchName)
		if err != nil {
			fmt.Printf("❌ Failed to fetch stats: %v\n", err)
			return
		}

		if stats.Total == 0 {
			fmt.Println("🎉 Workspace is empty! No statistics available.")
			return
		}

		percent := (stats.Done * 100) / stats.Total
		barWidth := 20
		filled := (stats.Done * barWidth) / stats.Total
		empty := barWidth - filled
		bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

		ctxLabel := "Global Workspace"
		if gitCtx.BranchName != "" {
			ctxLabel = fmt.Sprintf("Branch [%s]", gitCtx.BranchName)
		}

		fmt.Printf("\n📊 Kairo Analytics — %s\n", ctxLabel)
		fmt.Printf("--------------------------------------------\n")
		fmt.Printf("Progress:    [%s] %d%%\n", bar, percent)
		fmt.Printf("Total Tasks: %d\n", stats.Total)
		fmt.Printf("  ├─ Completed: %d\n", stats.Done)
		fmt.Printf("  └─ Pending:   %d\n", stats.Todo)
		fmt.Printf("\nPriority Breakdown:\n")
		fmt.Printf("  🚨 HIGH:   %d\n", stats.HighPrio)
		fmt.Printf("  ⚡ MEDIUM: %d\n", stats.MedPrio)
		fmt.Printf("  🧊 LOW:    %d\n", stats.LowPrio)
		fmt.Println("--------------------------------------------")
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}