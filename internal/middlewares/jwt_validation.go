package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	//
	"mcp-server-template/internal/globals"
)

type JWTValidationMiddlewareDependencies struct {
	AppCtx *globals.ApplicationContext
}

type JWTValidationMiddleware struct {
	dependencies JWTValidationMiddlewareDependencies

	// carried stuff
	jwks  *JWKS
	mutex sync.Mutex
}

func NewJWTValidationMiddleware(deps JWTValidationMiddlewareDependencies) *JWTValidationMiddleware {

	mw := &JWTValidationMiddleware{
		dependencies: deps,
	}

	// Launch JWKS worker only when requested
	if mw.dependencies.AppCtx.Config.Middleware.JWT.Enabled &&
		mw.dependencies.AppCtx.Config.Middleware.JWT.Validation.Strategy == "local" {
		go mw.cacheJWKS()
	}

	return mw
}

func (mw *JWTValidationMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if !mw.dependencies.AppCtx.Config.Middleware.JWT.Enabled {
			goto nextStage
		}

		switch mw.dependencies.AppCtx.Config.Middleware.JWT.Validation.Strategy {
		case "local":
			// 1. Extract token from header
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(rw, "Authorization header not found", http.StatusUnauthorized)
				return
			}
			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

			// Reject unauthorized requests
			_, err := mw.isTokenValid(tokenString)
			if err != nil {
				http.Error(rw, fmt.Sprintf("Invalid token: %v", err.Error()), http.StatusUnauthorized)
				return
			}

			// Put the JWT into the validated request header
			req.Header.Set(mw.dependencies.AppCtx.Config.Middleware.JWT.Validation.ForwardedHeader, tokenString)
		default:
			// Having a validated JWT into a specific header is the default behavior,
			// as having tools like Istio securing APIs is much more safe and reliable
			// When the token is already validated, do nothing.
		}

	nextStage:
		next.ServeHTTP(rw, req)
	})
}
