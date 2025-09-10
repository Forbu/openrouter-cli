# openrouter-cli

A fast, minimal command-line tool to chat with OpenRouter (and OpenAI-compatible) LLM APIs. Defaults to OpenRouter.

- Chat from args or stdin
- Streaming output with `--stream`
- List available models
- Simple config via flags and environment variables

## Requirements
- Go 1.21+ (module uses `go 1.21`)
- An API key for your provider (OpenRouter by default)

## Install / Build
```bash
# From the project root
go build -o openroutercli .

# Optional: install into GOPATH/bin or GOBIN
# go install ./...
```

## Quick Start
```bash
# 1) Set your API key (OpenRouter by default)
export OPENROUTER_API_KEY=YOUR_KEY

# 2) One-off chat
./openroutercli chat "Say hi in one word"

# 3) From stdin
echo "Explain RAG briefly" | ./openroutercli chat

# 4) Streaming output
./openroutercli --stream chat "Write a haiku about Go"

# 5) List models
./openroutercli models
```

## Configuration
You can configure the tool via flags or environment variables.

- Environment variables
  - `OPENROUTER_API_KEY`: API key (used by default)

- Global flags
  - `--base-url`: API base URL (default: `https://openrouter.ai/api/v1`)
  - `--api-key`: API key override (takes precedence over env var)
  - `--model`: Model name (default: `openai/gpt-4o-mini`)
  - `--temperature`: Sampling temperature (default: `0.7`)
  - `--stream`: Stream responses if supported (default: `false`)

Examples:
```bash
# Use a different model
./openroutercli --model openrouter/auto chat "What can you do?"

# Override API key from a different env var or secret manager
./openroutercli --api-key "$MY_TEMP_KEY" chat "Test"

# Use an OpenAI base URL and key
./openroutercli \
  --base-url https://api.openai.com/v1 \
  --api-key "$OPENAI_API_KEY" \
  --model gpt-4o-mini \
  chat "Hello"
```

## Commands
- `chat [prompt]` — Send a one-off prompt. If no prompt args are provided, reads from stdin.
- `models` — List available model IDs from the provider.

Run `--help` for details:
```bash
./openroutercli --help
./openroutercli chat --help
./openroutercli models --help
```

## Notes
- This CLI sets `HTTP-Referer` and `X-Title` headers for OpenRouter compatibility and analytics. Adjust in code if needed.
- Non-zero exit codes are returned on errors (useful for scripting).
- No conversation history/state is stored; each `chat` call is a single-turn request.

## Development
```bash
# Tidy modules
go mod tidy

# Run without building a separate binary
go run . chat "Ping"

# Build
go build -o openroutercli .
```

## Security
- Do not commit or echo your API keys. Prefer environment variables or your secret manager.
- Be cautious when piping prompts containing sensitive data.
