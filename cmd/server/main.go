package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/Simon-Weij/yamp/internal/db"
	"github.com/Simon-Weij/yamp/internal/db/sqlc"
	"github.com/Simon-Weij/yamp/internal/handler"
	"github.com/Simon-Weij/yamp/internal/router"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func buildDSN() string {
	host := requireEnv("POSTGRES_HOST")
	port := requireEnv("POSTGRES_PORT")
	user := requireEnv("POSTGRES_USER")
	pass := requireEnv("POSTGRES_PASSWORD")
	dbname := requireEnv("POSTGRES_DB")
	sslmode := requireEnv("POSTGRES_SSLMODE")

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   "/" + dbname,
	}
	q := u.Query()
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()

	return u.String()
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return v
}

func main() {
	dsn := buildDSN()

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer pool.Close()

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open sql db for migrations: %v", err)
	}

	if err := db.Migrate(sqlDB); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatalf("something went wrong closing connection: %v", err)
	}

	queries := sqlc.New(pool)

	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(queries)

	mux := router.New(*healthHandler, *authHandler)

	slog.Info("starting server...")

	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		fmt.Println(err)
	}
}
