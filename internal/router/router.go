package router

import (
	"net/http"

	"github.com/Simon-Weij/yamp/internal/handler"
	"github.com/Simon-Weij/yamp/internal/middleware"
)

func New(
	healtHandler handler.HealthHandler,
	authHandler handler.AuthHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle(
		"GET /health",
		middleware.New(middleware.WithLogging).Then(http.HandlerFunc(healtHandler.GetHealth)),
	)

	mux.Handle(
		"POST /auth/signup",
		middleware.New(
			middleware.WithLogging,
			middleware.WithTimeout,
		).Then(http.HandlerFunc(authHandler.HandleSignup)),
	)
	mux.Handle(
		"POST /auth/login",
		middleware.New(
			middleware.WithLogging,
			middleware.WithTimeout,
		).Then(http.HandlerFunc(authHandler.HandleLogin)),
	)
	mux.Handle(
		"POST /auth/refresh",
		middleware.New(
			middleware.WithLogging,
			middleware.WithTimeout,
		).Then(http.HandlerFunc(authHandler.HandleRefresh)),
	)
	mux.Handle(
		"POST /auth/logout",
		middleware.New(
			middleware.WithLogging,
			middleware.WithTimeout,
		).Then(http.HandlerFunc(authHandler.HandleLogout)),
	)

	return mux
}
