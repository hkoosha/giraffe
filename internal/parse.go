package internal

import (
	"github.com/hkoosha/giraffe/dialects"
	"github.com/hkoosha/giraffe/internal/inmem"
	"github.com/hkoosha/giraffe/internal/queryimpl"
	"github.com/hkoosha/giraffe/internal/queryimpl/dialectical"
	"github.com/hkoosha/giraffe/internal/queryimpl/gquery"
)

func parse(
	spec string,
) (dialectical.DialecticalQuery, error) {
	dq := dialectical.New()

	dialect, spec, err := dialects.Normalized(spec)
	if err != nil {
		return dq, err
	}

	var impl queryimpl.QueryImpl
	switch dialect {
	case dialects.Giraffe1v1:
		impl, err = gquery.Parse(queryimpl.MaxDepth, spec)
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
