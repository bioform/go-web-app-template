package action

import (
	"context"
	"slices"
)

type option int

const (
	SkipCache option = iota
	SkipTransaction
	NopIfDisabled
)

type ActionPerformer[A Action] struct {
	action         A
	isAllowedCach  *Cache
	IsEnabledCache *Cache
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

// Perform action in a transaction context.
// Returns: ok: A boolean indicating whether the action was successfully performed.
func (ap *ActionPerformer[A]) Perform() (ok bool, err error) {
	return ap.perform()
}

// Try performs the action in a transaction context, but does not return an
// error if the action is disabled. It returns a boolean indicating whether the
// action was successfully performed.
// Returns: ok: A boolean indicating whether the action was successfully performed.
func (ap *ActionPerformer[A]) Try() (ok bool, err error) {
	return ap.perform(NopIfDisabled)
}

// perform executes the action within a transaction context provided by the
// TransactionProvider. It first checks if the action is enabled and valid
// before performing it. If the action is disabled and NopIfDisabled option is provided,
// it will proceed without error.
//
// Parameters:
//   - ctx: The context for the operation, used for managing request-scoped
//     values (DB reference), cancellation, and deadlines.
//   - NopIfDisabled: A boolean flag indicating whether the action should proceed
//     even if it is disabled.
//
// Returns:
//   - ok: A boolean indicating whether the action was successfully performed.
//   - err: An error if any occurred during the transaction or action execution.
func (ap *ActionPerformer[A]) perform(opts ...option) (ok bool, err error) {
	tp := ap.action.TransactionProvider()
	ctx := ap.action.Context()

	err = tp.Transaction(ctx, func(transactionContext context.Context) error {
		ap.action.SetContext(transactionContext)
		defer ap.action.SetContext(ctx)

		// Check if the action is allowed and enabled.
		if ok, err = ap.IsPerformable(opts...); !ok || err != nil {
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

/*
 * The following methods are used to check if an action is allowed, enabled, or
 * performable. They cache the results of these checks to avoid redundant
 * computation.
 */

func (ap *ActionPerformer[A]) IsAllowed(opts ...option) (bool, error) {
	if !slices.Contains(opts, SkipCache) {
		cache := ap.isAllowedCach
		if cache != nil {
			// Return the cached result if it exists.
			return cache.ok, cache.err
		}
	}

	// Otherwise, check if the action is allowed.
	ok, err := ap.checkAllowed()
	ap.isAllowedCach = &Cache{
		ok:  ok,
		err: err,
	}
	return ok, err
}

func (ap *ActionPerformer[A]) IsEnabled(opts ...option) (bool, error) {
	if !slices.Contains(opts, SkipCache) {
		cache := ap.IsEnabledCache
		if cache != nil {
			// Return the cached result if it exists.
			return cache.ok, cache.err
		}
	}

	// Otherwise, check if the action is enabled.
	ok, err := ap.checkEnabled(opts...)
	ap.IsEnabledCache = &Cache{
		ok:  ok,
		err: err,
	}
	return ok, err
}

func (ap *ActionPerformer[A]) IsPerformable(opts ...option) (bool, error) {
	allowed, err := ap.IsAllowed(opts...)
	if !allowed || err != nil {
		return false, err
	}

	enabled, err := ap.IsEnabled(opts...)
	if !enabled || err != nil {
		return false, err
	}

	return true, nil
}

/*
 * The following methods are used to check if an action is allowed, enabled, or
 * valid. They don't use caching and perform the checks directly.
 * computation.
 */

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

func (ap *ActionPerformer[A]) checkEnabled(opts ...option) (bool, error) {
	ok, err := ap.action.IsEnabled()
	if ok {
		return true, nil
	}

	errMap, ok := err.(ErrorMap)
	if !ok {
		return false, err
	}

	if slices.Contains(opts, NopIfDisabled) {
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
