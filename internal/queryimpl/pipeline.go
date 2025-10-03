package queryimpl

import (
	"slices"
)

type Pipeline struct {
	// Debug DebugImpl

	queries []Query
	ref     string
}

func (q Pipeline) Segments() []Query {
	return slices.Clone(q.queries)
}

func (q Pipeline) String() string {
	return q.ref
}
