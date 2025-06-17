package stores

import (
	"context"
	"fmt"

	"github.com/tuantran1810/go-di-template/uberfx/internal/models"
)

type UserStore struct {
	*GenericStore[models.User]
}

func NewUserStore(repository *Repository) *UserStore {
	return &UserStore{
		GenericStore: NewGenericStore[models.User](repository),
	}
}

func (s *UserStore) Start(ctx context.Context) error {
	log.Infof("starting user store")

	s.repository.mutex.Lock()
	err := s.repository.db.AutoMigrate(&models.User{})
	s.repository.mutex.Unlock()
	if err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return s.Ping(timeoutCtx)
}

func (s *UserStore) Stop(_ context.Context) error {
	return nil
}

func (s *UserStore) FindByUsername(
	ctx context.Context,
	tx models.Transaction,
	username string,
) (*models.User, error) {
	if username == "" {
		return nil, fmt.Errorf("%w - input username is empty", models.ErrInvalid)
	}

	s.repository.mutex.RLock()
	defer s.repository.mutex.RUnlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	dbtx := s.repository.getTransaction(tx).WithContext(timeoutCtx)
	var user models.User
	if err := dbtx.First(&user, "username = ?", username).Error; err != nil {
		return nil, generateError("failed to find user by username", err)
	}

	return &user, nil
}
