package action

import "context"

type AddCallbackFunc func(callback AfterCommitCallback)
type addCallback string

var addCallbackKey addCallback = "addCallbackFunc"

func (ap *ActionPerformer[A]) AfterCommit() []error {
	var (
		errs []error
		act  Action          = ap.Action()
		ctx  context.Context = act.Context()
	)

	if callback := act.AfterCommitCallback(); callback != nil {
		ap.addCallback(callback)

		for _, callback := range ap.callbacks {
			if err := callback(ctx, act); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func (ap *ActionPerformer[A]) getAddCallbackFunc(ctx context.Context) (context.Context, AddCallbackFunc) {
	value := ctx.Value(addCallbackKey)

	if value != nil {
		if fn, ok := value.(AddCallbackFunc); ok {
			return ctx, fn
		}
	}

	fn := ap.addCallbackFunc
	// Set the root context value.
	ctx = context.WithValue(ctx, addCallbackKey, fn)
	return ctx, fn
}

func (ap *ActionPerformer[A]) addCallbackFunc(callback AfterCommitCallback) {
	ap.callbacks = append(ap.callbacks, callback)
}
