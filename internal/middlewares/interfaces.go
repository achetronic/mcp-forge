package middlewares

import (
	"net/http"

	//
	"github.com/mark3labs/mcp-go/server"
)

type ToolMiddleware interface {
	Middleware(next server.ToolHandlerFunc) server.ToolHandlerFunc
}

type HttpMiddleware interface {
	Middleware(next http.Handler) http.Handler
}
