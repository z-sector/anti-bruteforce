package grpc

import (
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"anti_bruteforce/internal/delivery/grpc/pb"
	"anti_bruteforce/pkg/gserver"
)

func NewGrpcServer(logger zerolog.Logger, port int, uc UseCaseI) *gserver.AppGRPCServer {
	server := gserver.NewAppGrpcServer(grpc.NewServer(), gserver.Addr("", port))
	handler := NewHandlerGrpc(logger, uc)

	pb.RegisterAntiBruteForceServiceServer(server, handler)
	return server
}
