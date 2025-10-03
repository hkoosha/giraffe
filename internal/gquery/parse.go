package gquery

import (
	"strings"

	"github.com/hkoosha/giraffe/internal/gquery/gqcmd"
)

// TODO: no 'overwrite' and 'maybe' at the same time.

// TODO: better conflicting cmd checks.

func Escaped(
	spec string,
) string {
	sb := strings.Builder{}
	sb.Grow(len(spec))

	for _, c := range spec {
		if _, ok := commands[c]; ok {
			sb.WriteRune(gqcmd.Escape)
		}

		sb.WriteRune(c)
	}

	return spec
}

func Parse(
	spec string,
) (Query, error) {
	return parse(spec)
}
