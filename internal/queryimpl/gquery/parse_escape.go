package gquery

import (
	"strings"

	"github.com/hkoosha/giraffe/qcmd"
)

var commands = qcmd.All()

func Escaped(
	spec string,
) string {
	sb := strings.Builder{}
	sb.Grow(len(spec))

	for _, c := range spec {
		// TODO cast fix
		if _, ok := commands[qcmd.Cmd(c)]; ok {
			sb.WriteByte(qcmd.Escape.Byte())
		}

		// TODO cast fix
		sb.WriteByte(byte(c))
	}

	return spec
}
