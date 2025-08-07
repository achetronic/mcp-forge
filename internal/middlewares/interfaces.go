package middlewares

import (
	"github.com/mark3labs/mcp-go/server"
	"net/http"
)

type ToolMiddleware interface {
	Middleware(next server.ToolHandlerFunc) server.ToolHandlerFunc
}

type HttpMiddleware interface {
	Middleware(next http.Handler) http.Handler
}
