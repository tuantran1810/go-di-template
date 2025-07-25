package stores

import (
	"context"

	"github.com/tuantran1810/go-di-template/uberfx/internal/models"
	"github.com/tuantran1810/go-di-template/uberfx/internal/stores/mysql"
)

type MessageStore struct {
	*mysql.GenericStore[models.Message]
}

func NewMessageStore(repository *mysql.Repository) *MessageStore {
	return &MessageStore{
		GenericStore: mysql.NewGenericStore[models.Message](repository),
	}
}

func (s *MessageStore) Start(ctx context.Context) error {
	log.Info("starting message store")

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	err := s.GenericStore.AutoMigrate(timeoutCtx)
	if err != nil {
		return err
	}

	return s.Ping(timeoutCtx)
}

func (s *MessageStore) Stop(_ context.Context) error {
	log.Info("stopping message store")
	return nil
}
