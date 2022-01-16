package app

import (
	"context"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/sirupsen/logrus"
)

type App struct {
	logger  *logrus.Logger
	storage *Storage
}

type Storage interface {
	Connect(ctx context.Context) (err error)
	Close() error
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, i int64) error

	ListEventsForDay(ctx context.Context, t time.Time) ([]storage.Event, error)
	ListEventsForWeek(ctx context.Context, t time.Time) ([]storage.Event, error)
	ListEventsForMonth(ctx context.Context, t time.Time) ([]storage.Event, error)
}

func New(logger *logrus.Logger, storage *Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
