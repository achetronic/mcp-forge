package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

type AccessLogsMiddlewareDependencies struct {
	Logger *slog.Logger
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

		mw.dependencies.Logger.Info("AccessLogsMiddleware output",
			"method", req.Method,
			"url", req.URL.String(),
			"remote_addr", req.RemoteAddr,
			"user_agent", req.UserAgent(),
			"headers", req.Header,
			"request_duration", duration.String(),
		)
	})
}
