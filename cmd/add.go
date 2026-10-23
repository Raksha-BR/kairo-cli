// cmd/add.go
package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

		priority, _ := cmd.Flags().GetString("priority")
		priority = strings.ToUpper(priority)
		if priority != "HIGH" && priority != "MEDIUM" && priority != "LOW" {
			priority = "MEDIUM"
		}

		dueStr, _ := cmd.Flags().GetString("due")
		dueDate := parseRelativeDate(dueStr)

		parentID, err := db.InsertTask(taskTitle, priority, dueDate, gitCtx.RepoPath, gitCtx.BranchName)
		if err != nil {
			fmt.Printf("❌ Failed to save task: %v\n", err)
			return
		}

		useAI, _ := cmd.Flags().GetBool("ai")
		if useAI {
			fmt.Printf("🤖 Asking AI to decompose: %q...\n", taskTitle)
			subtasks, err := ai.DecomposeTask(taskTitle)
			if err != nil {
				fmt.Printf("⚠️  AI Breakdown failed (%v). Saved standard task.\n", err)
				return
			}

			_ = db.InsertSubtasksForParent(parentID, subtasks, priority, gitCtx.RepoPath, gitCtx.BranchName)
			fmt.Printf("✨ Task decomposed into %d subtasks under parent ID [%s]!\n", len(subtasks), parentID[:8])
			for _, sub := range subtasks {
				fmt.Printf("   └─ %s\n", sub)
			}
			return
		}

		dueNotice := ""
		if dueDate != nil {
			dueNotice = fmt.Sprintf(" (Due: %s)", dueDate.Format("Jan 02"))
		}

		if gitCtx.BranchName != "" {
			fmt.Printf("✅ [%s] Task added to branch [%s]%s: %q\n", priority, gitCtx.BranchName, dueNotice, taskTitle)
		} else {
			fmt.Printf("✅ [%s] Global task added%s: %q\n", priority, dueNotice, taskTitle)
		}
	},
}

func parseRelativeDate(input string) *time.Time {
	if input == "" {
		return nil
	}
	input = strings.ToLower(strings.TrimSpace(input))
	now := time.Now()

	switch input {
	case "today":
		t := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		return &t
	case "tomorrow":
		t := time.Date(now.Year(), now.Month(), now.Day()+1, 23, 59, 59, 0, now.Location())
		return &t
	}

	if strings.HasSuffix(input, "d") {
		daysStr := strings.TrimSuffix(input, "d")
		if days, err := strconv.Atoi(daysStr); err == nil {
			t := now.AddDate(0, 0, days)
			return &t
		}
	}

	if strings.HasSuffix(input, "w") {
		weeksStr := strings.TrimSuffix(input, "w")
		if weeks, err := strconv.Atoi(weeksStr); err == nil {
			t := now.AddDate(0, 0, weeks*7)
			return &t
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("ai", "a", false, "Use AI to break down this task")
	addCmd.Flags().StringP("priority", "p", "MEDIUM", "Task priority: HIGH, MEDIUM, or LOW")
	addCmd.Flags().StringP("due", "u", "", "Due date: 'today', 'tomorrow', '3d', or '1w'")
}