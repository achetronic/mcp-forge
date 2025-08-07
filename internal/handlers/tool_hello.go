package handlers

import (
	"context"
	"fmt"
	"log"
	"tiny-mcp/internal/globals"

	//
	"github.com/mark3labs/mcp-go/mcp"
)

func HandleToolHello(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	name, ok := arguments["name"].(string)
	if !ok {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: name parameter is required and must be a string",
				},
			},
			IsError: true,
		}, nil
	}

	globals.Logger.Debug("showing parameters", "arguments", arguments, "request", request)

	log.Print("request: ", request.Request)
	log.Print("header: ", request.Header)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Hello, %s! 👋", name),
			},
		},
	}, nil
}
