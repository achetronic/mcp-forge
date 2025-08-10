package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	//
	"mcp-forge/internal/globals"

	//
	"github.com/google/cel-go/cel"
)

type JWTValidationMiddlewareDependencies struct {
	AppCtx *globals.ApplicationContext
}

type JWTValidationMiddleware struct {
	dependencies JWTValidationMiddlewareDependencies

	// Carried stuff
	jwks  *JWKS
	mutex sync.Mutex

	//
	celPrograms []*cel.Program
}

func NewJWTValidationMiddleware(deps JWTValidationMiddlewareDependencies) (*JWTValidationMiddleware, error) {

	mw := &JWTValidationMiddleware{
		dependencies: deps,
	}

	// Launch JWKS worker only when requested
	if mw.dependencies.AppCtx.Config.Middleware.JWT.Enabled &&
		mw.dependencies.AppCtx.Config.Middleware.JWT.Validation.Strategy == "local" {
		go mw.cacheJWKS()
	}

	// Precompile and check CEL expressions to fail-fast and safe resources.
	// They will be truly used later.
	allowConditionsEnv, err := cel.NewEnv(
		cel.Variable("payload", cel.DynType),
	)
	if err != nil {
		return nil, fmt.Errorf("CEL environment creation error: %s", err.Error())
	}

	for _, allowCondition := range mw.dependencies.AppCtx.Config.Middleware.JWT.Validation.Local.AllowConditions {

		// Compile and execute the code
		ast, issues := allowConditionsEnv.Compile(allowCondition.Expression)
		if issues != nil && issues.Err() != nil {
			return nil, fmt.Errorf("CEL expression compilation exited with error: %s", issues.Err())
		}

		prg, err := allowConditionsEnv.Program(ast)
		if err != nil {
			return nil, fmt.Errorf("CEL program construction error: %s", err.Error())
		}
		mw.celPrograms = append(mw.celPrograms, &prg)
	}

	return mw, nil
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
				http.Error(rw, "RBAC: Access Denied: Authorization header not found", http.StatusUnauthorized)
				return
			}
			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

			// Reject unauthorized requests
			_, err := mw.isTokenValid(tokenString)
			if err != nil {
				http.Error(rw, fmt.Sprintf("RBAC: Access Denied: Invalid token: %v", err.Error()), http.StatusUnauthorized)
				return
			}

			// Put the JWT into the validated request header
			req.Header.Set(mw.dependencies.AppCtx.Config.Middleware.JWT.Validation.ForwardedHeader, tokenString)

			// Extract the JWT payload
			tokenStringParts := strings.Split(tokenString, ".")

			// Decode it into a Go's structure for later
			tokenPayloadBytes, err := base64.RawURLEncoding.DecodeString(tokenStringParts[1])
			if err != nil {
				mw.dependencies.AppCtx.Logger.Error("error decoding JWT payload from base64", "error", err.Error())
				http.Error(rw, fmt.Sprintf("RBAC: Access Denied: JWT Payload can not be decoded"), http.StatusUnauthorized)
				return
			}

			tokenPayload := map[string]any{}
			err = json.Unmarshal(tokenPayloadBytes, &tokenPayload)
			if err != nil {
				mw.dependencies.AppCtx.Logger.Error("error decoding JWT payload from JSON", "error", err.Error())
				http.Error(rw, fmt.Sprintf("RBAC: Access Denied: Internal Issue"), http.StatusUnauthorized)
				return
			}

			// Check allowance conditions for the JWT
			// At this point, we assume the JWT is unmarshalled into a golang structure
			for _, celProgram := range mw.celPrograms {
				out, _, err := (*celProgram).Eval(map[string]interface{}{
					"payload": tokenPayload,
				})

				if err != nil {
					mw.dependencies.AppCtx.Logger.Error("CEL program evaluation error", "error", err.Error())
					http.Error(rw, fmt.Sprintf("RBAC: Access Denied: Internal Issue"), http.StatusUnauthorized)
					return
				}

				if out.Value() != true {
					http.Error(rw, fmt.Sprintf("RBAC: Access Denied: JWT does not meet conditions"), http.StatusUnauthorized)
					return
				}
			}

		default:
			// Having a validated JWT into a specific header is the default behavior,
			// as having tools like Istio securing APIs is much more safe and reliable
			// When the token is already validated, do nothing.
		}

	nextStage:
		next.ServeHTTP(rw, req)
	})
}
