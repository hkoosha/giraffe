package giraffe

import (
	"strings"

	. "github.com/hkoosha/giraffe/internal/dot"
	"github.com/hkoosha/giraffe/internal/gquery"
)

func modify(
	d *Datum,
	q gquery.Query,
	value any,
) error {
	switch {
	case value != nil && (q.Flags().IsDelete() || q.Flags().IsMove()):
		return newDataWriteUnexpectedValueError(q, value)

	case q.Flags().IsDelete():
		if err := del(d, q); err != nil {
			return err
		}

		return nil

	case q.Flags().IsMove():
		if err := move(d, q); err != nil {
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

func move(
	d *Datum,
	q gquery.Query,
) error {
	seg0, seg1, ok := q.Segments()
	if !ok {
		return newDataWriteMoveUnsegmentedQuery(q)
	}

	prev, err := d.get(seg0)
	if err != nil {
		return err
	}

	if dErr := del(d, seg0); dErr != nil {
		return dErr
	}

	newD, err := d.set(seg1, prev)
	if err != nil {
		return err
	}

	*d = newD

	return nil
}

func del(
	d *Datum,
	q gquery.Query,
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
		delete(obj, q.Named())
		a := any(obj)
		d.val = &a
	}

	return nil
}

// ==============================================================================.

func set(
	d *Datum,
	q gquery.Query,
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
	q gquery.Query,
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
		d.obj()[q.Named()] = item

	case qf.IsLeaf() && dt.isZero():
		*d = _newDatum(
			Obj, map[string]Datum{
				q.Named(): item,
			},
		)

	case !qf.IsLeaf() && dt.IsObj():
		dd := d.obj()
		ddI := dd[q.Named()]

		if err := set(&ddI, q.Next(), item); err != nil {
			return err
		}

		dd[q.Named()] = ddI

	case !qf.IsLeaf() && dt.isZero():
		dd := _newDatum(Type(0), nil)
		if err := set(&dd, q.Next(), item); err != nil {
			return err
		}

		*d = _newDatum(
			Obj, map[string]Datum{
				q.Named(): dd,
			},
		)

	default:
		panic(EF("unreachable: unknown case for set obj"))
	}

	return nil
}

func arrSet(
	d *Datum,
	q gquery.Query,
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
	query gquery.Query,
) error {
	return newDataWriteError(
		query,
		ErrCodeDataWriteImplicitOverwrite,
		"cannot implicitly overwrite data",
	)
}

func newDataWriteMissingKeyError(
	query gquery.Query,
) error {
	return newDataWriteError(
		query,
		ErrCodeDataWriteMissingKey,
		"missing key",
	)
}

func newDataWriteMoveUnsegmentedQuery(
	query gquery.Query,
) error {
	return newDataWriteError(
		query,
		ErrCodeDataWriteUnsegmentedQuery,
		"move query has no or too many segments",
	)
}

func newDataWriteError(
	query gquery.Query,
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
