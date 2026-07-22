package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Simon-Weij/yamp/internal/handler"
	"github.com/Simon-Weij/yamp/internal/router"
)

func main() {
	healthHandler := handler.NewHealthHandler()
	mux := router.New(*healthHandler)

	slog.Info("starting server...")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
	}
}
