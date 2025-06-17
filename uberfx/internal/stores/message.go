package stores

import (
	"context"

	"github.com/tuantran1810/go-di-template/uberfx/internal/models"
)

type MessageStore struct {
	*GenericStore[models.Message]
}

func NewMessageStore(repository *Repository) *MessageStore {
	return &MessageStore{
		GenericStore: NewGenericStore[models.Message](repository),
	}
}

func (s *MessageStore) Start(ctx context.Context) error {
	log.Infof("starting message store")

	s.repository.mutex.Lock()
	err := s.repository.db.AutoMigrate(&models.Message{})
	s.repository.mutex.Unlock()
	if err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return s.Ping(timeoutCtx)
}

func (s *MessageStore) Stop(_ context.Context) error {
	return nil
}
