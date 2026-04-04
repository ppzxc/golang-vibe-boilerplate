// Package main is the entrypoint for the API server.
package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ppzxc/golang-vibe-boilerplate/internal/adapter/httphandler"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/adapter/postgresrepo"
	apptodo "github.com/ppzxc/golang-vibe-boilerplate/internal/app/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/config"
)

func main() {
	// Structured logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := config.Load()

	// Open database connection
	db, err := postgresrepo.Open(cfg.Database.DSN())
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	if err := pingDB(db); err != nil {
		slog.Error("database ping failed", "error", err)
		os.Exit(1)
	}

	// Manual DI: wire up all dependencies
	todoRepo := postgresrepo.NewTodoRepository(db)
	todoSvc := apptodo.NewService(todoRepo)
	router := httphandler.NewRouter(todoSvc)

	// HTTP server
	srv := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start server in background, send errors via channel
	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case <-ctx.Done():
		slog.Info("server shutting down")
	case err := <-errCh:
		slog.Error("server error", "error", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}
	slog.Info("server stopped")
}

func pingDB(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}
