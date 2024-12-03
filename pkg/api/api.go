package api

import (
	"context"
	"log/slog"

	"github.com/bioform/go-web-app-template/pkg/action"
	"gorm.io/gorm"
)

type API interface {
	action.TransactionProvider
	DB() *gorm.DB
}

type api struct {
	API
	db *gorm.DB
}

func New(db *gorm.DB) *api {
	return &api{db: db}
}

func (a *api) DB() *gorm.DB {
	return a.db
}

func (a *api) Transaction(ctx context.Context, lambda func(newContext context.Context) error) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		newAPI := *a
		newAPI.db = tx
		return lambda(newAPI.AddTo(ctx))
	})
}
func (u *api) LogValue() slog.Value {
	return slog.GroupValue()
}
