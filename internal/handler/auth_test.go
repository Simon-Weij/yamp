package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Simon-Weij/yamp/internal/db/sqlc"
	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPassword = "password"

var testPasswordHash = mustHash(testPassword)

func mustHash(password string) string {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		panic(err)
	}
	return hash
}

type mockedUserCreator struct{}

func (f mockedUserCreator) CreateUser(
	ctx context.Context,
	arg sqlc.CreateUserParams,
) (sqlc.User, error) {
	return sqlc.User{
		Username:     arg.Username,
		PasswordHash: "aaaaa",
	}, nil
}

func (f mockedUserCreator) GetUserForLogin(
	ctx context.Context,
	username string,
) (sqlc.GetUserForLoginRow, error) {
	return sqlc.GetUserForLoginRow{
		ID:           5,
		PasswordHash: testPasswordHash,
	}, nil
}

func (f mockedUserCreator) GetUserByID(
	ctx context.Context,
	id int64,
) (sqlc.GetUserByIDRow, error) {
	return sqlc.GetUserByIDRow{
		ID:       id,
		Username: "testuser",
	}, nil
}

type tokenStore struct {
	byHash  map[string]sqlc.RefreshToken
	revoked map[int64]bool
}

func newTokenStore() *tokenStore {
	return &tokenStore{
		byHash:  make(map[string]sqlc.RefreshToken),
		revoked: make(map[int64]bool),
	}
}

func (s *tokenStore) CreateRefreshToken(
	ctx context.Context,
	arg sqlc.CreateRefreshTokenParams,
) (sqlc.RefreshToken, error) {
	t := sqlc.RefreshToken{
		ID:        int64(len(s.byHash) + 1),
		UserID:    arg.UserID,
		TokenHash: arg.TokenHash,
		ExpiresAt: arg.ExpiresAt,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}
	s.byHash[arg.TokenHash] = t
	return t, nil
}

func (s *tokenStore) GetRefreshTokenByHash(
	ctx context.Context,
	tokenHash string,
) (sqlc.RefreshToken, error) {
	t, ok := s.byHash[tokenHash]
	if !ok {
		return sqlc.RefreshToken{}, errNotFound
	}
	if s.revoked[t.ID] {
		t.RevokedAt = pgtype.Timestamp{Time: time.Now(), Valid: true}
	}
	return t, nil
}

func (s *tokenStore) RevokeRefreshToken(ctx context.Context, id int64) error {
	s.revoked[id] = true
	return nil
}

var errNotFound = &notFoundError{}

type notFoundError struct{}

func (*notFoundError) Error() string { return "not found" }

func Test_HandleSignup(t *testing.T) {
	tests := []struct {
		name               string
		json               string
		expectedStatusCode int
	}{
		{
			name:               "runs correctly",
			json:               `{"username":"testuser","password":"password"}`,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "empty username",
			json:               `{"username":"","password":"password"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "empty password",
			json:               `{"username":"testuser","password":""}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "empty username and password",
			json:               `{"username":"","password":""}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "emtpy json",
			json:               "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "extra field",
			json:               `{"username":"","password":"","newfield":"newvalue"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "name too short",
			json:               `{"username":"a","password":"password"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "password too short",
			json:               `{"username":"username","password":"pass"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "username too long",
			json:               `{"username":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","password":"password"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "password too long",
			json:               `{"username":"username","password":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewAuthHandler(mockedUserCreator{}, newTokenStore())

			req := httptest.NewRequest(
				http.MethodPost,
				"/auth/signup",
				strings.NewReader(tt.json),
			)

			rec := httptest.NewRecorder()

			handler.HandleSignup(rec, req)

			t.Logf("status code: %d", rec.Code)
			t.Logf("response body: %s", rec.Body.String())

			assert.Equal(t, tt.expectedStatusCode, rec.Code)
		})
	}
}

func Test_HandleLogin_setsRefreshTokenCookie(t *testing.T) {
	t.Setenv("JWT_TOKEN", "test-secret")

	handler := NewAuthHandler(mockedUserCreator{}, newTokenStore())

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		strings.NewReader(`{"username":"testuser","password":"password"}`),
	)
	rec := httptest.NewRecorder()

	handler.HandleLogin(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]string
	require.NoError(t, rec.Result().Body.Close())
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Contains(t, body, "access_token")
	assert.NotContains(t, body, "refresh_token")

	var foundCookie bool
	for _, c := range rec.Result().Cookies() {
		if c.Name == refreshTokenCookieName {
			foundCookie = true
			assert.True(t, c.HttpOnly)
			assert.True(t, c.Secure, "cookie should be Secure")
			assert.Equal(t, http.SameSiteStrictMode, c.SameSite)
			assert.NotEmpty(t, c.Value)
		}
	}
	assert.True(t, foundCookie, "refresh_token cookie should be set")
}

func Test_HandleRefresh_rotatesAndIssuesAccessToken(t *testing.T) {
	t.Setenv("JWT_TOKEN", "test-secret")

	store := newTokenStore()
	handler := NewAuthHandler(mockedUserCreator{}, store)

	raw := "valid-refresh-token"
	hash := hashRefreshToken(raw)
	_, err := store.CreateRefreshToken(context.Background(), sqlc.CreateRefreshTokenParams{
		UserID:    5,
		TokenHash: hash,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(time.Hour), Valid: true},
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: raw})
	rec := httptest.NewRecorder()

	handler.HandleRefresh(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]string
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Contains(t, body, "access_token")

	stored, err := store.GetRefreshTokenByHash(context.Background(), hash)
	require.NoError(t, err)
	assert.True(t, store.revoked[stored.ID], "old token should be revoked")

	var newCookie bool
	for _, c := range rec.Result().Cookies() {
		if c.Name == refreshTokenCookieName && c.Value != "" {
			newCookie = true
			assert.NotEqual(t, raw, c.Value)
		}
	}
	assert.True(t, newCookie, "a new refresh_token cookie should be set")
}

func Test_HandleRefresh_missingCookie(t *testing.T) {
	handler := NewAuthHandler(mockedUserCreator{}, newTokenStore())

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	rec := httptest.NewRecorder()

	handler.HandleRefresh(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func Test_HandleRefresh_revokedToken(t *testing.T) {
	t.Setenv("JWT_TOKEN", "test-secret")

	store := newTokenStore()
	handler := NewAuthHandler(mockedUserCreator{}, store)

	raw := "revoked-refresh-token"
	hash := hashRefreshToken(raw)
	tok, err := store.CreateRefreshToken(context.Background(), sqlc.CreateRefreshTokenParams{
		UserID:    5,
		TokenHash: hash,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(time.Hour), Valid: true},
	})
	require.NoError(t, err)
	require.NoError(t, store.RevokeRefreshToken(context.Background(), tok.ID))

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: raw})
	rec := httptest.NewRecorder()

	handler.HandleRefresh(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func Test_HandleRefresh_expiredToken(t *testing.T) {
	t.Setenv("JWT_TOKEN", "test-secret")

	store := newTokenStore()
	handler := NewAuthHandler(mockedUserCreator{}, store)

	raw := "expired-refresh-token"
	hash := hashRefreshToken(raw)
	_, err := store.CreateRefreshToken(context.Background(), sqlc.CreateRefreshTokenParams{
		UserID:    5,
		TokenHash: hash,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(-time.Hour), Valid: true},
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: raw})
	rec := httptest.NewRecorder()

	handler.HandleRefresh(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func Test_HandleLogout_revokesAndClearsCookie(t *testing.T) {
	store := newTokenStore()
	handler := NewAuthHandler(mockedUserCreator{}, store)

	raw := "logout-refresh-token"
	hash := hashRefreshToken(raw)
	tok, err := store.CreateRefreshToken(context.Background(), sqlc.CreateRefreshTokenParams{
		UserID:    5,
		TokenHash: hash,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(time.Hour), Valid: true},
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: raw})
	rec := httptest.NewRecorder()

	handler.HandleLogout(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.True(t, store.revoked[tok.ID], "token should be revoked")

	var cleared bool
	for _, c := range rec.Result().Cookies() {
		if c.Name == refreshTokenCookieName {
			cleared = true
			assert.Equal(t, "", c.Value)
			assert.LessOrEqual(t, c.MaxAge, 0)
		}
	}
	assert.True(t, cleared, "refresh_token cookie should be cleared")
}

func Test_HandleLogout_noCookie(t *testing.T) {
	handler := NewAuthHandler(mockedUserCreator{}, newTokenStore())

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rec := httptest.NewRecorder()

	handler.HandleLogout(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func Test_hashRefreshToken_stable(t *testing.T) {
	sum := sha256.Sum256([]byte("abc"))
	expected := hex.EncodeToString(sum[:])
	assert.Equal(t, expected, hashRefreshToken("abc"))
}
