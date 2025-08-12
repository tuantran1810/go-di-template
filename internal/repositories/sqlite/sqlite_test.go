package stores

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tuantran1810/go-di-template/internal/entities"
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
