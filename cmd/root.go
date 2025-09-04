package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-openrouter",
	Short: "A CLI to interact with the OpenRouter API",
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommand is provided
		fmt.Println("Welcome to go-openrouter CLI!")
		fmt.Println("Use 'go-openrouter --help' to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
