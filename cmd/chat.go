package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/revrost/go-openrouter"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an interactive chat session with a model",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")
		if apiKey == "" {
			apiKey = os.Getenv("OPENROUTER_API_KEY")
		}
		if apiKey == "" {
			log.Fatal("Error: API key not provided. Set OPENROUTER_API_KEY or use the --api-key flag.")
		}

		client := openrouter.NewClient(apiKey)
		scanner := bufio.NewScanner(os.Stdin)

		var messages []openrouter.ChatCompletionMessage

		fmt.Println("Starting interactive chat... (type 'exit' to quit)")

		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}
			userInput := scanner.Text()

			if strings.ToLower(userInput) == "exit" {
				fmt.Println("Goodbye!")
				break
			}

			userMessage := openrouter.ChatCompletionMessage{
				Role: "user",
				Content: openrouter.Content{
					Text: userInput,
				},
			}
			messages = append(messages, userMessage)

			req := openrouter.ChatCompletionRequest{
				Model:    "gryphe/mythomax-l2-13b",
				Messages: messages,
			}

			resp, err := client.CreateChatCompletion(cmd.Context(), req)
			if err != nil {
				log.Printf("Completion error: %v\n", err)
				messages = messages[:len(messages)-1]
				continue
			}

			assistantMessage := resp.Choices[0].Message
			fmt.Println("Assistant:", assistantMessage.Content.Text)
			messages = append(messages, assistantMessage)
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().String("api-key", "", "OpenRouter API key")
}
