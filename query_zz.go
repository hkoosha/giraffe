package giraffe

import (
	"reflect"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/gquery"
)

var (
	invalid = Query("")
	tQuery  = reflect.TypeOf((*Query)(nil)).Elem()
)

func (q Query) impl() queryimpl.Query {
	return M(queryimpl.Parse(string(q)))
}
