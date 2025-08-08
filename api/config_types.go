package api

import "time"

// ServerTransportHTTPConfig represents the HTTP transport configuration
type ServerTransportHTTPConfig struct {
	Host string `yaml:"host"`
}

// ServerTransportConfig represents the transport configuration
type ServerTransportConfig struct {
	Type string                    `yaml:"type"`
	HTTP ServerTransportHTTPConfig `yaml:"http,omitempty"`
}

// ServerConfig represents the server configuration section
type ServerConfig struct {
	Name      string                `yaml:"name"`
	Version   string                `yaml:"version"`
	Transport ServerTransportConfig `yaml:"transport,omitempty"`
}

// AccessLogsConfig represents the AccessLogs middleware configuration
type AccessLogsConfig struct {
	ExcludedHeaders []string `yaml:"excluded_headers"`
	RedactedHeaders []string `yaml:"redacted_headers"`
}

// JWTValidationLocalConfig represents the local JWT validation configuration
type JWTValidationLocalConfig struct {
	JWKSUri         string                        `yaml:"jwks_uri"`
	CacheInterval   time.Duration                 `yaml:"cache_interval"`
	AllowConditions []JWTValidationAllowCondition `yaml:"allow_conditions,omitempty"`
}

// JWTValidationAllowCondition represents a condition for allowing a request after the local JWT validation configuration
type JWTValidationAllowCondition struct {
	Expression string `yaml:"expression"`
}

// JWTValidationConfig represents the JWT validation configuration
type JWTValidationConfig struct {
	Strategy        string                   `yaml:"strategy"`
	ForwardedHeader string                   `yaml:"forwarded_header,omitempty"`
	Local           JWTValidationLocalConfig `yaml:"local,omitempty"`
}

// JWTConfig represents the JWT middleware configuration
type JWTConfig struct {
	Enabled    bool                `yaml:"enabled"`
	Validation JWTValidationConfig `yaml:"validation,omitempty"`
}

// MiddlewareConfig represents the middleware configuration section
type MiddlewareConfig struct {
	AccessLogs AccessLogsConfig `yaml:"access_logs"`
	JWT        JWTConfig        `yaml:"jwt,omitempty"`
}

// OAuthAuthorizationServer represents the OAuth Authorization Server configuration
type OAuthAuthorizationServer struct {
	Enabled   bool   `yaml:"enabled"`
	IssuerUri string `yaml:"issuer_uri"`
}

// OAuthProtectedResourceConfig represents the OAuth Protected Resource configuration
type OAuthProtectedResourceConfig struct {
	Enabled                               bool     `yaml:"enabled"`
	Resource                              string   `yaml:"resource"`
	AuthServers                           []string `yaml:"auth_servers"`
	JWKSUri                               string   `yaml:"jwks_uri"`
	ScopesSupported                       []string `yaml:"scopes_supported"`
	BearerMethodsSupported                []string `yaml:"bearer_methods_supported,omitempty"`
	ResourceSigningAlgValuesSupported     []string `yaml:"resource_signing_alg_values_supported,omitempty"`
	ResourceName                          string   `yaml:"resource_name,omitempty"`
	ResourceDocumentation                 string   `yaml:"resource_documentation,omitempty"`
	ResourcePolicyUri                     string   `yaml:"resource_policy_uri,omitempty"`
	ResourceTosUri                        string   `yaml:"resource_tos_uri,omitempty"`
	TLSClientCertificateBoundAccessTokens bool     `yaml:"tls_client_certificate_bound_access_tokens,omitempty"`
	AuthorizationDetailsTypesSupported    []string `yaml:"authorization_details_types_supported,omitempty"`
	DPoPSigningAlgValuesSupported         []string `yaml:"dpop_signing_alg_values_supported,omitempty"`
	DPoPBoundAccessTokensRequired         bool     `yaml:"dpop_bound_access_tokens_required,omitempty"`
}

// Configuration represents the complete configuration structure
type Configuration struct {
	Server                   ServerConfig                 `yaml:"server,omitempty"`
	Middleware               MiddlewareConfig             `yaml:"middleware,omitempty"`
	OAuthAuthorizationServer OAuthAuthorizationServer     `yaml:"oauth_authorization_server,omitempty"`
	OAuthProtectedResource   OAuthProtectedResourceConfig `yaml:"oauth_protected_resource,omitempty"`
}
