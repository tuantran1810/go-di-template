package controllers

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tuantran1810/go-di-template/internal/controllers/transformers"
	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/libs/utils"
	mocks "github.com/tuantran1810/go-di-template/mocks/controllers"
	pb "github.com/tuantran1810/go-di-template/pkg/go_di_template/v1"
)

func TestUserController_CreateUser(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()

	mockUserUsecase := mocks.NewMockIUserUsecase(t)
	mockUserUsecase.EXPECT().CreateUser(
		mock.Anything,
		&entities.User{
			Username: "test1",
			Password: "test1",
			Name:     "test1",
			Email:    &[]string{"test1@test.com"}[0],
		},
		[]entities.KeyValuePair{
			{
				Key:   "key1",
				Value: "value1",
			},
			{
				Key:   "key2",
				Value: "value2",
			},
		},
	).Return(
		&entities.User{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			Uuid:      "test1",
			Username:  "test1",
			Password:  "test1",
			Name:      "test1",
			Email:     &[]string{"test1@test.com"}[0],
		},
		[]entities.UserAttribute{
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
		nil,
	)

	mockUserUsecase.EXPECT().CreateUser(
		mock.Anything,
		&entities.User{
			Username: "test_failed",
			Password: "test_failed",
			Name:     "test_failed",
			Email:    &[]string{"test_failed@test.com"}[0],
		},
		[]entities.KeyValuePair{
			{
				Key:   "key1",
				Value: "value1",
			},
		},
	).Return(nil, nil, fmt.Errorf("fake error"))

	mockLoggingWorker := mocks.NewMockILoggingWorker(t)
	mockLoggingWorker.EXPECT().
		Inject(entities.Message{
			Key:   "user_created",
			Value: "user_id: 1, username: test1",
		}).
		Return()

	c := &UserController{
		userUsecase:              mockUserUsecase,
		loggingWorker:            mockLoggingWorker,
		userTransformer:          transformers.NewPbUserTransformer(),
		keyValuePairTransformer:  entities.NewBaseExtendedTransformer[pb.KeyValuePair, entities.KeyValuePair](),
		userAttributeTransformer: transformers.NewPbUserAttributesTransformer(),
	}

	tests := []struct {
		name    string
		req     *pb.CreateUserRequest
		want    *pb.CreateUserResponse
		wantErr bool
	}{
		{
			name: "success",
			req: &pb.CreateUserRequest{
				User: &pb.User{
					Username: "test1",
					Password: "test1",
					Name:     "test1",
					Email:    &[]string{"test1@test.com"}[0],
				},
				Attributes: []*pb.KeyValuePair{
					{
						Key:   "key1",
						Value: "value1",
					},
					{
						Key:   "key2",
						Value: "value2",
					},
				},
			},
			want: &pb.CreateUserResponse{
				User: &pb.User{
					Id:        1,
					CreatedAt: utils.ToTimepb(now),
					UpdatedAt: utils.ToTimepb(now),
					Uuid:      "test1",
					Username:  "test1",
					Password:  "test1",
					Name:      "test1",
					Email:     &[]string{"test1@test.com"}[0],
				},
				Attributes: []*pb.UserAttribute{
					{
						Id:        1,
						CreatedAt: utils.ToTimepb(now),
						UpdatedAt: utils.ToTimepb(now),
						UserId:    1,
						Key:       "key1",
						Value:     "value1",
					},
					{
						Id:        2,
						CreatedAt: utils.ToTimepb(now),
						UpdatedAt: utils.ToTimepb(now),
						UserId:    1,
						Key:       "key2",
						Value:     "value2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			req: &pb.CreateUserRequest{
				User: &pb.User{
					Username: "test1",
					Password: "this is a very long password, so that it should fail validation",
					Name:     "test1",
					Email:    &[]string{"test1@test.com"}[0],
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "failed to create user",
			req: &pb.CreateUserRequest{
				User: &pb.User{
					Username: "test_failed",
					Password: "test_failed",
					Name:     "test_failed",
					Email:    &[]string{"test_failed@test.com"}[0],
				},
				Attributes: []*pb.KeyValuePair{
					{
						Key:   "key1",
						Value: "value1",
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "nil user",
			req: &pb.CreateUserRequest{
				User: nil,
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := c.CreateUser(context.TODO(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserController.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserController.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserController_GetUserByUsername(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()

	mockUserUsecase := mocks.NewMockIUserUsecase(t)
	mockUserUsecase.EXPECT().
		GetUserByUsername(mock.Anything, "test1").
		Return(&entities.User{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			Username:  "test1",
			Password:  "test1",
			Uuid:      "test1",
			Name:      "test1",
			Email:     &[]string{"test1@test.com"}[0],
		}, []entities.UserAttribute{
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

	mockUserUsecase.EXPECT().
		GetUserByUsername(mock.Anything, "test_failed").
		Return(nil, nil, fmt.Errorf("fake error"))

	mockLoggingWorker := mocks.NewMockILoggingWorker(t)
	mockLoggingWorker.EXPECT().
		Inject(entities.Message{
			Key:   "user_get",
			Value: "user_id: 1, username: test1",
		}).
		Return()

	c := &UserController{
		userUsecase:              mockUserUsecase,
		loggingWorker:            mockLoggingWorker,
		userTransformer:          transformers.NewPbUserTransformer(),
		keyValuePairTransformer:  entities.NewBaseExtendedTransformer[pb.KeyValuePair, entities.KeyValuePair](),
		userAttributeTransformer: transformers.NewPbUserAttributesTransformer(),
	}

	tests := []struct {
		name    string
		req     *pb.GetUserByUsernameRequest
		want    *pb.GetUserByUsernameResponse
		wantErr bool
	}{
		{
			name: "success",
			req: &pb.GetUserByUsernameRequest{
				Username: "test1",
			},
			want: &pb.GetUserByUsernameResponse{
				User: &pb.User{
					Id:        1,
					CreatedAt: utils.ToTimepb(now),
					UpdatedAt: utils.ToTimepb(now),
					Uuid:      "test1",
					Username:  "test1",
					Password:  "test1",
					Name:      "test1",
					Email:     &[]string{"test1@test.com"}[0],
				},
				Attributes: []*pb.UserAttribute{
					{
						Id:        1,
						CreatedAt: utils.ToTimepb(now),
						UpdatedAt: utils.ToTimepb(now),
						UserId:    1,
						Key:       "key1",
						Value:     "value1",
					},
					{
						Id:        2,
						CreatedAt: utils.ToTimepb(now),
						UpdatedAt: utils.ToTimepb(now),
						UserId:    1,
						Key:       "key2",
						Value:     "value2",
					},
				},
			},
		},
		{
			name: "invalid request",
			req: &pb.GetUserByUsernameRequest{
				Username: "this is a very long username, so that it should fail validation, this is a very long username, so that it should fail validation.",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid request",
			req: &pb.GetUserByUsernameRequest{
				Username: "",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "failed to get user by username",
			req: &pb.GetUserByUsernameRequest{
				Username: "test_failed",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := c.GetUserByUsername(context.TODO(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserController.GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserController.GetUserByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
