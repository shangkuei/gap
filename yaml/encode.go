package yaml

import (
	"io"

	"github.com/goccy/go-yaml"
)

// EncodeOption is a type for functional options for the Encode function.
type EncodeOption struct {
	Indent int
}

// Encode encodes data to the writer with yaml.
func Encode[S any](writer io.Writer, data S, opts ...func(*EncodeOption)) error {
	var opt EncodeOption
	for _, fn := range opts {
		fn(&opt)
	}

	encoder := yaml.NewEncoder(writer, yaml.Indent(opt.Indent))
	return encoder.Encode(data)
}
