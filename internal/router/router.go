package router

import (
	"net/http"

	"github.com/Simon-Weij/yamp/internal/handler"
	"github.com/Simon-Weij/yamp/internal/middleware"
)

func New(healtHandler handler.HealthHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle(
		"GET /health",
		middleware.New(middleware.WithLogging).Then(http.HandlerFunc(healtHandler.GetHealth)),
	)

	return mux
}
