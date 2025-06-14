package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/kiryu2k/dns-server/internal/server"
	"github.com/pkg/errors"
)

const (
	host = "localhost"
	port = "2053"
)

func main() {
	defer func() {
		fmt.Printf("Total goroutines before exit: %d\n", runtime.NumGoroutine())
	}()

	var (
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelDebug,
			ReplaceAttr: nil,
		}))
		ctx = context.Background()
	)

	srv, err := server.NewDns(host, port, logger)
	if err != nil {
		logger.ErrorContext(ctx, errors.WithMessage(err, "new dns server").Error())
		return
	}

	go func() {
		logger.InfoContext(ctx, "starting dns server...")
		if err := srv.ListenAndServe(ctx); err != nil {
			logger.ErrorContext(ctx, errors.WithMessage(err, "listen and serve dns server").Error())
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.InfoContext(ctx, "gracefully shutting down...")

	srv.Close()
	logger.InfoContext(ctx, "finished")
}
