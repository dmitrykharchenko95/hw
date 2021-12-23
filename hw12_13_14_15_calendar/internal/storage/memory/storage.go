package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/sirupsen/logrus"
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
	s.mu.Lock()
	s.seq++
	event.ID = s.seq
	s.Data[s.seq] = event
	s.mu.Unlock()

	s.logg.Infof("event created: id %d", event.ID)
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	if _, ok := s.Data[event.ID]; !ok {
		return storage.ErrNoteNotFound
	}

	if event.Title == "" {
		event.Title = s.Data[event.ID].Title
	}

	if event.StartDate.IsZero() {
		event.StartDate = s.Data[event.ID].StartDate
	}

	if event.EndDate.IsZero() {
		event.EndDate = s.Data[event.ID].EndDate
	}

	if event.Content == "" {
		event.Content = s.Data[event.ID].Content
	}

	if event.UserID == 0 {
		event.UserID = s.Data[event.ID].UserID
	}

	if event.SendTime == 0 {
		event.SendTime = s.Data[event.ID].SendTime
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
		return storage.ErrNoteNotFound
	}
	delete(s.Data, id)
	s.mu.Unlock()

	s.logg.Infof("event deleted: id %d", id)
	return nil
}

func (s *Storage) ListEventsForDay(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return s.listEventsForPeriod(ctx, t, 24)
}

func (s *Storage) ListEventsForWeek(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return s.listEventsForPeriod(ctx, t, 168)
}

func (s *Storage) ListEventsForMonth(ctx context.Context, t time.Time) ([]storage.Event, error) {
	return s.listEventsForPeriod(ctx, t, 720)
}

func (s *Storage) listEventsForPeriod(c context.Context, t time.Time, d int) ([]storage.Event, error) { //nolint:unparam
	events := make([]storage.Event, 0, 7)

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, event := range s.Data {
		if !(event.StartDate.After(t.Add(time.Hour*time.Duration(d))) || event.EndDate.Before(t)) {
			events = append(events, event)
		}
	}

	s.logg.Infof("%d events from %v to %v listed", len(events), t.Format("2006-02-01 15:04"),
		t.Add(time.Hour*time.Duration(d)).Format("2006-02-01 15:04"))
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID > events[j].ID
	})
	return events, nil
}
