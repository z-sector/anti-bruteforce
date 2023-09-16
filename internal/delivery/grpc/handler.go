package grpc

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"

	"anti_bruteforce/internal/delivery/grpc/pb"
)

type UseCaseI interface {
	AddToBlackList(ctx context.Context, subnet string) error
	RemoveFromBlackList(ctx context.Context, subnet string) error
	AddToWhiteList(ctx context.Context, subnet string) error
	RemoveFromWhiteList(ctx context.Context, subnet string) error
	ClearLists(ctx context.Context) error
	CheckAuth(ctx context.Context, in *pb.AuthCheckRequest) (bool, error)
	Reset(ctx context.Context, in *pb.AuthCheckRequest) error
}

type HandlerGrpc struct {
	log zerolog.Logger
	uc  UseCaseI
	pb.UnimplementedAntiBruteForceServiceServer
}

func NewHandlerGrpc(log zerolog.Logger, uc UseCaseI) *HandlerGrpc {
	return &HandlerGrpc{log: log, uc: uc}
}

func (h *HandlerGrpc) AddToBlackList(ctx context.Context, in *pb.SubnetAddress) (*emptypb.Empty, error) {
	if err := h.uc.AddToBlackList(ctx, in.GetSubnetAddress()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *HandlerGrpc) RemoveFromBlackList(ctx context.Context, in *pb.SubnetAddress) (*emptypb.Empty, error) {
	if err := h.uc.RemoveFromBlackList(ctx, in.GetSubnetAddress()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *HandlerGrpc) AddToWhiteList(ctx context.Context, in *pb.SubnetAddress) (*emptypb.Empty, error) {
	if err := h.uc.AddToWhiteList(ctx, in.GetSubnetAddress()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *HandlerGrpc) RemoveFromWhiteList(ctx context.Context, in *pb.SubnetAddress) (*emptypb.Empty, error) {
	if err := h.uc.RemoveFromWhiteList(ctx, in.GetSubnetAddress()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *HandlerGrpc) ClearLists(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := h.uc.ClearLists(ctx); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *HandlerGrpc) AuthCheck(ctx context.Context, in *pb.AuthCheckRequest) (*pb.AuthCheckResponse, error) {
	ok, err := h.uc.CheckAuth(ctx, in)
	if err != nil {
		return nil, err
	}
	return &pb.AuthCheckResponse{Accepted: ok}, nil
}

func (h *HandlerGrpc) Reset(ctx context.Context, in *pb.AuthCheckRequest) (*emptypb.Empty, error) {
	if err := h.uc.Reset(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
