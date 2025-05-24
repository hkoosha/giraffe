package internal

import (
	"time"

	"github.com/hkoosha/giraffe/core/inmem"
	"github.com/hkoosha/giraffe/dialects"
	"github.com/hkoosha/giraffe/internal/queryimpl/gquery"
)

var cache = inmem.Make[gquery.GiraffeQuery](
	"github.com/hkoosha/giraffe|parse_query",
	7*24*time.Hour,
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
		return gquery.Parse(spec)

	case dialects.Unknown:
		return gquery.GiraffeQuery{}, dialects.ErrUnknown()

	default:
		return gquery.GiraffeQuery{}, dialects.ErrUnknown()
	}
}

func Parse(
	spec string,
) (gquery.GiraffeQuery, error) {
	cached, ok := cache.Get(spec)

	if !ok {
		query, err := parse(spec)
		cache.Set(spec, query, err)
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
