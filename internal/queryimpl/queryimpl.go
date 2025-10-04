package queryimpl

import (
	"fmt"

	"github.com/hkoosha/giraffe/dialects"
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/inmem"
	"github.com/hkoosha/giraffe/internal/queryimpl/gquery"
	"github.com/hkoosha/giraffe/internal/queryimpl/hquery"
	"github.com/hkoosha/giraffe/qflag"
)

var invalid = DialecticalQuery{}

type QueryImpl interface {
	fmt.Stringer

	Flags() qflag.QFlag

	Plus(other string) (QueryImpl, error)
}

type DialecticalQuery struct {
	Dialect dialects.Dialect
	QueryImpl
	impl QueryImpl
}

func (q DialecticalQuery) String() string {
	return q.impl.String()
}

func (q DialecticalQuery) Flags() qflag.QFlag {
	return q.impl.Flags()
}

func parse(
	spec string,
) (DialecticalQuery, error) {
	spec, err := dialects.Normalized(spec)
	if err != nil {
		return invalid, err
	}

	q := DialecticalQuery{
		Dialect: M(dialects.DialectOf(spec)),
	}

	switch q.Dialect {
	case dialects.Giraffe:
		q.impl, err = gquery.Parse(spec)

	case dialects.Http:
		q.impl, err = hquery.Parse(spec)
	}

	if err != nil {
		q = invalid
	}

	return q, err
}

func Parse(
	spec string,
) (DialecticalQuery, error) {
	cached, ok := inmem.Get[DialecticalQuery](spec)

	if !ok {
		query, err := parse(spec)
		inmem.Set(spec, query, err)
		return query, err
	}

	return cached.Unpack()
}
