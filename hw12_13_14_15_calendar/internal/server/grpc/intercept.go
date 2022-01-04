package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func UnaryServerRequestLoggerInterceptor(logg *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)
		if err != nil {
			logg.Errorf("UnaryServerRequestLoggerInterceptor: can not handle request: %v", err)
		}

		ip, err := getClientIP(ctx)
		if err != nil {
			logg.Errorf("can not get client IP: %v", err)
		}

		md, _ := metadata.FromIncomingContext(ctx)

		var dataLen int

		response, ok := resp.(*pb.Response)
		if ok {
			dataLen = len(response.String())
		}

		logg.Infof("%v %v %v %vbytes %vms %v",
			ip, info.FullMethod, md.Get(":authority"), dataLen, time.Since(start).Microseconds(), md.Get("user-agent"))

		return resp, err
	}
}

func getClientIP(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("couldn't parse client IP address")
	}
	return p.Addr.String(), nil
}
