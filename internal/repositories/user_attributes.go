package repositories

import (
	"context"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/internal/repositories/mysql"
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

type UserAttributeRepository struct {
	*mysql.GenericRepository[UserAttribute, entities.UserAttribute]
}

func NewUserAttributeRepository(repository *mysql.Repository) *UserAttributeRepository {
	transformer := entities.NewExtendedDataTransformer(&userAttributeTransformer{})
	return &UserAttributeRepository{
		GenericRepository: mysql.NewGenericRepository(repository, transformer),
	}
}

func (s *UserAttributeRepository) Start(ctx context.Context) error {
	log.Info("starting user attribute store")

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	err := s.AutoMigrate(timeoutCtx)
	if err != nil {
		return err
	}

	return s.Ping(timeoutCtx)
}

func (s *UserAttributeRepository) Stop(_ context.Context) error {
	log.Info("stopping user attribute store")
	return nil
}

func (s *UserAttributeRepository) GetByUserID(
	ctx context.Context,
	tx entities.Transaction,
	userID uint,
) ([]entities.UserAttribute, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return s.GetManyByCriterias(
		timeoutCtx, tx,
		nil,
		map[string]any{"user_id": userID},
		[]string{"id"},
		0, 0,
	)
}
