package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/adamkadda/arman/internal/cms"
	"github.com/adamkadda/arman/internal/cms/handler"
	"github.com/adamkadda/arman/pkg/database"
	"github.com/adamkadda/arman/pkg/logging"
	"github.com/adamkadda/arman/pkg/server"
	"github.com/caarlos0/env/v11"
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

	err := start(ctx)
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

func start(ctx context.Context) error {
	var cfg cms.Config
	if err := env.Parse(&cfg); err != nil {
		return err
	}

	// TODO: Initialize pkg

	db, err := database.NewWithConfig(ctx, cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close(ctx)

	server, err := server.New(cfg.Port)
	if err != nil {
		return err
	}

	router := handler.RegisterRoutes(db.Pool)

	return server.ServeHTTPHandler(ctx, router)
}
