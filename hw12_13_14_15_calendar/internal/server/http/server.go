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
	Logg  *logrus.Logger
	srv   *http.Server
	addr  string
	store *storage.Storage
}

func NewServer(logger *logrus.Logger, store *storage.Storage, host, port string) *Server {
	addr := net.JoinHostPort(host, port)

	s := &Server{
		Logg: logger,
		srv: &http.Server{
			Addr:         addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
		addr:  addr,
		store: store,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/create", s.createEvent())
	mux.HandleFunc("/update", s.updateEvent())
	mux.HandleFunc("/delete", s.deleteEvent())
	mux.HandleFunc("/list/day", s.listEvents("day"))
	mux.HandleFunc("/list/week", s.listEvents("week"))
	mux.HandleFunc("/list/month", s.listEvents("month"))

	s.srv.Handler = loggingMiddleware(mux, logger)

	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.Logg.Infof("Start http server on %s...", s.addr)

	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.Logg.Info("Stop http server...")

	return s.srv.Shutdown(ctx)
}
