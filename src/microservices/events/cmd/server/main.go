package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"events/internal/application"
	"events/internal/config"
	"events/internal/logger"
	"events/internal/router"
	"events/internal/worker"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conf, err := config.NewConfig(os.Args[1:])
	if err != nil {
		return fmt.Errorf("run: init config fail: %w", err)
	}

	zLog, err := logger.Build(conf.LogLevel)
	if err != nil {
		return fmt.Errorf("run: init logger fail: %w", err)
	}

	app := application.App{
		Log:  zLog,
		Conf: &conf,
	}

	g, gCtx := errgroup.WithContext(mainCtx)

	zLog.Info(
		"Server starting...",
		zap.String("address", conf.Address),
		zap.String("log level", conf.LogLevel),
	)

	server := http.Server{
		Addr:    conf.Address,
		Handler: router.New(&app),
		BaseContext: func(_ net.Listener) context.Context {
			return gCtx
		},
	}

	g.Go(server.ListenAndServe)

	g.Go(func() error {
		zLog.Info("Worker run")
		return worker.NewWorker(&app).Run(gCtx)
	})

	g.Go(func() error {
		<-mainCtx.Done()
		zLog.Info("Server stopped")
		return server.Shutdown(gCtx)
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("exit reason: %w", err)
	}

	return nil
}
