package dbaction

import (
	"context"

	"github.com/bioform/go-web-app-template/pkg/action"
	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/bioform/go-web-app-template/pkg/logging"
	"gorm.io/gorm"
)

type BaseAction struct {
	action.BaseAction
}

func (a BaseAction) DB(ctx context.Context) *gorm.DB {
	actionTransactionProvider := a.TransactionProvider()

	dbProvider, ok := actionTransactionProvider.(*database.DbProvider)
	if !ok {
		log := logging.Get(ctx)
		log.Error("DB provider is not a *database.DbProvider", "provider", actionTransactionProvider)
		return nil // or handle the error appropriately
	}

	return dbProvider.DB(ctx)
}
