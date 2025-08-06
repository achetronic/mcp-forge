package main

import (
	"log"
	"net/http"
	"time"

	//
	"tiny-mcp/internal/globals"
	"tiny-mcp/internal/handlers"
	"tiny-mcp/internal/tools"

	//
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// 1. Create a new MCP server
	mcpServer := server.NewMCPServer(
		"Tiny MCP Server",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// 2. Wrap MCP server in HTTP transport
	// It's possible to use another
	httpServer := server.NewStreamableHTTPServer(mcpServer,
		server.WithHeartbeatInterval(30*time.Second),
		server.WithStateLess(false))

	// 3. Register it under a path, then add custom endpoints.
	// Custom endpoints are needed as the library is not feature-complete according to MCP spec requirements (2025-06-16)
	// Ref: https://modelcontextprotocol.io/specification/2025-06-18/basic/authorization#overview
	mux := http.NewServeMux()
	mux.Handle("/mcp/", httpServer)
	mux.HandleFunc("/.well-known/oauth-protected-resource", handlers.HandleOauthProtectedResources)

	// 4. Add some useful magic in the form of tools to your MCP server
	// This is the most useful part
	tools.AddTools(mcpServer)

	// Start StreamableHTTP server
	globals.Logger.Info("Starting StreamableHTTP server on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}

}
