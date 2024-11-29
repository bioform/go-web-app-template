package action

import (
	"context"
)

type Action interface {
	SetTransactionProvider(TransactionProvider)
	TransactionProvider() TransactionProvider
	Perform(context.Context) error
	IsAllowed(context.Context) (bool, error)
	IsEnabled(context.Context) (bool, ErrorMap)
	IsValid(context.Context) (bool, ErrorMap)
}
