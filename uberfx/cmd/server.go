package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tuantran1810/go-di-template/uberfx/config"
	"go.uber.org/fx"
)

func newServerApp() *fx.App {
	cfg := config.MustLoadConfig[config.ServerConfig]()
	log.Infof("Starting server with config: %+v", cfg)

	return fx.New(
		fx.StartTimeout(fx.DefaultTimeout),
		fx.StopTimeout(fx.DefaultTimeout),
		fx.Supply(cfg),
	)
}

func startServer(_ *cobra.Command, _ []string) {
	newServerApp().Run()
}
