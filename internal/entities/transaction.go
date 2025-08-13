package entities

import (
	"context"
)

type Transaction interface {
	GetTransaction() any
}

type DBTxHandleFunc func(ctx context.Context, dbtx Transaction) (err error)
