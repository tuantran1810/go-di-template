package server

import (
	"context"
	"fmt"
)

type Server struct {
	grpc *grpcServer
	http *httpServer
}

func NewServer(conf *config) (*Server, error) {
	if conf == nil {
		return nil, fmt.Errorf("config is required")
	}

	grpcServer, err := newGRPCServer(conf)
	if err != nil {
		return nil, fmt.Errorf("new grpc server / %w", err)
	}

	httpServer, err := newHTTPServer(conf)
	if err != nil {
		return nil, fmt.Errorf("new http server / %w", err)
	}

	g := &Server{
		grpcServer,
		httpServer,
	}

	return g, nil
}

func (g *Server) Start() error {
	ch := make(chan error)
	defer close(ch)

	go func() {
		ch <- g.grpc.Serve()
	}()
	go func() {
		ch <- g.http.Serve()
	}()

	for i := 0; i < 2; i++ {
		if err := <-ch; err != nil {
			_ = g.Stop(context.Background())
			continue
		}
	}

	return nil
}

func (g *Server) StartBackground(ctx context.Context) error {
	go func() {
		g.Start()
	}()

	return nil
}

func (g *Server) Stop(ctx context.Context) error {
	if err := g.http.Stop(ctx); err != nil {
		return err
	}

	g.grpc.Stop(ctx)
	return nil
}
