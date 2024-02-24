//go:build windows
// +build windows

package log

import (
	"io"
	"os"

	"github.com/mattn/go-colorable"
)

func stderr() io.Writer {
	return colorable.NewColorable(os.Stderr)
}

func stdout() io.Writer {
	return colorable.NewColorable(os.Stdout)
}
