# go-sportsagent

## Overview

Agentic application using OpenAI function calling to interact with sports data services.

## Architecture

**Server**: HTTP server on `:8080` with `/query` endpoint

**Agent**: OpenAI GPT-4o with function calling - AI decides which service to call based on user query

**Services**:

- `rotoreader` (localhost:8001) - GET /feed - sports news feed
- `oddstracker` (localhost:8000) - GET /changes - betting odds changes

## Structure

```text
main.go                 # HTTP server entry
internal/
  handlers/            # HTTP request handlers
  services/            # Agent logic with OpenAI integration
  clients/             # HTTP clients for external services
  tools/               # OpenAI function definitions
```

## Configuration

Uses `.env` file (loaded via godotenv):

- `OPENAI_API_KEY` - required
- `ROTOREADER_URL` - optional, defaults to localhost:8001
- `ODDSTRACKER_URL` - optional, defaults to localhost:8000

## Testing

- **Unit tests**: Configuration validation, no external dependencies
- **Integration tests**: Run with `INTEGRATION_TESTS=1` to test actual service connectivity
- `.env` auto-loads in tests via `init()` function
- Focus on positive flows only

## Development Commands

```bash
make test              # unit tests
make test-integration  # integration tests (requires services running)
make run              # start server
make build            # build binary
```

## Function Calling

Functions currently have no parameters (services don't accept them yet).

When services add parameters (e.g., filtering by team/sport), update `internal/tools/definitions.go`:

```go
Parameters: openai.FunctionParameters{
    "type": "object",
    "properties": map[string]interface{}{
        "team": map[string]string{
            "type":        "string",
            "description": "Team name to filter for",
        },
    },
    "required": []string{"team"},
},
```

The agent service extracts parameters from function calls and constructs the appropriate HTTP requests (including path parameters if needed).
