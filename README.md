# binding library [![](https://godoc.org/github.com/seletskiy/binding-go?status.svg)](http://godoc.org/github.com/seletskiy/binding-go)

Binding library offers easy way of binding form field (like HTML form) into structs.

Typical usage with gin framework:

```go
var user struct {
    Age int
}

err := binding.Bind(&user, func (key string) interface{} {
    return context.PostForm(key)
})

errs := err.(binding.BindingErrors)

if err = errs.Field("Age"); err != nil {
    // Age field in POST form is malformed
}
```

Difference to built-in gin `context.Bind()` function is that there are:

* custom binding function parsers;
* per-field binding errors, not just generic ParseInt errors.
