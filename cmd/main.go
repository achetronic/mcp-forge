package main

import (
	"log"
	"net/http"
	"time"

	//
	"mcp-server-template/internal/globals"
	"mcp-server-template/internal/handlers"
	"mcp-server-template/internal/middlewares"
	"mcp-server-template/internal/tools"

	//
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// 0. Process the configuration
	environmentOptions, err := globals.GetEnvironmentOptions()
	if err != nil {
		globals.Logger.Error("error processing environment vars", "error", err.Error())
		return
	}
	globals.Environment = environmentOptions

	// 1. Initialize middlewares that need it
	accessLogsMw := middlewares.NewAccessLogsMiddleware(middlewares.AccessLogsMiddlewareDependencies{
		Logger: globals.Logger,
	})

	jwtValidationMw := middlewares.NewJWTValidationMiddleware(middlewares.JWTValidationMiddlewareDependencies{
		Logger: globals.Logger,
	})

	// 2. Create a new MCP server
	mcpServer := server.NewMCPServer(
		"MCP Server Template",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// 3. Add some useful magic in the form of tools to your MCP server
	// This is the most useful part
	tm := tools.NewToolsManager(tools.ToolsManagerDependencies{
		McpServer:   mcpServer,
		Middlewares: []middlewares.ToolMiddleware{},
	})
	tm.AddTools()

	// 3. Wrap MCP server in a transport (stdio, HTTP, SSE)
	switch globals.Environment.ServerTransport {
	case globals.ServerTransportHttp:
		httpServer := server.NewStreamableHTTPServer(mcpServer,
			server.WithHeartbeatInterval(30*time.Second),
			server.WithStateLess(false))

		// Register it under a path, then add custom endpoints.
		// Custom endpoints are needed as the library is not feature-complete according to MCP spec requirements (2025-06-16)
		// Ref: https://modelcontextprotocol.io/specification/2025-06-18/basic/authorization#overview
		mux := http.NewServeMux()
		mux.Handle("/mcp", accessLogsMw.Middleware(jwtValidationMw.Middleware(httpServer)))
		mux.Handle("/.well-known/oauth-protected-resource", accessLogsMw.Middleware(http.HandlerFunc(handlers.HandleOauthProtectedResources)))

		// Start StreamableHTTP server
		globals.Logger.Info("Starting StreamableHTTP server", "host", globals.Environment.ServerTransportHttpHost)
		err := http.ListenAndServe(globals.Environment.ServerTransportHttpHost, mux)
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
