package client

import (
	"context"
	"time"

	"github.com/tuantran1810/go-di-template/libs/logger"
	"github.com/tuantran1810/go-di-template/uberfx/internal/models"
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
	return nil
}

func (c *FakeClient) Stop(_ context.Context) error {
	return nil
}

func (c *FakeClient) Send(ctx context.Context, msg *models.Message) error {
	latencyCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*time.Duration(c.FakeClientConfig.LatencyMs),
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
