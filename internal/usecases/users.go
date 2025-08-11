package usecases

import (
	"context"
	"fmt"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/libs/utils"
)

type Users struct {
	repository         IRepository
	userStore          IUserStore
	userAttributeStore IUserAttributeStore
	uuidGenerator      IUUIDGenerator
}

func NewUsersUsecase(
	repository IRepository,
	userStore IUserStore,
	userAttributeStore IUserAttributeStore,
) *Users {
	return &Users{
		repository:         repository,
		userStore:          userStore,
		userAttributeStore: userAttributeStore,
		uuidGenerator:      &utils.UUIDGenerator{},
	}
}

func (u *Users) createUserImpl(
	ctx context.Context,
	dbtx entities.Transaction,
	user *entities.User,
	attributes []entities.KeyValuePair,
) (*entities.User, []entities.UserAttribute, error) {
	outUser, err := u.userStore.Create(ctx, dbtx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	atts := make([]entities.UserAttribute, len(attributes))
	if len(attributes) == 0 {
		return outUser, atts, nil
	}

	for i, attr := range attributes {
		atts[i] = entities.UserAttribute{
			UserID: outUser.ID,
			Key:    attr.Key,
			Value:  attr.Value,
		}
	}

	outAttributes, err := u.userAttributeStore.CreateMany(ctx, dbtx, atts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user attributes: %w", err)
	}

	return outUser, outAttributes, nil
}

func (u *Users) CreateUser(ctx context.Context, user *entities.User, attributes []entities.KeyValuePair) (*entities.User, []entities.UserAttribute, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	user.Uuid = u.uuidGenerator.MustNewUUID()

	var outUser *entities.User
	outAttributes := make([]entities.UserAttribute, 0)

	if err := u.repository.RunTx(
		timeoutCtx,
		func(ictx context.Context, dbtx entities.Transaction) error {
			var ierr error
			outUser, outAttributes, ierr = u.createUserImpl(ictx, dbtx, user, attributes)
			if ierr != nil {
				return fmt.Errorf("failed to create user: %w", ierr)
			}
			return nil
		},
	); err != nil {
		return nil, nil, err
	}

	return outUser, outAttributes, nil
}

func (u *Users) GetUserByUsername(ctx context.Context, username string) (*entities.User, []entities.UserAttribute, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	if username == "" {
		return nil, nil, fmt.Errorf("%w - input username is empty", entities.ErrInvalid)
	}

	user, err := u.userStore.FindByUsername(timeoutCtx, nil, username)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	atts, err := u.userAttributeStore.GetByUserID(timeoutCtx, nil, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user attributes: %w", err)
	}

	return user, atts, nil
}
