package usecases

import (
	"context"
	"fmt"

	"github.com/tuantran1810/go-di-template/internal/entities"
)

type Users struct {
	repository         IRepository
	userStore          IUserStore
	userAttributeStore IUserAttributeStore
}

func NewUsers(
	repository IRepository,
	userStore IUserStore,
	userAttributeStore IUserAttributeStore,
) *Users {
	return &Users{
		userStore:          userStore,
		userAttributeStore: userAttributeStore,
	}
}

func (u *Users) Create(ctx context.Context, user *entities.User, attributes []entities.KeyValuePair) (*entities.User, []entities.UserAttribute, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var outUser *entities.User
	outAttributes := make([]entities.UserAttribute, 0)

	err := u.repository.RunTx(
		timeoutCtx,
		func(ictx context.Context, dbtx entities.Transaction) error {
			var ierr error
			outUser, ierr = u.userStore.Create(ictx, dbtx, user)
			if ierr != nil {
				return fmt.Errorf("failed to create user: %w", ierr)
			}

			if len(attributes) == 0 {
				return nil
			}

			atts := make([]entities.UserAttribute, len(attributes))
			for i, attr := range attributes {
				atts[i] = entities.UserAttribute{
					UserID: outUser.ID,
					Key:    attr.Key,
					Value:  attr.Value,
				}
			}

			outAttributes, ierr = u.userAttributeStore.CreateMany(ictx, dbtx, atts)
			if ierr != nil {
				return fmt.Errorf("failed to create user attributes: %w", ierr)
			}

			return nil
		},
	)

	if err != nil {
		return nil, nil, err
	}

	return outUser, outAttributes, nil
}
