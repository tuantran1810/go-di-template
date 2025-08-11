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
	"github.com/tuantran1810/go-di-template/internal/controllers"
	"github.com/tuantran1810/go-di-template/internal/stores"
	"github.com/tuantran1810/go-di-template/internal/stores/mysql"
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
	_ usecases.IRepository         = &mysql.Repository{}
	_ usecases.IUserStore          = &stores.UserStore{}
	_ usecases.IUserAttributeStore = &stores.UserAttributeStore{}
	_ controllers.IUserUsecase     = &usecases.Users{}
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

func newUserStore(
	appLifecycle fx.Lifecycle,
	repository *mysql.Repository,
) *stores.UserStore {
	s := stores.NewUserStore(repository)
	appLifecycle.Append(fx.Hook{
		OnStart: s.Start,
		OnStop:  s.Stop,
	})
	return s
}

func newUserAttributeStore(
	appLifecycle fx.Lifecycle,
	repository *mysql.Repository,
) *stores.UserAttributeStore {
	s := stores.NewUserAttributeStore(repository)
	appLifecycle.Append(fx.Hook{
		OnStart: s.Start,
		OnStop:  s.Stop,
	})
	return s
}

func newUsersUsecase(
	repository usecases.IRepository,
	userStore usecases.IUserStore,
	userAttributeStore usecases.IUserAttributeStore,
) *usecases.Users {
	return usecases.NewUsersUsecase(repository, userStore, userAttributeStore)
}

func newController(
	usecase controllers.IUserUsecase,
) *controllers.UserController {
	return controllers.NewUserController(usecase)
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
		),
		fx.Provide(
			newRepository,
			func(r *mysql.Repository) usecases.IRepository {
				return r
			},
			fx.Annotate(
				newUserStore,
				fx.As(new(usecases.IUserStore)),
			),
			fx.Annotate(
				newUserAttributeStore,
				fx.As(new(usecases.IUserAttributeStore)),
			),
			fx.Annotate(
				newUsersUsecase,
				fx.As(new(controllers.IUserUsecase)),
			),
			newController,
		),
		fx.Invoke(startGrpcServer),
		fx.Invoke(startHttpServer),
	)
}

func startServer(_ *cobra.Command, _ []string) {
	newServerApp().Run()
}
