package pipelines

import (
	"context"
	"maps"
	"reflect"

	"github.com/hkoosha/giraffe"
)

func Mapping(
	m map[giraffe.Query]giraffe.Query,
) Fn {
	fn := mappingFn{
		mapping: maps.Clone(m),
	}

	return fn.Ekran
}

type mappingFn struct {
	mapping map[giraffe.Query]giraffe.Query
}

func (m *mappingFn) String() string {
	return reflect.TypeOf(m).Elem().String()
}

func (m *mappingFn) Ekran(
	_ context.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	mapping := giraffe.OfEmpty()

	for k, v := range m.mapping {
		if d, err := dat.Get(v); err != nil {
			return giraffe.OfErr(), err
		} else {
			mapping, err = mapping.Set(k, d)
			if err != nil {
				return giraffe.OfErr(), err
			}
		}
	}

	return mapping, nil
}
