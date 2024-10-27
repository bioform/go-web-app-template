package util

import "fmt"

func WrapError(errp *error, format string, args ...any) {
	if *errp != nil {
		s := fmt.Sprintf(format, args...)
		*errp = fmt.Errorf("%s: %w", s, *errp)
	}
}
