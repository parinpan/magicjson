package magicjson

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func walk(anything any, cb func(t reflect.Type, v reflect.Value, canMarshal bool, path string) error) error {
	return walker(reflect.TypeOf(anything), reflect.ValueOf(anything), "", cb)
}

func walker(t reflect.Type, v reflect.Value, path string, cb func(t reflect.Type, v reflect.Value, canMarshal bool, path string) error) error {
	if v.CanInterface() && isMarshaler(v) && t.Kind() != reflect.Ptr {
		vPtr := reflect.New(t)
		vPtr.Elem().Set(v)
		return cb(vPtr.Type(), vPtr, true, path)
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
		if canMarshal(t.Elem()) {
			return cb(t, v, true, path)
		}
		return walker(t.Elem(), v.Elem(), path, cb)
	default:
		return cb(t, v, false, path)
	}

	return nil
}

func canMarshal(t reflect.Type) bool {
	v := reflect.New(t)
	return isMarshaler(v)
}

func isMarshaler(v reflect.Value) bool {
	_, ok := v.Interface().(json.Marshaler)
	return ok
}

func addPath(path, suffix string) string {
	if len(path) > 0 {
		return fmt.Sprint(path, ".", suffix)
	}
	return suffix
}
