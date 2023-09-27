package gserver

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	defaultAddr            = ":8080"
	defaultShutdownTimeout = 20 * time.Second
)

type AppGRPCServer struct {
	origin          *grpc.Server
	notify          chan error
	shutdownTimeout time.Duration
	addr            string
}

func NewAppGrpcServer(origin *grpc.Server, opts ...Option) *AppGRPCServer {
	s := &AppGRPCServer{
		origin:          origin,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
		addr:            defaultAddr,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *AppGRPCServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.origin.RegisterService(desc, impl)
}

func (s *AppGRPCServer) GetAddr() string {
	return s.addr
}

func (s *AppGRPCServer) Run() {
	go func() {
		defer close(s.notify)
		lis, err := net.Listen("tcp", s.addr)
		if err != nil {
			s.notify <- err
			return
		}
		s.notify <- s.origin.Serve(lis)
	}()
}

func (s *AppGRPCServer) Notify() <-chan error {
	return s.notify
}

func (s *AppGRPCServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	ok := make(chan struct{})
	go func() {
		s.origin.GracefulStop()
		close(ok)
	}()

	select {
	case <-ok:
		return nil
	case <-ctx.Done():
		s.origin.Stop()
		return ctx.Err()
	}
}
