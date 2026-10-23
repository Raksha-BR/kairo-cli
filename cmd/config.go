// cmd/config.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Raksha-BR/kairo-cli/internal/config"
	"github.com/Raksha-BR/kairo-cli/internal/git"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"info", "settings"},
	Short:   "Display active Kairo configuration and workspace info",
	Run: func(cmd *cobra.Command, args []string) {
		gitCtx := git.GetCurrentContext()

		fmt.Println("\n⚙️  Kairo System Configuration")
		fmt.Println("--------------------------------------------")
		fmt.Printf(" AI Provider:   %s\n", config.AppConfig.AI.Provider)
		fmt.Printf(" AI Model:      %s\n", config.AppConfig.AI.Model)
		fmt.Printf(" Ollama URL:    %s\n", config.AppConfig.AI.OllamaURL)
		if config.AppConfig.AI.OpenAIKey != "" {
			fmt.Println(" OpenAI Key:    [CONFIGURED]")
		} else {
			fmt.Println(" OpenAI Key:    [NOT SET]")
		}
		fmt.Println("--------------------------------------------")
		fmt.Println("📍 Active Workspace Context")
		if gitCtx.BranchName != "" {
			fmt.Printf(" Git Repo:      %s\n", gitCtx.RepoPath)
			fmt.Printf(" Git Branch:    %s\n", gitCtx.BranchName)
		} else {
			fmt.Println(" Git Context:   Global Workspace (Not inside a Git repo)")
		}
		fmt.Println(" Config File:   ~/.kairo.yaml")
		fmt.Println(" SQLite DB:     ~/.kairo.db")
		fmt.Println("--------------------------------------------")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}