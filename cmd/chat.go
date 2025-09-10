package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	FinishReason string      `json:"finish_reason"`
	Message      ChatMessage `json:"message"`
	Delta        *ChatMessage `json:"delta,omitempty"`
}

type ChatResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
}

var chatCmd = &cobra.Command{
	Use:   "chat [prompt]",
	Short: "Send a chat prompt to the model",
	Long:  "Send a one-off prompt to the model. Provide [prompt] or pipe via stdin.",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := strings.TrimSpace(strings.Join(args, " "))
		if prompt == "" {
			// Read from stdin if no args
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				b, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				prompt = strings.TrimSpace(string(b))
			}
		}
		if prompt == "" {
			return errors.New("no prompt provided: pass as args or pipe via stdin")
		}

		if cfgAPIKey == "" {
			return errors.New("missing API key: set OPENROUTER_API_KEY or --api-key")
		}

		ctx, cancel := context.WithTimeout(cmd.Context(), 2*time.Minute)
		defer cancel()

		messages := []ChatMessage{{Role: "user", Content: prompt}}
		reqBody := ChatRequest{
			Model:       cfgModel,
			Messages:    messages,
			Temperature: cfgTemp,
			Stream:      cfgStream,
		}

		b, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}

		url := strings.TrimRight(cfgAPIBase, "/") + "/chat/completions"
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+cfgAPIKey)
		// Optional but recommended by OpenRouter
		req.Header.Set("HTTP-Referer", "https://github.com/adrienbufort/openrouter-cli")
		req.Header.Set("X-Title", "openrouter-cli")

		client := &http.Client{Timeout: 0}
		if !cfgStream {
			client.Timeout = 60 * time.Second
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
		}

		if !cfgStream {
			var cr ChatResponse
			if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
				return err
			}
			if len(cr.Choices) == 0 {
				return errors.New("no choices returned")
			}
			fmt.Println(strings.TrimSpace(cr.Choices[0].Message.Content))
			return nil
		}

		// Stream mode (SSE-like: data: {json}\n\n)
		r := bufio.NewReader(resp.Body)
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}
			line = strings.TrimSpace(line)
			if line == "" || !strings.HasPrefix(line, "data:") {
				continue
			}
			payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if payload == "[DONE]" {
				break
			}
			var chunk ChatResponse
			if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
				continue
			}
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
				fmt.Print(chunk.Choices[0].Delta.Content)
			}
		}
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
