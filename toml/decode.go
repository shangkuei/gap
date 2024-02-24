package toml

import (
	"io"

	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml/v2"
)

// Decode decodes yaml encoded data from the reader and stores the result in the value pointed to by result.
func Decode[S any](reader io.Reader, result *S, hooks ...mapstructure.DecodeHookFunc) error {
	var data any
	if err := toml.NewDecoder(reader).Decode(&data); err != nil {
		return err
	}

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(hooks...),
		Result:     result,
	})
	return decoder.Decode(data)
}
