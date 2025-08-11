package usecases

import (
	"context"
	"time"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/libs/logger"
)

var log = logger.MustNamedLogger("usecases")

const defaultTimeout = 10 * time.Second

type IUserStore interface {
	Create(ctx context.Context, tx entities.Transaction, user *entities.User) (*entities.User, error)
	FindByUsername(ctx context.Context, tx entities.Transaction, username string) (*entities.User, error)
}

type IUserAttributeStore interface {
	CreateMany(ctx context.Context, tx entities.Transaction, userAttributes []entities.UserAttribute) ([]entities.UserAttribute, error)
	GetByUserID(ctx context.Context, tx entities.Transaction, userID uint) ([]entities.UserAttribute, error)
}

type IMessageStore interface {
	CreateMany(ctx context.Context, tx entities.Transaction, messages []entities.Message) ([]entities.Message, error)
}

type IRepository interface {
	RunTx(ctx context.Context, funcs ...entities.DBTxHandleFunc) error
}

type IUUIDGenerator interface {
	MustNewUUID() string
}
