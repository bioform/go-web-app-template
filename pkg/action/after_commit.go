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
	}

	for _, callback := range ap.callbacks {
		if err := callback(ctx, act); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (ap *ActionPerformer[A]) setAddCallback(ctx context.Context) context.Context {
	if fn, ok := ctx.Value(addCallbackKey).(AddCallbackFunc); ok {
		ap.addCallback = fn
	} else {
		ap.addCallback = func(callback AfterCommitCallback) {
			ap.callbacks = append(ap.callbacks, callback)
		}
		ctx = context.WithValue(ctx, addCallbackKey, ap.addCallback)
	}
	return ctx
}
