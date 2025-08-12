package clients

import (
	"context"
	"time"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/libs/logger"
)

var log = logger.MustNamedLogger("client")

type FakeClientConfig struct {
	LatencyMs uint
}

type FakeClient struct {
	FakeClientConfig
}

func NewFakeClient(config FakeClientConfig) *FakeClient {
	return &FakeClient{
		FakeClientConfig: config,
	}
}

func (c *FakeClient) Start(_ context.Context) error {
	log.Info("starting fake client")
	return nil
}

func (c *FakeClient) Stop(_ context.Context) error {
	log.Info("stopping fake client")
	return nil
}

func (c *FakeClient) Send(ctx context.Context, msg *entities.Message) error {
	latencyCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*time.Duration(c.LatencyMs),
	)
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-latencyCtx.Done():
		log.Infof("message sent: %s", msg.Key)
		return nil
	}
}
