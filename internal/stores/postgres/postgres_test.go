package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/libs/utils"
	"gorm.io/gorm"
)

func Test_isInvalidInputError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "no error",
			err:  nil,
			want: false,
		},
		{
			name: "constraint failed",
			err:  errors.New("xxx constraint failed"),
			want: true,
		},
		{
			name: "invalid data",
			err:  gorm.ErrInvalidData,
			want: true,
		},
		{
			name: "invalid data",
			err:  gorm.ErrInvalidField,
			want: true,
		},
		{
			name: "invalid data",
			err:  gorm.ErrInvalidValue,
			want: true,
		},
		{
			name: "invalid data",
			err:  gorm.ErrInvalidValueOfLength,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isInvalidInputError(tt.err); got != tt.want {
				t.Errorf("isInvalidInputError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isNotFoundError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "no error",
			err:  nil,
			want: false,
		},
		{
			name: "no row error",
			err:  sql.ErrNoRows,
			want: true,
		},
		{
			name: "record not found error",
			err:  gorm.ErrRecordNotFound,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isNotFoundError(tt.err); got != tt.want {
				t.Errorf("isNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isCanceledError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "no error",
			err:  nil,
			want: false,
		},
		{
			name: "context canceled",
			err:  context.Canceled,
			want: true,
		},
		{
			name: "operation was canceled",
			err:  errors.New("xxx operation was canceled"),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isCanceledError(tt.err); got != tt.want {
				t.Errorf("isCanceledError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEntityError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want error
	}{
		{
			name: "no error",
			err:  nil,
			want: nil,
		},
		{
			name: "canceled error",
			err:  context.Canceled,
			want: entities.ErrCanceled,
		},
		{
			name: "invalid error",
			err:  gorm.ErrInvalidField,
			want: entities.ErrInvalid,
		},
		{
			name: "not found error",
			err:  gorm.ErrRecordNotFound,
			want: entities.ErrNotFound,
		},
		{
			name: "unknown error",
			err:  errors.New("error"),
			want: entities.ErrDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := getEntityError(tt.err); !errors.Is(got, tt.want) {
				t.Errorf("getEntityError() error = %v, wantErr %v", got, tt.want)
			}
		})
	}
}

func TestRepository_Start(t *testing.T) {
	postgresContainer, err := postgres.Run(context.Background(),
		"postgres:latest",
		postgres.WithDatabase("test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Errorf("failed to start mysql container: %v", err)
		return
	}

	port, err := postgresContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		t.Errorf("failed to get postgres port: %v", err)
		return
	}

	defer postgresContainer.Terminate(context.Background())

	config := RepositoryConfig{
		Host:                   "localhost",
		Port:                   port.Int(),
		Username:               "postgres",
		Password:               "postgres",
		Database:               "test",
		MaxOpenConns:           10,
		MaxIdleConns:           10,
		ConnMaxLifeTimeSeconds: 1800,
	}

	t.Run("Start", func(t *testing.T) {
		r := MustNewRepository(config)
		defer r.Stop(context.Background())

		if err := r.Start(context.Background()); err != nil {
			t.Errorf("failed to start repository: %v", err)
		}
	})
}

func TestRepositoryConfig_DSN(t *testing.T) {
	tests := []struct {
		config RepositoryConfig
		want   string
	}{
		{
			config: RepositoryConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "postgres",
				Database: "test",
				SSLMode:  utils.Pointer("disable"),
				Timezone: utils.Pointer("UTC"),
				Params:   map[string]string{"application_name": "test"},
			},
			want: "host=localhost port=5432 user=postgres password=postgres dbname=test sslmode=disable TimeZone=UTC application_name=test",
		},
	}
	for _, tt := range tests {
		t.Run("TestRepositoryConfig_DSN", func(t *testing.T) {
			if got := tt.config.DSN(); got != tt.want {
				t.Errorf("RepositoryConfig.DSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
