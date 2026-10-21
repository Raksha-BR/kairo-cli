// cmd/rm.go
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/Raksha-BR/kairo-cli/internal/db"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

var rmCmd = &cobra.Command{
	Use:     "rm [task_number]",
	Aliases: []string{"delete", "remove"},
	Short:   "Delete a task from Kairo",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil || index < 1 {
			fmt.Println("❌ Please provide a valid task number (e.g., 'kairo rm 1').")
			return
		}

		gitCtx := git.GetCurrentContext()
		tasks, err := db.GetTasks(gitCtx.RepoPath, gitCtx.BranchName)
		if err != nil {
			fmt.Printf("❌ Failed to fetch tasks: %v\n", err)
			return
		}

		if index > len(tasks) {
			fmt.Printf("❌ Invalid task number %d. You only have %d tasks in this context.\n", index, len(tasks))
			return
		}

		targetTask := tasks[index-1]
		if err := db.DeleteTask(targetTask.ID); err != nil {
			fmt.Printf("❌ Failed to delete task: %v\n", err)
			return
		}

		fmt.Printf("🗑️  Deleted task #%d: %q\n", index, targetTask.Title)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}