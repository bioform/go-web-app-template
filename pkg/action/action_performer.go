package action

import (
	"context"
	"log/slog"
)

type ActionPerformer[A Action] struct {
	action    A
	performer any
}

func New[A Action](action A) *ActionPerformer[A] {
	if action.TransactionProvider() == nil {
		action.SetTransactionProvider(transactionProvider)
	}

	return &ActionPerformer[A]{
		action: action,
	}
}

func (ap *ActionPerformer[A]) Action() A {
	return ap.action
}

// As sets the performer for the action.
func (ap *ActionPerformer[A]) As(performer any) *ActionPerformer[A] {
	ap.performer = performer
	return ap
}

func (ap *ActionPerformer[A]) Performer() any {
	return ap.performer
}

func (ap *ActionPerformer[A]) Perform(ctx context.Context) (ok bool, err error) {
	return ap.perform(ctx, false)
}

func (ap *ActionPerformer[A]) PerformIfEnabled(ctx context.Context) (ok bool, err error) {
	return ap.perform(ctx, true)
}

// perform executes the action within a transaction context provided by the
// TransactionProvider. It first checks if the action is enabled and valid
// before performing it. If the action is disabled and nopIfDisabled is true,
// it will proceed without error.
//
// Parameters:
//   - ctx: The context for the operation, used for managing request-scoped
//     values (DB reference), cancellation, and deadlines.
//   - nopIfDisabled: A boolean flag indicating whether the action should proceed
//     even if it is disabled.
//
// Returns:
//   - ok: A boolean indicating whether the action was successfully performed.
//   - err: An error if any occurred during the transaction or action execution.
func (ap *ActionPerformer[A]) perform(ctx context.Context, nopIfDisabled bool) (ok bool, err error) {
	provider := ap.action.TransactionProvider()

	err = provider.Transaction(ctx, func(transactionContext context.Context) error {
		slog.Debug("Shared logic before perform")

		if ok, err = ap.checkAllowed(transactionContext); !ok || err != nil {
			return err
		}

		if ok, err = ap.checkEnabled(transactionContext, nopIfDisabled); !ok || err != nil {
			return err
		}

		if ok, err = ap.checkValid(transactionContext); !ok || err != nil {
			return err
		}

		if err = ap.action.Perform(transactionContext); err != nil {
			ok = false
			return err
		}

		slog.Debug("Shared logic after perform")
		return nil
	})

	return ok, err
}

func (ap *ActionPerformer[A]) checkAllowed(ctx context.Context) (bool, error) {
	ok, err := ap.action.IsAllowed(ctx)
	if err != nil {
		return false, err

	}
	if !ok {
		return false, NewAuthorizationError(ap.performer)
	}
	return true, nil
}

func (ap *ActionPerformer[A]) checkEnabled(ctx context.Context, nopIfDisabled bool) (bool, error) {
	ok, errMap := ap.action.IsEnabled(ctx)
	if ok {
		return true, nil
	}
	if nopIfDisabled {
		return false, nil
	}
	return false, NewDisabledError(ap.Performer(), errMap)
}

func (ap *ActionPerformer[A]) checkValid(ctx context.Context) (bool, error) {
	ok, errMap := ap.action.IsValid(ctx)
	if !ok {
		return false, NewValidationError(ap.Performer(), errMap)
	}
	return true, nil
}
