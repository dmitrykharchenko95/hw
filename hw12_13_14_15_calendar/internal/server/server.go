package server

import (
	"context"
	"sync"

	internalgrpc "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/server/http"
)

type Server struct {
	HTTP *internalhttp.Server
	GRPC *internalgrpc.Server
}

func NewServer(httpServer *internalhttp.Server, grpcServer *internalgrpc.Server) *Server {
	return &Server{
		HTTP: httpServer,
		GRPC: grpcServer,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var err error
	var wg sync.WaitGroup

	wg.Add(2)

	go func(wg *sync.WaitGroup) error {
		defer wg.Done()
		err = s.HTTP.Start(ctx)
		if err != nil {
			return err
		}
		return nil
	}(&wg)

	go func(wg *sync.WaitGroup) error {
		defer wg.Done()
		err = s.GRPC.Start(ctx)
		if err != nil {
			return err
		}
		return nil
	}(&wg)

	wg.Wait()
	return err
}
