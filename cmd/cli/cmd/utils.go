package cmd

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"anti_bruteforce/internal/delivery/grpc/pb"
)

func getGRPCClient(host string) (pb.AntiBruteForceServiceClient, error) {
	conn, err := grpc.Dial(
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return pb.NewAntiBruteForceServiceClient(conn), nil
}
