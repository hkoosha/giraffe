package giraffe

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/hkoosha/giraffe/internal"
	. "github.com/hkoosha/giraffe/internal/dot"
	"github.com/hkoosha/giraffe/internal/gdatum"
	"github.com/hkoosha/giraffe/internal/gquery"
	"github.com/hkoosha/giraffe/zebra/z"
)

// TODO schemas.

var (
	errD     = _newDatum(Type(0), nil)
	emptyObj = _newDatum(Obj, map[string]Datum{})
	emptyArr = _newDatum(Obj, []Datum{})

	_datumPtrType = reflect.TypeOf((*Datum)(nil))
	_prohibited   = map[reflect.Type]z.NA{
		reflect.TypeOf((*Query)(nil)).Elem():         z.None,
		reflect.TypeOf((*Type)(nil)).Elem():          z.None,
		reflect.TypeOf((*Tuple)(nil)).Elem():         z.None,
		reflect.TypeOf((*reflect.Value)(nil)).Elem(): z.None,
		reflect.TypeOf((*reflect.Type)(nil)).Elem():  z.None,
		reflect.TypeOf(reflect.Type(nil)):            z.None,
	}

	jsonNumberType = reflect.TypeOf((*json.Number)(nil)).Elem()
)

func of(
	v any,
) (Datum, error) {
	r := reflect.ValueOf(v)

	switch {
	// TODO is this the correct way to detect reflect.ValueOf(nil)?
	case !r.IsValid() && v == nil:
		return _newDatum(Nil, nil), nil

	case !r.IsValid():
		return OfErr(), newDataMakeInvalidTypeError(r.Type())

	case r.Type() == _datumPtrType && r.IsNil():
		return OfErr(), newNilError()

	// MUST make a copy: return value WILL BE MODIFIED by internal code.
	case r.Type() == _datumPtrType:
		d, ok := v.(*Datum)
		if !ok {
			panic(EF("unreachable: unexpected value for datum"))
		}

		val, typ, err := _ofAny(d.deref(), reflect.ValueOf(d.deref()))
		if err != nil {
			return OfErr(), err
		}

		if typ == Err {
			panic(EF("unreachable, invalid giraffe type: %s", typ.String()))
		}

		return _newDatum(typ, val), nil

	case r.Kind() == reflect.Map && r.Type().Key() == tQuery:
		return _ofExpandable(v, r)

	default:
		val, typ, err := _ofAny(v, r)
		if err != nil {
			return OfErr(), err
		}

		if typ == Err {
			panic(EF("unreachable, invalid giraffe type: %s", typ.String()))
		}

		return _newDatum(typ, val), nil
	}
}

func _newDatum(
	typ Type,
	val any,
) Datum {
	return Datum{
		val:   &val,
		typ:   typ,
		Debug: gdatum.NewDatumDebug(),
	}
}

func _ofExpandable(
	_ any,
	r reflect.Value,
) (Datum, error) {
	d := OfEmpty()

	it := r.MapRange()
	for it.Next() {
		var q gquery.Query
		if qCast, ok := it.Key().Interface().(Query); ok {
			q = qCast.impl()
		} else if qCast, ok := it.Key().Interface().(Query); ok {
			q = qCast.impl()
		} else {
			panic(EF(
				"unreachable, not a Query or Write: %s",
				it.Key().String(),
			))
		}

		q = q.WithMake()
		v := it.Value().Interface()

		if dd, err := d.set(q, v); err != nil {
			return OfErr(), err
		} else {
			d = dd
		}
	}

	return d, nil
}

func _ofAny(
	v any,
	r reflect.Value,
) (any, Type, error) {
	switch r.Kind() {
	case reflect.Pointer:
		return _ofPtr(v, r)

	case reflect.Slice, reflect.Array:
		return _ofSeq(v, r)

	case reflect.Map:
		return _ofMap(r)

	case reflect.Struct:
		return _ofStruct(v, r)

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		b := big.NewInt(0)
		b.SetUint64(r.Uint())

		return b, Int, nil

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return big.NewInt(r.Int()), Int, nil

	case reflect.Float32, reflect.Float64:
		return big.NewFloat(r.Float()), Flt, nil

	case reflect.Bool:
		return r.Bool(), Bln, nil

	//nolint:nestif
	case reflect.String:
		if r.Type() == jsonNumberType {
			if num, ok := r.Interface().(json.Number); ok {
				if i, err := num.Int64(); err == nil {
					return big.NewInt(i), Int, nil
				} else if f, err := num.Float64(); err == nil {
					return big.NewFloat(f), Flt, nil
				} else {
					return nil, Err, newDataMakeUnimplementedType(r.Type())
				}
			} else {
				panic(EF("unreachable, unexpected type for json.Number"))
			}
		}

		return r.String(), Str, nil

	case reflect.Func,
		reflect.Complex64,
		reflect.Complex128:
		return nil, Err, newDataMakeUnimplementedType(r.Type())

	case reflect.Uintptr,
		reflect.Interface,
		reflect.Chan,
		reflect.UnsafePointer:
		return nil, Err, newDataMakeUnimplementedType(r.Type())

	case reflect.Invalid:
		panic(EF("unreachable, invalid should be already handled"))

	default:
		panic(EF("unreachable, golang type unknown: %s", r.Kind().String()))
	}
}

func _ofStruct(
	v any,
	r reflect.Value,
) (any, Type, error) {
	if dat, ok := v.(Datum); ok {
		datAgain, err := of(dat.deref())
		if err != nil {
			return nil, Err, newDataMakeMarshalError(err)
		}

		return datAgain.deref(), datAgain.typ, nil
	}

	if isProhibited(r.Type()) {
		return nil, Err, newDataMakeProhibitedTypeError(r.Type())
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil, Err, newDataMakeMarshalError(err)
	}

	var conv any
	if err1 := json.Unmarshal(b, &conv); err1 != nil {
		return nil, Err, newDataMakeUnmarshalError(err1)
	}

	asMap, ok := conv.(map[string]any)
	if !ok {
		return nil, Err, newDataMakeUnimplementedType(r.Type())
	}

	return _ofMap(reflect.ValueOf(asMap))
}

//nolint:nestif
func _ofMap(
	r reflect.Value,
) (any, Type, error) {
	dat := _newDatum(Obj, make(map[string]Datum, r.Len()))

	for _, key := range r.MapKeys() {
		d, err := of(r.MapIndex(key).Interface())
		if err != nil {
			return nil, Err, err
		}

		if key.Kind() == reflect.String {
			o := dat.obj()
			o[key.String()] = d
			d.val = Ref(any(o))
		} else if key.CanInterface() && key.CanConvert(tQuery) {
			var q gquery.Query
			if qCast, ok := key.Interface().(Query); ok {
				q = qCast.impl()
			} else if qCast, ok := key.Interface().(Query); ok {
				q = qCast.impl()
			} else {
				panic(EF(
					"unreachable, not a Query or Write: %s",
					key.String(),
				))
			}

			dat, err = dat.set(q, d)
			if err != nil {
				return nil, Err, newDataMakeUnmarshalError(err)
			}
		} else {
			return nil, Err, newDataMakeUnexpectedTypeError(
				key.Type(),
				internal.TStr,
				tQuery,
			)
		}
	}

	return dat.obj(), Obj, nil
}

func _ofSeq(
	_ any,
	r reflect.Value,
) (any, Type, error) {
	casted := make([]Datum, r.Len())

	for i := range r.Len() {
		d, err := of(r.Index(i).Interface())
		if err != nil {
			return nil, Err, err
		}

		casted[i] = d
	}

	return casted, Arr, nil
}

func _ofPtr(
	v any,
	r reflect.Value,
) (any, Type, error) {
	switch vv := v.(type) {
	case *big.Int:
		cp := big.NewInt(0)
		cp.Set(vv)

		return cp, Int, nil

	case *big.Float:
		cp := big.NewFloat(0)
		cp.Set(vv)

		return cp, Flt, nil

	case *Datum:
		d := M(of(vv.deref()))

		return d.deref(), d.typ, nil
	}

	if r.Kind() != reflect.Pointer {
		panic(EF("unreachable, expecting a pointer, got: %s", r.Kind().String()))
	}

	d, err := of(r.Elem())
	if err != nil {
		return nil, Err, err
	}

	return d.deref(), d.typ, nil
}

// ==============================================================================.

func isProhibited0(
	t reflect.Type,
	seen map[reflect.Type]z.NA,
) bool {
	if _, ok := _prohibited[t]; ok {
		return true
	}

	if t.Kind() != reflect.Pointer {
		return false
	}

	if _, ok := seen[t.Elem()]; ok {
		return false
	}

	seen[t] = z.None

	return isProhibited0(t.Elem(), seen)
}

func isProhibited(
	t reflect.Type,
) bool {
	return isProhibited0(t.Elem(), make(map[reflect.Type]z.NA))
}

// =====================================.

func newDataMakeUnexpectedTypeError(
	actual reflect.Type,
	expecting ...reflect.Type,
) error {
	msg := fmt.Sprintf("unexpected type: actual=%v", actual)
	if len(expecting) > 0 {
		msg += fmt.Sprintf(", expecting=%v", expecting)
	}

	return newDataMakeError(ErrCodeDataMakeUnexpectedType, msg)
}

func newDataMakeProhibitedTypeError(
	t reflect.Type,
) error {
	return newDataMakeError(
		ErrCodeDataMakeProhibitedType,
		"prohibited type: "+t.String(),
	)
}

func newDataMakeUnimplementedType(
	t reflect.Type,
) error {
	return newDataMakeError(
		ErrCodeDataMakeUnimplementedType,
		"type not implemented yet: "+t.String(),
	)
}

func newDataMakeInvalidTypeError(
	t reflect.Type,
) error {
	return newDataMakeError(
		ErrCodeDataMakeInvalidType,
		"invalid type: "+t.String(),
	)
}

func newDataMakeMarshalError(
	err error,
) error {
	return E(err, newDataMakeError(
		ErrCodeDataMakeSerializationFailure,
		"serialization failure",
	))
}

func newDataMakeUnmarshalError(
	err error,
) error {
	return E(err, newDataMakeError(
		ErrCodeDataMakeDeserializationFailure,
		"deserialization failure",
	))
}

func newDataMakeDuplicatedKeyError(
	q Query,
) error {
	return E(newDataMakeError(
		ErrCodeDataMakeDuplicateKey,
		"duplicated key: "+q.String(),
	))
}
