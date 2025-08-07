package middlewares

import (
	"log/slog"
	"net/http"

	//
	"tiny-mcp/internal/globals"
)

type JWTValidationMiddlewareDependencies struct {
	Logger *slog.Logger
}

type JWTValidationMiddleware struct {
	dependencies JWTValidationMiddlewareDependencies
}

func NewJWTValidationMiddleware(dependencies JWTValidationMiddlewareDependencies) *JWTValidationMiddleware {
	return &JWTValidationMiddleware{
		dependencies: dependencies,
	}
}

func (mw *JWTValidationMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		switch globals.Environment.ServerTransportHttpJwtValidationStrategy {

		case globals.ServerTransportHttpJwtValidationStrategyLocal:
			// Local validation against an external OIDP is pending review.
			// The code is done, but not reviewed yet.

			// 1. Validate the JWT against the JWKS
			// 2. Reject unauthorized requests
			// 3. Put the JWT into the validated request header
		default:
			// Having a validated JWT into a specific header is the default
			// as having tools like Istio securing APIs is much more safe and reliable

			// When the token is already validated, do nothing
		}
		next.ServeHTTP(rw, req)
	})
}
