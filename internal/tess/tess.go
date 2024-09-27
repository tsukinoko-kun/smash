package tess

import (
	"fmt"
	"reflect"
)

// Field unsafely accesses a struct field with the given name
func Field[T any](s T, name string) any {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		// maybe a pointer to a struct?
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			if v.Kind() != reflect.Struct {
				panic("Field: s must be a struct")
			}
		} else {
			panic("Field: s must be a struct")
		}
	}

	field := v.FieldByName(name)
	if !field.IsValid() {
		panic(fmt.Sprintf("Field: no such field %s in struct", name))
	}

	return field.Interface()
}

func UnsafeCast[T any](v any) T {
	if cs, ok := v.(T); ok {
		return cs
	}
	panic(fmt.Sprintf("UnsafeCast: failed to cast %T", v))
}
