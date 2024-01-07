package magicjson

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type callbackFn func(t reflect.Type, v reflect.Value, marshaller bool, path string) error

func walk(anything any, cb callbackFn) error {
	return walker(reflect.TypeOf(anything), reflect.ValueOf(anything), "", cb)
}

func walker(t reflect.Type, v reflect.Value, path string, cb callbackFn) error {
	// check if the value is the type of marshaler
	if isMarshaler(v) {
		ref := toRef(v)
		return cb(ref.Type(), ref, true, path)
	}

	// de-reference the value when it's a pointer - a value can be a type of marshaler
	if t.Kind() == reflect.Ptr && isMarshaler(v.Elem()) {
		return cb(t, v, true, path)
	}

	switch t.Kind() {
	case reflect.Struct:
		for idx := 0; idx < t.NumField(); idx++ {
			field := t.Field(idx)

			if err := walker(field.Type, v.Field(idx), addPath(path, field.Name), cb); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for idx := 0; idx < v.Len(); idx++ {
			if err := walker(v.Index(idx).Type(), v.Index(idx), addPath(path, fmt.Sprint(idx)), cb); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if err := walker(v.MapIndex(key).Type(), v.MapIndex(key), addPath(path, key.String()), cb); err != nil {
				return err
			}
		}
	case reflect.Ptr:
		return walker(t.Elem(), v.Elem(), path, cb)
	default:
		return cb(t, v, false, path)
	}

	return nil
}

func toRef(v reflect.Value) reflect.Value {
	ref := reflect.New(v.Type())
	ref.Elem().Set(v)
	return ref
}

func isMarshaler(v reflect.Value) bool {
	_, ok := reflect.New(v.Type()).Interface().(json.Marshaler)
	return ok
}

func addPath(path, suffix string) string {
	if len(path) > 0 {
		return fmt.Sprint(path, ".", suffix)
	}
	return suffix
}
