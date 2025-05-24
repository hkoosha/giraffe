package giraffe

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/hkoosha/giraffe/cmd"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/internal"
	"github.com/hkoosha/giraffe/internal/gstrings"
)

func (d Datum) hasShallow(
	q queryT,
) bool {
	qf := q.Flags()
	dt := d.typ

	if qf.IsSelf() {
		panic(EF("unreachable: self query on hasShallow"))
	}

	if dt.isZero() {
		return false
	}

	switch {
	case qf.IsArr():
		return dt.IsArr() && 0 <= q.Index() && q.Index() < len(d.arr())

	case qf.IsObj():
		if !dt.IsObj() {
			return false
		}

		_, ok := d.obj()[q.Attr()]

		return ok

	default:
		panic(EF("unreachable, unknown query type: %s", q.String()))
	}
}

func (d Datum) tryObj() (map[string]Datum, error) {
	dt := d.typ

	switch {
	case dt.IsNil():
		return nil, newNilError()

	case !dt.IsObj():
		return nil, newTypeCastError(dt, Obj)

	default:
		return cast[map[string]Datum](d), nil
	}
}

func (d Datum) tryArr() ([]Datum, error) {
	dt := d.typ

	switch {
	case dt.IsNil():
		return nil, newNilError()

	case !dt.IsArr():
		return nil, newTypeCastError(dt, Arr)

	default:
		return cast[[]Datum](d), nil
	}
}

func (d Datum) tryLen() (int, error) {
	dt := d.typ

	switch {
	case dt.IsNil():
		return -1, newNilError()

	case dt.IsArr():
		return len(d.arr()), nil

	case dt.IsObj():
		return len(d.obj()), nil

	default:
		return -1, newTypeCastError(dt, Arr)
	}
}

func (d Datum) obj() map[string]Datum {
	return M(d.tryObj())
}

func (d Datum) arr() []Datum {
	return M(d.tryArr())
}

func (d Datum) len() int {
	return M(d.tryLen())
}

func (d Datum) get(
	q queryT,
) (Datum, error) {
	// 	q, err := q.Resolved(d.resolver)

	qf := q.Flags()
	dt := d.typ

	switch {
	case qf.IsIndeterministic():
		return OfErr(), newDataReadOnlyError(q)

	case !qf.IsReadonly():
		return OfErr(), newDataReadIndeterministicQueryError(q)

	case qf.IsSelf():
		if qf.IsLeaf() {
			return d, nil
		}

		return d.get(q.Next())

	case qf.IsArr() && dt.IsArr():
		v := d.arr()
		i := q.Index()

		switch {
		case i >= len(v):
			return OfErr(), newDataReadOutOfBoundsError(q)
		case qf.IsLeaf():
			return v[i], nil
		default:
			return v[i].get(q.Next())
		}

	case qf.IsObj() && dt.IsObj():
		v, ok := d.obj()[q.Attr()]
		if !ok {
			return OfErr(), newDataReadMissingKeyError(q)
		} else if qf.IsLeaf() {
			return v, nil
		} else {
			return v.get(q.Next())
		}

	case qf.IsArr():
		return OfErr(), newDataReadUnexpectedTypeError(q, Arr, dt)

	case qf.IsObj():
		return OfErr(), newDataReadUnexpectedTypeError(q, Obj, dt)

	default:
		panic(EF("unreachable, cannot handler query for item get: %s", q.String()))
	}
}

func (d Datum) resolver(
	q string,
) (string, error) {
	qp, err := internal.Parse(q)
	if err != nil {
		return "", err
	}

	dd, err := d.get(qp)
	if err != nil {
		return "", err
	}

	return dd.Str()
}

func (d Datum) tree() []queryT {
	var tr []queryT
	tree(&tr, &d, []string{})

	return tr
}

func tree(
	tr *[]queryT,
	d *Datum,
	path []string,
) bool {
	dt := d.typ

	switch {
	case !dt.IsObj() && len(path) == 0:
		return false

	// Once the first property of the object is handled, we have the path for
	// this object, and we don't want to repeat the current path for each
	// property (which would result in the same path). Hence, we signal the
	// parent call to bail out.
	case !dt.IsObj():
		q := Q(strings.Join(path, cmd.Sep.String())).impl()
		*tr = append(*tr, q)

		return true

	case d.len() == 0 && len(path) == 0:
		return false

	case d.len() == 0:
		q := Q(strings.Join(path, cmd.Sep.String())).impl()
		*tr = append(*tr, q)

		return false

	default:
		for k, v := range d.obj() {
			if tree(tr, &v, append(slices.Clone(path), internal.Escaped(k))) {
				return false
			}
		}

		return false
	}
}

// ==============================================================================.

func newDataReadOnlyError(
	query queryT,
) error {
	return newDataReadError(
		query,
		ErrCodeDataReadOnly,
		"query is not readonly",
	)
}

func newDataReadIndeterministicQueryError(
	query queryT,
) error {
	return newDataReadError(
		query,
		ErrCodeDataReadIndeterministicQuery,
		"indeterministic query",
	)
}

func newDataReadOutOfBoundsError(
	query queryT,
) error {
	return newDataReadError(
		query,
		ErrCodeDataReadIndexOutOfBounds,
		"index out of bounds",
	)
}

func newDataReadMissingKeyError(
	query queryT,
) error {
	return newDataReadError(
		query,
		ErrCodeDataReadMissingKey,
		"missing key",
	)
}

func newDataReadIntegerOverflowError(
	need reflect.Type,
	extra ...string,
) error {
	return newGiraffeError(
		ErrCodeOverflowError,
		fmt.Sprintf(
			"integer does not fit: target=%s%s",
			need.String(),
			gstrings.Joined(extra),
		),
	)
}

func newDataReadUnexpectedTypeError(
	query queryT,
	expecting Type,
	actual Type,
) error {
	return newDataReadError(
		query,
		ErrCodeDataReadUnexpectedType,
		"unexpected type",
		"expecting=",
		expecting.String(),
		"actual=",
		actual.String(),
	)
}

func newDataReadError(
	query queryT,
	code uint64,
	msg string,
	extra ...string,
) error {
	sb := strings.Builder{}
	sb.WriteString("data read error, ")
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
