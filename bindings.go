package binding

import (
	"fmt"
	"strconv"
	"strings"
)

// Bindings is a map of binding function to it's name in `binding` tag.
type Bindings map[string]BindFunc

// BindFunc is a binding function signature which is used as parser for every
// mapped value.
//
// First argument is mapped value to be parsed. Only strings are supported for
// now.
//
// Second argument is optional argument string that can control binding
// function execution (like set bitness for ints), which is specified after
// `:` char in the `binding` tag.
type BindFunc func(interface{}, string) (interface{}, error)

func bindInt(data interface{}, opts string) (interface{}, error) {
	var (
		bits = 0
		base = 10
	)

	_, err := fmt.Sscanf(opts, "%d,%d", &bits, &base)
	if err != nil && !strings.HasSuffix(err.Error(), "EOF") {
		return nil, InvalidBindingError(err.Error())
	}

	if _, ok := data.(string); !ok {
		return nil, InvalidBindingError(
			fmt.Sprintf("only strings are supported, but %T given", data),
		)
	}

	result, err := strconv.ParseInt(data.(string), base, bits)
	if err != nil {
		return nil, err
	}

	switch bits {
	case 8:
		return int8(result), nil
	case 16:
		return int16(result), nil
	case 32:
		return int32(result), nil
	case 64:
		return int64(result), nil
	default:
		return int(result), nil
	}
}

func bindFloat(data interface{}, opts string) (interface{}, error) {
	var (
		bits = 32
	)

	_, err := fmt.Sscanf(opts, "%d", &bits)
	if err != nil && !strings.HasSuffix(err.Error(), "EOF") {
		return nil, InvalidBindingError(err.Error())
	}

	if _, ok := data.(string); !ok {
		return nil, InvalidBindingError(
			fmt.Sprintf("only strings are supported, but %T given", data),
		)
	}

	result, err := strconv.ParseFloat(data.(string), bits)
	if err != nil {
		return nil, err
	}

	switch bits {
	case 32:
		return float32(result), nil
	case 64:
		return float64(result), nil
	default:
		return float32(result), nil
	}
}

func bindString(data interface{}, _ string) (interface{}, error) {
	return data, nil
}
