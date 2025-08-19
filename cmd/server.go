package cmd

import (
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/tuantran1810/go-di-template/config"
	"github.com/tuantran1810/go-di-template/internal/controllers"
	"github.com/tuantran1810/go-di-template/internal/inbound"
	"github.com/tuantran1810/go-di-template/internal/inbound/server"
	"github.com/tuantran1810/go-di-template/internal/outbound"
	"github.com/tuantran1810/go-di-template/internal/repositories"
	"github.com/tuantran1810/go-di-template/internal/repositories/mysql"
	"github.com/tuantran1810/go-di-template/internal/usecases"
	"github.com/tuantran1810/go-di-template/libs/middlewares"
	pb "github.com/tuantran1810/go-di-template/pkg/go_di_template/v1"
	"go.uber.org/fx"
	"google.golang.org/grpc"
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
	client *outbound.FakeClient,
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
	config outbound.FakeClientConfig,
) *outbound.FakeClient {
	c := outbound.NewFakeClient(config)
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
) *inbound.FakeConsumer {
	c := inbound.NewFakeConsumer(
		inbound.FakeConsumerConfig{
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

func startInboundServer(
	appLifecycle fx.Lifecycle,
	cfg config.ServerConfig,
	userController *controllers.UserController,
) *server.Server {
	serverConfig := server.NewServerConfig().
		SetLogger(log).
		SetGRPCAddr(fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort)).
		SetHTTPAddr(fmt.Sprintf("0.0.0.0:%d", cfg.HttpPort)).
		SetGRPCReflection(true).
		RegisterGRPC(func(s *grpc.Server) {
			pb.RegisterUserServiceServer(s, userController)
		}).
		RegisterHTTP(func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
			if err := pb.RegisterUserServiceHandlerServer(globalContext, mux, userController); err != nil {
				log.Fatalln("Failed to register server:", err)
			}
		}).AddInterceptor(
		middlewares.HandleErrorCodes,
	)

	server, err := server.NewServer(serverConfig)
	if err != nil {
		log.Fatalln("Failed to create server:", err)
	}

	appLifecycle.Append(fx.Hook{
		OnStart: server.StartBackground,
		OnStop:  server.Stop,
	})

	return server
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
			outbound.FakeClientConfig{
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
		fx.Invoke(startInboundServer),
		fx.Invoke(newFakeConsumer),
	)
}

func startServer(_ *cobra.Command, _ []string) {
	newServerApp().Run()
}
