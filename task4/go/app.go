package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	db     *sql.DB
	log    *slog.Logger
	server *http.Server
}

func NewApp(log *slog.Logger) *App {
	db, err := initDB()
	if err != nil {
		log.Error("failed to initialize database", slog.String("error", err.Error()))
		panic(err)
	}

	handler := NewHandler(db, log)

	mux := http.NewServeMux()
	mux.HandleFunc("/device", handler.DeviceHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		db:     db,
		log:    log,
		server: server,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	a.log.Info("server started", slog.String("address", ":8080"))
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.log.Error("failed to start server")
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (a *App) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.log.Error("an error occurred while stopping the server", slog.String("error", err.Error()))
	}

	if err := a.db.Close(); err != nil {
		a.log.Error("an error occurred while closing the connection to the database", slog.String("error", err.Error()))
	}
}
