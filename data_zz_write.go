package giraffe

import (
	"strings"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

func modify(
	d *Datum,
	q queryT,
	value any,
) error {
	switch {
	case value != nil && q.Flags().IsDelete():
		return newDataWriteUnexpectedValueError(q, value)

	case q.Flags().IsDelete():
		if err := del(d, q); err != nil {
			return err
		}

		return nil

	default:
		v, err := of(value)
		if err != nil {
			return err
		}

		if err := set(d, q, v); err != nil {
			return err
		}

		return nil
	}
}

func del(
	d *Datum,
	q queryT,
) error {
	qf := q.Flags()
	dt := d.typ

	switch {
	case dt.IsObj() != qf.IsObj():
		return newQueryTypeCastError(dt, qf)

	case dt.IsArr() && q.Index() >= d.len():
		return newDataWriteMissingKeyError(q)
	}

	switch {
	case !d.hasShallow(q):
		// Do nothing.

	case dt.IsArr():
		arr := any(append(d.arr()[:q.Index()], d.arr()[q.Index()+1:]...))
		d.val = &arr

	default:
		obj := d.obj()
		delete(obj, q.Attr())
		a := any(obj)
		d.val = &a
	}

	return nil
}

// ==============================================================================.

func set(
	d *Datum,
	q queryT,
	value Datum,
) error {
	switch {
	case q.Flags().IsObj():
		return setObj(d, q, value)

	default:
		return arrSet(d, q, value)
	}
}

func setObj(
	d *Datum,
	q queryT,
	item Datum,
) error {
	qf := q.Flags()
	dt := d.typ

	has := d.hasShallow(q)

	switch {
	case qf.IsMaybe():
		panic(EF("TODO: implement maybe for set obj"))

	case !dt.isZero() && dt.IsObj() != q.Flags().IsObj():
		return newQueryTypeCastError(dt, q.Flags())

	case !qf.IsLeaf() && !qf.IsMake() && !has:
		return newDataWriteMissingKeyError(q)

	case qf.IsLeaf() && !qf.IsOverwrite() && has:
		return newDataWriteOverwriteErr(q)
	}

	switch {
	case qf.IsLeaf() && dt.IsObj():
		d.obj()[q.Attr()] = item

	case qf.IsLeaf() && dt.isZero():
		*d = _newDatum(
			Obj, map[string]Datum{
				q.Attr(): item,
			},
		)

	case !qf.IsLeaf() && dt.IsObj():
		dd := d.obj()
		ddI := dd[q.Attr()]

		if err := set(&ddI, q.Next(), item); err != nil {
			return err
		}

		dd[q.Attr()] = ddI

	case !qf.IsLeaf() && dt.isZero():
		dd := _newDatum(Type(0), nil)
		if err := set(&dd, q.Next(), item); err != nil {
			return err
		}

		*d = _newDatum(
			Obj, map[string]Datum{
				q.Attr(): dd,
			},
		)

	default:
		panic(EF("unreachable: unknown case for set obj"))
	}

	return nil
}

func arrSet(
	d *Datum,
	q queryT,
	item Datum,
) error {
	qf := q.Flags()
	dt := d.typ

	switch {
	case qf.IsMaybe():
		panic(EF("TODO: implement maybe for set arr"))

	case !dt.isZero() && dt.IsObj() != q.Flags().IsObj():
		return newQueryTypeCastError(dt, qf)

	case !qf.IsMake() && dt.isZero():
		return newDataWriteMissingKeyError(q)

	case !qf.IsAppend() && d.len() <= q.Index():
		return newDataWriteMissingKeyError(q)
	}

	Assert(!dt.isZero() || qf.IsMake() || qf.Val() == 0)

	switch {
	case d.hasShallow(q):
		d.arr()[q.Index()] = item

	case dt.isZero() && !qf.IsLeaf():
		dd := _newDatum(Type(0), nil)
		if err := set(&dd, q.Next(), item); err != nil {
			return err
		}

		*d = _newDatum(Arr, []Datum{dd})

	case dt.isZero() && qf.IsLeaf():
		*d = _newDatum(Arr, []Datum{item})

	case qf.IsAppend():
		dd := d.arr()
		dd = append(dd, item)
		a := any(dd)
		d.val = &a

	default:
		panic(EF("unreachable: unknown case for set arr"))
	}

	return nil
}

// ==============================================================================.

func newDataWriteOverwriteErr(
	query queryT,
) error {
	return newDataWriteError(
		query,
		ErrCodeDataWriteImplicitOverwrite,
		"cannot implicitly overwrite data",
	)
}

func newDataWriteMissingKeyError(
	query queryT,
) error {
	return newDataWriteError(
		query,
		ErrCodeDataWriteMissingKey,
		"missing key",
	)
}

func newDataWriteError(
	query queryT,
	code uint64,
	msg string,
	extra ...string,
) error {
	sb := strings.Builder{}
	sb.WriteString("data write error, ")
	sb.WriteString(msg)
	sb.WriteString(": ")
	sb.WriteString(query.String())

	for _, e := range extra {
		sb.WriteString(", ")
		sb.WriteString(e)
	}

	return E(newGiraffeError(
		code,
		sb.String(),
	))
}
