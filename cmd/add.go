// cmd/add.go
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Raksha-BR/kairo-cli/internal/db"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new task to Kairo",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskTitle := strings.Join(args, " ")
		
		useAI, _ := cmd.Flags().GetBool("ai")
		if useAI {
			fmt.Printf("🤖 AI Breakdown triggered for: %q\n", taskTitle)
			fmt.Println("   (Note: AI connection will be implemented in Phase 5)")
			return
		}

		// Get the current Git context
		gitCtx := git.GetCurrentContext()
		
		// Save to Database with Git context
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