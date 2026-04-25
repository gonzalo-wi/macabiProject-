package database

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

func withTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// TxFromCtx returns the active transaction from ctx, or fallback if none.
func TxFromCtx(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return fallback
}

type GORMTransactor struct {
	db *gorm.DB
}

func NewGORMTransactor(db *gorm.DB) *GORMTransactor {
	return &GORMTransactor{db: db}
}

func (t *GORMTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(withTx(ctx, tx))
	})
}
