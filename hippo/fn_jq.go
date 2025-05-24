package hippo

import (
	"errors"

	"github.com/itchyny/gojq"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
)

func MkJqFn(
	jq string,
) (JqFn, error) {
	return JqFn{jq: nil}.WithJq(jq)
}

type JqFn struct {
	jq *gojq.Query
}

func (e JqFn) Fn() *Fn {
	return FnOf(e.exe)
}

func (e JqFn) WithJq(
	jq string,
) (JqFn, error) {
	parsed, err := gojq.Parse(jq)
	if err != nil {
		return JqFn{jq: nil}, E(err)
	}

	return JqFn{jq: parsed}, nil
}

func (e JqFn) exe(
	ctx gtx.Context,
	call Call,
) (giraffe.Datum, error) {
	dat := call.Data()

	plain, err := dat.Plain()
	if err != nil {
		return dErr, err
	}

	collect := giraffe.OfEmpty()

	it := e.jq.RunWithContext(ctx, plain)
	v, ok := it.Next()
	for ok {
		switch vt := v.(type) {
		case error:
			var hErr *gojq.HaltError
			if errors.As(vt, &hErr) && hErr.Value() == nil {
				break
			}
			return dErr, E(vt)

		default:
			vd := giraffe.OfJsonable(vt)
			collect, err = collect.Merge(vd)
			if err != nil {
				return dErr, err
			}
		}
		v, ok = it.Next()
	}

	return collect, nil
}
