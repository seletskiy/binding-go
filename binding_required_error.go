package binding

import (
	"fmt"
)

type RequiredError struct {
	name string
}

func (err RequiredError) Name() string {
	return err.name
}

func (err RequiredError) Error() string {
	return fmt.Sprintf(
		`%s — field required but not specified`,
		err.Name(),
	)
}
