package middlewares

import (
	"mcp-forge/internal/globals"
	"net/http"
	"time"
)

type AccessLogsMiddlewareDependencies struct {
	AppCtx *globals.ApplicationContext
}

type AccessLogsMiddleware struct {
	dependencies AccessLogsMiddlewareDependencies
}

func NewAccessLogsMiddleware(dependencies AccessLogsMiddlewareDependencies) *AccessLogsMiddleware {
	return &AccessLogsMiddleware{
		dependencies: dependencies,
	}
}

func (mw *AccessLogsMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		start := time.Now()
		next.ServeHTTP(rw, req)
		duration := time.Since(start)

		filteredHeaders := req.Header.Clone()
		// Redact selected headers
		for _, redactedHeader := range mw.dependencies.AppCtx.Config.Middleware.AccessLogs.RedactedHeaders {
			tmpHeader := filteredHeaders.Get(redactedHeader)

			if len(tmpHeader) >= 10 {
				filteredHeaders.Set(redactedHeader, tmpHeader[:10]+"***")
				continue
			}
			filteredHeaders.Set(redactedHeader, "***")
		}

		// Exclude selected headers
		for _, excludedHeader := range mw.dependencies.AppCtx.Config.Middleware.AccessLogs.ExcludedHeaders {
			filteredHeaders.Del(excludedHeader)
		}

		mw.dependencies.AppCtx.Logger.Info("AccessLogsMiddleware output",
			"method", req.Method,
			"url", req.URL.String(),
			"remote_addr", req.RemoteAddr,
			"user_agent", req.UserAgent(),
			"headers", filteredHeaders,
			"request_duration", duration.String(),
		)
	})
}
