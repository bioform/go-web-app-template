package action

import (
	"context"
)

type Action interface {
	SetTransactionProvider(TransactionProvider)
	TransactionProvider() TransactionProvider
	Perform(context.Context) error
	IsAllowed(context.Context) (ok bool, errorMap ErrorMap)
	IsEnabled(context.Context) (ok bool, errorMap ErrorMap)
	IsValid(context.Context) (ok bool, errorMap ErrorMap)
}
