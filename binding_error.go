package binding

import (
	"fmt"
)

// BindingError will be part of BindingErrors slice to describe binding error
// of specific field.
type BindingError struct {
	name  string
	cause error
}

func (err BindingError) Name() string {
	return err.name
}

func (err BindingError) Cause() error {
	return err.cause
}

func (err BindingError) Error() string {
	return fmt.Sprintf(
		`%s — %s`,
		err.Name(),
		err.Cause(),
	)
}
