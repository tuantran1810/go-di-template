package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tuantran1810/go-di-template/uberfx/config"
	"github.com/tuantran1810/go-di-template/uberfx/internal/client"
	"github.com/tuantran1810/go-di-template/uberfx/internal/consumer"
	"github.com/tuantran1810/go-di-template/uberfx/internal/processor"
	"github.com/tuantran1810/go-di-template/uberfx/internal/stores"
	"go.uber.org/fx"
)

var (
	_ processor.DataRepository = &stores.Repository{}
	_ processor.DataWriter     = &stores.MessageStore{}
	_ processor.Sender         = &client.FakeClient{}
	_ consumer.Processor       = &processor.Processor{}
)

func newRepository(
	appLifecycle fx.Lifecycle,
	config stores.RepositoryConfig,
) *stores.Repository {
	r := stores.MustNewRepository(config)
	appLifecycle.Append(fx.Hook{
		OnStart: r.Start,
		OnStop:  r.Stop,
	})
	return r
}

func newMessageStore(
	appLifecycle fx.Lifecycle,
	repository *stores.Repository,
) *stores.MessageStore {
	s := stores.NewMessageStore(repository)
	appLifecycle.Append(fx.Hook{
		OnStart: s.Start,
		OnStop:  s.Stop,
	})
	return s
}

func newFakeClient(
	appLifecycle fx.Lifecycle,
	config client.FakeClientConfig,
) *client.FakeClient {
	c := client.NewFakeClient(config)
	appLifecycle.Append(fx.Hook{
		OnStart: c.Start,
		OnStop:  c.Stop,
	})
	return c
}

func newFakeConsumer(
	appLifecycle fx.Lifecycle,
	config consumer.FakeConsumerConfig,
	processor consumer.Processor,
) *consumer.FakeConsumer {
	c := consumer.NewFakeConsumer(config, processor)
	appLifecycle.Append(fx.Hook{
		OnStart: c.Start,
		OnStop:  c.Stop,
	})
	return c
}

func newProcessor(
	appLifecycle fx.Lifecycle,
	client processor.Sender,
	writer processor.DataWriter,
	repo processor.DataRepository,
) *processor.Processor {
	p := processor.NewProcessor(client, writer, repo)
	appLifecycle.Append(fx.Hook{
		OnStart: p.Start,
		OnStop:  p.Stop,
	})
	return p
}

func newServerApp() *fx.App {
	cfg := config.MustLoadConfig[config.ServerConfig]()
	log.Infof("Starting server with config: %+v", cfg)

	return fx.New(
		fx.StartTimeout(fx.DefaultTimeout),
		fx.StopTimeout(fx.DefaultTimeout),
		fx.Supply(
			cfg,
			client.FakeClientConfig{
				LatencyMs: cfg.FakeClient.LatencyMs,
			},
			consumer.FakeConsumerConfig{
				PerMs: cfg.FakeConsumer.PerMs,
			},
			stores.RepositoryConfig{
				DatabasePath: cfg.Sqlite.DatabasePath,
			},
		),
		fx.Provide(
			newRepository,
			fx.Annotate(
				func(r *stores.Repository) *stores.Repository {
					return r
				},
				fx.As(new(processor.DataRepository)),
			),
			fx.Annotate(
				newMessageStore,
				fx.As(new(processor.DataWriter)),
			),
			fx.Annotate(
				newFakeClient,
				fx.As(new(processor.Sender)),
			),
			fx.Annotate(
				newProcessor,
				fx.As(new(consumer.Processor)),
			),
		),
		fx.Invoke(newFakeConsumer),
	)
}

func startServer(_ *cobra.Command, _ []string) {
	newServerApp().Run()
}
