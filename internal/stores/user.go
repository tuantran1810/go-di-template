package stores

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/internal/stores/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex,size:32"`
	Password string
	Uuid     string
	Name     string
	Email    sql.NullString
}

type userTransformer struct{}

func (t *userTransformer) ToEntity(data *User) (*entities.User, error) {
	var email *string
	if data.Email.Valid {
		email = &data.Email.String
	}
	return &entities.User{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Username:  data.Username,
		Password:  data.Password,
		Uuid:      data.Uuid,
		Name:      data.Name,
		Email:     email,
	}, nil
}

func (t *userTransformer) FromEntity(entity *entities.User) (*User, error) {
	var email sql.NullString
	if entity.Email != nil {
		email = sql.NullString{String: *entity.Email, Valid: true}
	}
	return &User{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		Username: entity.Username,
		Password: entity.Password,
		Uuid:     entity.Uuid,
		Name:     entity.Name,
		Email:    email,
	}, nil
}

type UserStore struct {
	*mysql.GenericStore[User, entities.User]
}

func NewUserStore(repository *mysql.Repository) *UserStore {
	transformer := entities.NewExtendedDataTransformer(&userTransformer{})
	return &UserStore{
		GenericStore: mysql.NewGenericStore(repository, transformer),
	}
}

func (s *UserStore) Start(ctx context.Context) error {
	log.Info("starting user store")
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	err := s.AutoMigrate(timeoutCtx)
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
	tx entities.Transaction,
	username string,
) (*entities.User, error) {
	if username == "" {
		return nil, fmt.Errorf("%w - input username is empty", entities.ErrInvalid)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	user, err := s.GetByCriterias(
		timeoutCtx, tx,
		nil,
		map[string]any{"username": username},
		[]string{"id"},
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
