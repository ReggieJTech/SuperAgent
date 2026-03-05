package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

// Authenticator handles webhook authentication
type Authenticator struct {
	config AuthConfig
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator(config AuthConfig) *Authenticator {
	return &Authenticator{
		config: config,
	}
}

// Authenticate validates the incoming request
func (a *Authenticator) Authenticate(r *http.Request, body []byte) error {
	switch strings.ToLower(a.config.Type) {
	case "none":
		return nil
	case "bearer":
		return a.authenticateBearer(r)
	case "apikey":
		return a.authenticateAPIKey(r)
	case "basic":
		return a.authenticateBasic(r)
	case "hmac":
		return a.authenticateHMAC(r, body)
	default:
		return fmt.Errorf("unsupported auth type: %s", a.config.Type)
	}
}

// authenticateBearer validates bearer token authentication
func (a *Authenticator) authenticateBearer(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return fmt.Errorf("missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return fmt.Errorf("invalid Authorization header format")
	}

	if parts[1] != a.config.Token {
		return fmt.Errorf("invalid bearer token")
	}

	return nil
}

// authenticateAPIKey validates API key authentication
func (a *Authenticator) authenticateAPIKey(r *http.Request) error {
	header := a.config.Header
	if header == "" {
		header = "X-API-Key"
	}

	apiKey := r.Header.Get(header)
	if apiKey == "" {
		return fmt.Errorf("missing API key header: %s", header)
	}

	if apiKey != a.config.Key {
		return fmt.Errorf("invalid API key")
	}

	return nil
}

// authenticateBasic validates basic authentication
func (a *Authenticator) authenticateBasic(r *http.Request) error {
	username, password, ok := r.BasicAuth()
	if !ok {
		return fmt.Errorf("missing or invalid Basic auth")
	}

	if username != a.config.Username || password != a.config.Password {
		return fmt.Errorf("invalid credentials")
	}

	return nil
}

// authenticateHMAC validates HMAC signature
func (a *Authenticator) authenticateHMAC(r *http.Request, body []byte) error {
	header := a.config.Header
	if header == "" {
		header = "X-Signature"
	}

	signature := r.Header.Get(header)
	if signature == "" {
		return fmt.Errorf("missing HMAC signature header: %s", header)
	}

	// Calculate expected signature
	algorithm := strings.ToLower(a.config.Algorithm)
	var expectedSig string

	switch algorithm {
	case "sha256", "":
		mac := hmac.New(sha256.New, []byte(a.config.Secret))
		mac.Write(body)
		expectedSig = hex.EncodeToString(mac.Sum(nil))
	default:
		return fmt.Errorf("unsupported HMAC algorithm: %s", algorithm)
	}

	// Remove "sha256=" prefix if present
	signature = strings.TrimPrefix(signature, "sha256=")

	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return fmt.Errorf("invalid HMAC signature")
	}

	return nil
}
