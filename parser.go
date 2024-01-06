package magicjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func parse(t reflect.Type, v reflect.Value, canMarshal bool) ([]byte, error) {
	if !canMarshal {
		return defaultParser(t, v)
	}
	return jsonMarshal(t, v)
}

func jsonMarshal(t reflect.Type, v reflect.Value) ([]byte, error) {
	vPtr := reflect.NewAt(t.Elem(), v.UnsafePointer())
	results := vPtr.Elem().MethodByName("MarshalJSON").Call([]reflect.Value{})

	if len(results) != 2 {
		return nil, fmt.Errorf("invalid length of results")
	}

	b, ok := results[0].Interface().([]byte)
	if !ok {
		return nil, fmt.Errorf("expected []byte type, got %s", results[0].Type())
	}

	if results[1].Interface() == nil {
		return b, nil
	}

	err, ok := results[1].Interface().(error)
	if !ok {
		return nil, fmt.Errorf("expected error type, got %s", results[1].Type())
	}

	return b, err
}

func defaultParser(t reflect.Type, v reflect.Value) ([]byte, error) {
	result, err := resolve(t, v)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func resolve(t reflect.Type, v reflect.Value) (any, error) {
	vKind := v.Kind()
	s := fmt.Sprint(reflect.Indirect(v))

	if isBytes(t) {
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

func matchKind(kind reflect.Kind, kinds ...reflect.Kind) bool {
	for i := range kinds {
		if kind == kinds[i] {
			return true
		}
	}
	return false
}
