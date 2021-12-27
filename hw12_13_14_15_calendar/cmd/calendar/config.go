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
		DSN      string `config:"db_dsn"`
		InMemory bool   `config:"in_memory"`
	}
	HTTPServer struct {
		Host string `config:"http_host"`
		Port string `config:"http_port"`
	}
	GRPCServer struct {
		Addr string `config:"grpc_addr"`
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
