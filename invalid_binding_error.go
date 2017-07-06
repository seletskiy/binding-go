package binding

// InvalidBindingError represents malformed binding specification which
// is not caused by malformed mapped value.
type InvalidBindingError string

func (err InvalidBindingError) Error() string {
	return string(err)
}
