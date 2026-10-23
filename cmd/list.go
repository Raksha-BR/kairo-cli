// cmd/list.go
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Raksha-BR/kairo-cli/internal/db"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List tasks for the current context",
	Run: func(cmd *cobra.Command, args []string) {
		prioFilter, _ := cmd.Flags().GetString("priority")
		prioFilter = strings.ToUpper(prioFilter)

		gitCtx := git.GetCurrentContext()
		tasks, err := db.GetTasks(gitCtx.RepoPath, gitCtx.BranchName)
		if err != nil {
			fmt.Printf("❌ Failed to fetch tasks: %v\n", err)
			return
		}

		if prioFilter != "" {
			var filtered []db.Task
			for _, t := range tasks {
				if t.Priority == prioFilter {
					filtered = append(filtered, t)
				}
			}
			tasks = filtered
		}

		if len(tasks) == 0 {
			fmt.Println("🎉 No tasks found matching your context/filter!")
			return
		}

		doneCount := 0
		for _, t := range tasks {
			if t.Status == "DONE" {
				doneCount++
			}
		}

		if gitCtx.BranchName != "" {
			fmt.Printf("📋 Tasks for branch [%s]", gitCtx.BranchName)
		} else {
			fmt.Printf("📋 Global tasks")
		}
		if prioFilter != "" {
			fmt.Printf(" (%s priority)", prioFilter)
		}
		fmt.Printf(" • %s\n\n", renderProgressBar(doneCount, len(tasks)))

		for i, t := range tasks {
			checkbox := "[ ]"
			if t.Status == "DONE" {
				checkbox = "[x]"
			}
			fmt.Printf("  %s %d. [%s] %s\n", checkbox, i+1, t.Priority, t.Title)
		}
	},
}

func renderProgressBar(doneCount, totalCount int) string {
	if totalCount == 0 {
		return ""
	}
	percent := (doneCount * 100) / totalCount
	width := 10
	filled := (doneCount * width) / totalCount
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return fmt.Sprintf("[%s] %d%% (%d/%d done)", bar, percent, doneCount, totalCount)
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("priority", "p", "", "Filter tasks by priority: HIGH, MEDIUM, or LOW")
}