package action

import (
	"context"
	"log/slog"
)

type ActionPerformer struct {
	action Action
}

func New(action Action) *ActionPerformer {
	if action.TransactionProvider() == nil {
		action.SetTransactionProvider(transactionProvider)
	}

	return &ActionPerformer{
		action: action,
	}
}

func (ap *ActionPerformer) Perform(ctx context.Context) (ok bool, err error) {
	return ap.perform(ctx, false)
}

func (ap *ActionPerformer) PerformIfEnabled(ctx context.Context) (ok bool, err error) {
	return ap.perform(ctx, true)
}

// perform executes the action within a transaction context provided by the
// TransactionProvider. It first checks if the action is enabled and valid
// before performing it. If the action is disabled and okIfDisabled is true,
// it will proceed without error.
//
// Parameters:
//   - ctx: The context for the operation, used for managing request-scoped
//     values (DB reference), cancellation, and deadlines.
//   - okIfDisabled: A boolean flag indicating whether the action should proceed
//     even if it is disabled.
//
// Returns:
//   - ok: A boolean indicating whether the action was successfully performed.
//   - err: An error if any occurred during the transaction or action execution.
func (ap *ActionPerformer) perform(ctx context.Context, okIfDisabled bool) (ok bool, err error) {
	provider := ap.action.TransactionProvider()

	err = provider.Transaction(ctx, func(transactionContext context.Context) error {
		slog.Debug("Shared logic before perform")

		if ok, err = ap.checkEnabled(transactionContext, okIfDisabled); !ok || err != nil {
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

func (ap *ActionPerformer) checkEnabled(ctx context.Context, okIfDisabled bool) (bool, error) {
	ok, errMap := ap.action.IsEnabled(ctx)
	if !ok {
		return okIfDisabled, NewDisabledError(errMap)
	}
	return true, nil
}

func (ap *ActionPerformer) checkValid(ctx context.Context) (bool, error) {
	ok, errMap := ap.action.IsValid(ctx)
	if !ok {
		return false, NewValidationError(errMap)
	}
	return true, nil
}
