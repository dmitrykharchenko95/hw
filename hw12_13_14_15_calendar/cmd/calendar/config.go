package main

import (
	"context"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"github.com/sirupsen/logrus"
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

func InitStore(storeType, dsn string, logg *logrus.Logger) (storage.Storage, error) {
	var store storage.Storage

	switch storeType {
	case "in-memory":
		logg.Info("Use in-memory storage")
		store = memorystorage.New(logg)
	case "sql":
		logg.Info("Use sql storage")
		store = sqlstorage.New(dsn, logg)
	default:
		logg.Errorf("wrong srorage type:%v", storeType)
		return nil, storage.ErrWrongStorageType
	}

	return store, nil
}
