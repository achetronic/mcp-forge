package tools

import (
	"mcp-forge/internal/globals"
	"mcp-forge/internal/middlewares"

	//
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ToolsManagerDependencies struct {
	AppCtx *globals.ApplicationContext

	McpServer   *server.MCPServer
	Middlewares []middlewares.ToolMiddleware
}

type ToolsManager struct {
	dependencies ToolsManagerDependencies
}

func NewToolsManager(deps ToolsManagerDependencies) *ToolsManager {
	return &ToolsManager{
		dependencies: deps,
	}
}

func (tm *ToolsManager) AddTools() {

	// 1. Describe a tool, then add it
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
	tm.dependencies.McpServer.AddTool(tool, tm.HandleToolHello)

	// 2. Describe and add another tool
	tool = mcp.NewTool("whoami",
		mcp.WithDescription("Expose information about the user"),
	)
	tm.dependencies.McpServer.AddTool(tool, tm.HandleToolWhoami)
}
