package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	grpc_health "google.golang.org/grpc/health/grpc_health_v1"
)

type HttpServer struct {
	server *http.Server
	conn   *grpc.ClientConn
}

type httpServerConfig struct {
	addr            string
	registerFunc    func(mux *runtime.ServeMux, conn *grpc.ClientConn)
	serveMuxOptions []runtime.ServeMuxOption
	middlewares     []func(http.Handler) http.Handler
}

func NewHTTPServer(conf *config) (*HttpServer, error) {
	if conf.http.addr == "" {
		conf.http.addr = "localhost:8080"
	}

	conn, err := grpc.Dial(conf.grpc.addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial grpc / %w", err)
	}

	options := []runtime.ServeMuxOption{
		runtime.WithHealthEndpointAt(grpc_health.NewHealthClient(conn), "/health"),
	}
	options = append(options, conf.http.serveMuxOptions...)
	options = append(options, runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONBuiltin{}))
	mux := runtime.NewServeMux(options...)

	promHandler := promhttp.Handler()
	mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		promHandler.ServeHTTP(w, r)
	})

	handler := http.Handler(mux)

	if len(conf.http.middlewares) > 0 {
		for _, middleware := range conf.http.middlewares {
			handler = middleware(handler)
		}
	}

	if conf.http.registerFunc != nil {
		conf.http.registerFunc(mux, conn)
	}

	s := &HttpServer{
		server: &http.Server{
			Addr:    conf.http.addr,
			Handler: handler,
		},
		conn: conn,
	}

	return s, nil
}

func (s *HttpServer) Serve() error {
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http listen and serve / %w", err)
	}

	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	defer s.conn.Close()
	return s.server.Shutdown(ctx)
}
