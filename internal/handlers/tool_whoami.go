package handlers

import (
	"context"
	"fmt"
	"log"

	//
	"github.com/mark3labs/mcp-go/mcp"
)

func HandleToolWhoami(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	log.Print("request: ", request.Request)
	log.Print("header: ", request.Header)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("I am gilbertito"),
			},
		},
	}, nil
}
