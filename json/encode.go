package json

import (
	"encoding/json"
	"io"
)

// EncodeOption is a type for functional options for the Encode function.
type EncodeOption struct {
	EscapeHTML   bool
	IndentPrefix string
	IndentValue  string
}

// Encode encodes data to the writer with json.
func Encode[S any](writer io.Writer, data S, opts ...func(*EncodeOption)) error {
	var opt EncodeOption
	for _, fn := range opts {
		fn(&opt)
	}

	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(opt.EscapeHTML)
	encoder.SetIndent(opt.IndentPrefix, opt.IndentValue)
	return encoder.Encode(data)
}
