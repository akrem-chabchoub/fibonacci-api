package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"

	"github.com/akrem-chabchoub/fibonacci-api/internal/server"
)

func main() {
	listen := flag.String("listen", ":8080", "HTTP listen address")
	shutdownTimeout := flag.Duration("shutdown-timeout", 10*time.Second, "time to finish in-flight requests")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync() //nolint:errcheck

	_, _ = maxprocs.Set(maxprocs.Logger(logger.Sugar().Infof))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := server.New(server.Options{
		Addr:            *listen,
		ShutdownTimeout: *shutdownTimeout,
		Logger:          logger,
	})

	if err := srv.Run(ctx); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
