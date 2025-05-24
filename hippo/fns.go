package hippo

import (
	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
)

func RunQuery() *Fn {
	// TODO
	panic("unimplemented")
}

// =====================================

func AssertNoError() *Fn {
	q := giraffe.Q("error")

	return FnOf(func(
		_ gtx.Context,
		c Call,
	) (giraffe.Datum, error) {
		subject, err := c.Data().QBln(q)
		if err != nil {
			return dErr, err
		}

		if subject {
			return dErr, EF("unexpected error")
		}

		return giraffe.OfEmpty(), nil
	})
}

func AssertEq() *Fn {
	return FnOf(func(
		_ gtx.Context,
		c Call,
	) (giraffe.Datum, error) {
		q, err := c.Args().QStr("query")
		if err != nil {
			return dErr, err
		}

		value, err := c.Args().Get("value")
		if err != nil {
			return dErr, err
		}

		qq, err := giraffe.GQParse(q)
		if err != nil {
			return dErr, err
		}

		subject, err := c.Data().Get(qq)
		if err != nil {
			return dErr, err
		}

		if !subject.Eq(value) {
			return dErr, EF(
				"data mismatch, key=%s, expecting=%v, have=%v",
				q,
				value,
				subject,
			)
		}

		return giraffe.OfEmpty(), nil
	})
}

// =====================================

func Data() *Fn {
	return FnOf(func(
		_ gtx.Context,
		c Call,
	) (giraffe.Datum, error) {
		return c.Data().Merge(c.Args())
	})
}

// =====================================

func Static(
	dat giraffe.Datum,
) *Fn {
	return FnOf(func(
		gtx.Context,
		Call,
	) (giraffe.Datum, error) {
		return dat, nil
	})
}

func StaticOf(
	pairs ...giraffe.Tuple,
) (*Fn, error) {
	dat, err := giraffe.OfN(pairs...)
	if err != nil {
		return nil, err
	}

	return Static(dat), nil
}

// =====================================

func SelectRand() *Fn {
	return FnOf(selectRand)
}

func selectRand(
	ctx gtx.Context,
	call Call,
) (giraffe.Datum, error) {
	l, err := call.Data().Len()
	if err != nil {
		return dErr, err
	}

	// TODO how to sanely select obj keys? Probably never.
	// TODO chk int range.

	i := ctx.Rand().StdV2().UintN(uint(l)) //nolint:gosec // TODO fix later
	return call.Data().At(int(i))          //nolint:gosec // TODO fix later
}
