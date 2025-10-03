package qcmd

import (
	"maps"
)

const (
	Dialect   = '~'
	Pipe      = '|'
	Extern    = '^'
	Append    = '+'
	Overwrite = '='
	Delete    = '!'
	Make      = '$'
	Maybe     = '?'
	Sep       = '.'
	Escape    = '\\'
	At        = '@'
	Self      = '#'
	Move      = '>'
)

var all = map[rune]struct{}{
	Dialect:   {},
	Pipe:      {},
	Extern:    {},
	Append:    {},
	Overwrite: {},
	Delete:    {},
	Make:      {},
	Maybe:     {},
	Sep:       {},
	Escape:    {},
	At:        {},
	Self:      {},
	Move:      {},
}

func All() map[rune]struct{} {
	return maps.Clone(all)
}
