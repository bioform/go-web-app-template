package action

import (
	"context"
	"log/slog"
)

type ActionPerformer[A Action] struct {
	action A
}

func New[A Action](ctx context.Context, action A) *ActionPerformer[A] {
	action.SetContext(ctx)

	return &ActionPerformer[A]{
		action: action,
	}
}

func (ap *ActionPerformer[A]) Action() A {
	return ap.action
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
	tp := ap.action.TransactionProvider()

	err = tp.Transaction(ctx, func(transactionContext context.Context) error {
		defer ap.action.SetContext(ap.action.Context())

		ap.action.SetContext(transactionContext)

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
		return false, NewAuthorizationError(ap.action)
	}
	return true, nil
}

func (ap *ActionPerformer[A]) checkEnabled(ctx context.Context, nopIfDisabled bool) (bool, error) {
	ok, err := ap.action.IsEnabled(ctx)
	if ok {
		return true, nil
	}

	errMap, ok := err.(ErrorMap)
	if !ok {
		return false, err
	}

	if nopIfDisabled {
		return false, nil
	}
	return false, NewDisabledError(ap.action, ErrorMap(errMap))
}

func (ap *ActionPerformer[A]) checkValid(ctx context.Context) (bool, error) {
	ok, err := ap.action.IsValid(ctx)
	if !ok {
		errMap, ok := err.(ErrorMap)
		if !ok {
			return false, err
		}
		return false, NewValidationError(ap.action, errMap)
	}
	return true, nil
}
