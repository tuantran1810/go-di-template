package middlewares

import (
	"fmt"
	"testing"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_HandleError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		errType error
		outErr  error
	}{
		{
			errType: nil,
			outErr:  nil,
		},
		{
			errType: entities.ErrCanceled,
			outErr:  status.Error(codes.Canceled, "canceled - test err"),
		},
		{
			errType: entities.ErrTooManyRequests,
			outErr:  status.Error(codes.ResourceExhausted, "too many requests - test err"),
		},
		{
			errType: entities.ErrDatabase,
			outErr:  status.Error(codes.Unavailable, "database error - test err"),
		},
		{
			errType: entities.ErrInvalid,
			outErr:  status.Error(codes.InvalidArgument, "invalid input - test err"),
		},
		{
			errType: entities.ErrNotFound,
			outErr:  status.Error(codes.NotFound, "not found - test err"),
		},
		{
			errType: entities.ErrUnauthorized,
			outErr:  status.Error(codes.Unauthenticated, "unauthorized - test err"),
		},
		{
			errType: entities.ErrInternal,
			outErr:  status.Error(codes.Internal, "internal error - test err"),
		},
	}
	for _, tt := range tests {
		t.Run("Test_HandleError", func(t *testing.T) {
			t.Parallel()

			var ierr error
			if tt.errType != nil {
				ierr = fmt.Errorf("%w - test err", tt.errType)
			}

			err := HandleError(ierr)
			if err == nil && tt.outErr == nil {
				return
			}

			if err.Error() != tt.outErr.Error() {
				t.Errorf("HandleError() = %v, want %v", err, tt.outErr)
			}
		})
	}
}
