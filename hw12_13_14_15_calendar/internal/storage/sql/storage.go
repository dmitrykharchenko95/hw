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

const (
	day   = time.Hour * 24
	week  = time.Hour * 168
	month = time.Hour * 720
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
	s.logg.Info("sql storage connected")
	return nil
}

func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}

	logrus.Info("sql storage closed")
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	err := storage.NewEventValidate(event)
	if err != nil {
		return err
	}

	_, err = s.db.NamedExecContext(ctx, `
		INSERT INTO events (title, start_date, end_date, content, user_id, send_time)
		VALUES (:title, :start_date, :end_date, :content, :user_id, :send_time)`, &event)
	if err != nil {
		return err
	}

	s.logg.Infof("event created")
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE events 
		SET title=:title, start_date=:start_date, end_date=:end_date, content=:content, user_id=:user_id, send_time=:send_time
		WHERE id=:id`,
		event)
	if err != nil {
		s.logg.Errorf("can not update event id: %v", err)
		return err
	}

	s.logg.Infof("event updated: id %d", event.ID)
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	if id < 1 {
		return storage.ErrWrongID
	}

	res, err := s.db.ExecContext(ctx, `
		DELETE FROM events
		WHERE id=$1 `, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return storage.ErrEventNotFound
	}

	s.logg.Infof("event deleted: id %d", id)
	return nil
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

	s.logg.Infof("%d events from %v to %v listed", len(events), t.Format("2006-02-01 15:04"),
		t.Add(d).Format("2006-02-01 15:04"))
	return events, err
}

func (s *Storage) ListEventsForDay(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return s.listEventsForPeriod(ctx, t, day)
}

func (s *Storage) ListEventsForWeek(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return s.listEventsForPeriod(ctx, t, week)
}

func (s *Storage) ListEventsForMonth(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return s.listEventsForPeriod(ctx, t, month)
}
