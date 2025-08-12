package consumer

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/libs/logger"
)

var log = logger.MustNamedLogger("consumer")

type ILoggingWorker interface {
	LogAndSend(ctx context.Context, msg entities.Message) error
}

type FakeConsumerConfig struct {
	PerMs uint
}

type FakeConsumer struct {
	FakeConsumerConfig
	loggingWorker ILoggingWorker
	cancel        context.CancelFunc
}

func NewFakeConsumer(config FakeConsumerConfig, loggingWorker ILoggingWorker) *FakeConsumer {
	return &FakeConsumer{
		FakeConsumerConfig: config,
		loggingWorker:      loggingWorker,
	}
}

func (c *FakeConsumer) generateContentWorker(ctx context.Context) {
	timer := time.NewTicker(time.Millisecond * time.Duration(c.FakeConsumerConfig.PerMs))
	for {
		select {
		case <-timer.C:
			log.Infof("message generated")
			if err := c.loggingWorker.LogAndSend(
				ctx,
				entities.Message{
					Key:   "consumer_messages",
					Value: gofakeit.Sentence(20),
				},
			); err != nil {
				log.Error("failed to log and send message", "error", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *FakeConsumer) Start(_ context.Context) error {
	log.Info("starting fake consumer")
	cancelCtx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	go c.generateContentWorker(cancelCtx)
	return nil
}

func (c *FakeConsumer) Stop(_ context.Context) error {
	log.Info("stopping fake consumer")
	c.cancel()
	return nil
}
