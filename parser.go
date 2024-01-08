package magicjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func parse(v reflect.Value, marshaller bool) ([]byte, error) {
	if !marshaller {
		return defaultParser(v)
	}
	return marshalJSON(v)
}

func marshalJSON(v reflect.Value) ([]byte, error) {
	vPtr := reflect.NewAt(v.Elem().Type(), v.UnsafePointer())
	results := vPtr.Elem().MethodByName("MarshalJSON").Call([]reflect.Value{})

	bytes, _ := results[0].Interface().([]byte)
	err, _ := results[1].Interface().(error)

	return bytes, err
}

func defaultParser(v reflect.Value) ([]byte, error) {
	result, err := resolve(v)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func resolve(v reflect.Value) (any, error) {
	vKind := v.Kind()
	s := fmt.Sprint(reflect.Indirect(v))

	if isBytes(v.Type()) {
		return json.Marshal(s)
	}

	if matchKind(vKind,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64) {
		return strconv.Atoi(s)
	}

	if matchKind(vKind,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64) {
		return strconv.ParseUint(s, 10, 64)
	}

	if matchKind(vKind, reflect.Float32, reflect.Float64) {
		return strconv.ParseFloat(s, 64)
	}

	if matchKind(vKind, reflect.Complex64, reflect.Complex128) {
		return strconv.ParseComplex(s, 128)
	}

	return s, nil
}

func isBytes(t reflect.Type) bool {
	return fmt.Sprint(t) == `[]uint8`
}

func isEmptySlice(anything any) bool {
	v := reflect.ValueOf(anything)

	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return v.Len() == 0
	}

	return false
}

func matchKind(kind reflect.Kind, kinds ...reflect.Kind) bool {
	for i := range kinds {
		if kind == kinds[i] {
			return true
		}
	}
	return false
}
