package binding

import (
	"strings"
)

// BindingErrors will be returned from Bind function if mapper provides values
// that can't be successfully bind to specified struct.
type BindingErrors []error

func (errors BindingErrors) Error() string {
	messages := []string{}

	for _, err := range errors {
		messages = append(messages, err.Error())
	}

	return strings.Join(messages, "; ")
}

// Field returns error for specific field name if any.
func (errors BindingErrors) Field(name string) error {
	for _, err := range errors {
		switch err := err.(type) {
		case RequiredError:
			if err.Name() == name {
				return err
			}
		case BindingError:
			if err.Name() == name {
				return err
			}
		}
	}

	return nil
}
