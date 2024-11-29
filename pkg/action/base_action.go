package action

import (
	"context"
	//validator "github.com/rezakhademix/govalidator/v2"
)

// Define the base type with the shared template method.
type BaseAction struct {
	transactionProvider TransactionProvider
}

func (ba *BaseAction) SetTransactionProvider(provider TransactionProvider) {
	ba.transactionProvider = provider
}

func (ba *BaseAction) TransactionProvider() TransactionProvider {
	return ba.transactionProvider
}

func (ba *BaseAction) IsAllowed(ctx context.Context) (bool, error) {
	// See https://github.com/rezakhademix/govalidator
	// For example:
	// v := validator.New()
	// return v.IsPassed(), v.Errors()
	return true, nil
}

func (ba *BaseAction) IsEnabled(ctx context.Context) (bool, ErrorMap) {
	// See https://github.com/rezakhademix/govalidator
	// For example:
	// v := validator.New()
	// return v.IsPassed(), v.Errors()
	return true, nil
}

func (ba *BaseAction) IsValid(ctx context.Context) (bool, ErrorMap) {
	// See https://github.com/rezakhademix/govalidator
	// For example:
	// v := validator.New()
	// return v.IsPassed(), v.Errors()
	return true, nil
}
