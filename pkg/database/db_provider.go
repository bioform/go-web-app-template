package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

var dbKey = contextDBKey("DB")

type contextDBKey string

type DbProvider struct {
	db *gorm.DB
}

func (provider *DbProvider) Transaction(ctx context.Context, lambda func(newContext context.Context) error) error {
	return provider.db.Transaction(func(tx *gorm.DB) error {
		return lambda(provider.SetDB(ctx, tx))
	})
}

func (provider *DbProvider) DB(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(dbKey).(*gorm.DB)
	if ok {
		return db
	}

	return provider.db
}

func (dp *DbProvider) SetDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

func (dp *DbProvider) String() string {
	db := dp.db
	if db == nil {
		return fmt.Sprintf("%T(nil)", dp)
	}

	dialect := dp.db.Dialector.Name()
	return fmt.Sprintf("%T(db: %s)", dp, dialect)
}
