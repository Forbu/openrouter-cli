package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgAPIBase string
	cfgAPIKey  string
	cfgModel   string
	cfgTemp    float64
	cfgStream  bool
)

var rootCmd = &cobra.Command{
	Use:   "openroutercli",
	Short: "CLI to interact with OpenRouter-compatible LLM APIs",
	Long:  "A fast, ergonomic CLI to chat with LLMs via OpenRouter-compatible APIs.",
}

func init() {
	// Defaults aligned with user's preferences
	rootCmd.PersistentFlags().StringVar(&cfgAPIBase, "base-url", "https://openrouter.ai/api/v1", "API base URL (OpenRouter-compatible)")
	rootCmd.PersistentFlags().StringVar(&cfgAPIKey, "api-key", "", "API key (overrides env OPENROUTER_API_KEY)")
	rootCmd.PersistentFlags().StringVar(&cfgModel, "model", "openai/gpt-4o-mini", "Default model")
	rootCmd.PersistentFlags().Float64Var(&cfgTemp, "temperature", 0.7, "Sampling temperature")
	rootCmd.PersistentFlags().BoolVar(&cfgStream, "stream", false, "Stream responses when supported")

	// Env fallback
	if v := os.Getenv("OPENROUTER_API_KEY"); v != "" {
		cfgAPIKey = v
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
