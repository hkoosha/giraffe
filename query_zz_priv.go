package giraffe

import (
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/internal"
	"github.com/hkoosha/giraffe/internal/queryimpl"
)

// Uncomment, anc change return types in query.go from Query to
// queryimpl.QueryImpl to do a manual check.
// var _ queryimpl.QueryImpl = (*Query)(nil)

func (q Query) impl() queryimpl.QueryImpl {
	return M(internal.Parse(string(q)))
}
