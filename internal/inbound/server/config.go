package server

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc_logger "github.com/tuantran1810/go-di-template/internal/inbound/server/logger"
	"google.golang.org/grpc"
)

type config struct {
	logger grpc_logger.Logger
	grpc   grpcServerConfig
	http   httpServerConfig
}

func NewServerConfig() *config {
	conf := &config{
		grpc: grpcServerConfig{
			addr:          "localhost:9090",
			unaryInts:     make([]grpc.UnaryServerInterceptor, 0),
			serverOptions: make([]grpc.ServerOption, 0),
			registerFunc:  nil,
			reflection:    true,
			loggerOptions: grpc_logger.DefaultOptions,
		},
		http: httpServerConfig{
			addr:         "localhost:8080",
			registerFunc: nil,
			middlewares:  make([]func(http.Handler) http.Handler, 0),
		},
	}

	return conf
}

// SetLogger will enable logger unary interceptor
func (c *config) SetLogger(logger grpc_logger.Logger) *config {
	c.logger = logger
	return c
}

func (c *config) SetGRPCLoggerOptions(options grpc_logger.Options) *config {
	c.grpc.loggerOptions = options
	return c
}

func (c *config) SetGRPCAddr(addr string) *config {
	c.grpc.addr = addr
	return c
}

func (c *config) SetHTTPAddr(addr string) *config {
	c.http.addr = addr
	return c
}

func (c *config) SetGRPCReflection(enable bool) *config {
	c.grpc.reflection = enable
	return c
}

// AddInterceptor is deprecated, use ChainUnaryInterceptor instead.
func (c *config) AddInterceptor(interceptors ...grpc.UnaryServerInterceptor) *config {
	return c.ChainUnaryInterceptors(interceptors...)
}

func (c *config) ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) *config {
	c.grpc.unaryInts = append(c.grpc.unaryInts, interceptors...)
	return c
}

func (c *config) ChainStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) *config {
	c.grpc.streamInts = append(c.grpc.streamInts, interceptors...)
	return c
}

func (c *config) AddGRPCServerOptions(options ...grpc.ServerOption) *config {
	c.grpc.serverOptions = append(c.grpc.serverOptions, options...)
	return c
}

func (c *config) AddHTTPServerOptions(options ...runtime.ServeMuxOption) *config {
	c.http.serveMuxOptions = append(c.http.serveMuxOptions, options...)
	return c
}

func (c *config) RegisterGRPC(handler func(s *grpc.Server)) *config {
	c.grpc.registerFunc = func(s *grpc.Server) {
		handler(s)
	}
	return c
}

func (c *config) RegisterHTTP(handler func(mux *runtime.ServeMux, conn *grpc.ClientConn)) *config {
	c.http.registerFunc = func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
		handler(mux, conn)
	}
	return c
}

func (c *config) AddMiddleware(middlewares ...func(http.Handler) http.Handler) *config {
	c.http.middlewares = append(c.http.middlewares, middlewares...)
	return c
}
