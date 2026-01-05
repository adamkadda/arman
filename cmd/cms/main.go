package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/adamkadda/arman/pkg/logging"
)

func main() {
	ctx, done := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	logger := logging.NewLoggerFromEnv()
	ctx = logging.WithLogger(ctx, logger)

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Error(
				"application panic",
				slog.Any("panic", r),
				slog.String("stack", string(debug.Stack())),
			)
			panic(r)
		}
	}()

	err := realMain(ctx)
	done()

	if err != nil {
		logger.Error(
			"application error",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	logger.Info("successful shutdown")
}

func realMain(ctx context.Context) error {
	// TODO: Load env

	// TODO: Setup DB connection

	// TODO: Setup middleware stack

	// TODO: Register routes

	// TODO: Start listening
	return nil
}
