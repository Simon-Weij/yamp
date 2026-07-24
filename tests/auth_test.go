//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/Simon-Weij/yamp/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func setupApp(t *testing.T, ctx context.Context) string {
	t.Helper()

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:18.4-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = pgContainer.Terminate(context.Background())
	})

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	port, err := pgContainer.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	t.Setenv("POSTGRES_HOST", host)
	t.Setenv("POSTGRES_PORT", port.Port())
	t.Setenv("POSTGRES_DB", "test")
	t.Setenv("POSTGRES_USER", "user")
	t.Setenv("POSTGRES_PASSWORD", "password")
	t.Setenv("POSTGRES_SSLMODE", "disable")

	ln, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	baseURL := fmt.Sprintf("http://localhost:%d", ln.Addr().(*net.TCPAddr).Port)

	go func() { require.NoError(t, app.Run(ctx, ln)) }()

	require.Eventually(t, func() bool {
		resp, err := http.Get(baseURL + "/health")
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 10*time.Second, 200*time.Millisecond)

	return baseURL
}

func TestAuthFullFlow(t *testing.T) {
	ctx := t.Context()
	baseURL := setupApp(t, ctx)

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := &http.Client{Jar: jar}

	const username = "integrationuser"
	const password = "securepass"

	t.Run("signup missing password", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"username": "u1"})
		resp, err := client.Post(baseURL+"/auth/signup", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("signup invalid username", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"username": "u4", "password": "password"})
		resp, err := client.Post(baseURL+"/auth/signup", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("signup valid", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"username": username, "password": password})
		resp, err := client.Post(baseURL+"/auth/signup", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("signup duplicate", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"username": username, "password": password})
		resp, err := client.Post(baseURL+"/auth/signup", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	var accessToken string
	t.Run("login", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"username": username, "password": password})
		resp, err := client.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var lr loginResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&lr))
		require.NotEmpty(t, lr.AccessToken)
		accessToken = lr.AccessToken

		require.Len(t, resp.Cookies(), 1)
		cookie := resp.Cookies()[0]
		assert.Equal(t, "refresh_token", cookie.Name)
		assert.True(t, cookie.HttpOnly)
		assert.True(t, cookie.Secure)
		assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	})

	t.Run("me with valid token", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, username, body["username"])
		assert.NotZero(t, body["id"])
	})

	t.Run("me without token", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/auth/me")
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("me with invalid token", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/auth/me", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	var newAccessToken string
	t.Run("refresh", func(t *testing.T) {
		resp, err := client.Post(baseURL+"/auth/refresh", "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var lr loginResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&lr))
		require.NotEmpty(t, lr.AccessToken)
		newAccessToken = lr.AccessToken
	})

	t.Run("me with refreshed token", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+newAccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, username, body["username"])
	})

	t.Run("refresh with old token is rejected", func(t *testing.T) {
		noJarClient := &http.Client{}
		resp, err := noJarClient.Post(baseURL+"/auth/refresh", "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("logout", func(t *testing.T) {
		resp, err := client.Post(baseURL+"/auth/logout", "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("refresh after logout is rejected", func(t *testing.T) {
		resp, err := client.Post(baseURL+"/auth/refresh", "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
