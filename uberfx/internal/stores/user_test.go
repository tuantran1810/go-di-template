package stores

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/tuantran1810/go-di-template/uberfx/internal/models"
	"gorm.io/gorm"
)

func getUserTestData(t *testing.T) []models.User {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Second)
	return []models.User{
		{
			Model: gorm.Model{
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username: "user1",
		},
		{
			Model: gorm.Model{
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username: "user2",
		},
		{
			Model: gorm.Model{
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username: "user3",
		},
	}
}

func createUserTestData(t *testing.T, store *UserStore) {
	t.Helper()

	if _, err := store.CreateMany(context.Background(), nil, getUserTestData(t)); err != nil {
		t.Errorf("failed to create data: %v", err)
		return
	}
}

func setupUserTest(t *testing.T) (*Repository, error) {
	t.Helper()
	r, err := NewRepository(RepositoryConfig{DatabasePath: ":memory:"})
	if err != nil {
		return nil, err
	}

	if err := r.db.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}
	return r, nil
}

func TestUserStore_FindByUsername(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC().Truncate(time.Second)
	repository, err := setupUserTest(t)
	if err != nil {
		t.Errorf("failed to setup repository: %v", err)
		return
	}
	defer repository.Stop(context.Background())

	s := NewUserStore(repository)
	createUserTestData(t, s)

	tests := []struct {
		name     string
		username string
		want     *models.User
		wantErr  bool
	}{
		{
			name:     "user1",
			username: "user1",
			want: &models.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Username: "user1",
			},
			wantErr: false,
		},
		{
			name:     "not found",
			username: "user10",
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.FindByUsername(context.TODO(), nil, tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.FindByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserStore.FindByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
