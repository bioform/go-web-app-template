package action

type option int

const (
	SkipCache option = iota
	SkipTransaction
	NopIfDisabled
)

type Options struct {
	SkipCache       bool
	SkipTransaction bool
	NopIfDisabled   bool
}

func options(opts []option) Options {
	var o Options
	for _, opt := range opts {
		switch opt {
		case SkipCache:
			o.SkipCache = true
		case SkipTransaction:
			o.SkipTransaction = true
		case NopIfDisabled:
			o.NopIfDisabled = true
		}
	}
	return o
}
