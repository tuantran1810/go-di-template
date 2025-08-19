package server

import (
	"context"
	"fmt"
	"net"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpc_logger "github.com/tuantran1810/go-di-template/internal/inbound/server/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	server *grpc.Server
	lis    net.Listener
}

type grpcServerConfig struct {
	addr          string
	registerFunc  func(s *grpc.Server)
	unaryInts     []grpc.UnaryServerInterceptor
	streamInts    []grpc.StreamServerInterceptor
	serverOptions []grpc.ServerOption
	loggerOptions grpc_logger.Options
	reflection    bool
}

func newGRPCServer(conf *config) (*grpcServer, error) {
	if conf.grpc.registerFunc == nil {
		return nil, fmt.Errorf(".RegisterGRPCHandler() is required to be able to handle grpc messages")
	}
	if conf.grpc.addr == "" {
		conf.grpc.addr = "localhost:9090"
	}

	// unary interceptors
	unaryInts := make([]grpc.UnaryServerInterceptor, 0)
	if conf.logger != nil {
		unaryInts = append(unaryInts, grpc_logger.UnaryServerInterceptor(conf.logger, conf.grpc.loggerOptions))
	}
	unaryInts = append(unaryInts,
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor(),
		grpc_validator.UnaryServerInterceptor(),
	)
	unaryInts = append(unaryInts, conf.grpc.unaryInts...)

	// stream interceptors
	streamInts := make([]grpc.StreamServerInterceptor, 0)
	if conf.logger != nil {
		streamInts = append(streamInts, grpc_logger.StreamServerInterceptor(conf.logger))
	}
	streamInts = append(streamInts,
		grpc_prometheus.StreamServerInterceptor,
		grpc_recovery.StreamServerInterceptor(),
		grpc_validator.StreamServerInterceptor(),
	)

	// options
	options := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
		grpc.ChainStreamInterceptor(streamInts...),
	}
	options = append(options, conf.grpc.serverOptions...)

	// setup server
	server := grpc.NewServer(options...)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	grpc_prometheus.Register(server)
	grpc_prometheus.EnableHandlingTimeHistogram()
	conf.grpc.registerFunc(server)
	if conf.grpc.reflection {
		reflection.Register(server)
	}

	lis, err := net.Listen("tcp", conf.grpc.addr)
	if err != nil {
		return nil, fmt.Errorf("grpc tcp listen / %w", err)
	}

	s := &grpcServer{
		server: server,
		lis:    lis,
	}

	return s, nil
}

func (s *grpcServer) Serve() error {
	if err := s.server.Serve(s.lis); err != nil {
		return fmt.Errorf("grpc serve / %w", err)
	}

	return nil
}

func (s *grpcServer) Stop(ctx context.Context) {
	s.server.GracefulStop()
}
