package handlers

import (
	"encoding/json"
	"net/http"

	//
	"tiny-mcp/internal/globals"

	//
	"github.com/sethvargo/go-envconfig"
)

// OauthProtectedResourceResponse represents the response returned by '.well-known/oauth-protected-resource' endpoint
// According to the RFC9728 (Section 2)
// Ref: https://datatracker.ietf.org/doc/rfc9728/
type OauthProtectedResourceResponse struct {

	// Essential: these are commonly included
	Resource                          string   `json:"resource"                                           env:"OAUTH_PR_RESOURCE,required"`                     // Required
	AuthorizationServers              []string `json:"authorization_servers,omitempty"                    env:"OAUTH_PR_AUTH_SERVERS"`                          // Optional
	JwksUri                           string   `json:"jwks_uri,omitempty"                                 env:"OAUTH_PR_JWKS_URI"`                              // Optional
	ScopesSupported                   []string `json:"scopes_supported,omitempty"                         env:"OAUTH_PR_SCOPES_SUPPORTED"`                      // Optional
	BearerMethodsSupported            []string `json:"bearer_methods_supported,omitempty"                 env:"OAUTH_PR_BEARER_METHODS_SUPPORTED"`              // Optional
	ResourceSigningAlgValuesSupported []string `json:"resource_signing_alg_values_supported,omitempty"    env:"OAUTH_PR_RESOURCE_SIGNING_ALG_VALUES_SUPPORTED"` // Optional

	// Extra: these are commonly omitted
	// For reading
	ResourceName          string `json:"resource_name,omitempty"          env:"OAUTH_PR_RESOURCE_NAME"`          // Recommended
	ResourceDocumentation string `json:"resource_documentation,omitempty" env:"OAUTH_PR_RESOURCE_DOCUMENTATION"` // Optional
	ResourcePolicyUri     string `json:"resource_policy_uri,omitempty"    env:"OAUTH_PR_RESOURCE_POLICY_URI"`    // Optional
	ResourceTosUri        string `json:"resource_tos_uri,omitempty"       env:"OAUTH_PR_RESOURCE_TOS_URI"`       // Optional

	// For advanced security
	TlsClientCertificateBoundAccessTokens bool     `json:"tls_client_certificate_bound_access_tokens,omitempty"   env:"OAUTH_PR_TLS_CLIENT_CERTIFICATE_BOUND_ACCESS_TOKENS"` // Optional
	AuthorizationDetailsTypesSupported    []string `json:"authorization_details_types_supported,omitempty"        env:"OAUTH_PR_AUTHORIZATION_DETAILS_TYPES_SUPPORTED"`      // Optional
	DpopSigningAlgValuesSupported         []string `json:"dpop_signing_alg_values_supported,omitempty"            env:"OAUTH_PR_DPOP_SIGNING_ALG_VALUES_SUPPORTED"`          // Optional
	DpopBoundAccessTokensRequired         bool     `json:"dpop_bound_access_tokens_required,omitempty"            env:"OAUTH_PR_DPOP_BOUND_ACCESS_TOKENS_REQUIRED"`          // Optional
}

func HandleOauthProtectedResources(response http.ResponseWriter, request *http.Request) {

	//
	ResponseObject := &OauthProtectedResourceResponse{}
	if err := envconfig.Process(globals.Context, ResponseObject); err != nil {
		globals.Logger.Error("error processing environment vars", "error", err.Error())
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Transform into JSON
	ResponseObjectBytes, err := json.Marshal(ResponseObject)
	if err != nil {
		globals.Logger.Error("error converting response into json", "error", err.Error())
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "max-age=3600")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET")          // FIXME: TOO STRICT
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type") // FIXME: TOO STRICT

	_, err = response.Write(ResponseObjectBytes)
	if err != nil {
		globals.Logger.Error("error sending response to client", "error", err.Error())
		return
	}
}
