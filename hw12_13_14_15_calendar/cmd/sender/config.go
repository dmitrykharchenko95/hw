package main

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

type Config struct {
	Logger struct {
		Level string `config:"level"`
	}
	RabbitMQ struct {
		Host         string `config:"rmq_host"`
		Port         int    `config:"rmq_port"`
		User         string `config:"rmq_user"`
		Password     string `config:"rmq_password"`
		ExchangeName string `config:"rmq_exchangename"`
		ExchangeType string `config:"rmq_exchangetype"`
		Queue        string `config:"rmq_queue"`
		BindingKey   string `config:"rmq_bindingkey"`
		Tag          string `config:"rmq_tag"`
	}
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	l := confita.NewLoader(file.NewBackend(configFile))

	if err := l.Load(context.Background(), cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
