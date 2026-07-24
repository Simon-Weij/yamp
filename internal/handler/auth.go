package handler

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Simon-Weij/yamp/internal/db/sqlc"
	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	accessTokenDuration  = time.Minute * 7
	refreshTokenDuration = time.Hour * 24 * 7

	refreshTokenCookieName = "refresh_token"
)

var (
	errRefreshTokenMissing = errors.New("missing refresh token cookie")
	errRefreshTokenRevoked = errors.New("refresh token has been revoked")
	errRefreshTokenExpired = errors.New("refresh token has expired")
)

type AuthHandler struct {
	UserRepository         UserRepository
	RefreshTokenRepository RefreshTokenRepository
}

func NewAuthHandler(
	userRepository UserRepository,
	refreshTokenRepository RefreshTokenRepository,
) *AuthHandler {
	return &AuthHandler{
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

type UserRepository interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	GetUserForLogin(ctx context.Context, username string) (sqlc.GetUserForLoginRow, error)
}

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, arg sqlc.CreateRefreshTokenParams) (sqlc.RefreshToken, error)
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (sqlc.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id int64) error
}

type authRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func validateAuthRequest(w http.ResponseWriter, r *http.Request) (*authRequest, error) {
	validate := validator.New()

	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid body: %w", err)
	}
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &req, nil
}

func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	req, err := validateAuthRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passwordHash, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		http.Error(w, "something went wrong while hashing", http.StatusInternalServerError)
		return
	}

	_, err = h.UserRepository.CreateUser(r.Context(), sqlc.CreateUserParams{
		Username:     req.Username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("user already exists %s", req.Username), http.StatusConflict)
		return
	}
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	req, err := validateAuthRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.UserRepository.GetUserForLogin(r.Context(), req.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("something went wrong while getting user for %s: %v", req.Username, err), http.StatusUnauthorized)
		return
	}

	passwordMatches, err := argon2id.ComparePasswordAndHash(req.Password, user.PasswordHash)
	if err != nil {
		http.Error(w, "could not compare passwords: "+err.Error(), http.StatusInternalServerError)
	}
	if !passwordMatches {
		http.Error(w, "password does not match!", http.StatusUnauthorized)
		return
	}

	accessToken, err := createAccessToken(user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("something went wrong signing: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.issueRefreshToken(w, r, user.ID); err != nil {
		http.Error(w, fmt.Sprintf("something went wrong issuing refresh token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		http.Error(w, errRefreshTokenMissing.Error(), http.StatusUnauthorized)
		return
	}

	stored, err := h.RefreshTokenRepository.GetRefreshTokenByHash(
		r.Context(),
		hashRefreshToken(cookie.Value),
	)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	if stored.RevokedAt.Valid {
		http.Error(w, errRefreshTokenRevoked.Error(), http.StatusUnauthorized)
		return
	}

	if !stored.ExpiresAt.Valid || stored.ExpiresAt.Time.Before(time.Now()) {
		http.Error(w, errRefreshTokenExpired.Error(), http.StatusUnauthorized)
		return
	}

	if err := h.RefreshTokenRepository.RevokeRefreshToken(r.Context(), stored.ID); err != nil {
		http.Error(w, "something went wrong while rotating refresh token", http.StatusInternalServerError)
		return
	}

	if err := h.issueRefreshToken(w, r, stored.UserID); err != nil {
		http.Error(w, fmt.Sprintf("something went wrong issuing refresh token: %v", err), http.StatusInternalServerError)
		return
	}

	accessToken, err := createAccessToken(stored.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("something went wrong signing: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		setRefreshTokenCookie(w, "", true)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	stored, err := h.RefreshTokenRepository.GetRefreshTokenByHash(
		r.Context(),
		hashRefreshToken(cookie.Value),
	)
	if err == nil {
		_ = h.RefreshTokenRepository.RevokeRefreshToken(r.Context(), stored.ID)
	}

	setRefreshTokenCookie(w, "", true)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) issueRefreshToken(w http.ResponseWriter, r *http.Request, userID int64) error {
	raw, err := generateRefreshToken()
	if err != nil {
		return fmt.Errorf("generate refresh token: %w", err)
	}

	expiresAt := time.Now().Add(refreshTokenDuration)

	_, err = h.RefreshTokenRepository.CreateRefreshToken(r.Context(), sqlc.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: hashRefreshToken(raw),
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("store refresh token: %w", err)
	}

	setRefreshTokenCookie(w, raw, false)

	return nil
}

func setRefreshTokenCookie(w http.ResponseWriter, value string, expired bool) {
	cookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	if expired {
		cookie.Value = ""
		cookie.Expires = time.Unix(0, 0)
		cookie.MaxAge = -1
	}
	http.SetCookie(w, cookie)
}

func createAccessToken(userID int64) (string, error) {
	key := os.Getenv("JWT_TOKEN")
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(accessTokenDuration).Unix(),
			"iat":     time.Now().Unix(),
		})
	token, err := t.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return token, nil
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
