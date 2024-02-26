package yaml

import (
	"encoding/json"
	"reflect"

	"github.com/goccy/go-yaml"
)

// DecodeYAMLUnmarshalFunc is a DecodeHookFunc that converts a map to a struct using json.Unmarshal.
func DecodeYAMLUnmarshalFunc(from reflect.Type, to reflect.Type, d interface{}) (interface{}, error) {
	if !reflect.PointerTo(to).Implements(reflect.TypeOf((*yaml.BytesUnmarshaler)(nil)).Elem()) {
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
	unmarshaller, ok := result.(yaml.BytesUnmarshaler)
	if !ok {
		return d, nil
	}

	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	if err := unmarshaller.UnmarshalYAML(data); err != nil {
		return nil, err
	}
	return result, nil
}
