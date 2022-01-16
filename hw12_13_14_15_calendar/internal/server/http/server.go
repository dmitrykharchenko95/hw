package internalhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/sirupsen/logrus"
)

type Server struct {
	logg  *logrus.Logger
	srv   *http.Server
	addr  string
	store *storage.Storage
}

func NewServer(logger *logrus.Logger, store *storage.Storage, host, port string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/create", createEvent(store, logger))
	mux.HandleFunc("/update", updateEvent(store, logger))
	mux.HandleFunc("/delete", deleteEvent(store, logger))
	mux.HandleFunc("/list/day", listEvents(store, logger, "day"))
	mux.HandleFunc("/list/week", listEvents(store, logger, "week"))
	mux.HandleFunc("/list/month", listEvents(store, logger, "month"))

	wrappedMux := loggingMiddleware(mux, logger)

	addr := net.JoinHostPort(host, port)
	return &Server{
		logg: logger,
		srv: &http.Server{
			Addr:         addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			Handler:      wrappedMux,
		},
		addr:  addr,
		store: store,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Infof("Start http server on %s...", s.addr)

	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("Stop http server...")

	return s.srv.Shutdown(ctx)
}
