package globals

import "github.com/sethvargo/go-envconfig"

const (
	ServerTransportHttp                           = "http"
	ServerTransportHttpJwtValidationStrategyLocal = "local"
)

type EnvironmentOptions struct {
	// 1. Server related
	ServerTransport                          string `env:"SERVER_TRANSPORT"`
	ServerTransportHttpHost                  string `env:"SERVER_TRANSPORT_HTTP_HOST"`
	ServerTransportHttpJwtValidatedHeader    string `env:"SERVER_TRANSPORT_HTTP_JWT_VALIDATED_HEADER"`
	ServerTransportHttpJwtValidationStrategy string `env:"SERVER_TRANSPORT_HTTP_JWT_VALIDATION_STRATEGY"`

	// 2. Oauth Protected Resource related
	// Essential: these are commonly included
	Resource                          string   `env:"OAUTH_PR_RESOURCE,required"`                     // Required
	AuthorizationServers              []string `env:"OAUTH_PR_AUTH_SERVERS"`                          // Optional
	JwksUri                           string   `env:"OAUTH_PR_JWKS_URI"`                              // Optional
	ScopesSupported                   []string `env:"OAUTH_PR_SCOPES_SUPPORTED"`                      // Optional
	BearerMethodsSupported            []string `env:"OAUTH_PR_BEARER_METHODS_SUPPORTED"`              // Optional
	ResourceSigningAlgValuesSupported []string `env:"OAUTH_PR_RESOURCE_SIGNING_ALG_VALUES_SUPPORTED"` // Optional

	// Extra: these are commonly omitted
	// For reading
	ResourceName          string `env:"OAUTH_PR_RESOURCE_NAME"`          // Recommended
	ResourceDocumentation string `env:"OAUTH_PR_RESOURCE_DOCUMENTATION"` // Optional
	ResourcePolicyUri     string `env:"OAUTH_PR_RESOURCE_POLICY_URI"`    // Optional
	ResourceTosUri        string `env:"OAUTH_PR_RESOURCE_TOS_URI"`       // Optional

	// For advanced security
	TlsClientCertificateBoundAccessTokens bool     `env:"OAUTH_PR_TLS_CLIENT_CERTIFICATE_BOUND_ACCESS_TOKENS"` // Optional
	AuthorizationDetailsTypesSupported    []string `env:"OAUTH_PR_AUTHORIZATION_DETAILS_TYPES_SUPPORTED"`      // Optional
	DpopSigningAlgValuesSupported         []string `env:"OAUTH_PR_DPOP_SIGNING_ALG_VALUES_SUPPORTED"`          // Optional
	DpopBoundAccessTokensRequired         bool     `env:"OAUTH_PR_DPOP_BOUND_ACCESS_TOKENS_REQUIRED"`          // Optional
}

func GetEnvironmentOptions() (*EnvironmentOptions, error) {
	EnvironmentObject := &EnvironmentOptions{}
	if err := envconfig.Process(Context, EnvironmentObject); err != nil {
		return nil, err
	}

	return EnvironmentObject, nil
}
