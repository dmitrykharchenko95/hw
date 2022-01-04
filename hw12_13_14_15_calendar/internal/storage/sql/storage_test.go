package sqlstorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestStorage_Connect(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		testEvent := storage.Event{
			Title:     "Test event",
			StartDate: time.Date(2021, time.December, 20, 11, 0o0, 0o0, 0o0, time.Local),
			EndDate:   time.Date(2021, time.December, 30, 11, 0o0, 0o0, 0o0, time.Local),
			Content:   "This is test event!",
			UserID:    1,
			SendTime:  time.Hour,
		}

		s := New("host=localhost port=5432 user=username password=userpassword dbname=calendar sslmode=disable",
			&logrus.Logger{})

		ctx, _ := context.WithTimeout(context.Background(), time.Second*5) //nolint

		err := s.Connect(ctx)
		require.NoError(t, err)

		err = s.CreateEvent(ctx, testEvent)
		require.NoError(t, err)

		event, err := s.ListEventsForDay(ctx, time.Date(2021, time.December, 20, 10, 0o0, 0o0, 0o0, time.Local))
		require.NoError(t, err)
		for i, e := range event {
			fmt.Println("events", i, e)
		}

		err = s.DeleteEvent(ctx, 2)
		require.NoError(t, err)

		err = s.Close()
		require.NoError(t, err)
	})
}
