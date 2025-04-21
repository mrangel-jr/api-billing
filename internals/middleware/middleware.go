package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct{}
type contextKey string

const ClaimsContextKey = contextKey("claims")

func GetUser(r *http.Request) (string, error) {
	claims, ok := r.Context().Value(ClaimsContextKey).(*CustomClaims)
	if !ok {
		return "", errors.New("no claims found in context")
	}

	return claims.Username, nil
}

// JWTMiddleware returns an HTTP middleware that checks the JWT token
func (am *AuthMiddleware) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Get signing key from environment variable
		signingKey := os.Getenv("JWT_SIGNING_KEY")
		if signingKey == "" {
			http.Error(w, "Missing signing key", http.StatusInternalServerError)
			return
		}
		// Extract token from Authorization header

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(signingKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		if claims, ok := token.Claims.(*CustomClaims); ok {
			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
	})
}
