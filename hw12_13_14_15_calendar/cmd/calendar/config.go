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
	Storage struct {
		DSN   string `config:"db_dsn"`
		Store string `config:"db_store"`
	}
	HTTPServer struct {
		Host string `config:"http_host"`
		Port string `config:"http_port"`
	}
	GRPCServer struct {
		Host string `config:"grpc_host"`
		Port string `config:"grpc_port"`
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
