package action

import (
	"context"
)

type Action interface {
	SetContext(context.Context)
	Context() context.Context
	TransactionProvider() TransactionProvider
	Performer() any
	Perform(context.Context) error
	IsAllowed(context.Context) (bool, error)
	IsEnabled(context.Context) (bool, error)
	IsValid(context.Context) (bool, error)
}
