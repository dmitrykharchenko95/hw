package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/server"
	internalgrpc "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/migrations"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.json", "Path to configuration file")
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

	var store storage.Storage

	if config.Storage.InMemory {
		store = sqlstorage.New(config.Storage.DSN, logg)
	} else {
		store = memorystorage.New(logg)
	}

	calendar := app.New(logg, store)
	_ = calendar
	httpServer := internalhttp.NewServer(logg, config.HTTPServer.Host, config.HTTPServer.Port)
	grpcServer := internalgrpc.NewSever(logg, config.GRPCServer.Addr)

	srv := server.NewServer(httpServer, grpcServer)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := srv.HTTP.Stop(ctx); err != nil {
			logg.Error("failed to stop http srv: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := srv.HTTP.Start(ctx); err != nil {
		logg.Error("failed to start http srv: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
