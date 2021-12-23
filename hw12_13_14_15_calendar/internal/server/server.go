package server

import (
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
