package tools

import (
	"context"
	"fmt"

	//
	"github.com/mark3labs/mcp-go/mcp"
)

func (tm *ToolsManager) HandleToolWhoami(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	validatedJwt := request.Header.Get(tm.dependencies.AppCtx.Config.Middleware.JWT.Validation.ForwardedHeader)

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
