package middlewares

import (
	"context"
	"errors"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleError(err error) error {
	if err == nil {
		return nil
	}

	errString := err.Error()

	switch {
	case errors.Is(err, entities.ErrCanceled):
		return status.Error(codes.Canceled, errString)
	case errors.Is(err, entities.ErrTooManyRequests):
		return status.Error(codes.ResourceExhausted, errString)
	case errors.Is(err, entities.ErrDatabase):
		return status.Error(codes.Unavailable, errString)
	case errors.Is(err, entities.ErrInvalid):
		return status.Error(codes.InvalidArgument, errString)
	case errors.Is(err, entities.ErrNotFound):
		return status.Error(codes.NotFound, errString)
	case errors.Is(err, entities.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, errString)
	case errors.Is(err, entities.ErrInternal):
		return status.Error(codes.Internal, errString)
	default:
		return status.Error(codes.Unknown, errString)
	}
}

func HandleErrorCodes(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	out, err := handler(ctx, req)

	return out, HandleError(err)
}
