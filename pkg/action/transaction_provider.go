package action

import (
	"context"
)

var transactionProvider TransactionProvider

type TransactionProvider interface {
	Transaction(currentContext context.Context, executeInTransaction func(newContext context.Context) error) error
}

func SetTransactionProvider(provider TransactionProvider) {
	transactionProvider = provider
}
