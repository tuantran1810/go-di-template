package models

import (
	"context"

	"gorm.io/gorm"
)

type Transaction interface {
	GetTransaction() any
}

type GormTransaction struct {
	Tx *gorm.DB
}

func NewGormTransaction(tx *gorm.DB) *GormTransaction {
	return &GormTransaction{Tx: tx}
}

func (t *GormTransaction) GetTransaction() any {
	return t.Tx
}

type DBTxHandleFunc func(ctx context.Context, dbtx Transaction, data any) (out any, cont bool, err error)
