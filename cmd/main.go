package main

import (
	"github.com/sethvargo/go-envconfig"
	"log"
	"net/http"
	"time"
	"tiny-mcp/internal/tools"

	//
	"tiny-mcp/internal/globals"
	"tiny-mcp/internal/handlers"
	//
	"github.com/mark3labs/mcp-go/server"
)

const ServerTransportHttp = "http"

type ServerOptions struct {
	Transport         string `env:"SERVER_TRANSPORT"`
	TransportHttpHost string `env:"SERVER_TRANSPORT_HTTP_HOST"`

	TransportHttpForwardedJwtHeader    string `env:"SERVER_TRANSPORT_HTTP_FORWARDED_JWT_HEADER"`
	TransportHttpJwtValidationStrategy string `env:"SERVER_TRANSPORT_HTTP_JWT_VALIDATION_STRATEGY"`
}

func main() {

	// 1. Create a new MCP server
	mcpServer := server.NewMCPServer(
		"Tiny MCP Server",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// 2. Add some useful magic in the form of tools to your MCP server
	// This is the most useful part
	tools.AddTools(mcpServer)

	// 3. Wrap MCP server in a transport (stdio, HTTP, SSE)
	ServerOptionsObject := &ServerOptions{}
	if err := envconfig.Process(globals.Context, ServerOptionsObject); err != nil {
		globals.Logger.Error("error processing environment vars", "error", err.Error())
		return
	}

	switch ServerOptionsObject.Transport {
	case ServerTransportHttp:
		httpServer := server.NewStreamableHTTPServer(mcpServer,
			server.WithHeartbeatInterval(30*time.Second),
			server.WithStateLess(false))

		// Register it under a path, then add custom endpoints.
		// Custom endpoints are needed as the library is not feature-complete according to MCP spec requirements (2025-06-16)
		// Ref: https://modelcontextprotocol.io/specification/2025-06-18/basic/authorization#overview
		mux := http.NewServeMux()
		mux.Handle("/mcp", httpServer)
		mux.HandleFunc("/.well-known/oauth-protected-resource", handlers.HandleOauthProtectedResources)

		// Start StreamableHTTP server
		globals.Logger.Info("Starting StreamableHTTP server", "host", ServerOptionsObject.TransportHttpHost)
		err := http.ListenAndServe(ServerOptionsObject.TransportHttpHost, mux)
		if err != nil {
			log.Fatal(err)
		}

	default:
		// Start stdio server
		globals.Logger.Info("Starting stdio server")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatal(err)
		}
	}
}
