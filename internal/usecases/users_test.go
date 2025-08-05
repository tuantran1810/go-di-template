package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tuantran1810/go-di-template/internal/entities"
	mockUsecases "github.com/tuantran1810/go-di-template/mocks/usecases"
)

type mockRepository struct{}

func (m *mockRepository) RunTx(ctx context.Context, funcs ...entities.DBTxHandleFunc) error {
	for _, f := range funcs {
		if f != nil {
			if err := f(ctx, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func TestUsers_createUserImpl(t *testing.T) {
	t.Parallel()
	now := time.Now()

	mockRepository := &mockRepository{}
	mockUserStore := mockUsecases.NewMockIUserStore(t)
	mockUserAttributeStore := mockUsecases.NewMockIUserAttributeStore(t)

	mockUserStore.EXPECT().
		Create(mock.Anything, mock.Anything, &entities.User{
			Username: "test1",
			Password: "test1",
			Uuid:     "test1",
			Name:     "test1",
			Email:    &[]string{"test1@test.com"}[0],
		}).
		Return(&entities.User{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			Username:  "test1",
			Password:  "test1",
			Uuid:      "test1",
			Name:      "test1",
			Email:     &[]string{"test1@test.com"}[0],
		}, nil)

	mockUserStore.EXPECT().
		Create(mock.Anything, mock.Anything, &entities.User{
			Username: "test_failed",
		}).
		Return(nil, errors.New("fake error"))

	mockUserAttributeStore.EXPECT().
		CreateMany(mock.Anything, mock.Anything, []entities.UserAttribute{
			{
				UserID: 1,
				Key:    "key1",
				Value:  "value1",
			},
			{
				UserID: 1,
				Key:    "key2",
				Value:  "value2",
			},
			{
				UserID: 1,
				Key:    "key3",
				Value:  "value3",
			},
		}).
		Return([]entities.UserAttribute{
			{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    1,
				Key:       "key1",
				Value:     "value1",
			},
			{
				ID:        2,
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    1,
				Key:       "key2",
				Value:     "value2",
			},
			{
				ID:        3,
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    1,
				Key:       "key3",
				Value:     "value3",
			},
		}, nil)

	mockUserAttributeStore.EXPECT().
		CreateMany(mock.Anything, mock.Anything, []entities.UserAttribute{
			{
				UserID: 1,
				Key:    "failed_key",
				Value:  "failed_value",
			},
		}).
		Return(nil, errors.New("fake error"))

	u := &Users{
		repository:         mockRepository,
		userStore:          mockUserStore,
		userAttributeStore: mockUserAttributeStore,
	}

	tests := []struct {
		name           string
		user           *entities.User
		attributes     []entities.KeyValuePair
		wantUser       *entities.User
		wantAttributes []entities.UserAttribute
		wantErr        bool
	}{
		{
			name: "success",
			user: &entities.User{
				Username: "test1",
				Password: "test1",
				Uuid:     "test1",
				Name:     "test1",
				Email:    &[]string{"test1@test.com"}[0],
			},
			attributes: []entities.KeyValuePair{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
				{
					Key:   "key3",
					Value: "value3",
				},
			},
			wantUser: &entities.User{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				Username:  "test1",
				Password:  "test1",
				Uuid:      "test1",
				Name:      "test1",
				Email:     &[]string{"test1@test.com"}[0],
			},
			wantAttributes: []entities.UserAttribute{
				{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
					UserID:    1,
					Key:       "key1",
					Value:     "value1",
				},
				{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
					UserID:    1,
					Key:       "key2",
					Value:     "value2",
				},
				{
					ID:        3,
					CreatedAt: now,
					UpdatedAt: now,
					UserID:    1,
					Key:       "key3",
					Value:     "value3",
				},
			},
			wantErr: false,
		},
		{
			name: "failed to create user",
			user: &entities.User{
				Username: "test_failed",
			},
			attributes: []entities.KeyValuePair{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
			wantUser:       nil,
			wantAttributes: nil,
			wantErr:        true,
		},
		{
			name: "create user with no attributes",
			user: &entities.User{
				Username: "test1",
				Password: "test1",
				Uuid:     "test1",
				Name:     "test1",
				Email:    &[]string{"test1@test.com"}[0],
			},
			attributes: []entities.KeyValuePair{},
			wantUser: &entities.User{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				Username:  "test1",
				Password:  "test1",
				Uuid:      "test1",
				Name:      "test1",
				Email:     &[]string{"test1@test.com"}[0],
			},
			wantAttributes: []entities.UserAttribute{},
			wantErr:        false,
		},
		{
			name: "failed to create user attributes",
			user: &entities.User{
				Username: "test1",
				Password: "test1",
				Uuid:     "test1",
				Name:     "test1",
				Email:    &[]string{"test1@test.com"}[0],
			},
			attributes: []entities.KeyValuePair{
				{
					Key:   "failed_key",
					Value: "failed_value",
				},
			},
			wantUser:       nil,
			wantAttributes: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotUser, gotAttributes, err := u.createUserImpl(context.TODO(), nil, tt.user, tt.attributes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.createUserImpl() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("Users.createUserImpl() gotUser = %v, want %v", gotUser, tt.user)
			}
			if !reflect.DeepEqual(gotAttributes, tt.wantAttributes) {
				t.Errorf("Users.createUserImpl() gotAttributes = %v, want %v", gotAttributes, tt.wantAttributes)
			}
		})
	}
}

func TestUsers_CreateUser(t *testing.T) {
	t.Parallel()
	now := time.Now()

	mockRepository := &mockRepository{}
	mockUserStore := mockUsecases.NewMockIUserStore(t)
	mockUserAttributeStore := mockUsecases.NewMockIUserAttributeStore(t)
	mockUUIDGenerator := mockUsecases.NewMockIUUIDGenerator(t)

	mockUserStore.EXPECT().
		Create(mock.Anything, mock.Anything, &entities.User{
			Username: "test1",
			Password: "test1",
			Uuid:     "test1",
			Name:     "test1",
			Email:    &[]string{"test1@test.com"}[0],
		}).
		Return(&entities.User{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			Username:  "test1",
			Password:  "test1",
			Uuid:      "test1",
			Name:      "test1",
			Email:     &[]string{"test1@test.com"}[0],
		}, nil)

	mockUUIDGenerator.EXPECT().
		MustNewUUID().
		Return("test1")

	mockUserStore.EXPECT().
		Create(mock.Anything, mock.Anything, &entities.User{
			Username: "test_failed",
			Uuid:     "test1",
		}).
		Return(nil, errors.New("fake error"))

	mockUserAttributeStore.EXPECT().
		CreateMany(mock.Anything, mock.Anything, []entities.UserAttribute{
			{
				UserID: 1,
				Key:    "key1",
				Value:  "value1",
			},
			{
				UserID: 1,
				Key:    "key2",
				Value:  "value2",
			},
			{
				UserID: 1,
				Key:    "key3",
				Value:  "value3",
			},
		}).
		Return([]entities.UserAttribute{
			{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    1,
				Key:       "key1",
				Value:     "value1",
			},
			{
				ID:        2,
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    1,
				Key:       "key2",
				Value:     "value2",
			},
			{
				ID:        3,
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    1,
				Key:       "key3",
				Value:     "value3",
			},
		}, nil)

	mockUserAttributeStore.EXPECT().
		CreateMany(mock.Anything, mock.Anything, []entities.UserAttribute{
			{
				UserID: 1,
				Key:    "failed_key",
				Value:  "failed_value",
			},
		}).
		Return(nil, errors.New("fake error"))

	u := &Users{
		repository:         mockRepository,
		userStore:          mockUserStore,
		userAttributeStore: mockUserAttributeStore,
		uuidGenerator:      mockUUIDGenerator,
	}

	tests := []struct {
		name           string
		user           *entities.User
		attributes     []entities.KeyValuePair
		wantUser       *entities.User
		wantAttributes []entities.UserAttribute
		wantErr        bool
	}{
		{
			name: "success",
			user: &entities.User{
				Username: "test1",
				Password: "test1",
				Name:     "test1",
				Email:    &[]string{"test1@test.com"}[0],
			},
			attributes: []entities.KeyValuePair{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
				{
					Key:   "key3",
					Value: "value3",
				},
			},
			wantUser: &entities.User{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				Username:  "test1",
				Password:  "test1",
				Uuid:      "test1",
				Name:      "test1",
				Email:     &[]string{"test1@test.com"}[0],
			},
			wantAttributes: []entities.UserAttribute{
				{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
					UserID:    1,
					Key:       "key1",
					Value:     "value1",
				},
				{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
					UserID:    1,
					Key:       "key2",
					Value:     "value2",
				},
				{
					ID:        3,
					CreatedAt: now,
					UpdatedAt: now,
					UserID:    1,
					Key:       "key3",
					Value:     "value3",
				},
			},
			wantErr: false,
		},
		{
			name: "failed to create user",
			user: &entities.User{
				Username: "test_failed",
			},
			attributes: []entities.KeyValuePair{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
			wantUser:       nil,
			wantAttributes: nil,
			wantErr:        true,
		},
		{
			name: "create user with no attributes",
			user: &entities.User{
				Username: "test1",
				Password: "test1",
				Name:     "test1",
				Email:    &[]string{"test1@test.com"}[0],
			},
			attributes: []entities.KeyValuePair{},
			wantUser: &entities.User{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				Username:  "test1",
				Password:  "test1",
				Uuid:      "test1",
				Name:      "test1",
				Email:     &[]string{"test1@test.com"}[0],
			},
			wantAttributes: []entities.UserAttribute{},
			wantErr:        false,
		},
		{
			name: "failed to create user attributes",
			user: &entities.User{
				Username: "test1",
				Password: "test1",
				Name:     "test1",
				Email:    &[]string{"test1@test.com"}[0],
			},
			attributes: []entities.KeyValuePair{
				{
					Key:   "failed_key",
					Value: "failed_value",
				},
			},
			wantUser:       nil,
			wantAttributes: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotUser, gotAttributes, err := u.CreateUser(context.TODO(), tt.user, tt.attributes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("Users.CreateUser() gotUser = %v, want %v", gotUser, tt.user)
			}
			if !reflect.DeepEqual(gotAttributes, tt.wantAttributes) {
				t.Errorf("Users.CreateUser() gotAttributes = %v, want %v", gotAttributes, tt.wantAttributes)
			}
		})
	}
}
