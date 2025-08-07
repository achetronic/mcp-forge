package middlewares

import (
	"context"

	//
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type NoopMiddlewareDependencies struct{}

type NoopMiddleware struct {
	dependencies NoopMiddlewareDependencies
}

func NewNoopMiddleware(dependencies NoopMiddlewareDependencies) *NoopMiddleware {
	return &NoopMiddleware{
		dependencies: dependencies,
	}
}

func (mw *NoopMiddleware) ToolMiddleware(next server.ToolHandlerFunc) server.ToolHandlerFunc {

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return next(ctx, request)
	}
}
