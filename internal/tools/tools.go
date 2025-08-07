package tools

import (
	"tiny-mcp/internal/handlers"

	//
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func AddTools(mcpServer *server.MCPServer) {

	// 1. Describe a tool, then add it
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
	mcpServer.AddTool(tool, handlers.HandleToolHello)

	// 2. Describe and add another tool
	tool = mcp.NewTool("whoami",
		mcp.WithDescription("Expose information about the user"),
	)
	mcpServer.AddTool(tool, handlers.HandleToolWhoami)
}
