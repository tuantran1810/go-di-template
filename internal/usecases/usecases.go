package usecases

import (
	"context"
	"time"

	"github.com/tuantran1810/go-di-template/internal/entities"
)

const defaultTimeout = 10 * time.Second

type IUserStore interface {
	Create(ctx context.Context, tx entities.Transaction, user *entities.User) (*entities.User, error)
	FindByUsername(ctx context.Context, tx entities.Transaction, username string) (*entities.User, error)
}

type IUserAttributeStore interface {
	CreateMany(ctx context.Context, tx entities.Transaction, userAttributes []entities.UserAttribute) ([]entities.UserAttribute, error)
	GetByUserID(ctx context.Context, tx entities.Transaction, userID uint) ([]entities.UserAttribute, error)
}

type IRepository interface {
	RunTx(ctx context.Context, funcs ...entities.DBTxHandleFunc) error
}

type IUUIDGenerator interface {
	MustNewUUID() string
}
