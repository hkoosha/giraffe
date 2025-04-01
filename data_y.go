package giraffe

import (
	"fmt"
	"reflect"
	"strings"

	. "github.com/hkoosha/giraffe/internal/dot"
	"github.com/hkoosha/giraffe/internal/gquery"
	"github.com/hkoosha/giraffe/zebra/z"
)

func (d Datum) String0() string {
	switch {
	case d.typ.IsErr():
		return getDataErrRepr()

	case d.typ.IsArr():
		sb := strings.Builder{}
		sb.WriteByte('[')

		for i := range d.len() {
			if i >= 3 {
				sb.WriteString(", ...")

				break
			}

			if i > 0 {
				sb.WriteString(", ")
			}

			item := M(d.At(i))
			sb.WriteString(item.shallowStr())
		}

		sb.WriteByte(']')

		return sb.String()

	case d.typ.IsObj():
		sb := strings.Builder{}
		sb.WriteByte('{')

		i := 0
		for k := range d.obj() {
			if i >= 3 {
				sb.WriteString(", ...")

				break
			}

			if i > 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(k)

			i++
		}

		sb.WriteByte('}')

		return sb.String()

	default:
		c, err := d.raw()
		if err != nil {
			return getDataErrRepr()
		}

		return fmt.Sprintf("%v", c)
	}
}

func (d Datum) shallowStr() string {
	switch {
	case d.typ.IsArr() && d.len() == 0:
		return "[]"
	case d.typ.IsArr():
		return "[...]"
	case d.typ.IsObj() && d.len() == 0:
		return "{}"
	case d.typ.IsObj():
		return "{...}"
	default:
		return d.String0()
	}
}

func (d Datum) raw() (any, error) {
	switch {
	case d.typ.IsErr():
		return nil, newInvalidDatumError()

	case d.typ.IsNil():
		//nolint:nilnil
		return nil, nil

	case d.typ.IsStr():
		return d.Str()

	case d.typ.IsBln():
		return d.Bln()

	case d.typ.IsFlt():
		return d.Flt()

	case d.typ.IsInt():
		return d.Int()

	case d.typ.IsArr():
		return z.TryApplied(d.arr(), func(it Datum) (any, error) {
			return it.raw()
		})

	case d.typ.IsObj():
		return z.TryApplied2(d.obj(), func(_ string, it Datum) (any, error) {
			return it.raw()
		})

	case d.typ.isZero():
		panic(EF("unreachable, datum is zero"))

	default:
		panic(EF("unreachable, unknown datum type: %s", d.typ.String()))
	}
}

func (d Datum) eq(
	other Datum,
) bool {
	if d.typ.IsErr() || other.typ.IsErr() {
		return false
	}

	if d.typ != other.typ {
		return false
	}

	switch {
	case d.typ.IsObj():
		dObj := d.obj()
		oObj := other.obj()
		if len(dObj) != len(oObj) {
			return false
		}
		for k, v := range dObj {
			if oV, ok := oObj[k]; !ok || !v.eq(oV) {
				return false
			}
		}

		return true

	case d.typ.IsArr():
		dArr := d.arr()
		oArr := other.arr()
		if len(dArr) != len(oArr) {
			return false
		}
		for i, v := range dArr {
			if !v.eq(oArr[i]) {
				return false
			}
		}

		return true

	default:
		return reflect.DeepEqual(d.val, other.val)
	}
}

func (d Datum) cmp(
	other Datum,
) (int, error) {
	dI, err := d.Int()
	if err != nil {
		return -1, err
	}

	oI, err := other.Int()
	if err != nil {
		return -1, err
	}

	return dI.Cmp(oI), nil
}

func (d Datum) merge(
	right Datum,
	path []string,
) (Datum, error) {
	if d.typ != right.typ {
		return OfErr(), newMergeIncompatibleTypesError(d.typ, right.typ)
	}

	switch {
	case d.typ.IsObj():
		fin := Of(d)
		obj := fin.obj()

		for k, v := range right.obj() {
			if existing, ok := obj[k]; ok {
				nPath := Appended(path, k)
				var err error
				if obj[k], err = existing.merge(v, nPath); err != nil {
					return OfErr(), err
				}
			} else {
				obj[k] = v
			}
		}

		return fin, nil

	case d.typ.IsArr():
		panic("todo")

	default:
		if !d.eq(right) {
			return OfErr(), newMergeClashingKeysError(path)
		}

		return d, nil
	}
}

func (d Datum) nest(
	q Query,
) (Datum, error) {
	if q.impl().Flags().IsNonDeterministic() {
		panic("TODO: fix non deterministic flag for nest")
	}

	obj, err := d.tryObj()
	if err != nil {
		return errD, err
	}

	nested := OfEmpty()
	for k, v := range obj {
		qK, err := Parse(k)
		if err != nil {
			return errD, err
		}
		qK = q.Plus(qK)

		nested, err = nested.Set(qK, v)
		if err != nil {
			return errD, err
		}
	}

	return nested, nil
}

func (d Datum) set(
	q gquery.Query,
	value any,
) (Datum, error) {
	if !d.typ.IsArr() && !d.typ.IsObj() {
		panic("TODO unimplemented, set for non-container types: " + d.typ.String())
	}

	cp := M(of(d.deref()))
	if err := modify(&cp, q, value); err != nil {
		return errD, err
	}

	return cp, nil
}

func (d Datum) has(
	q gquery.Query,
) bool {
	switch {
	case q.Flags().IsObj() && d.typ.IsObj():
		v, ok := d.obj()[q.Named()]
		if !ok {
			return false
		}

		return q.Flags().IsLeaf() || v.has(q.Next())

	case q.Flags().IsArr() && d.typ.IsArr():
		arr := d.arr()
		if q.Index() >= len(arr) {
			return false
		}

		return q.Flags().IsLeaf() || arr[q.Index()].has(q.Next())

	default:
		return false
	}
}

func (d Datum) deref() any {
	return *d.val
}

func cast[T any](d Datum) T {
	if d.typ.IsNil() {
		panic(newNilError())
	}

	t, ok := d.deref().(T)
	if !ok {
		var zero T

		panic(newReflectiveTypeCastError(d.typ, zero))
	}

	return t
}

func copyToArr[T any](
	d Datum,
	castFn func(Datum) (T, error),
) ([]T, error) {
	v, err := d.tryArr()
	if err != nil {
		return nil, err
	}

	val := make([]T, len(v))

	for i, j := range v {
		r, err1 := castFn(j)
		if err1 != nil {
			return nil, err1
		}

		val[i] = r
	}

	return val, nil
}

// =============================================================================.

func newMergeIncompatibleTypesError(
	left Type,
	right Type,
) error {
	return newDataMakeError(
		ErrCodeDataMergeIncompatibleTypes,
		fmt.Sprintf(
			"incompatible types on merge: left=%s, right=%s",
			left.String(),
			right.String(),
		),
	)
}

func newMergeClashingKeysError(
	key []string,
) error {
	q := z.Applied(key, Escaped)

	return newDataMakeError(
		ErrCodeDataMergeClashingKeys,
		"clashing keys: ["+strings.Join(q, CmdSep)+"]",
	)
}

func newDataWriteUnexpectedValueError(
	q gquery.Query,
	v any,
) error {
	return newDataWriteError(
		q,
		ErrCodeDataModifyOperationTakesNoValue,
		"the modification takes no value, but a value was provided: "+fmt.Sprint(v),
	)
}
