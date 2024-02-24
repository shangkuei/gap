package yaml

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.Nil(t, Encode(&buf, tt.got))
			if err := Decode(&buf, &tt.object); !tt.wantErr {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.got, tt.object)
		})
	}
}
