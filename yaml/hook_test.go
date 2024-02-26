package yaml

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/goccy/go-yaml"
	helper "github.com/shangkuei/gap/testhelper"
)

type unmarshalerObject struct {
	Number int `json:"number" mapstructure:"number"`
}

func (u *unmarshalerObject) UnmarshalYAML(data []byte) error {
	var obj map[string]any
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return err
	}
	if number, ok := obj["number"].(uint64); !ok {
		return fmt.Errorf("number field is incorrect")
	} else {
		u.Number = int(number)
	}
	return nil
}

type unmarshalerSliceObject []string

func (u *unmarshalerSliceObject) UnmarshalYAML(data []byte) error {
	var obj []string
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return err
	}
	*u = obj
	return nil
}

type unmarshalerMapObject map[string]string

func (u *unmarshalerMapObject) UnmarshalYAML(data []byte) error {
	var obj map[string]string
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return err
	}
	*u = obj
	return nil
}

type unmarshalerErrorObject struct{}

func (u *unmarshalerErrorObject) UnmarshalYAML(data []byte) error {
	return errors.New("unmarshaler error")
}

func TestDecodeYAMLUnmarshalFunc(t *testing.T) {
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
			err := Encode(&buf, tt.got)
			if diff, ok := helper.Equal(err, error(nil)); !ok {
				t.Error(helper.Message(t, "unexpected error", diff))
			}
			err = Decode(&buf, &tt.object, DecodeYAMLUnmarshalFunc)
			if got := err != nil; got != tt.wantErr {
				t.Error(helper.Message(t, "unexpected error", fmt.Sprintf("Err: %v", got)))
			}
			if diff, ok := helper.Equal(tt.got, tt.object); !ok {
				t.Error(helper.Message(t, "unexpected object", diff))
			}
		})
	}
}
