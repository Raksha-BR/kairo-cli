// cmd/clear.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Raksha-BR/kairo-cli/internal/db"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

var clearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"purge", "clean"},
	Short:   "Clear all completed tasks in the active workspace",
	Run: func(cmd *cobra.Command, args []string) {
		gitCtx := git.GetCurrentContext()
		rowsDeleted, err := db.PurgeDoneTasks(gitCtx.RepoPath, gitCtx.BranchName)
		if err != nil {
			fmt.Printf("❌ Failed to clear tasks: %v\n", err)
			return
		}

		if rowsDeleted == 0 {
			fmt.Println("ℹ️ No completed tasks to clear.")
			return
		}

		fmt.Printf("🧹 Cleared %d completed task(s) from this workspace!\n", rowsDeleted)
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}