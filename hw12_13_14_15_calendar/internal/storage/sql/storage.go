package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v4/stdlib" //nolint
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	db   *sqlx.DB
	dsn  string
	logg *logrus.Logger
}

func New(dsn string, logg *logrus.Logger) *Storage {
	return &Storage{
		db:   new(sqlx.DB),
		dsn:  dsn,
		logg: logg,
	}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.ConnectContext(ctx, "pgx", s.dsn)
	if err != nil {
		return err
	}

	err = s.db.Ping()
	if err != nil {
		err = fmt.Errorf("cannot ping with db: %w", err)
		return err
	}

	s.logg.Info("storage connected")
	return nil
}

func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}

	logrus.Info("storage closed")
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	_, err := s.db.NamedExecContext(ctx, `
		INSERT INTO events (title, start_date, end_date, content, user_id, send_time)
		VALUES (:title, :start_date, :end_date, :content, :user_id, :send_time)`, &event)
	if err != nil {
		return err
	}

	s.logg.Info("event created")
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	oldEvent := &storage.Event{}
	err = tx.QueryRowxContext(ctx, `
		SELECT title, start_date, end_date, content, user_id, send_tame
		FROM events
		WHERE id= $1`,
		event.ID).StructScan(&oldEvent)

	if err != nil {
		return err
	}

	if event.Title != "" {
		oldEvent.Title = event.Title
	}

	if !event.StartDate.IsZero() {
		oldEvent.StartDate = event.StartDate
	}

	if !event.EndDate.IsZero() {
		oldEvent.EndDate = event.EndDate
	}

	if event.Content != "" {
		oldEvent.Content = event.Content
	}

	if event.UserID != 0 {
		oldEvent.UserID = event.UserID
	}

	if event.SendTime != 0 {
		oldEvent.SendTime = event.SendTime
	}

	_, err = s.db.ExecContext(ctx, `
		UPDATE events 
		SET title=:title, start_date=:start_date, end_date=:end_date, content=:content, user_id=:user_id, send_time=:send_time
		WHERE id=:id`,
		oldEvent)

	if err != nil {
		return err
	}

	s.logg.WithField("event_id", event.ID).Info("event update")
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	if id < 1 {
		return storage.ErrWrongID
	}

	_, err := s.db.ExecContext(ctx, `
		DELETE FROM events
		WHERE id=$1 `, id)
	if err != nil {
		return err
	}

	s.logg.WithField("event_id", id).Info("event deleted")
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id int64) (*storage.Event, error) { // удалить
	event := new(storage.Event)

	err := s.db.QueryRowxContext(ctx, `
		SELECT * FROM events 
		WHERE id=$1`, id).StructScan(event)
	if err != nil {
		return nil, err
	}

	s.logg.WithField("event_id", id).Info("event got")
	return event, nil
}

func (s *Storage) listEventsForPeriod(ctx context.Context, t time.Time, d time.Duration) ([]storage.Event, error) {
	events := make([]storage.Event, 0, 1)

	rows, err := s.db.QueryxContext(ctx, //nolint:sqlclosecheck
		`SELECT * FROM events 
		WHERE NOT 
		(end_date < $1 AND start_date > $2)`,
		t, t.Add(d))
	if err != nil {
		return events, err
	}

	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("cannot close rows: %w", err)
			s.logg.Error(err)
		}
	}(rows)

	for rows.Next() {
		var event storage.Event
		if err := rows.StructScan(&event); err != nil {
			return events, err
		}
		events = append(events, event)
	}

	s.logg.Info("events for day are got")
	return events, err
}

func (s *Storage) ListEventsForDay(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return s.listEventsForPeriod(ctx, t, time.Hour*64)
}

func (s *Storage) ListEventsForWeek(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return s.listEventsForPeriod(ctx, t, time.Hour*168)
}

func (s *Storage) ListEventsForMonth(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return s.listEventsForPeriod(ctx, t, time.Hour*720)
}
