package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Simon-Weij/yamp/internal/db/sqlc"
	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	UserRepository UserRepository
}

func NewAuthHandler(userRepository UserRepository) *AuthHandler {
	return &AuthHandler{
		UserRepository: userRepository,
	}
}

type UserRepository interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	GetUserForLogin(ctx context.Context, username string) (sqlc.GetUserForLoginRow, error)
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
		http.Error(w, "something went wrong while creating user", http.StatusInternalServerError)
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
	if !passwordMatches {
		http.Error(w, "password does not match!", http.StatusUnauthorized)
		return
	}

	key := os.Getenv("JWT_TOKEN")
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(time.Minute * 7).Unix(),
			"iat":     time.Now().Unix(),
		})
	token, err := t.SignedString([]byte(key))
	if err != nil {
		http.Error(w, fmt.Sprintf("something went wrong signing: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	})
}
