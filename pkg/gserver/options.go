package gserver

import (
	"net"
	"strconv"
)

type Option func(*AppGRPCServer)

func Addr(host string, port int) Option {
	return func(s *AppGRPCServer) {
		s.addr = net.JoinHostPort(host, strconv.Itoa(port))
	}
}
