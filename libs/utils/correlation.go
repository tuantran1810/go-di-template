package utils

import (
	"context"

	"github.com/google/uuid"
)

type CorrelationIDKey string

const XCorrelationID CorrelationIDKey = "x-correlation-id"

func GenerateCorrelationID() string {
	u := uuid.New()
	return u.String()
}

func GetCorrelationID(ctx context.Context) string {
	correlationID := ctx.Value(XCorrelationID)
	if correlationID == nil {
		return GenerateCorrelationID()
	}

	id, ok := correlationID.(string)
	if !ok {
		return GenerateCorrelationID()
	}

	return id
}

func InjectCorrelationIDToContext(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, XCorrelationID, correlationID)
}
