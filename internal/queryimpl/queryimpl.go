package queryimpl

import (
	"fmt"

	"github.com/hkoosha/giraffe/dialects"
	"github.com/hkoosha/giraffe/internal/inmem"
	"github.com/hkoosha/giraffe/internal/queryimpl/dialectical"
	"github.com/hkoosha/giraffe/internal/queryimpl/gquery"
	"github.com/hkoosha/giraffe/qflag"
)

// MaxDepth must fit in the gqflag.QFlag in the sequence part, i.e., 8 bits.
const MaxDepth = 255

func parse(
	spec string,
) (dialectical.DialecticalQuery, error) {
	dq := dialectical.New()

	dialect, spec, err := dialects.Normalized(spec)
	if err != nil {
		return dq, err
	}

	var impl QueryImpl
	switch dialect {
	case dialects.Giraffe1v1:
		impl, err = gquery.Parse(spec)
	}

	if err != nil {
		return dq, err
	}

	dq = dq.WithDialect(dialect).WithImpl(impl)
	return dq, nil
}

func Parse(
	spec string,
) (dialectical.DialecticalQuery, error) {
	cached, ok := inmem.Get[dialectical.DialecticalQuery](spec)

	if !ok {
		query, err := parse(spec)
		inmem.Set(spec, query, err)
		return query, err
	}

	return cached.Unpack()
}

type QueryImpl interface {
	fmt.Stringer

	Flags() qflag.QFlag

	Plus(query string) (QueryImpl, error)

	Attr() string
	Index() int

	Root() QueryImpl
	Leaf() QueryImpl
	Prev() QueryImpl
	Next() QueryImpl

	WithoutOverwrite() QueryImpl
	WithOverwrite() QueryImpl
	WithMake() QueryImpl
}
