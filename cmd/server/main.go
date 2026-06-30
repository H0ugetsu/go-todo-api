package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/h0ugetsu/todo-api/internal/server"
	"github.com/h0ugetsu/todo-api/internal/todo"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	todoRepo := todo.NewRepository()
	todoService := todo.NewService(todoRepo)
	todoHandler := todo.NewHandler(todoService, logger)

	mux := server.NewRouter(todoHandler)

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", addr),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Info("🚀 Starting server...", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("😭 Failed to start server", "error", err)
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.Info("🌇 Shutting down server...")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("😭 Failed to shutdown server", "error", err)
			return err
		}
	case err := <-errCh:
		return err
	}

	return nil
}
