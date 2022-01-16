package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/server"
	internalgrpc "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/server/http"
	_ "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/migrations"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/calendar_config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig()
	if err != nil {
		log.Fatal(err, config)
	}

	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatal("cannot create new logger", err)
	}

	store, err := InitStore(config.Storage.Store, config.Storage.DSN, logg)
	if err != nil {
		logg.Warn(err)
		return
	}

	err = store.Connect(context.Background())
	if err != nil {
		logg.Warnf("can not connect with DB: %v", err)
		return
	}

	httpServer := internalhttp.NewServer(logg, &store, config.HTTPServer.Host, config.HTTPServer.Port)
	grpcServer := internalgrpc.NewSever(logg, &store, config.GRPCServer.Host, config.GRPCServer.Port)

	srv := server.NewServer(httpServer, grpcServer)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := srv.HTTP.Stop(ctx); err != nil {
			logg.Errorf("failed to stop http srv: %v", err)
		}

		if err := srv.GRPC.Stop(ctx); err != nil {
			logg.Errorf("failed to stop grpc srv: %v", err)
		}
	}()

	logg.Info("calendar is running...")

	if err := srv.Start(ctx); err != nil {
		logg.Errorf("failed to start srv: %v", err)
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
