package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type unmarshalerObject struct {
	Number int `json:"number" mapstructure:"number"`
}

func (u *unmarshalerObject) UnmarshalJSON(data []byte) error {
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if number, ok := obj["number"].(float64); !ok {
		return fmt.Errorf("number field is incorrect")
	} else {
		u.Number = int(number)
	}
	return nil
}

type unmarshalerSliceObject []string

func (u *unmarshalerSliceObject) UnmarshalJSON(data []byte) error {
	var obj []string
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*u = obj
	return nil
}

type unmarshalerMapObject map[string]string

func (u *unmarshalerMapObject) UnmarshalJSON(data []byte) error {
	var obj map[string]string
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*u = obj
	return nil
}

type unmarshalerErrorObject struct{}

func (u *unmarshalerErrorObject) UnmarshalJSON(data []byte) error {
	return errors.New("unmarshaler error")
}

func TestDecodeJSONUnmarshalFunc(t *testing.T) {
	tests := []struct {
		name    string
		got     any
		object  any
		wantErr bool
	}{
		{
			name:   "happy path",
			got:    unmarshalerObject{Number: 1},
			object: unmarshalerObject{},
		},
		{
			name:   "slice",
			got:    unmarshalerSliceObject([]string{"arg1"}),
			object: unmarshalerSliceObject(nil),
		},
		{
			name:   "map",
			got:    unmarshalerMapObject(map[string]string{"arg1": "value"}),
			object: unmarshalerMapObject(nil),
		},
		{
			name:    "unmarshaler error",
			got:     unmarshalerErrorObject{},
			object:  unmarshalerErrorObject{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			assert.Nil(t, Encode(&buf, tt.got))
			if err := Decode(&buf, &tt.object, DecodeJSONUnmarshalFunc); !tt.wantErr {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.got, tt.object)
		})
	}
}
