package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrangel-jr/api-billing/internals/middleware"
	"github.com/stretchr/testify/assert"
)

// Helper function to generate a valid JWT token
func generateToken(expired bool) string {
	exp := time.Now().Add(time.Hour).Unix()
	if expired {
		exp = time.Now().Add(-time.Hour).Unix()
	}
	secret := os.Getenv("JWT_SIGNING_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      "test-user",
		"exp":      exp,
		"username": "tenant_a",
		"iss":      "my-app",
		"iat":      1745377200,
		"aud":      "my-audience",
		"jti":      "unique-id-123",
	})

	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func TestJWTMiddleware_TableDriven(t *testing.T) {
	os.Setenv("JWT_SIGNING_KEY", "Du005rSSzktWnX5n6KEzGYAI4XYgcUPF8hknf3Urn+0=")
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authorized"))
	})

	handlerToTest := (&middleware.AuthMiddleware{}).JWTMiddleware(finalHandler)

	tests := []struct {
		name         string
		token        string
		expectedCode int
		expectedBody string
		setHeader    bool
	}{
		{
			name:         "Valid Token",
			token:        generateToken(false),
			expectedCode: http.StatusOK,
			expectedBody: "Authorized",
			setHeader:    true,
		},
		{
			name:         "Missing Authorization Header",
			token:        "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Missing Authorization header",
			setHeader:    false,
		},
		{
			name:         "Invalid Token Format",
			token:        "this.is.not.valid",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Invalid token",
			setHeader:    true,
		},
		{
			name:         "Expired Token",
			token:        generateToken(true),
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Invalid token",
			setHeader:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/consume/202504/summary", nil)
			if tc.setHeader {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}

			rr := httptest.NewRecorder()
			handlerToTest.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tc.expectedBody)
		})
	}
}
