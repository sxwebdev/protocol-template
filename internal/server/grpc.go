package server

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func (s *Server) newGRPC() error {
	s.grpc = grpc.NewServer()

	grpc_health_v1.RegisterHealthServer(s.grpc, health.NewServer())

	reflection.Register(s.grpc)

	listen, err := net.Listen("tcp", s.config.GrpcAddr)
	if err != nil {
		return err
	}

	s.logger.Infof("GRPC server start successfully on addr %s", s.config.GrpcAddr)

	if err := s.grpc.Serve(listen); err != nil {
		return err
	}

	return nil
}
