package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEventNotFound     = errors.New("event not found")
	ErrWrongID           = errors.New("wrong event id")
	ErrEmptyFieldOfEvent = errors.New("one of field (title, start date, end date or user id) is empty")
)

func NewEventValidate(e Event) error {
	if e.Title == "" || e.StartDate.IsZero() || e.EndDate.IsZero() || e.UserID == 0 {
		return ErrEmptyFieldOfEvent
	}
	return nil
}

type Event struct {
	ID        int64         `db:"id"`
	Title     string        `db:"title"`
	StartDate time.Time     `db:"start_date"`
	EndDate   time.Time     `db:"end_date"`
	Content   string        `db:"content"`
	UserID    int64         `db:"user_id"`
	SendTime  time.Duration `db:"send_time"`
}

type Storage interface {
	Connect(ctx context.Context) (err error)
	Close() error
	CreateEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, id int64) error

	ListEventsForDay(ctx context.Context, t time.Time) ([]Event, error)
	ListEventsForWeek(ctx context.Context, t time.Time) ([]Event, error)
	ListEventsForMonth(ctx context.Context, t time.Time) ([]Event, error)
}
