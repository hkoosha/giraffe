package queryimpl

import (
	"strings"

	"github.com/hkoosha/giraffe/qcmd"
)

// TODO: no 'overwrite' and 'maybe' at the same time.
// TODO: better conflicting cmd checks.

var commands = qcmd.All()

func Escaped(
	spec string,
) string {
	sb := strings.Builder{}
	sb.Grow(len(spec))

	for _, c := range spec {
		if _, ok := commands[c]; ok {
			sb.WriteRune(qcmd.Escape)
		}

		sb.WriteRune(c)
	}

	return spec
}
