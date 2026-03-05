package webui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// AuthConfig contains authentication configuration
type AuthConfig struct {
	Enabled       bool
	JWTSecret     string
	SessionExpiry time.Duration
}

// Claims represents JWT claims
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

// UserInfo represents user information
type UserInfo struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

// handleLogin handles user login
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// Validate credentials (simplified - should use proper user store)
	if !s.validateCredentials(req.Username, req.Password) {
		s.writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Generate JWT token
	token, expiresAt, err := s.generateToken(req.Username, "admin")
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		s.writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	response := LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: UserInfo{
			Username: req.Username,
			Role:     "admin",
		},
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleLogout handles user logout
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	// In a stateless JWT system, logout is handled client-side by discarding the token
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "logged out successfully",
	})
}

// handleRefreshToken handles token refresh
func (s *Server) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	// Extract token from header
	tokenString := s.extractToken(r)
	if tokenString == "" {
		s.writeError(w, http.StatusUnauthorized, "no token provided")
		return
	}

	// Parse and validate token
	claims, err := s.validateToken(tokenString)
	if err != nil {
		s.writeError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	// Generate new token
	token, expiresAt, err := s.generateToken(claims.Username, claims.Role)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	response := LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: UserInfo{
			Username: claims.Username,
			Role:     claims.Role,
		},
	}

	s.writeJSON(w, http.StatusOK, response)
}

// generateToken generates a JWT token
func (s *Server) generateToken(username, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "bigpanda-super-agent",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret-key-change-me")) // Should use config
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// validateToken validates a JWT token
func (s *Server) validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret-key-change-me"), nil // Should use config
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// extractToken extracts JWT token from request
func (s *Server) extractToken(r *http.Request) string {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Check query parameter
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	return ""
}

// validateCredentials validates user credentials
func (s *Server) validateCredentials(username, password string) bool {
	// Simplified validation - should use proper user store and password hashing
	// For demo purposes, accept any non-empty credentials
	return username != "" && password != ""
}

// requireAuth is middleware that requires authentication
func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := s.extractToken(r)
		if tokenString == "" {
			s.writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		claims, err := s.validateToken(tokenString)
		if err != nil {
			s.writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// Add claims to request context (optional)
		log.Debug().Str("user", claims.Username).Msg("Authenticated request")

		next.ServeHTTP(w, r)
	})
}
