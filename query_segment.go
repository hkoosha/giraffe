package giraffe

import (
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/queryimpl"
	"github.com/hkoosha/giraffe/internal/queryimpl/dialectical"
)

// Segment NEVER INSTANTIATE DIRECTLY. NEVER CAST TO. NEVER CAST FROM.
type Segment string

func (q Segment) impl() dialectical.DialecticalQuery {
	return M(queryimpl.Parse(string(q)))
}

func (q Segment) String() string {
	return q.impl().String()
}
