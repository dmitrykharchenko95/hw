package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/sirupsen/logrus"
)

func hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, _ = w.Write([]byte("Welcome to Calendar!"))
}

type Response struct {
	Info  string
	Data  []storage.Event
	Error string
}

func WriteResponse(w http.ResponseWriter, resp *Response, logg *logrus.Logger) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		logg.Errorf("response marshal error: %s", err)
	}
	_, err = w.Write(resBuf)
	if err != nil {
		logg.Errorf("response write error: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (s *Server) createEvent() func(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{}
		if r.Method != http.MethodPost {
			resp.Error = fmt.Sprintf("method %s not supported on uri %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusMethodNotAllowed)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("unsupported method %s, must be POST", r.Method)
			return
		}

		event := storage.Event{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&event)

		if err != nil && !errors.Is(err, io.EOF) {
			resp.Error = fmt.Sprintf("wrong format input:%v", err)
			w.WriteHeader(http.StatusBadRequest)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("can not decode request body: %v", err)
			return
		}

		err = (*s.store).CreateEvent(context.Background(), event)
		if err != nil {
			resp.Error = fmt.Sprintf("can not create event:%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("can not create event: %v", err)
			return
		}

		resp.Info = "Event created" //nolint:goconst
		WriteResponse(w, resp, s.Logg)
	}
}

func (s *Server) listEvents(period string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{}
		if r.Method != http.MethodGet {
			resp.Error = fmt.Sprintf("method %s not supported on uri %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusMethodNotAllowed)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("unsupported method %s, must be GET", r.Method)
			return
		}

		in := make([]byte, r.ContentLength)
		_, err := r.Body.Read(in)
		if err != nil && !errors.Is(err, io.EOF) {
			resp.Error = fmt.Sprintf("can not read request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("can not read request body: %v", err)
			return
		}

		date, err := time.Parse("2006-01-02", string(in))
		if err != nil && !errors.Is(err, io.EOF) {
			resp.Error = "wrong date format, must be 'YYYY-MM-DD'"
			w.WriteHeader(http.StatusBadRequest)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("can not parse date: %v", err)
			return
		}

		events := make([]storage.Event, 0, 1)

		switch period {
		case "day":
			events, err = (*s.store).ListEventsForDay(context.Background(), date)
		case "week":
			events, err = (*s.store).ListEventsForWeek(context.Background(), date)
		case "month":
			events, err = (*s.store).ListEventsForMonth(context.Background(), date)
		}

		if err != nil {
			resp.Error = fmt.Sprintf("events search error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("events search error: %v", err)
			return
		}

		resp.Data = events
		WriteResponse(w, resp, s.Logg)
	}
}

func (s *Server) updateEvent( //nolint:dupl
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{}
		if r.Method != http.MethodPut {
			resp.Error = fmt.Sprintf("method %s not supported on uri %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusMethodNotAllowed)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("unsupported method %s, must be PUT", r.Method)
			return
		}

		event := storage.Event{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&event)

		if err != nil && !errors.Is(err, io.EOF) {
			resp.Error = fmt.Sprintf("wrong format input:%v", err)
			w.WriteHeader(http.StatusBadRequest)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("can not decode request body: %v", err)
			return
		}

		err = (*s.store).UpdateEvent(context.Background(), event)
		if err != nil {
			resp.Error = fmt.Sprintf("can not update event:%v", err)
			w.WriteHeader(http.StatusInternalServerError)

			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("can not update event: %v", err)
			return
		}

		resp.Info = "Event updated"
		WriteResponse(w, resp, s.Logg)
	}
}

func (s *Server) deleteEvent() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{}
		if r.Method != http.MethodDelete {
			resp.Error = fmt.Sprintf("method %s not supported on uri %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusMethodNotAllowed)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("unsupported method %s, must be DELETE", r.Method)
			return
		}

		in := make([]byte, r.ContentLength)
		_, err := r.Body.Read(in)
		if err != nil && !errors.Is(err, io.EOF) {
			resp.Error = fmt.Sprintf("can not read request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("can not read request body: %v", err)
			return
		}

		noteID, err := strconv.Atoi(string(in))
		if err != nil {
			resp.Error = "event id must be digit"
			w.WriteHeader(http.StatusBadRequest)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("can not convert request body to int: %v", err)
			return
		}

		err = (*s.store).DeleteEvent(context.Background(), int64(noteID))
		if err != nil {
			resp.Error = fmt.Sprintf("can not delete event:%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			WriteResponse(w, resp, s.Logg)
			s.Logg.Errorf("can not delete event: %v", err)
			return
		}

		resp.Info = "Event deleted"
		WriteResponse(w, resp, s.Logg)
	}
}
