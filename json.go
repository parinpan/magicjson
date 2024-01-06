package magicjson

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/tidwall/sjson"
)

func Marshal(anything any) (payload []byte, err error) {
	t := reflect.TypeOf(anything)

	if anything == nil || isBytes(t) {
		return json.Marshal(anything)
	}

	err = walk(anything, func(t reflect.Type, v reflect.Value, canMarshal bool, path string) error {
		b, err := parse(t, v, canMarshal)
		if err != nil {
			return err
		}

		if len(path) == 0 {
			payload = b
			return nil
		}

		if payload, err = sjson.SetRawBytes(payload, path, b); err != nil {
			return err
		}

		return nil
	})

	return payload, err
}

func isBytes(t reflect.Type) bool {
	return fmt.Sprint(t) == `[]uint8`
}
