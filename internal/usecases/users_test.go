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
	mockUserRepository := mockUsecases.NewMockIUserRepository(t)
	mockUserAttributeRepository := mockUsecases.NewMockIUserAttributeRepository(t)

	mockUserRepository.EXPECT().
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

	mockUserRepository.EXPECT().
		Create(mock.Anything, mock.Anything, &entities.User{
			Username: "test_failed",
		}).
		Return(nil, errors.New("fake error"))

	mockUserAttributeRepository.EXPECT().
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

	mockUserAttributeRepository.EXPECT().
		CreateMany(mock.Anything, mock.Anything, []entities.UserAttribute{
			{
				UserID: 1,
				Key:    "failed_key",
				Value:  "failed_value",
			},
		}).
		Return(nil, errors.New("fake error"))

	u := &Users{
		repository:              mockRepository,
		userRepository:          mockUserRepository,
		userAttributeRepository: mockUserAttributeRepository,
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
	mockUserRepository := mockUsecases.NewMockIUserRepository(t)
	mockUserAttributeRepository := mockUsecases.NewMockIUserAttributeRepository(t)
	mockUUIDGenerator := mockUsecases.NewMockIUUIDGenerator(t)

	mockUserRepository.EXPECT().
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

	mockUserRepository.EXPECT().
		Create(mock.Anything, mock.Anything, &entities.User{
			Username: "test_failed",
			Uuid:     "test1",
		}).
		Return(nil, errors.New("fake error"))

	mockUserAttributeRepository.EXPECT().
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

	mockUserAttributeRepository.EXPECT().
		CreateMany(mock.Anything, mock.Anything, []entities.UserAttribute{
			{
				UserID: 1,
				Key:    "failed_key",
				Value:  "failed_value",
			},
		}).
		Return(nil, errors.New("fake error"))

	u := &Users{
		repository:              mockRepository,
		userRepository:          mockUserRepository,
		userAttributeRepository: mockUserAttributeRepository,
		uuidGenerator:           mockUUIDGenerator,
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

func TestUsers_GetUserByUsername(t *testing.T) {
	t.Parallel()
	now := time.Now()

	mockUserRepository := mockUsecases.NewMockIUserRepository(t)
	mockUserAttributeRepository := mockUsecases.NewMockIUserAttributeRepository(t)

	mockUserRepository.EXPECT().
		FindByUsername(mock.Anything, mock.Anything, "test1").
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

	mockUserRepository.EXPECT().
		FindByUsername(mock.Anything, mock.Anything, "test_failed").
		Return(nil, errors.New("fake error"))

	mockUserRepository.EXPECT().
		FindByUsername(mock.Anything, mock.Anything, "failed_atts").
		Return(&entities.User{
			ID:        2,
			CreatedAt: now,
			UpdatedAt: now,
			Username:  "failed_atts",
		}, nil)
	mockUserRepository.EXPECT().
		FindByUsername(mock.Anything, mock.Anything, "no_atts").
		Return(&entities.User{
			ID:        3,
			CreatedAt: now,
			UpdatedAt: now,
			Username:  "no_atts",
			Password:  "no_atts",
			Uuid:      "no_atts",
			Name:      "no_atts",
			Email:     &[]string{"no_atts@test.com"}[0],
		}, nil)

	mockUserAttributeRepository.EXPECT().
		GetByUserID(mock.Anything, mock.Anything, uint(1)).
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
		}, nil)
	mockUserAttributeRepository.EXPECT().
		GetByUserID(mock.Anything, mock.Anything, uint(2)).
		Return(nil, errors.New("fake error"))

	mockUserAttributeRepository.EXPECT().
		GetByUserID(mock.Anything, mock.Anything, uint(3)).
		Return([]entities.UserAttribute{}, nil)

	u := &Users{
		userRepository:          mockUserRepository,
		userAttributeRepository: mockUserAttributeRepository,
	}

	tests := []struct {
		name     string
		username string
		want     *entities.User
		want1    []entities.UserAttribute
		wantErr  bool
	}{
		{
			name:     "success",
			username: "test1",
			want: &entities.User{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				Username:  "test1",
				Password:  "test1",
				Uuid:      "test1",
				Name:      "test1",
				Email:     &[]string{"test1@test.com"}[0],
			},
			want1: []entities.UserAttribute{
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
			},
		},
		{
			name:     "failed to find user",
			username: "test_failed",
			want:     nil,
			want1:    nil,
			wantErr:  true,
		},
		{
			name:     "failed to get user attributes",
			username: "failed_atts",
			want:     nil,
			want1:    nil,
			wantErr:  true,
		},
		{
			name:     "no attributes",
			username: "no_atts",
			want: &entities.User{
				ID:        3,
				CreatedAt: now,
				UpdatedAt: now,
				Username:  "no_atts",
				Password:  "no_atts",
				Uuid:      "no_atts",
				Name:      "no_atts",
				Email:     &[]string{"no_atts@test.com"}[0],
			},
			want1:   []entities.UserAttribute{},
			wantErr: false,
		},
		{
			name:     "empty username",
			username: "",
			want:     nil,
			want1:    nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, err := u.GetUserByUsername(context.TODO(), tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.GetUserByUsername() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Users.GetUserByUsername() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
