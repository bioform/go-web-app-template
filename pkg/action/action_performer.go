package action

import (
	"context"
	"errors"
	"fmt"
	"slices"
)

type ActionPerformer[A Action] struct {
	action A

	callbacks   []AfterCommitCallback
	addCallback AddCallbackFunc

	isAllowedCache *Cache
	isEnabledCache *Cache
}

func New[A Action](ctx context.Context, action A) *ActionPerformer[A] {
	ap := &ActionPerformer[A]{action: action}

	action.SetContext(ap.setAddCallback(ctx))
	action.Init()

	return ap
}

func (ap *ActionPerformer[A]) As(performer Performer) *ActionPerformer[A] {
	ap.action.SetPerformer(performer)
	return ap
}

func (ap *ActionPerformer[A]) AsSystem() *ActionPerformer[A] {
	ap.action.SetPerformer(SystemPerformer)
	return ap
}

func (ap *ActionPerformer[A]) Action() A {
	return ap.action
}

// Perform executes the action within a transaction context.
func (ap *ActionPerformer[A]) Perform() (bool, error) {
	return ap.perform()
}

// Try executes the action but does not return an error if the action is disabled.
func (ap *ActionPerformer[A]) Try() (bool, error) {
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
func (ap *ActionPerformer[A]) perform(opts ...option) (bool, error) {
	fn := func() (bool, error) {
		if ok, err := ap.IsPerformable(append(opts, SkipTransaction)...); !ok || err != nil {
			return ok, err
		}
		if ok, err := ap.checkValid(); !ok || err != nil {
			return ok, err
		}
		if err := ap.action.Perform(); err != nil {
			return false, err
		}
		return true, nil
	}
	ok, err := ap.transaction(fn)

	if ok && err == nil {
		if errs := ap.AfterCommit(); len(errs) > 0 {
			ok = false
			err = NewActionError(ap.action, fmt.Errorf("after commit: %w", errors.Join(errs...)))
		}
	}

	err = ap.wrapError(err)

	return ok, err
}

// IsAllowed checks if the action is allowed and caches the result.
func (ap *ActionPerformer[A]) IsAllowed(opts ...option) (bool, error) {
	o := options(opts)

	if !o.SkipCache && ap.isAllowedCache != nil {
		return ap.isAllowedCache.ok, ap.isAllowedCache.err
	}
	fn := func() (bool, error) {
		ok, err := ap.checkAllowed()
		ap.isAllowedCache = &Cache{ok: ok, err: err}
		return ok, err
	}

	if o.SkipTransaction {
		return fn()
	}
	return ap.transaction(fn)
}

// IsEnabled checks if the action is enabled and caches the result.
func (ap *ActionPerformer[A]) IsEnabled(opts ...option) (bool, error) {
	o := options(opts)

	if !o.SkipCache && ap.isEnabledCache != nil {
		return ap.isEnabledCache.ok, ap.isEnabledCache.err
	}
	fn := func() (bool, error) {
		ok, err := ap.checkEnabled(opts...)
		ap.isEnabledCache = &Cache{ok: ok, err: err}
		return ok, err
	}
	if o.SkipTransaction {
		return fn()
	}
	return ap.transaction(fn)
}

// IsPerformable checks if the action is both allowed and enabled.
func (ap *ActionPerformer[A]) IsPerformable(opts ...option) (bool, error) {
	if ok, err := ap.IsAllowed(opts...); !ok || err != nil {
		return false, err
	}
	if ok, err := ap.IsEnabled(opts...); !ok || err != nil {
		return false, err
	}
	return true, nil
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

func (ap *ActionPerformer[A]) checkEnabled(opts ...option) (bool, error) {
	ok, err := ap.action.IsEnabled()
	if !ok {
		if errMap, ok := err.(ErrorMap); ok {
			if slices.Contains(opts, NopIfDisabled) {
				return false, nil
			}
			return false, NewDisabledError(ap.action, errMap)
		}
		return false, err
	}
	return true, nil
}

func (ap *ActionPerformer[A]) checkValid() (bool, error) {
	ok, err := ap.action.IsValid()
	if !ok {
		if errMap, ok := err.(ErrorMap); ok {
			return false, NewValidationError(ap.action, errMap)
		}
		return false, err
	}
	return true, nil
}

func (ap *ActionPerformer[A]) transaction(fn func() (bool, error)) (ok bool, err error) {
	tp := ap.action.TransactionProvider()
	ctx := ap.action.Context()

	err = tp.Transaction(ctx, func(txCtx context.Context) error {
		ap.action.SetContext(txCtx)
		defer ap.action.SetContext(ctx)

		ok, err = fn()
		return err
	})
	return ok, err
}

func (ap *ActionPerformer[A]) wrapError(err error) error {
	if err != nil {
		handledError := ap.action.ErrorHandler(err)
		if handledError != err {
			return handledError
		}

		var actionErr ActionError
		if ok := errors.As(err, &actionErr); ok {
			var action Action = ap.action
			if actionErr.action != action {
				err = NewActionError(ap.action, err)
			}
		} else {
			err = NewActionError(ap.action, err)
		}
	}
	return err
}
