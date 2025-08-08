package middlewares

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	//
	"github.com/golang-jwt/jwt/v5"
)

// JWKS represents a set (group) of several JWK
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represent a single JSON Web Key
type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`

	N string `json:"n,omitempty"` // RSA modulus
	E string `json:"e,omitempty"` // RSA exponent

	Crv string `json:"crv,omitempty"` // EC curve
	X   string `json:"x,omitempty"`   // EC x coordinate
	Y   string `json:"y,omitempty"`   // EC y coordinate

	K   string `json:"k,omitempty"` // Symmetric key (for HMAC)
	Alg string `json:"alg"`
	Use string `json:"use"`
}

// cacheJWKS obtains JWKS keys from remote, from time to time,
// and keep internal cache reasonable up-to-date
func (mw *JWTValidationMiddleware) cacheJWKS() {

	// Bypass the cache thread when middleware is disabled by config
	if !mw.dependencies.AppCtx.Config.Middleware.JWT.Enabled {
		return
	}

	mw.dependencies.AppCtx.Logger.Info("JWKS cache daemon running for JWT auth middleware")

	for {
		var jwks JWKS

		//
		resp, err := http.Get(mw.dependencies.AppCtx.Config.Middleware.JWT.Validation.Local.JWKSUri)
		if err != nil {
			mw.dependencies.AppCtx.Logger.Error("failed getting JWKS from remote", "error", err.Error())
			goto haveANap
		}
		defer resp.Body.Close()

		//
		if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
			mw.dependencies.AppCtx.Logger.Error("failed decoding JWKS from remote", "error", err.Error())
			goto haveANap
		}

		//
		mw.mutex.Lock()
		mw.jwks = &jwks
		mw.mutex.Unlock()

		// Don't be greedy, man
	haveANap:
		time.Sleep(mw.dependencies.AppCtx.Config.Middleware.JWT.Validation.Local.CacheInterval)
	}
}

func (mw *JWTValidationMiddleware) isTokenValid(token string) (bool, error) {
	// Get JWT header
	header, err := parseJWTHeader(token)
	if err != nil {
		return false, fmt.Errorf("error parsing token: %v", err)
	}

	// Retrieve 'Kid' and 'Alg' from token's header
	kid, ok := header["kid"].(string)
	if !ok {
		return false, fmt.Errorf("jwt header 'kid' field not found")
	}

	alg, ok := header["alg"].(string)
	if !ok {
		return false, fmt.Errorf("jwt header 'alg' field not found")
	}

	// Obtain JWKS from the middleware cache
	mw.mutex.Lock()
	jwks := mw.jwks
	mw.mutex.Unlock()

	// Look for the published key with the same Kid as the token
	var matchingKey *JWK
	for _, key := range jwks.Keys {
		if key.Kid == kid && (key.Use == "" || key.Use == "sig") {
			matchingKey = &key
			break
		}
	}

	if matchingKey == nil {
		return false, fmt.Errorf("no matching 'kid' in JWKS")
	}

	// Algorithm must match
	if matchingKey.Alg != "" && matchingKey.Alg != alg {
		return false, fmt.Errorf("algorithm missmatch")
	}

	// Convert JWK to a public key of corresponding type (RSA, EC, etc.)
	publicKey, err := jwkToKey(matchingKey)
	if err != nil {
		return false, fmt.Errorf("error converting JWK to public key")
	}

	// Validate the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		//
		expectedMethod, localErr := getSigningMethod(alg)
		if localErr != nil {
			return nil, localErr
		}

		if token.Method != expectedMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})

	if err != nil || !parsedToken.Valid {
		return false, fmt.Errorf("invalid token: %v", err)
	}

	return true, nil
}

// parseJWTHeader extracts the header of a JWT without verifying the signature
// This is used to infer algorithm to be used and the key from the JWKS
func parseJWTHeader(tokenString string) (map[string]interface{}, error) {
	//
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed token: It must be like header.payload.signature")
	}

	// Extract the header (first part)
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("error decoding header: %s", err.Error())
	}

	//
	var header map[string]interface{}
	if err = json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("error parsing JSON header: %s", err.Error())
	}

	return header, nil
}

// jwkToKey calculate corresponding real key (RSA, EC, etc.) from params present in the JWK
func jwkToKey(jwk *JWK) (interface{}, error) {
	switch jwk.Kty {
	case "RSA":
		return jwkToRSAPublicKey(jwk)
	case "EC":
		return jwkToECPublicKey(jwk)
	case "oct": // Symmetric keys
		return jwkToSymmetricKey(jwk)
	default:
		return nil, fmt.Errorf("unsupported key type: %s", jwk.Kty)
	}
}

// jwkToRSAPublicKey converts a JWK into a public RSA key
func jwkToRSAPublicKey(jwk *JWK) (*rsa.PublicKey, error) {
	if jwk.N == "" || jwk.E == "" {
		return nil, fmt.Errorf("incomplete RSA key data")
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("error decoding modulus: %v", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("error decoding exponent: %v", err)
	}

	n := new(big.Int)
	n.SetBytes(nBytes)

	var e int
	for i := 0; i < len(eBytes); i++ {
		e = e<<8 + int(eBytes[i])
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

// jwkToECPublicKey converts a JWK into a public ECDSA key
func jwkToECPublicKey(jwk *JWK) (*ecdsa.PublicKey, error) {
	if jwk.X == "" || jwk.Y == "" || jwk.Crv == "" {
		return nil, fmt.Errorf("incomplete EC key data")
	}

	var curve elliptic.Curve
	switch jwk.Crv {
	case "P-256":
		curve = elliptic.P256()
	case "P-384":
		curve = elliptic.P384()
	case "P-521":
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("unsupported curve: %s", jwk.Crv)
	}

	xBytes, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, fmt.Errorf("error decoding X coordinate: %v", err)
	}

	yBytes, err := base64.RawURLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return nil, fmt.Errorf("error decoding Y coordinate: %v", err)
	}

	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}

// jwkToSymmetricKey converts a JWK into a simetric key (for HMAC)
func jwkToSymmetricKey(jwk *JWK) ([]byte, error) {
	if jwk.K == "" {
		return nil, fmt.Errorf("incomplete symmetric key data")
	}

	k, err := base64.RawURLEncoding.DecodeString(jwk.K)
	if err != nil {
		return nil, fmt.Errorf("error decoding symmetric key: %v", err)
	}

	return k, nil
}

// getSigningMethod returns suitable signing method according to the algorithm
func getSigningMethod(alg string) (jwt.SigningMethod, error) {
	switch alg {
	case "RS256":
		return jwt.SigningMethodRS256, nil
	case "RS384":
		return jwt.SigningMethodRS384, nil
	case "RS512":
		return jwt.SigningMethodRS512, nil
	case "ES256":
		return jwt.SigningMethodES256, nil
	case "ES384":
		return jwt.SigningMethodES384, nil
	case "ES512":
		return jwt.SigningMethodES512, nil
	case "HS256":
		return jwt.SigningMethodHS256, nil
	case "HS384":
		return jwt.SigningMethodHS384, nil
	case "HS512":
		return jwt.SigningMethodHS512, nil
	default:
		return nil, fmt.Errorf("unsupported signing method: %s", alg)
	}
}
