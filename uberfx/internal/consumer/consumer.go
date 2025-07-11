package consumer

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/tuantran1810/go-di-template/libs/logger"
	"github.com/tuantran1810/go-di-template/uberfx/internal/models"
)

var log = logger.MustNamedLogger("consumer")

type Processor interface {
	Process(ctx context.Context, msg *models.Message) error
}

type FakeConsumerConfig struct {
	PerMs uint
}

type FakeConsumer struct {
	FakeConsumerConfig
	processor Processor
	cancel    context.CancelFunc
}

func NewFakeConsumer(config FakeConsumerConfig, processor Processor) *FakeConsumer {
	return &FakeConsumer{
		FakeConsumerConfig: config,
		processor:          processor,
	}
}

func (c *FakeConsumer) generateContentWorker(ctx context.Context) {
	timer := time.NewTicker(time.Millisecond * time.Duration(c.FakeConsumerConfig.PerMs))
	for {
		select {
		case <-timer.C:
			log.Infof("message generated")
			msg := &models.Message{
				Key:   gofakeit.UUID(),
				Value: gofakeit.Sentence(20),
			}
			if err := c.processor.Process(ctx, msg); err != nil {
				log.Errorf("failed to process message: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *FakeConsumer) Start(_ context.Context) error {
	cancelCtx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	go c.generateContentWorker(cancelCtx)
	return nil
}

func (c *FakeConsumer) Stop(_ context.Context) error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}
