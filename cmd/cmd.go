package cmd

const (
	// Write

	Append    Cmd = '+'
	Overwrite Cmd = '='
	Delete    Cmd = '!'
	Make      Cmd = '$'

	// Read

	Maybe  Cmd = '?'
	Sep    Cmd = '.'
	Escape Cmd = '\\'
	At     Cmd = '@'
	Self   Cmd = '#'

	// Control

	Dialect Cmd = '~'
	Pipe    Cmd = '|'
)

type Cmd byte

func (c Cmd) String() string {
	return string(c)
}

func (c Cmd) Byte() byte {
	return byte(c)
}

var all = map[Cmd]struct{}{
	Dialect:   {},
	Pipe:      {},
	Append:    {},
	Overwrite: {},
	Delete:    {},
	Make:      {},
	Maybe:     {},
	Sep:       {},
	Escape:    {},
	At:        {},
	Self:      {},
}
