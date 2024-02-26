package yaml

import (
	"bytes"
	"fmt"
	"testing"

	helper "github.com/shangkuei/gap/testhelper"
)

type unmarshalerObject struct {
	Number int `yaml:"number" mapstructure:"number"`
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Encode(&buf, tt.got)
			if diff, ok := helper.Equal(err, error(nil)); !ok {
				t.Error(helper.Message(t, "unexpected error", diff))
			}
			err = Decode(&buf, &tt.object)
			if got := err != nil; got != tt.wantErr {
				t.Error(helper.Message(t, "unexpected error", fmt.Sprintf("Err: %v", got)))
			}
			if diff, ok := helper.Equal(tt.got, tt.object); !ok {
				t.Error(helper.Message(t, "unexpected object", diff))
			}
		})
	}
}
