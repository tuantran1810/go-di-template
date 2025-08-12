package stores

import (
	"context"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/internal/stores/mysql"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Key   string
	Value string
}

type MessageStore struct {
	*mysql.GenericStore[Message, entities.Message]
}

func NewMessageStore(repository *mysql.Repository) *MessageStore {
	transformer := entities.NewBaseExtendedTransformer[Message, entities.Message]()
	return &MessageStore{
		GenericStore: mysql.NewGenericStore(repository, transformer),
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
