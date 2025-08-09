package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RochaKaique/forecastgo/cmd/server"
)

func main() {
	start := time.Now()
	cfg, err := server.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	configureLog(cfg.GetString("application.name"))

	srv := server.CreateServer(cfg)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("application started",
		"startup_time", fmt.Sprintf("%dms", time.Since(start).Milliseconds()),
		"port", cfg.GetString("server.port"))

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Start() }()

	select {
	case <-ctx.Done():
		slog.Warn("received shutdown signal, shutting down...")
	case err := <-errCh:
		if err != nil {
			slog.Error("server stopped with error", "err", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}

func configureLog(appName string) {
	serverName, _ := os.Hostname()

	handlerLevel := slog.LevelInfo
	logLevel := os.Getenv("LOG_LEVEL")

	switch logLevel {
	case "DEBUG":
		handlerLevel = slog.LevelDebug
	case "WARN":
		handlerLevel = slog.LevelWarn
	case "ERROR":
		handlerLevel = slog.LevelError
	default:
		handlerLevel = slog.LevelInfo
	}

	attrs := []slog.Attr{
		slog.String("application", appName),
		slog.String("serverName", serverName),
	}

	traceHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Any("_time", a.Value)
			}
			return a
		},
		AddSource: true,
		Level:     handlerLevel,
	})

	slog.SetDefault(slog.New(traceHandler.WithAttrs(attrs)))
}
