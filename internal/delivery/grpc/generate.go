package grpc

//go:generate protoc --proto_path=./../../../ --go_out=. --go-grpc_out=. ./../../../api/proto/api.proto
//go:generate mockgen -source=handler.go -destination=mocks/handler.go
