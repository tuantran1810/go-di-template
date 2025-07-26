package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/tuantran1810/go-di-template/libs/logger"
	"go.uber.org/zap"
)

var (
	log           *zap.SugaredLogger = logger.MustNamedLogger("cmd")
	globalContext                    = context.Background()

	RootCmd = &cobra.Command{
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				log.Fatalf("Failed to run command: %v", err)
			}
		},
	}
)

func init() {
	startServerCmd := &cobra.Command{
		Use:   "start-server",
		Short: "Starts the server",
		Long:  `Starts the server`,
		Run:   startServer,
	}

	startConsumerCmd := &cobra.Command{
		Use:   "start-consumer",
		Short: "Starts the consumer",
		Long:  `Starts the consumer`,
		Run:   startConsumer,
	}

	RootCmd.AddCommand(startServerCmd)
	RootCmd.AddCommand(startConsumerCmd)
}
