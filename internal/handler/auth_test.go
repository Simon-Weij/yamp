package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Simon-Weij/yamp/internal/db/sqlc"
	"github.com/stretchr/testify/assert"
)

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
		PasswordHash: "AAAAA",
	}, nil
}

func Test_HandleSignup(t *testing.T) {
	tests := []struct {
		name               string
		json               string
		expectedStatusCode int
	}{
		{
			name:               "runs correctly",
			json:               `{"username":"testuser","password":"password"}`,
			expectedStatusCode: http.StatusOK,
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
			handler := NewAuthHandler(mockedUserCreator{})

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
