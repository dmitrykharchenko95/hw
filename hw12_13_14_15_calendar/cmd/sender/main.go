package main

import (
	"context"
	"flag"
	"log"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/rabbitmq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender_config.json", "Path to scheduler configuration file")
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

	consumer := rabbitmq.NewConsumer(
		config.RabbitMQ.Port,
		config.RabbitMQ.Host,
		config.RabbitMQ.User,
		config.RabbitMQ.Password,
		config.RabbitMQ.ExchangeName,
		config.RabbitMQ.ExchangeType,
		config.RabbitMQ.Queue,
		config.RabbitMQ.BindingKey,
		config.RabbitMQ.Tag)

	sender := NewSender(logg, consumer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = sender.ListenQueue(ctx)
	if err != nil {
		sender.Logger.Warnf("can not run sender:%v", err)
	}
}
