package action

import (
	"context"
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

func (ap *ActionPerformer[A]) Perform() (ok bool, err error) {
	return ap.perform(false)
}

func (ap *ActionPerformer[A]) PerformIfEnabled() (ok bool, err error) {
	return ap.perform(true)
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
func (ap *ActionPerformer[A]) perform(nopIfDisabled bool) (ok bool, err error) {
	tp := ap.action.TransactionProvider()
	ctx := ap.action.Context()

	err = tp.Transaction(ctx, func(transactionContext context.Context) error {
		ap.action.SetContext(transactionContext)
		defer ap.action.SetContext(ctx)

		if ok, err = ap.checkAllowed(); !ok || err != nil {
			return err
		}

		if ok, err = ap.checkEnabled(nopIfDisabled); !ok || err != nil {
			return err
		}

		if ok, err = ap.checkValid(); !ok || err != nil {
			return err
		}

		if err = ap.action.Perform(); err != nil {
			ok = false
			return err
		}

		return nil
	})

	return ok, err
}

func (ap *ActionPerformer[A]) checkAllowed() (bool, error) {
	ok, err := ap.action.IsAllowed()
	if err != nil {
		return false, err

	}
	if !ok {
		return false, NewAuthorizationError(ap.action)
	}
	return true, nil
}

func (ap *ActionPerformer[A]) checkEnabled(nopIfDisabled bool) (bool, error) {
	ok, err := ap.action.IsEnabled()
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

func (ap *ActionPerformer[A]) checkValid() (bool, error) {
	ok, err := ap.action.IsValid()
	if !ok {
		errMap, ok := err.(ErrorMap)
		if !ok {
			return false, err
		}
		return false, NewValidationError(ap.action, errMap)
	}
	return true, nil
}
