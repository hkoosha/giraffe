package pipelines

import (
	"context"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/dot"
)

//goland:noinspection GoUnusedExportedFunction
func Step(
	name string,
	fn Fn,
	requires ...giraffe.Query,
) Fn {
	if len(requires) == 0 {
		panic(EF("no requirement for step provided"))
	}

	stepFn := stepFn{
		name:     name,
		fn:       fn,
		requires: requires,
	}

	return stepFn.Ekran
}

type stepFn struct {
	name     string
	fn       Fn
	requires []giraffe.Query
}

func (m *stepFn) String() string {
	return "Step[" + m.name + "]"
}

func newMissingKeysError(queries []giraffe.Query) error {
	return newPipelineQueriesError(
		ErrCodeMissingKeys,
		"missing keys",
		queries...,
	)
}

func (m *stepFn) Ekran(
	ctx context.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	var missing []giraffe.Query

	for _, k := range m.requires {
		if !dat.Has(k) {
			missing = append(missing, k)
		}
	}

	if len(missing) > 0 {
		return giraffe.OfErr(), newMissingKeysError(missing)
	}

	result, err := m.fn(ctx, dat)
	if err != nil {
		return giraffe.OfErr(), err
	}

	return result, nil
}
