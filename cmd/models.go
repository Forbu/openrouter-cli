package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type ModelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List available models",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgAPIKey == "" {
			return errors.New("missing API key: set OPENROUTER_API_KEY or --api-key")
		}

		ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
		defer cancel()

		url := strings.TrimRight(cfgAPIBase, "/") + "/models"
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewReader(nil))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+cfgAPIKey)
		req.Header.Set("HTTP-Referer", "https://github.com/adrienbufort/openrouter-cli")
		req.Header.Set("X-Title", "openrouter-cli")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
		}

		var mr ModelsResponse
		if err := json.NewDecoder(resp.Body).Decode(&mr); err != nil {
			return err
		}
		for _, m := range mr.Data {
			fmt.Println(m.ID)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(modelsCmd)
}
