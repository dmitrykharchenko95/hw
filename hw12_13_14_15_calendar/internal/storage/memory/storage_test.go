package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	testLogger = logrus.New()
	wg         sync.WaitGroup
	testEvent1 = storage.Event{
		ID:        1,
		Title:     "Title 1",
		StartDate: time.Date(2021, time.November, 10, 23, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2021, time.November, 11, 23, 0, 0, 0, time.UTC),
		Content:   "Test content #1",
		UserID:    1,
		SendTime:  time.Hour,
	}
	testEvent2 = storage.Event{
		ID:        2,
		Title:     "Title 2",
		StartDate: time.Date(2021, time.November, 16, 23, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2021, time.November, 17, 23, 0, 0, 0, time.UTC),
		Content:   "Test content #2",
		UserID:    2,
		SendTime:  time.Hour,
	}
	testEvent3 = storage.Event{
		ID:        3,
		Title:     "Title 3",
		StartDate: time.Date(2021, time.November, 29, 23, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2021, time.November, 30, 23, 0, 0, 0, time.UTC),
		Content:   "Test content #3",
		UserID:    3,
		SendTime:  time.Hour,
	}
)

func TestStorage_CreateEvent(t *testing.T) {
	testStorage := New(testLogger)
	tests := []storage.Event{testEvent1, testEvent2, testEvent3}

	wg.Add(3)
	for i, tt := range tests {
		tt := tt
		i := i
		go t.Run(fmt.Sprintf("Parallel #%d", i), func(t *testing.T) {
			defer wg.Done()
			err := testStorage.CreateEvent(context.Background(), tt)
			require.NoError(t, err)
		})
	}
	wg.Wait()
	require.Len(t, testStorage.Data, 3)

	testStorage = New(testLogger)

	t.Run("base", func(t *testing.T) {
		err := testStorage.CreateEvent(context.Background(), testEvent1)

		require.NoError(t, err)
		require.Equal(t, &Storage{seq: 1, Data: map[int64]storage.Event{1: testEvent1}, logg: testLogger}, testStorage)
	})
}

func TestStorage_DeleteEvent(t *testing.T) {
	testStorage := New(testLogger)
	testStorage.Data = map[int64]storage.Event{1: testEvent1, 2: testEvent2, 3: testEvent3}
	wg.Add(3)
	for i := 1; i < 4; i++ {
		i := i
		go t.Run(fmt.Sprintf("parallel #%d", i), func(t *testing.T) {
			defer wg.Done()
			err := testStorage.DeleteEvent(context.Background(), int64(i))
			require.NoError(t, err)
		})
	}

	wg.Wait()
	require.Equal(t, map[int64]storage.Event{}, testStorage.Data)

	testStorage = &Storage{
		seq:  2,
		Data: map[int64]storage.Event{1: testEvent1, 2: testEvent2},
		logg: testLogger,
	}

	t.Run("base", func(t *testing.T) {
		err := testStorage.DeleteEvent(context.Background(), 2)
		require.NoError(t, err)
		require.Equal(t, &Storage{seq: 2, Data: map[int64]storage.Event{1: testEvent1}, logg: testLogger}, testStorage)
	})

	t.Run("wrong id", func(t *testing.T) {
		err := testStorage.DeleteEvent(context.Background(), -1)
		require.ErrorIs(t, err, storage.ErrWrongID)
	})

	t.Run("note not found", func(t *testing.T) {
		err := testStorage.DeleteEvent(context.Background(), 5)
		require.ErrorIs(t, err, storage.ErrNoteNotFound)
	})
}

func TestStorage_UpdateEvent(t *testing.T) {
	testStorage := &Storage{
		seq:  3,
		Data: map[int64]storage.Event{1: testEvent1, 2: testEvent2, 3: testEvent3},
		logg: testLogger,
	}

	newEvent := storage.Event{
		ID:        2,
		Title:     "New Title 2",
		StartDate: time.Date(2021, time.October, 10, 23, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2021, time.October, 11, 23, 0, 0, 0, time.UTC),
		Content:   "New Test content #2",
		UserID:    1,
		SendTime:  time.Hour * 4,
	}

	expStorage := &Storage{
		seq:  3,
		Data: map[int64]storage.Event{1: testEvent1, 2: newEvent, 3: testEvent3},
		logg: testLogger,
	}

	t.Run("1 event", func(t *testing.T) {
		err := testStorage.UpdateEvent(context.Background(), newEvent)
		require.NoError(t, err)
		require.Equal(t, expStorage, testStorage)
	})

	t.Run("event not exist", func(t *testing.T) {
		err := testStorage.UpdateEvent(context.Background(), storage.Event{ID: 4})
		require.ErrorIs(t, err, storage.ErrNoteNotFound)
	})
}

func TestStorage_ListEventsForDay(t *testing.T) {
	testStorage := &Storage{
		seq:  3,
		Data: map[int64]storage.Event{1: testEvent1, 2: testEvent2, 3: testEvent3},
		logg: testLogger,
	}

	testDates := []time.Time{
		time.Date(2021, time.November, 10, 23, 0, 0, 0, time.UTC),
		time.Date(2021, time.November, 16, 23, 0, 0, 0, time.UTC),
		time.Date(2021, time.November, 29, 23, 0, 0, 0, time.UTC),
	}

	for i, td := range testDates {
		td := td
		i := i
		t.Run(fmt.Sprintf("base #%d", i), func(t *testing.T) {
			t.Parallel()

			e, err := testStorage.ListEventsForDay(context.Background(), td)
			require.NoError(t, err)
			require.Equal(t, testStorage.Data[int64(i+1)], e[0])
		})
	}

	t.Run("note not exist", func(t *testing.T) {
		e, err := testStorage.ListEventsForDay(context.Background(),
			time.Date(2021, time.December, 10, 23, 0, 0, 0, time.UTC))
		require.NoError(t, err)
		require.Equal(t, []storage.Event{}, e)
	})
}

func TestStorage_ListEventsForWeek(t *testing.T) {
	testStorage := &Storage{
		seq:  3,
		Data: map[int64]storage.Event{1: testEvent1, 2: testEvent2, 3: testEvent3},
		logg: testLogger,
	}

	testDates := []time.Time{
		time.Date(2021, time.November, 16, 23, 0, 0, 0, time.UTC),
		time.Date(2021, time.November, 10, 23, 0, 0, 0, time.UTC),
		time.Date(2021, time.December, 29, 23, 0, 0, 0, time.UTC),
	}

	t.Run("1 event", func(t *testing.T) {
		e, err := testStorage.ListEventsForWeek(context.Background(), testDates[0])
		require.NoError(t, err)
		require.Equal(t, []storage.Event{testEvent2}, e)
	})

	t.Run("several events", func(t *testing.T) {
		e, err := testStorage.ListEventsForWeek(context.Background(), testDates[1])
		require.NoError(t, err)
		require.Equal(t, []storage.Event{testEvent2, testEvent1}, e)
	})

	t.Run("note not exist", func(t *testing.T) {
		e, err := testStorage.ListEventsForWeek(context.Background(), testDates[2])
		require.NoError(t, err)
		require.Equal(t, []storage.Event{}, e)
	})
}

func TestStorage_ListEventsForMonth(t *testing.T) {
	testStorage := &Storage{
		seq:  3,
		Data: map[int64]storage.Event{1: testEvent1, 2: testEvent2, 3: testEvent3},
		logg: testLogger,
	}

	testDates := []time.Time{
		time.Date(2021, time.November, 30, 23, 0, 0, 0, time.UTC),
		time.Date(2021, time.November, 10, 23, 0, 0, 0, time.UTC),
		time.Date(2021, time.December, 16, 23, 0, 0, 0, time.UTC),
	}

	t.Run("1 event", func(t *testing.T) {
		e, err := testStorage.ListEventsForMonth(context.Background(), testDates[0])
		require.NoError(t, err)
		require.Equal(t, []storage.Event{testEvent3}, e)
	})

	t.Run("several events", func(t *testing.T) {
		e, err := testStorage.ListEventsForMonth(context.Background(), testDates[1])
		require.NoError(t, err)
		require.Equal(t, []storage.Event{testEvent3, testEvent2, testEvent1}, e)
	})

	t.Run("note not exist", func(t *testing.T) {
		e, err := testStorage.ListEventsForMonth(context.Background(),
			time.Date(2021, time.December, 10, 23, 0, 0, 0, time.UTC))
		require.NoError(t, err)
		require.Equal(t, []storage.Event{}, e)
	})
}
