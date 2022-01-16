package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/rabbitmq"
	sqlstorage "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/scheduler_config.json", "Path to scheduler configuration file")
}

func main() {
	flag.Parse()

	config, err := NewConfig()
	if err != nil {
		log.Fatal(err, config)
	}

	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatal("cannot create new logger", err)
	}

	store := sqlstorage.New(config.Storage.DSN, logg)

	producer := rabbitmq.NewProducer(
		config.RabbitMQ.Port,
		config.RabbitMQ.Host,
		config.RabbitMQ.User,
		config.RabbitMQ.Password,
		config.RabbitMQ.ExchangeName,
		config.RabbitMQ.ExchangeType,
		config.RabbitMQ.Queue,
		config.RabbitMQ.BindingKey)

	period, err := time.ParseDuration(config.Scheduler.Period)
	if err != nil {
		logg.Errorf("can not parse period: %v", err)
		logg.Info("use default period - 30s")
		period = time.Second * 30
	}

	scheduler := NewScheduler(producer, store, logg, period)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = scheduler.SearchEventForNotification(ctx)
	if err != nil {
		logg.Fatalf("can not run scheduler:%v", err)
	}
}
