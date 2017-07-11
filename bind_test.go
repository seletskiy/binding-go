package binding

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBind_CanBindStringIdentically(t *testing.T) {
	test := assert.New(t)

	var user struct {
		Name string
	}

	err := Bind(&user, func(string) interface{} {
		return "John Doe"
	})

	test.NoError(err)
	test.Equal("John Doe", user.Name)
}

func TestBind_CanBindInts(t *testing.T) {
	test := assert.New(t)

	var user struct {
		Age   int
		Age8  int8
		Age16 int16
		Age32 int32
		Age64 int64
	}

	err := Bind(&user, func(key string) interface{} {
		switch key {
		case "Age":
			return fmt.Sprint(math.MaxInt32)
		case "Age8":
			return fmt.Sprint(math.MaxInt8)
		case "Age16":
			return fmt.Sprint(math.MaxInt16)
		case "Age32":
			return fmt.Sprint(math.MaxInt32)
		case "Age64":
			return fmt.Sprint(math.MaxInt64)
		default:
			return "XXX"
		}
	})

	test.NoError(err)
	test.Equal(int8(math.MaxInt8), user.Age8)
	test.Equal(int16(math.MaxInt16), user.Age16)
	test.Equal(int32(math.MaxInt32), user.Age32)
	test.Equal(int64(math.MaxInt64), user.Age64)
}

func TestBind_CanBindFloats(t *testing.T) {
	test := assert.New(t)

	var point struct {
		Distance32 float32
		Distance64 float64
	}

	err := Bind(&point, func(key string) interface{} {
		switch key {
		case "Distance32":
			return fmt.Sprint(math.MaxFloat32)
		case "Distance64":
			return fmt.Sprint(math.MaxFloat64)
		default:
			return "XXX"
		}
	})

	test.NoError(err)
	test.Equal(float32(math.MaxFloat32), point.Distance32)
	test.Equal(float64(math.MaxFloat64), point.Distance64)
}

func TestBind_CanUseCustomFieldName(t *testing.T) {
	test := assert.New(t)

	var profile struct {
		UserAge int `form:"user_age"`
	}

	err := Bind(&profile, func(key string) interface{} {
		switch key {
		case "user_age":
			return "88"
		default:
			return "XXX"
		}
	})

	test.NoError(err)
	test.Equal(88, profile.UserAge)
}

func TestBind_CanCheckRequiredFields(t *testing.T) {
	test := assert.New(t)

	var user struct {
		Age    int    `required:"true"`
		Name   string `required:"true"`
		Height int
	}

	err := Bind(&user, func(key string) interface{} {
		return nil
	})

	test.Equal(BindingErrors{RequiredError{"Age"}, RequiredError{"Name"}}, err)
	test.NotNil(err.(BindingErrors).Field("Age"))
	test.NotNil(err.(BindingErrors).Field("Name"))
	test.Nil(err.(BindingErrors).Field("Height"))

	test.Equal(0, user.Age)
	test.Equal("", user.Name)
	test.Equal(0, user.Height)
}

func TestBind_PreservesAlreadySetValues(t *testing.T) {
	test := assert.New(t)

	var user struct {
		Age  int
		Name string
	}

	user.Age = 1
	user.Name = "John Doe Jr."

	err := Bind(&user, func(key string) interface{} {
		return nil
	})

	test.NoError(err)
	test.Equal(1, user.Age)
	test.Equal("John Doe Jr.", user.Name)
}

func TestBind_CanUseCustomBindFunc(t *testing.T) {
	test := assert.New(t)

	var contract struct {
		ExpiresIn time.Duration `binding:"duration"`
	}

	var bindDuration = func(data interface{}, _ string) (interface{}, error) {
		return time.ParseDuration(data.(string))
	}

	err := Bind(&contract, func(key string) interface{} {
		return "1h30m"
	}, Bindings{"duration": bindDuration})

	test.NoError(err)
	test.Equal("1h30m0s", contract.ExpiresIn.String())
}

func TestBin_CanUseCustomFieldNameFunc(t *testing.T) {
	test := assert.New(t)

	var user struct {
		Name string
		Age  int `name:"age"`
	}

	err := Bind(&user, func(key string) interface{} {
		return "27"
	}, FieldNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("name")
	}))

	test.NoError(err)
	test.Empty(user.Name)
	test.Equal(27, user.Age)
}
