package yaml

import (
	"io"

	"github.com/goccy/go-yaml"
	"github.com/mitchellh/mapstructure"
)

// Decode decodes yaml encoded data from the reader and stores the result in the value pointed to by result.
func Decode[S any](reader io.Reader, result *S, hooks ...mapstructure.DecodeHookFunc) error {
	var data any
	if err := yaml.NewDecoder(reader).Decode(&data); err != nil {
		return err
	}

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(hooks...),
		Result:     result,
	})
	return decoder.Decode(data)
}
