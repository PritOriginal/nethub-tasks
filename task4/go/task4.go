package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// Device простая модель устройства
type Device struct {
	ID       int64
	Hostname string
	IP       string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	app := NewApp(logger)

	go func() {
		app.MustRun()
	}()

	// Graceful shutdown

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done

	app.Stop()

	logger.Info("server stopped")
}
