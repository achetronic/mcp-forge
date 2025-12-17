package middlewares

import (
	"net/http"
	"strings"
)

// getRequestScheme returns the scheme of a given request.
// It is inferred from headers and connection in cascade.
func getRequestScheme(req *http.Request) string {

	// Priority 1: X-Forwarded-Proto (Istio/Envoy/Nginx)
	if proto := req.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}

	// Priority 2: Forwarded header (RFC 7239)
	if fwd := req.Header.Get("Forwarded"); fwd != "" {
		if strings.Contains(fwd, "proto=https") {
			return "https"
		}
	}

	// Priority 3: Direct TLS connection
	if req.TLS != nil {
		return "https"
	}

	return "http"
}
