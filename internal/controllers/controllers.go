package controllers

import (
	"context"

	"github.com/tuantran1810/go-di-template/internal/entities"
)

type IUserUsecase interface {
	CreateUser(ctx context.Context, user *entities.User, attributes []entities.KeyValuePair) (*entities.User, []entities.UserAttribute, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, []entities.UserAttribute, error)
}
