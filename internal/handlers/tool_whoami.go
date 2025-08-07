package handlers

import (
	"context"
	"fmt"

	//
	"mcp-server-template/internal/globals"

	//
	"github.com/mark3labs/mcp-go/mcp"
)

func HandleToolWhoami(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	validatedJwt := request.Header.Get(globals.Environment.ServerTransportHttpJwtValidatedHeader)

	if validatedJwt == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error: JWT is empty. Information is not available"),
				},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Success! Data are in the following JWT. You have to decode it first: %s", validatedJwt),
			},
		},
	}, nil
}
