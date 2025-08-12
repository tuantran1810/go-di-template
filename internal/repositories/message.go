package repositories

import (
	"context"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/internal/repositories/mysql"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Key   string
	Value string
}

type MessageRepository struct {
	*mysql.GenericRepository[Message, entities.Message]
}

func NewMessageRepository(repository *mysql.Repository) *MessageRepository {
	transformer := entities.NewBaseExtendedTransformer[Message, entities.Message]()
	return &MessageRepository{
		GenericRepository: mysql.NewGenericRepository(repository, transformer),
	}
}

func (s *MessageRepository) Start(ctx context.Context) error {
	log.Info("starting message store")

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	err := s.AutoMigrate(timeoutCtx)
	if err != nil {
		return err
	}

	return s.Ping(timeoutCtx)
}

func (s *MessageRepository) Stop(_ context.Context) error {
	log.Info("stopping message store")
	return nil
}
