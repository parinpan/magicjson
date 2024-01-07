package magicjson

import (
	"encoding/json"
	"reflect"

	"github.com/tidwall/sjson"
)

func Marshal(anything any) (payload []byte, err error) {
	t := reflect.TypeOf(anything)

	if anything == nil || isBytes(t) || isEmptySlice(anything) {
		return json.Marshal(anything)
	}

	err = walk(anything, func(t reflect.Type, v reflect.Value, marshaller bool, path string) error {
		bytes, err := parse(t, v, marshaller)
		if err != nil {
			return err
		}

		if payload, err = sjson.SetRawBytes(payload, path, bytes); err != nil {
			return err
		}

		return nil
	})

	return payload, err
}
