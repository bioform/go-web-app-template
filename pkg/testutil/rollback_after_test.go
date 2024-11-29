package testutil

import (
	"context"

	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/onsi/ginkgo/v2"
	"gorm.io/gorm"
)

func RollbackAfterTest(ctx context.Context) *gorm.DB {
	tx := database.Get(ctx).Begin()

	ginkgo.DeferCleanup(tx.Rollback)

	return tx
}
