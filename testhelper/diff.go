package helper

import (
	"strings"

	"github.com/google/go-cmp/cmp"
)

// Equal compares got and want and returns a formatted error message and a bool
func Equal(got, want any, opts ...cmp.Option) (string, bool) {
	if cmp.Equal(got, want, opts...) {
		return "", true
	}
	var messages []string
	messages = append(messages, "Diff(-got,+want):")
	for _, line := range strings.Split(cmp.Diff(got, want, opts...), "\n") {
		if line == "" {
			continue
		}
		messages = append(messages, "\t"+line)
	}
	return strings.Join(messages, "\n"), false
}
