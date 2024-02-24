package json

import (
	"encoding/json"
	"reflect"
)

// DecodeJSONUnmarshalFunc is a DecodeHookFunc that converts a map to a struct using json.Unmarshal.
func DecodeJSONUnmarshalFunc(from reflect.Type, to reflect.Type, d interface{}) (interface{}, error) {
	if !reflect.PointerTo(to).Implements(reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()) {
		return d, nil
	}

	var result any
	if to.Kind() == reflect.Slice {
		result = reflect.MakeSlice(to, 0, 0).Interface()
	} else if to.Kind() == reflect.Map {
		result = reflect.MakeMap(to).Interface()
	} else {
		result = reflect.New(to).Interface()
	}
	unmarshaller, ok := result.(json.Unmarshaler)
	if !ok {
		return d, nil
	}

	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	if err := unmarshaller.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return result, nil
}
