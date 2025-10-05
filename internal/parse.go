package internal

import (
	"github.com/hkoosha/giraffe/dialects"
	"github.com/hkoosha/giraffe/internal/inmem"
	"github.com/hkoosha/giraffe/internal/queryimpl"
	"github.com/hkoosha/giraffe/internal/queryimpl/gquery"
)

func parse(
	spec string,
) (gquery.GiraffeQuery, error) {
	dialect, spec, err := dialects.Normalized(spec)
	if err != nil {
		return gquery.GiraffeQuery{}, err
	}

	switch dialect {
	case dialects.Giraffe1v1:
		return gquery.Parse(queryimpl.MaxDepth, spec)

	case dialects.Unknown:
		return gquery.GiraffeQuery{}, dialects.ErrUnknown()

	default:
		return gquery.GiraffeQuery{}, dialects.ErrUnknown()
	}
}

func Parse(
	spec string,
) (gquery.GiraffeQuery, error) {
	cached, ok := inmem.Get[gquery.GiraffeQuery](inmem.BucketParseQuery, spec)

	if !ok {
		query, err := parse(spec)
		inmem.Set(inmem.BucketParseQuery, spec, query, err)
		return query, err
	}

	return cached.Unpack()
}

func Escaped(
	spec string,
) string {
	// TODO :D
	return spec
}
