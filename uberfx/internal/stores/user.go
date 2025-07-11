package stores

import (
	"context"
	"fmt"

	"github.com/tuantran1810/go-di-template/uberfx/internal/models"
	"github.com/tuantran1810/go-di-template/uberfx/internal/stores/mysql"
)

type UserStore struct {
	*mysql.GenericStore[models.User]
}

func NewUserStore(repository *mysql.Repository) *UserStore {
	return &UserStore{
		GenericStore: mysql.NewGenericStore[models.User](repository),
	}
}

func (s *UserStore) Start(ctx context.Context) error {
	log.Info("starting user store")
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	err := s.GenericStore.AutoMigrate(timeoutCtx)
	if err != nil {
		return err
	}

	return s.Ping(timeoutCtx)
}

func (s *UserStore) Stop(_ context.Context) error {
	log.Info("stopping user store")
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

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return s.GenericStore.GetByCriterias(
		timeoutCtx, tx,
		nil,
		map[string]any{"username": username},
		[]string{"id"},
	)
}
