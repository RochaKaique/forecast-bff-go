package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/RochaKaique/forecastgo/cmd/server"
)

func main() {
	ctx := context.Background()

	configs, err := server.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	srv := server.CreateServer(configs)
	ConfigureLog(configs.GetString("application.name"))

	srv.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	slog.WarnContext(ctx, "received Interrupt signal... shutting down server")

	ctx, cancelFn := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFn()

	srv.Shutdown(ctx)
}

func ConfigureLog(appName string) {
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
