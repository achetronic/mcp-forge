package handlers

import (
	"encoding/json"
	"net/http"

	//
	"tiny-mcp/internal/globals"
)

// OauthProtectedResourceResponse represents the response returned by '.well-known/oauth-protected-resource' endpoint
// According to the RFC9728 (Section 2)
// Ref: https://datatracker.ietf.org/doc/rfc9728/
type OauthProtectedResourceResponse struct {

	// Essential: these are commonly included
	Resource                          string   `json:"resource"`                                        // Required
	AuthorizationServers              []string `json:"authorization_servers,omitempty"`                 // Optional
	JwksUri                           string   `json:"jwks_uri,omitempty"`                              // Optional
	ScopesSupported                   []string `json:"scopes_supported,omitempty"`                      // Optional
	BearerMethodsSupported            []string `json:"bearer_methods_supported,omitempty"`              // Optional
	ResourceSigningAlgValuesSupported []string `json:"resource_signing_alg_values_supported,omitempty"` // Optional

	// Extra: these are commonly omitted
	// For reading
	ResourceName          string `json:"resource_name,omitempty"`          // Recommended
	ResourceDocumentation string `json:"resource_documentation,omitempty"` // Optional
	ResourcePolicyUri     string `json:"resource_policy_uri,omitempty"`    // Optional
	ResourceTosUri        string `json:"resource_tos_uri,omitempty"`       // Optional

	// For advanced security
	TlsClientCertificateBoundAccessTokens bool     `json:"tls_client_certificate_bound_access_tokens,omitempty"` // Optional
	AuthorizationDetailsTypesSupported    []string `json:"authorization_details_types_supported,omitempty"`      // Optional
	DpopSigningAlgValuesSupported         []string `json:"dpop_signing_alg_values_supported,omitempty"`          // Optional
	DpopBoundAccessTokensRequired         bool     `json:"dpop_bound_access_tokens_required,omitempty"`          // Optional
}

func HandleOauthProtectedResources(response http.ResponseWriter, request *http.Request) {

	//
	ResponseObject := &OauthProtectedResourceResponse{
		Resource:                              globals.Environment.Resource,
		AuthorizationServers:                  globals.Environment.AuthorizationServers,
		JwksUri:                               globals.Environment.JwksUri,
		ScopesSupported:                       globals.Environment.ScopesSupported,
		BearerMethodsSupported:                globals.Environment.BearerMethodsSupported,
		ResourceSigningAlgValuesSupported:     globals.Environment.ResourceSigningAlgValuesSupported,
		ResourceName:                          globals.Environment.ResourceName,
		ResourceDocumentation:                 globals.Environment.ResourceDocumentation,
		ResourcePolicyUri:                     globals.Environment.ResourcePolicyUri,
		ResourceTosUri:                        globals.Environment.ResourceTosUri,
		TlsClientCertificateBoundAccessTokens: globals.Environment.TlsClientCertificateBoundAccessTokens,
		AuthorizationDetailsTypesSupported:    globals.Environment.AuthorizationDetailsTypesSupported,
		DpopSigningAlgValuesSupported:         globals.Environment.DpopSigningAlgValuesSupported,
		DpopBoundAccessTokensRequired:         globals.Environment.DpopBoundAccessTokensRequired,
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
