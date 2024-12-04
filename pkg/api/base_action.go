package api

import (
	"context"
	"fmt"

	"github.com/bioform/go-web-app-template/pkg/action"
	"gorm.io/gorm"
)

type BaseAction struct {
	action.BaseAction
	api API
}

func (ba *BaseAction) SetContext(ctx context.Context) {
	ba.BaseAction.SetContext(ctx)

	api, err := From(ctx)
	if err != nil {
		panic(fmt.Errorf("set api: %w", err))
	}

	ba.api = api
}

func (ba *BaseAction) TransactionProvider() action.TransactionProvider {
	return ba.api
}

func (ba BaseAction) API() API {
	return ba.api
}

func (ba BaseAction) DB() (*gorm.DB, error) {
	return ba.api.DB(), nil
}
