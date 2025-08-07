# Tiny MCP Server

## Description
TBD

## Motivation
TBD

## What to expect
TBD

## How to deploy
TBD

## How to develop
TBD

```json5
// file: claude_desktop_config.json

{
  "mcpServers": {
    
    // After executing: make build
    "stdio": {
      "command": "/home/example/tiny-mcp/bin/tiny-mcp-linux-amd64",
      "env": {
        "PLACEHOLDER": "placeholder"
      }
    },
    
    // After executing: npm i mcp-remote && source .env && make run
    "local-proxy-remote": {
      "command": "npx",
      "args": [
        "mcp-remote",
        "http://localhost:8080/mcp",
        "--transport",
        "http-only",
        "--header",
        "Authorization: Bearer ${JWT}"
      ],
      "env": {
        "JWT": "eyJhbGciOiJSUzI1NiIsImtpZCI6..."
      }
    }
  }
}
```

## References
1. https://www.npmjs.com/package/mcp-remote
2. https://datatracker.ietf.org/doc/rfc9728/
3. https://mcp-go.dev/getting-started
4. 