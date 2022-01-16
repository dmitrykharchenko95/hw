package memorystorage

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/sirupsen/logrus"
)

const (
	day   = time.Hour * 24
	week  = time.Hour * 168
	month = time.Hour * 720
)

type Storage struct {
	seq  int64
	Data map[int64]storage.Event
	mu   sync.RWMutex
	logg *logrus.Logger
}

func New(logg *logrus.Logger) *Storage {
	return &Storage{
		Data: make(map[int64]storage.Event),
		logg: logg,
	}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.logg.Info("storage in-memory started")
	return
}

func (s *Storage) Close() error {
	s.logg.Info("storage in-memory finished")
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	err := storage.NewEventValidate(event)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.seq++
	event.ID = s.seq
	s.Data[s.seq] = event
	s.mu.Unlock()

	s.logg.Infof("event created: id %d\n", event.ID)
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	if _, ok := s.Data[event.ID]; !ok {
		return storage.ErrEventNotFound
	}
	s.Data[event.ID] = event
	s.mu.Unlock()

	s.logg.Infof("event updated: id %d", event.ID)
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	if id < 1 {
		return storage.ErrWrongID
	}

	s.mu.Lock()
	if _, ok := s.Data[id]; !ok {
		return storage.ErrEventNotFound
	}
	delete(s.Data, id)
	s.mu.Unlock()

	s.logg.Infof("event deleted: id %d", id)
	return nil
}

func (s *Storage) listEventsForPeriod(
	c context.Context, //nolint:unparam
	t time.Time,
	d time.Duration,
) ([]storage.Event, error) {
	events := make([]storage.Event, 0, 2)

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, event := range s.Data {
		if !(event.StartDate.After(t.Add(d)) || event.EndDate.Before(t)) {
			events = append(events, event)
		}
	}

	s.logg.Infof("%d events from %v to %v listed", len(events), t.Format("2006-02-01 15:04"),
		t.Add(d).Format("2006-02-01 15:04"))
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID > events[j].ID
	})
	return events, nil
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

func (s *Storage) GetEventForNotification(ctx context.Context, t time.Duration) ([]storage.Event, error) {
	now := time.Now()
	events := make([]storage.Event, 0, 1)

	for _, e := range s.Data {
		if e.StartDate.Add(-e.SendTime).After(now) && e.StartDate.Add(-e.SendTime).Before(now.Add(t)) {
			events = append(events, e)
		}
	}
	s.logg.Infof("Got %d events for notification", len(events))
	return events, nil
}

func (s *Storage) DeleteOldEvents(ctx context.Context) error {
	now := time.Now()

	for _, e := range s.Data {
		if e.EndDate.Before(now.AddDate(-1, 0, 0)) {
			err := s.DeleteEvent(ctx, e.ID)
			if err != nil {
				s.logg.Errorf("can not delete event: %v", err)
				return fmt.Errorf("can not delete event: %w", err)
			}
		}
	}

	return nil
}
