// Package binding offers easy way of binding form-like sources into structs.
//
// It's particularly useful with web-frameworks like gin.
//
// Package offers rich-structured errors which can be easily integrated into
// UI error reports (like HTML page).
package binding

import (
	"fmt"
	"reflect"
	"strings"
)

// FieldNameFunc represents function that retrieves field name by given
// reflect type of field.
type FieldNameFunc func(field reflect.StructField) string

// MapFunc is a signature for function that maps field name into raw
// representation. Only string return values are supported for now.
type MapFunc func(name string) interface{}

// Bind binds values provided by mapper function into output struct.
//
// In simplest use it will try to populate every exported by value, provided
// by mapper function converting it's return value from string to appropriate
// struct's field type.
//
// Additionally, struct's tags can be used to control binding. Following tags
// will be inspected by Bind function: `binding`, `form` and `required`.
//
// Tag `binding` used to override binding function which will be used for
// converting value returned by mapper function to struct's field type.
//
// There are three built-in functions: `int`, `float`, `string`. They used to
// parse mapped value into int, int8, int16, int32, int64, float32, float64 and
// string types accordingly.
//
// Binding `int` accepts two arguments in the form of `int:<bits>,<base>`,
// which are optional and can be used to override automatically detected
// bitness of resulting int and base of 10.
//
// Binding `float` accepts one argument in the form of `float:<bits>`.
//
// Binding `string` has no arguments and do not apply any parsing to mapped
// value.
//
// Tag `required` used to specify, that field should have mapped value and
// error will be reported otherwise. Tag should be specified as
// `required:"true"`.
//
// Tag `form` can be used to override field name that will be passed into
// mapper function to obtain value. Bind will also inspect `json`, `bson`,
// `yaml` and `toml` tags if `form` tag is not specified. If no known tags
// specify mapped name, then field's name will be used.
//
// To customize binding behavior, third variable argument can be used:
//
// To specify binding functions, pass functions in the form of
// `Bindings{"<name>": <function>}`.
//
// To specify function that maps field to it's name, specify it as
// `FieldNameFunc(<func>)`.
func Bind(output interface{}, mapper MapFunc, options ...interface{}) error {
	var bindings = Bindings{
		"int":    bindInt,
		"float":  bindFloat,
		"string": bindString,
	}

	var fieldNameFunc = getFieldName

	for _, option := range options {
		switch option := option.(type) {
		case Bindings:
			for key, binding := range option {
				bindings[key] = binding
			}
		case FieldNameFunc:
			fieldNameFunc = option
		}
	}

	if reflect.ValueOf(output).Kind() != reflect.Ptr {
		return InvalidBindingError("specified output is not a pointer")
	}

	var (
		structValue = reflect.Indirect(reflect.ValueOf(output))
		structType  = structValue.Type()
	)

	if structType.Kind() != reflect.Struct {
		return InvalidBindingError(
			fmt.Sprintf(
				`output should be struct type, but %s is given`,
				structType,
			),
		)
	}

	if !structValue.CanSet() {
		return InvalidBindingError(`output can not be set`)
	}

	var errors BindingErrors

	for i := 0; i < structType.NumField(); i++ {
		var (
			field = structType.Field(i)
			name  = fieldNameFunc(field)
		)

		if name == "" {
			continue
		}

		if binding, ok := getBinding(field, bindings); !ok {
			return InvalidBindingError(
				fmt.Sprintf(
					`binding for %s.%s is specified but not registered`,
					structType,
					field.Name,
				),
			)
		} else {
			data := mapper(name)

			if data == nil {
				if isRequired(field) {
					errors = append(errors, RequiredError{name: name})
				}

				continue
			}

			if _, ok := data.(string); !ok {
				return InvalidBindingError(
					fmt.Sprintf(
						`binding values of type %T (%s.%s) is not supported`,
						data,
						structType,
						field.Name,
					),
				)
			}

			value, err := binding(data.(string))
			if err != nil {
				errors = append(errors, BindingError{
					name:  name,
					cause: err,
				})

				continue
			}

			structField := structValue.Field(i)
			if !structField.CanSet() {
				return InvalidBindingError(
					fmt.Sprintf(
						`field %s.%s is unexported and can not be set`,
						structType.Name(),
						field.Name,
					),
				)
			}

			structField.Set(reflect.ValueOf(value))
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func getFieldName(field reflect.StructField) string {
	for _, key := range []string{"form", "json", "bson", "yaml", "toml"} {
		if name, ok := field.Tag.Lookup(key); ok {
			name = strings.Split(name, ",")[0]
			if name != "" {
				return name
			}
		}
	}

	return field.Name
}

func isRequired(field reflect.StructField) bool {
	value, ok := field.Tag.Lookup("required")

	return ok && value == "true"
}

func getBinding(
	field reflect.StructField,
	bindings map[string]BindFunc,
) (func(string) (interface{}, error), bool) {
	tag, _ := field.Tag.Lookup("binding")
	if tag == "" {
		tag = getDefaultBindingTag(field)
	}

	var (
		args = strings.SplitN(tag, ":", 2)
		name = args[0]
		opts = ""
	)

	if len(args) == 2 {
		opts = args[1]
	}

	if binding, ok := bindings[name]; ok {
		return func(data string) (interface{}, error) {
			return binding(data, opts)
		}, true
	} else {
		return nil, false
	}

	return nil, true
}

func getDefaultBindingTag(field reflect.StructField) string {
	var defaults = map[reflect.Kind]string{
		reflect.Int:   "int",
		reflect.Int8:  "int:8",
		reflect.Int16: "int:16",
		reflect.Int32: "int:32",
		reflect.Int64: "int:64",

		reflect.Float32: "float:32",
		reflect.Float64: "float:64",

		reflect.String: "string",
	}

	return defaults[field.Type.Kind()]
}
