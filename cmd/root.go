// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kairo",
	Short: "Kairo is a context-aware terminal task engine",
	Long: `Kairo (Kairos - the opportune moment) is a blazingly fast, 
context-aware, and AI-powered terminal task engine for developers.
It automatically tracks your Git branches and uses AI to break down complex tasks.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.kairo.yaml)")
}