// cmd/export.go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Raksha-BR/kairo-cli/internal/db"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export workspace tasks to Markdown checklist or JSON",
	Long:  "Export active tasks for the current Git workspace formatted as a Markdown checklist for PRs or documentation.",
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		outFile, _ := cmd.Flags().GetString("out")

		gitCtx := git.GetCurrentContext()
		tasks, err := db.GetTasks(gitCtx.RepoPath, gitCtx.BranchName)
		if err != nil {
			fmt.Printf("❌ Failed to fetch tasks: %v\n", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("🎉 No tasks to export for this context!")
			return
		}

		var output string

		if strings.ToLower(format) == "json" {
			data, err := json.MarshalIndent(tasks, "", "  ")
			if err != nil {
				fmt.Printf("❌ JSON export error: %v\n", err)
				return
			}
			output = string(data)
		} else {
			// Markdown format (Default)
			header := "### 📋 Workspace Tasks"
			if gitCtx.BranchName != "" {
				header = fmt.Sprintf("### 📋 Workspace Tasks (`%s`)", gitCtx.BranchName)
			}
			var builder strings.Builder
			builder.WriteString(header + "\n\n")

			for _, t := range tasks {
				check := " "
				if t.Status == "DONE" {
					check = "x"
				}
				builder.WriteString(fmt.Sprintf("- [%s] **[%s]** %s\n", check, t.Priority, t.Title))
			}
			output = builder.String()
		}

		if outFile != "" {
			err := os.WriteFile(outFile, []byte(output), 0644)
			if err != nil {
				fmt.Printf("❌ Failed to save export file: %v\n", err)
				return
			}
			fmt.Printf("✅ Exported %d tasks to %s!\n", len(tasks), outFile)
		} else {
			fmt.Println("\n--- Exported Markdown Output ---")
			fmt.Println(output)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("format", "f", "markdown", "Export format: 'markdown' or 'json'")
	exportCmd.Flags().StringP("out", "o", "", "Output file path (e.g., PR_CHECKLIST.md)")
}