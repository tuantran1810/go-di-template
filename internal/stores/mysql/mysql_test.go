package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tuantran1810/go-di-template/internal/models"
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
			want: models.ErrCanceled,
		},
		{
			name: "invalid error",
			err:  gorm.ErrInvalidField,
			want: models.ErrInvalid,
		},
		{
			name: "not found error",
			err:  gorm.ErrRecordNotFound,
			want: models.ErrNotFound,
		},
		{
			name: "unknown error",
			err:  errors.New("error"),
			want: models.ErrDatabase,
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
	mysqlContainer, err := mysql.Run(context.Background(),
		"mysql:lts",
		mysql.WithDatabase("test"),
		mysql.WithUsername("root"),
		mysql.WithPassword("secret"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(30*time.Second),
			wait.ForListeningPort("3306/tcp").WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Errorf("failed to start mysql container: %v", err)
		return
	}

	port, err := mysqlContainer.MappedPort(context.Background(), "3306")
	if err != nil {
		t.Errorf("failed to get mysql port: %v", err)
		return
	}

	defer mysqlContainer.Terminate(context.Background())

	config := RepositoryConfig{
		Username:  "root",
		Password:  "secret",
		Protocol:  "tcp",
		Address:   fmt.Sprintf("127.0.0.1:%d", port.Int()),
		Database:  "test",
		Params:    map[string]string{},
		Collation: "utf8mb4_general_ci",
		Loc:       time.Local,
		TLSConfig: "",

		Timeout:                 10 * time.Second,
		ReadTimeout:             10 * time.Second,
		WriteTimeout:            10 * time.Second,
		AllowAllFiles:           false,
		AllowCleartextPasswords: false,
		AllowOldPasswords:       false,
		ClientFoundRows:         false,
		ColumnsWithAlias:        false,
		InterpolateParams:       false,
		MultiStatements:         false,
		ParseTime:               true,

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
