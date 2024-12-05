package action

import (
	"context"
)

type Action interface {
	SetContext(context.Context)
	Context() context.Context
	TransactionProvider() TransactionProvider
	Performer() any
	Perform() error
	IsAllowed() (bool, error)
	IsEnabled() (bool, error)
	IsValid() (bool, error)
	AfterCommitCallback() AfterCommitCallback
}

type AfterCommitCallback func(context.Context, Action) error
