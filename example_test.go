package binding_test

import (
	"fmt"
	"time"

	"github.com/seletskiy/binding-go"
)

func Example_simpleBinding() {
	var user struct {
		Age int
	}

	binding.Bind(&user, func(key string) interface{} {
		switch key {
		case "Age":
			return "27"
		default:
			return "??"
		}
	})

	fmt.Printf("Age: %d\n", user.Age)

	// Output:
	// Age: 27
}

func Example_customFieldName() {
	var user struct {
		Age int `form:"how_old"`
	}

	binding.Bind(&user, func(key string) interface{} {
		switch key {
		case "how_old":
			return "27"
		default:
			return "??"
		}
	})

	fmt.Printf("Age: %d\n", user.Age)

	// Output:
	// Age: 27
}

func Example_customBindingFunction() {
	var contract struct {
		Duration time.Duration `binding:"duration"`
	}

	var bindDuration = func(data interface{}, _ string) (interface{}, error) {
		return time.ParseDuration(data.(string))
	}

	binding.Bind(&contract, func(key string) interface{} {
		return "1h23m45s"
	}, binding.Bindings{"duration": bindDuration})

	fmt.Printf("Duration: %s\n", contract.Duration.String())

	// Output:
	// Duration: 1h23m45s
}

func Example_perFieldErrors() {
	var user struct {
		Age    int
		Name   string `required:"true"`
		Height int
	}

	err := binding.Bind(&user, func(key string) interface{} {
		switch key {
		case "Age":
			// return invalid int
			return "???"
		case "Name":
			// return no data for name
			return nil
		}

		return nil
	})

	errors := err.(binding.BindingErrors)

	fmt.Printf("Errors (%d):\n", len(errors))
	fmt.Printf("* %s\n", errors[0])
	fmt.Printf("* %s\n", errors[1])
	fmt.Printf("Age Error: %T\n", errors.Field("Age"))
	fmt.Printf("Name Error: %T\n", errors.Field("Name"))
	fmt.Printf("Height Error: %T\n", errors.Field("Height"))

	// Output:
	// Errors (2):
	// * Age — strconv.ParseInt: parsing "???": invalid syntax
	// * Name — field required but not specified
	// Age Error: binding.BindingError
	// Name Error: binding.RequiredError
	// Height Error: <nil>
}
