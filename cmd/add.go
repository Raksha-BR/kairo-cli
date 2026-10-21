// cmd/add.go
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Raksha-BR/kairo-cli/internal/ai"
	"github.com/Raksha-BR/kairo-cli/internal/db"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new task to Kairo",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskTitle := strings.Join(args, " ")
		gitCtx := git.GetCurrentContext()

		useAI, _ := cmd.Flags().GetBool("ai")
		if useAI {
			fmt.Printf("🤖 Asking AI to decompose: %q...\n", taskTitle)
			subtasks, err := ai.DecomposeTask(taskTitle)
			if err != nil {
				fmt.Printf("⚠️  AI Breakdown failed (%v). Saving as standard task.\n", err)
				_ = db.InsertTask(taskTitle, gitCtx.RepoPath, gitCtx.BranchName)
				return
			}

			err = db.InsertTaskWithSubtasks(taskTitle, gitCtx.RepoPath, gitCtx.BranchName, subtasks)
			if err != nil {
				fmt.Printf("❌ Failed to save task with subtasks: %v\n", err)
				return
			}

			fmt.Printf("✨ Task decomposed into %d actionable subtasks!\n", len(subtasks))
			for _, sub := range subtasks {
				fmt.Printf("   └─ %s\n", sub)
			}
			return
		}

		// Save Standard Task
		err := db.InsertTask(taskTitle, gitCtx.RepoPath, gitCtx.BranchName)
		if err != nil {
			fmt.Printf("❌ Failed to save task: %v\n", err)
			return
		}

		if gitCtx.BranchName != "" {
			fmt.Printf("✅ Task added to branch [%s]: %q\n", gitCtx.BranchName, taskTitle)
		} else {
			fmt.Printf("✅ Global task added: %q\n", taskTitle)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("ai", "a", false, "Use AI to break down this task")
}