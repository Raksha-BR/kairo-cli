// cmd/list.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Raksha-BR/kairo-cli/internal/db"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all tasks for the current context",
	Run: func(cmd *cobra.Command, args []string) {
		
		gitCtx := git.GetCurrentContext()
		tasks, err := db.GetTasks(gitCtx.RepoPath, gitCtx.BranchName)
		if err != nil {
			fmt.Printf("❌ Failed to fetch tasks: %v\n", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("🎉 No tasks found for this context! You're all caught up.")
			return
		}

		if gitCtx.BranchName != "" {
			fmt.Printf("📋 Tasks for branch [%s]:\n", gitCtx.BranchName)
		} else {
			fmt.Printf("📋 Your global tasks:\n")
		}
		
		for i, t := range tasks {
			checkbox := "[ ]"
			if t.Status == "DONE" {
				checkbox = "[x]"
			}
			fmt.Printf("  %s %d. %s\n", checkbox, i+1, t.Title)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}