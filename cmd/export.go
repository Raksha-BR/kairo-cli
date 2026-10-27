// cmd/export.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Raksha-BR/kairo-cli/internal/db"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

var exportCmd = &cobra.Command{
	Use:     "export [filename]",
	Aliases: []string{"pr", "md"},
	Short:   "Export active workspace tasks as a Markdown PR checklist",
	Long:    `Export tasks and AI subtasks for the current Git branch formatted as a Markdown checklist for GitHub PR descriptions.`,
	Run: func(cmd *cobra.Command, args []string) {
		gitCtx := git.GetCurrentContext()
		tasks, err := db.GetTasks(gitCtx.RepoPath, gitCtx.BranchName)
		if err != nil {
			fmt.Printf("❌ Failed to fetch tasks: %v\n", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("ℹ️  No tasks found in this workspace context to export.")
			return
		}

		md := "## 📋 PR Checklist & Task Breakdown\n\n"
		if gitCtx.BranchName != "" {
			md += fmt.Sprintf("> **Workspace Context:** `%s` (Branch: `%s`)\n\n", gitCtx.RepoPath, gitCtx.BranchName)
		} else {
			md += "> **Workspace Context:** Global Workspace\n\n"
		}

		for _, t := range tasks {
			box := "[ ]"
			if t.Status == "DONE" {
				box = "[x]"
			}

			indent := ""
			if t.ParentID.Valid {
				indent = "  "
			}

			md += fmt.Sprintf("%s- %s **[%s]** %s\n", indent, box, t.Priority, t.Title)
		}

		filename, _ := cmd.Flags().GetString("output")
		if len(args) > 0 {
			filename = args[0]
		}

		if filename != "" {
			err := os.WriteFile(filename, []byte(md), 0644)
			if err != nil {
				fmt.Printf("❌ Failed to write export file: %v\n", err)
				return
			}
			fmt.Printf("✅ Exported PR checklist to: %s\n", filename)
		} else {
			fmt.Println(md)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("output", "o", "", "Output markdown filename (e.g. PR_CHECKLIST.md)")
}