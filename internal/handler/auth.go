package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Simon-Weij/yamp/internal/db/sqlc"
	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	userCreator UserCreator
}

func NewAuthHandler(userCreator UserCreator) *AuthHandler {
	return &AuthHandler{
		userCreator: userCreator,
	}
}

type UserCreator interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
}

type signupRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()

	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, "validation failed", http.StatusBadRequest)
		return
	}

	passwordHash, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		http.Error(w, "something went wrong while hashing", http.StatusInternalServerError)
	}

	_, err = h.userCreator.CreateUser(r.Context(), sqlc.CreateUserParams{
		Username:     req.Username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		http.Error(w, "something went wrong while creating user", http.StatusInternalServerError)
	}
}
