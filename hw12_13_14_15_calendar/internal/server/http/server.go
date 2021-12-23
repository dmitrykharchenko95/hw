package internalhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	logg *logrus.Logger
	srv  *http.Server
	addr string
}

func NewServer(logger *logrus.Logger, host, port string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	wrappedMux := loggingMiddleware(mux, logger)

	return &Server{
		logg: logger,
		srv: &http.Server{
			Addr:         net.JoinHostPort(host, port),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			Handler:      wrappedMux,
		},
		addr: net.JoinHostPort(host, port),
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

	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
