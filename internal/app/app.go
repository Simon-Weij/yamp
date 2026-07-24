package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Simon-Weij/yamp/internal/db"
	"github.com/Simon-Weij/yamp/internal/db/sqlc"
	"github.com/Simon-Weij/yamp/internal/handler"
	"github.com/Simon-Weij/yamp/internal/router"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func buildDSN() (string, error) {
	host, err := requireEnv("POSTGRES_HOST")
	if err != nil {
		return "", err
	}
	port, err := requireEnv("POSTGRES_PORT")
	if err != nil {
		return "", err
	}
	user, err := requireEnv("POSTGRES_USER")
	if err != nil {
		return "", err
	}
	pass, err := requireEnv("POSTGRES_PASSWORD")
	if err != nil {
		return "", err
	}
	dbname, err := requireEnv("POSTGRES_DB")
	if err != nil {
		return "", err
	}
	sslmode, err := requireEnv("POSTGRES_SSLMODE")
	if err != nil {
		return "", err
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   "/" + dbname,
	}
	q := u.Query()
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func requireEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("missing required env var: %s", key)
	}
	return v, nil
}

func waitForDB(ctx context.Context, db *sql.DB) error {
	const maxRetries = 30
	for i := range maxRetries {
		if err := db.PingContext(ctx); err != nil {
			slog.Warn("waiting for db", "attempt", i+1, "err", err)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Second):
			}
			continue
		}
		return nil
	}
	return fmt.Errorf("db did not become ready after %d retries", maxRetries)
}

func Run(ctx context.Context, ln net.Listener) error {
	dsn, err := buildDSN()
	if err != nil {
		return fmt.Errorf("build dsn: %w", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer pool.Close()

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open sql db for migrations: %w", err)
	}

	if err := waitForDB(ctx, sqlDB); err != nil {
		pool.Close()
		return fmt.Errorf("wait for db: %w", err)
	}

	if err := db.Migrate(sqlDB); err != nil {
		pool.Close()
		return fmt.Errorf("migrate: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("something went wrong closing connection: %w", err)
	}

	queries := sqlc.New(pool)
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(queries, queries)
	mux := router.New(*healthHandler, *authHandler)

	srv := &http.Server{
		Handler: mux,
	}

	serveErrCh := make(chan error, 1)
	go func() {
		slog.Info("starting server...", "addr", ln.Addr().String())
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErrCh <- fmt.Errorf("serve: %w", err)
			return
		}
		serveErrCh <- nil
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
		return <-serveErrCh
	case err := <-serveErrCh:
		return err
	}
}
