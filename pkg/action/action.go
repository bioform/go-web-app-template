package action

import (
	"context"
)

type Action interface {
	SetTransactionProvider(TransactionProvider)
	TransactionProvider() TransactionProvider
	Perform(context.Context) error
	IsEnabled(context.Context) (ok bool, errorMap ErrorMap)
	IsValid(context.Context) (ok bool, errorMap ErrorMap)
}
