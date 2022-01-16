package internalgrpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	logg  *logrus.Logger
	store *storage.Storage
	srv   *grpc.Server
	addr  string
	*pb.UnimplementedStorageServer
}

func NewSever(logg *logrus.Logger, store *storage.Storage, host, port string) *Server {
	return &Server{
		logg:  logg,
		store: store,
		srv: grpc.NewServer(
			grpc.ChainUnaryInterceptor(UnaryServerRequestLoggerInterceptor(logg)),
		),
		addr: net.JoinHostPort(host, port),
	}
}

func (s *Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	pb.RegisterStorageServer(s.srv, s)
	s.logg.Infof("Start grpc server on %s...", s.addr)
	if err := s.srv.Serve(lsn); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("Stop grpc server...")
	s.srv.Stop()
	return nil
}

func (s *Server) CreateEvent(ctx context.Context, e *pb.Event) (*pb.Response, error) {
	event := storage.Event{
		ID:        e.Id,
		Title:     e.Title,
		StartDate: e.StartDate.AsTime(),
		EndDate:   e.EndDate.AsTime(),
		Content:   e.Content,
		UserID:    e.UserId,
		SendTime:  e.SendTime.AsDuration(),
	}
	res := pb.Response{}

	err := (*s.store).CreateEvent(ctx, event)
	if err != nil {
		s.logg.Errorf("can not create event: %v", err)
		res.Error = fmt.Sprintf("can not create event: %v", err)
		return &res, err
	}

	res.Info = "Event created"
	return &res, nil
}

func (s *Server) UpdateEvent(ctx context.Context, e *pb.Event) (*pb.Response, error) {
	event := storage.Event{
		ID:        e.Id,
		Title:     e.Title,
		StartDate: e.StartDate.AsTime(),
		EndDate:   e.EndDate.AsTime(),
		Content:   e.Content,
		UserID:    e.UserId,
		SendTime:  e.SendTime.AsDuration(),
	}
	res := pb.Response{}

	err := (*s.store).UpdateEvent(ctx, event)
	if err != nil {
		s.logg.Errorf("can not update event: %v", err)
		res.Error = fmt.Sprintf("can not update event: %v", err)
		return &res, err
	}

	res.Info = "Event updated"
	return &res, nil
}

func (s *Server) DeleteEvent(ctx context.Context, id *pb.DeleteEventReq) (*pb.Response, error) {
	res := pb.Response{}

	err := (*s.store).DeleteEvent(ctx, id.GetId())
	if err != nil {
		s.logg.Errorf("can not delete event: %v", err)
		res.Error = fmt.Sprintf("can not delete event: %v", err)
		return &res, err
	}

	res.Info = "Event deleted"
	return &res, nil
}

func (s *Server) ListEventsForDay(ctx context.Context, t *timestamp.Timestamp) (*pb.Response, error) {
	return s.listEventsForPeriod(ctx, t, "day")
}

func (s *Server) ListEventsForWeek(ctx context.Context, t *timestamp.Timestamp) (*pb.Response, error) {
	return s.listEventsForPeriod(ctx, t, "week")
}

func (s *Server) ListEventsForMonth(ctx context.Context, t *timestamp.Timestamp) (*pb.Response, error) {
	return s.listEventsForPeriod(ctx, t, "month")
}

func (s *Server) listEventsForPeriod(ctx context.Context, t *timestamp.Timestamp, period string) (*pb.Response, error) {
	res := pb.Response{}
	events := make([]storage.Event, 0)
	var err error

	switch period {
	case "day":
		events, err = (*s.store).ListEventsForDay(ctx, t.AsTime())
	case "week":
		events, err = (*s.store).ListEventsForWeek(ctx, t.AsTime())
	case "month":
		events, err = (*s.store).ListEventsForMonth(ctx, t.AsTime())
	}

	if err != nil {
		s.logg.Errorf("can not list events: %v", err)
		res.Error = fmt.Sprintf("can not list events: %v", err)
		return &res, err
	}

	if len(events) > 0 {
		for _, e := range events {
			res.Events = append(res.Events, parseEvents(e))
		}
	}

	return &res, nil
}

func parseEvents(in storage.Event) *pb.Event {
	return &pb.Event{
		Id:    in.ID,
		Title: in.Title,
		StartDate: &timestamp.Timestamp{
			Seconds: in.StartDate.Unix(),
			Nanos:   int32(in.StartDate.Nanosecond()),
		},
		EndDate: &timestamp.Timestamp{
			Seconds: in.EndDate.Unix(),
			Nanos:   int32(in.EndDate.Nanosecond()),
		},
		Content: in.Content,
		UserId:  in.UserID,
		SendTime: &duration.Duration{
			Seconds: int64(in.SendTime.Seconds()),
			Nanos:   int32(in.SendTime.Nanoseconds()),
		},
	}
}
