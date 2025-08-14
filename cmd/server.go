package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/tuantran1810/go-di-template/config"
	"github.com/tuantran1810/go-di-template/internal/clients"
	"github.com/tuantran1810/go-di-template/internal/consumers"
	"github.com/tuantran1810/go-di-template/internal/controllers"
	"github.com/tuantran1810/go-di-template/internal/repositories"
	"github.com/tuantran1810/go-di-template/internal/repositories/mysql"
	"github.com/tuantran1810/go-di-template/internal/usecases"
	"github.com/tuantran1810/go-di-template/libs/middlewares"
	pb "github.com/tuantran1810/go-di-template/pkg/go_di_template/v1"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	_ usecases.IRepository              = &mysql.Repository{}
	_ usecases.IUserRepository          = &repositories.UserRepository{}
	_ usecases.IUserAttributeRepository = &repositories.UserAttributeRepository{}
	_ usecases.IMessageRepository       = &repositories.MessageRepository{}
	_ controllers.IUserUsecase          = &usecases.Users{}
	_ controllers.ILoggingWorker        = &usecases.LoggingWorker{}
)

func newRepository(
	appLifecycle fx.Lifecycle,
	config mysql.RepositoryConfig,
) *mysql.Repository {
	r := mysql.MustNewRepository(config)
	appLifecycle.Append(fx.Hook{
		OnStart: r.Start,
		OnStop:  r.Stop,
	})
	return r
}

func newUserRepository(
	appLifecycle fx.Lifecycle,
	repository *mysql.Repository,
) *repositories.UserRepository {
	s := repositories.NewUserRepository(repository)
	appLifecycle.Append(fx.Hook{
		OnStart: s.Start,
		OnStop:  s.Stop,
	})
	return s
}

func newUserAttributeRepository(
	appLifecycle fx.Lifecycle,
	repository *mysql.Repository,
) *repositories.UserAttributeRepository {
	s := repositories.NewUserAttributeRepository(repository)
	appLifecycle.Append(fx.Hook{
		OnStart: s.Start,
		OnStop:  s.Stop,
	})
	return s
}

func newMessageRepository(
	appLifecycle fx.Lifecycle,
	repository *mysql.Repository,
) *repositories.MessageRepository {
	s := repositories.NewMessageRepository(repository)
	appLifecycle.Append(fx.Hook{
		OnStart: s.Start,
		OnStop:  s.Stop,
	})
	return s
}

func newUsersUsecase(
	repository *mysql.Repository,
	userRepository *repositories.UserRepository,
	userAttributeRepository *repositories.UserAttributeRepository,
) *usecases.Users {
	return usecases.NewUsersUsecase(repository, userRepository, userAttributeRepository)
}

func newLoggingWorker(
	cfg usecases.LoggingWorkerConfig,
	appLifecycle fx.Lifecycle,
	messageRepository *repositories.MessageRepository,
	client *clients.FakeClient,
) *usecases.LoggingWorker {
	w := usecases.NewLoggingWorker(cfg, messageRepository, client)
	appLifecycle.Append(fx.Hook{
		OnStart: w.Start,
		OnStop:  w.Stop,
	})
	return w
}

func newController(
	usecase *usecases.Users,
	loggingWorker *usecases.LoggingWorker,
) *controllers.UserController {
	return controllers.NewUserController(usecase, loggingWorker)
}

func newFakeClient(
	appLifecycle fx.Lifecycle,
	config clients.FakeClientConfig,
) *clients.FakeClient {
	c := clients.NewFakeClient(config)
	appLifecycle.Append(fx.Hook{
		OnStart: c.Start,
		OnStop:  c.Stop,
	})
	return c
}

func newFakeConsumer(
	appLifecycle fx.Lifecycle,
	config config.ConsumerConfig,
	loggingWorker *usecases.LoggingWorker,
) *consumers.FakeConsumer {
	c := consumers.NewFakeConsumer(
		consumers.FakeConsumerConfig{
			PerMs: config.PerMs,
		},
		loggingWorker,
	)

	appLifecycle.Append(fx.Hook{
		OnStart: c.Start,
		OnStop:  c.Stop,
	})

	return c
}

func startGrpcServer(
	appLifecycle fx.Lifecycle,
	cfg config.ServerConfig,
	userController *controllers.UserController,
) *grpc.Server {
	grpcHost := fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort)

	listen, err := net.Listen("tcp", grpcHost)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middlewares.HandleErrorCodes,
		),
	)
	pb.RegisterUserServiceServer(server, userController)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	reflection.Register(server)

	appLifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Infof("Starting grpc server on %s", grpcHost)
			go func() {
				if err := server.Serve(listen); err != nil {
					log.Fatalln("Failed to serve grpc:", err)
				}
			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Infof("Stopping grpc server on %s", grpcHost)
			server.GracefulStop()

			return nil
		},
	})

	return server
}

func startHttpServer(
	appLifecycle fx.Lifecycle,
	cfg config.ServerConfig,
	userController *controllers.UserController,
) *http.Server {
	httpHost := fmt.Sprintf("0.0.0.0:%d", cfg.HttpPort)

	gwmux := runtime.NewServeMux(
		runtime.WithErrorHandler(
			func(
				ctx context.Context,
				mux *runtime.ServeMux,
				marshaler runtime.Marshaler,
				w http.ResponseWriter,
				r *http.Request,
				err error,
			) {
				runtime.DefaultHTTPErrorHandler(
					ctx, mux, marshaler, w, r,
					middlewares.HandleError(err),
				)
			},
		),
	)
	if err := pb.RegisterUserServiceHandlerServer(
		globalContext, gwmux, userController,
	); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:        httpHost,
		ReadTimeout: cfg.HttpServerReadTimeout,
		Handler:     gwmux,
	}

	appLifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Infof("Starting http server on %s", httpHost)
			go func() {
				if err := gwServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalln("Failed to serve http:", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Infof("Stopping http server on %s", httpHost)
			if err := gwServer.Shutdown(ctx); err != nil {
				log.Fatalln("Failed to shutdown http server:", err)
			}

			return nil
		},
	})

	return gwServer
}

func newServerApp() *fx.App {
	cfg := config.MustLoadConfig[config.ServerConfig]()
	log.Infof("Starting server with config: %+v", cfg)

	return fx.New(
		fx.StartTimeout(fx.DefaultTimeout),
		fx.StopTimeout(fx.DefaultTimeout),
		fx.Supply(
			cfg,
			mysql.RepositoryConfig{
				Username:  cfg.MySql.Username,
				Password:  cfg.MySql.Password,
				Protocol:  cfg.MySql.Protocol,
				Address:   cfg.MySql.Address,
				Database:  cfg.MySql.Database,
				ParseTime: true,
			},
			usecases.LoggingWorkerConfig{
				BufferCapacity: cfg.LoggingWorker.BufferCapacity,
				FlushInterval:  cfg.LoggingWorker.FlushInterval,
			},
			config.ConsumerConfig{
				PerMs: cfg.Consumer.PerMs,
			},
			clients.FakeClientConfig{
				LatencyMs: cfg.Client.LatencyMs,
			},
		),
		fx.Provide(
			newRepository,
			newMessageRepository,
			newUserRepository,
			newUserAttributeRepository,
			newUsersUsecase,
			newFakeClient,
			newLoggingWorker,
			newController,
		),
		fx.Invoke(startGrpcServer),
		fx.Invoke(startHttpServer),
		fx.Invoke(newFakeConsumer),
	)
}

func startServer(_ *cobra.Command, _ []string) {
	newServerApp().Run()
}
