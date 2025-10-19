package giraffe

import (
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/internal/pipeimpl"
)

// Pipe NEVER INSTANTIATE DIRECTLY. NEVER CAST TO. NEVER CAST FROM.
type Pipe string

func (p Pipe) Steps() []Query {
	steps := p.impl().Steps()
	q := make([]Query, len(steps))
	for i, s := range steps {
		q[i] = Query(s.String())
	}
	return q
}

func (p Pipe) String() string {
	return p.impl().String()
}

func (p Pipe) impl() pipeimpl.PipeImpl {
	return M(pipeimpl.Parse(string(p)))
}
