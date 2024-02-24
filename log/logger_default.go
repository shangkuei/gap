//go:build !windows
// +build !windows

package log

import (
	"io"
	"os"
)

func stderr() io.Writer {
	return os.Stderr
}

func stdout() io.Writer {
	return os.Stdout
}
