package toml

import (
	"io"

	"github.com/pelletier/go-toml/v2"
)

// EncodeOption is a type for functional options for the Encode function.
type EncodeOption struct {
	TablesInline    bool
	IndentSymbol    string
	IndentTables    bool
	ArraysMultiline bool
}

// Encode encodes data to the writer with toml.
func Encode[S any](writer io.Writer, data S, opts ...func(*EncodeOption)) error {
	opt := EncodeOption{IndentSymbol: "  "}
	for _, fn := range opts {
		fn(&opt)
	}

	encoder := toml.NewEncoder(writer)
	encoder.SetTablesInline(opt.TablesInline)
	encoder.SetIndentSymbol(opt.IndentSymbol)
	encoder.SetIndentTables(opt.IndentTables)
	encoder.SetArraysMultiline(opt.ArraysMultiline)
	return encoder.Encode(data)
}
