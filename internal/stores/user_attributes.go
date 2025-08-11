package stores

import (
	"context"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/internal/stores/mysql"
	"gorm.io/gorm"
)

type UserAttribute struct {
	gorm.Model
	UserID uint `gorm:"index"`
	Key    string
	Value  string
}

type userAttributeTransformer struct{}

func (t *userAttributeTransformer) ToEntity(data *UserAttribute) (*entities.UserAttribute, error) {
	return &entities.UserAttribute{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		UserID:    data.UserID,
		Key:       data.Key,
		Value:     data.Value,
	}, nil
}

func (t *userAttributeTransformer) FromEntity(entity *entities.UserAttribute) (*UserAttribute, error) {
	return &UserAttribute{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		UserID: entity.UserID,
		Key:    entity.Key,
		Value:  entity.Value,
	}, nil
}

type UserAttributeStore struct {
	*mysql.GenericStore[UserAttribute, entities.UserAttribute]
}

func NewUserAttributeStore(repository *mysql.Repository) *UserAttributeStore {
	transformer := entities.NewExtendedDataTransformer(&userAttributeTransformer{})
	return &UserAttributeStore{
		GenericStore: mysql.NewGenericStore(repository, transformer),
	}
}

func (s *UserAttributeStore) Start(ctx context.Context) error {
	log.Info("starting user attribute store")

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	err := s.GenericStore.AutoMigrate(timeoutCtx)
	if err != nil {
		return err
	}

	return s.Ping(timeoutCtx)
}

func (s *UserAttributeStore) Stop(_ context.Context) error {
	log.Info("stopping user attribute store")
	return nil
}

func (s *UserAttributeStore) GetByUserID(
	ctx context.Context,
	tx entities.Transaction,
	userID uint,
) ([]entities.UserAttribute, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return s.GenericStore.GetManyByCriterias(
		timeoutCtx, tx,
		nil,
		map[string]any{"user_id": userID},
		[]string{"id"},
		0, 0,
	)
}
