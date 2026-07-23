package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WithAuth(t *testing.T) {
	secret := "test-secret"
	t.Setenv("JWT_TOKEN", secret)

	makeToken := func(claims jwt.MapClaims, signingKey []byte, method jwt.SigningMethod) string {
		token := jwt.NewWithClaims(method, claims)
		signed, err := token.SignedString(signingKey)
		require.NoError(t, err)
		return signed
	}

	validToken := makeToken(jwt.MapClaims{
		"user_id": "user-123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}, []byte(secret), jwt.SigningMethodHS256)

	expiredToken := makeToken(jwt.MapClaims{
		"user_id": "user-123",
		"exp":     time.Now().Add(-time.Hour).Unix(),
	}, []byte(secret), jwt.SigningMethodHS256)

	noUserIDToken := makeToken(jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	}, []byte(secret), jwt.SigningMethodHS256)

	wrongSecretToken := makeToken(jwt.MapClaims{
		"user_id": "user-123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}, []byte("wrong-secret"), jwt.SigningMethodHS256)

	tests := []struct {
		name               string
		authHeader         string
		expectedStatusCode int
	}{
		{
			name:               "missing header",
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "malformed header",
			authHeader:         "Bearer",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "garbage token",
			authHeader:         "Bearer aaaa",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "expired token",
			authHeader:         "Bearer " + expiredToken,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "wrong secret",
			authHeader:         "Bearer " + wrongSecretToken,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "missing user_id claim",
			authHeader:         "Bearer " + noUserIDToken,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "valid token",
			authHeader:         "Bearer " + validToken,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("ok"))
				require.NoError(t, err)
			})
			handler := WithAuth(nextHandler)
			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedStatusCode, rec.Code)
		})
	}
}
