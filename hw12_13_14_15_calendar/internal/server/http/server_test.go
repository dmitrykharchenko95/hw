package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"testing"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	testEvent1 = `
		{
		  "ID": 1,
		  "Title": "New year",
		  "StartDate": "2022-01-01T00:00:00Z",
		  "EndDate": "2022-01-02T00:00:00Z",
		  "Content": "Happy New Year!!!",
		  "UserID": 1,
		  "SendTime": 12
		}`

	testEvent1_1 = `
		{
		  "ID": 2,
		  "Title": "New year",
		  "StartDate": "2022-01-01T00:00:00Z",
		  "EndDate": "2022-01-02T00:00:00Z",
		  "Content": "Happy New Year!!!",
		  "UserID": 1,
		  "SendTime": 12
		}`

	testEvent2 = `
		{
		  "ID": 1,
		  "Title": "Christmas",
		  "StartDate": "2022-01-07T00:00:00Z",
		  "EndDate": "2022-01-08T00:00:00Z",
		  "Content": "Merry Christmas!!!",
		  "UserID": 3,
		  "SendTime": 24
		}`
	logg                                   = logrus.New()
	store                  storage.Storage = memorystorage.New(logg)
	expectedRes, actualRes Response
	expectedEvent          = storage.Event{}
	buf                    = &bytes.Buffer{}
	resBody                = make([]byte, 0, 20)
	ghost                  = &bytes.Buffer{}
)

func TestServer(t *testing.T) { //nolint
	logg.SetOutput(ghost)
	server := NewServer(logg, &store, "0.0.0.0", "8080")

	go func() {
		err := server.Start(context.Background())
		require.NoError(t, err)
	}()

	defer func() {
		err := server.Stop(context.Background())
		require.NoError(t, err)
	}()
	t.Run("create", func(t *testing.T) {
		createCmd := exec.Command("curl", "-X", "POST", "-i", "0.0.0.0:8080/create", `--data`, testEvent1)

		createCmd.Stdout = buf
		err := createCmd.Run()
		require.NoError(t, err)

		for !errors.Is(err, io.EOF) {
			resBody, err = buf.ReadBytes(10)
			if !errors.Is(err, io.EOF) {
				require.NoError(t, err)
			}
		}

		expectedRes.Info = "Event created"

		err = json.Unmarshal(resBody, &actualRes)

		require.NoError(t, err)
		require.Equal(t, expectedRes, actualRes)
	})

	t.Run("list events for day", func(t *testing.T) {
		listEventsDayCmd := exec.Command("curl", "-X", "GET", "-i", "0.0.0.0:8080/list/day", "--data", "2022-01-01")

		listEventsDayCmd.Stdout = buf
		err := listEventsDayCmd.Run()
		require.NoError(t, err)

		for !errors.Is(err, io.EOF) {
			resBody, err = buf.ReadBytes(10)
			if !errors.Is(err, io.EOF) {
				require.NoError(t, err)
			}
		}

		err = json.Unmarshal([]byte(testEvent1), &expectedEvent)
		require.NoError(t, err)

		expectedRes.Info = ""
		expectedRes.Data = append(expectedRes.Data, expectedEvent)

		err = json.Unmarshal(resBody, &actualRes)
		require.NoError(t, err)
		require.Equal(t, expectedRes, actualRes)
	})
	t.Run("update event", func(t *testing.T) {
		updateEventCmd := exec.Command("curl", "-X", "PUT", "-i", "0.0.0.0:8080/update", "--data", testEvent2)

		updateEventCmd.Stdout = buf
		err := updateEventCmd.Run()
		require.NoError(t, err)

		for !errors.Is(err, io.EOF) {
			resBody, err = buf.ReadBytes(10)
			if !errors.Is(err, io.EOF) {
				require.NoError(t, err)
			}
		}

		expectedRes.Info = "Event updated"
		expectedRes.Data = nil

		err = json.Unmarshal(resBody, &actualRes)
		require.NoError(t, err)
		require.Equal(t, expectedRes, actualRes)
	})
	t.Run("list events for week", func(t *testing.T) {
		listEventsWeekCmd := exec.Command("curl", "-X", "GET", "-i", "0.0.0.0:8080/list/week", "--data", "2022-01-01")

		listEventsWeekCmd.Stdout = buf
		err := listEventsWeekCmd.Run()
		require.NoError(t, err)

		for !errors.Is(err, io.EOF) {
			resBody, err = buf.ReadBytes(10)
			if !errors.Is(err, io.EOF) {
				require.NoError(t, err)
			}
		}

		err = json.Unmarshal([]byte(testEvent2), &expectedEvent)
		require.NoError(t, err)

		expectedRes.Info = ""
		expectedRes.Data = append(expectedRes.Data, expectedEvent)

		err = json.Unmarshal(resBody, &actualRes)
		require.NoError(t, err)
		require.Equal(t, expectedRes, actualRes)
	})
	t.Run("create new event", func(t *testing.T) {
		createCmd := exec.Command("curl", "-X", "POST", "-i", "0.0.0.0:8080/create", `--data`, testEvent1)

		createCmd.Stdout = buf
		err := createCmd.Run()
		require.NoError(t, err)

		for !errors.Is(err, io.EOF) {
			resBody, err = buf.ReadBytes(10)
			if !errors.Is(err, io.EOF) {
				require.NoError(t, err)
			}
		}
		expectedRes.Info = "Event created"
		expectedRes.Data = nil

		err = json.Unmarshal(resBody, &actualRes)
		require.NoError(t, err)
		require.Equal(t, expectedRes, actualRes)
	})
	t.Run("list events for month", func(t *testing.T) {
		listEventsMonthCmd := exec.Command("curl", "-X", "GET", "-i", "0.0.0.0:8080/list/month", "--data", "2021-12-06")

		listEventsMonthCmd.Stdout = buf
		err := listEventsMonthCmd.Run()
		require.NoError(t, err)

		for !errors.Is(err, io.EOF) {
			resBody, err = buf.ReadBytes(10)
			if !errors.Is(err, io.EOF) {
				require.NoError(t, err)
			}
		}

		err = json.Unmarshal([]byte(testEvent1_1), &expectedEvent)
		require.NoError(t, err)

		expectedRes.Info = ""
		expectedRes.Data = append(expectedRes.Data, expectedEvent)

		err = json.Unmarshal(resBody, &actualRes)
		require.NoError(t, err)
		require.Equal(t, expectedRes, actualRes)
	})
	t.Run("delete event", func(t *testing.T) {
		deleteCmd := exec.Command("curl", "-X", "DELETE", "-i", "0.0.0.0:8080/delete", `--data`, "1")

		deleteCmd.Stdout = buf
		err := deleteCmd.Run()
		require.NoError(t, err)

		for !errors.Is(err, io.EOF) {
			resBody, err = buf.ReadBytes(10)
			if !errors.Is(err, io.EOF) {
				require.NoError(t, err)
			}
		}

		expectedRes.Info = "Event deleted"
		expectedRes.Data = nil

		err = json.Unmarshal(resBody, &actualRes)
		require.NoError(t, err)
		require.Equal(t, expectedRes, actualRes)
	})
	t.Run("list events for month", func(t *testing.T) {
		listEventsMonthCmd := exec.Command("curl", "-X", "GET", "-i", "0.0.0.0:8080/list/month", "--data", "2022-01-01")

		listEventsMonthCmd.Stdout = buf
		err := listEventsMonthCmd.Run()
		require.NoError(t, err)

		for !errors.Is(err, io.EOF) {
			resBody, err = buf.ReadBytes(10)
			if !errors.Is(err, io.EOF) {
				require.NoError(t, err)
			}
		}

		err = json.Unmarshal([]byte(testEvent1_1), &expectedEvent)
		require.NoError(t, err)

		expectedRes.Info = ""
		expectedRes.Data = append(expectedRes.Data, expectedEvent)

		err = json.Unmarshal(resBody, &actualRes)
		require.NoError(t, err)
		require.Equal(t, expectedRes, actualRes)
	})
}
