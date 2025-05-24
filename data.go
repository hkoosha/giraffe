package giraffe

import (
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"maps"
	"math"
	"math/big"
	"slices"
	"strconv"

	"github.com/hkoosha/giraffe/cmd"
	"github.com/hkoosha/giraffe/core/serdes/gson"
	"github.com/hkoosha/giraffe/core/t11y"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/internal"
	"github.com/hkoosha/giraffe/internal/gdatum"
	"github.com/hkoosha/giraffe/internal/reflected"
	"github.com/hkoosha/giraffe/zebra/z"
)

func From[V any](v V) (Datum, error) {
	return of(v)
}

func FromJsonable(v any) (Datum, error) {
	return ofJsonable(v)
}

func OfJsonable(v any) Datum {
	return M(FromJsonable(v))
}

func Of[V Safe](v V) Datum {
	return M(From(v))
}

func Of1[V Safe](
	q Query,
	v V,
) Datum {
	m := map[Query]V{
		q.WithMake(): v,
	}

	return M(of(m))
}

func OfN(
	pairs ...Tuple,
) (Datum, error) {
	m := make(Implode, len(pairs))

	for _, pair := range pairs {
		if _, ok := m[pair.Query]; ok {
			return OfErr(), newDataMakeDuplicatedKeyError(pair.Query)
		}
		m[pair.Query] = pair.Dat
	}

	return of(m)
}

func OfKV[V Safe](
	k string,
	v V,
) Datum {
	return OfJsonable(map[string]V{
		k: v,
	})
}

func OfEmpty() Datum {
	return emptyObj
}

func OfEmptyArr() Datum {
	return emptyArr
}

func OfErr() Datum {
	return errD
}

type Num interface {
	int |
		uint |
		int8 |
		int16 |
		int32 |
		int64 |
		uint8 |
		uint16 |
		uint32 |
		uint64
}

type Ord interface {
	Num | string
}

type Basic interface {
	Ord | bool
}

type Seq interface {
	[]bool |
		[]string |
		[]int |
		[]int64 |
		[]uint64 |
		[][]bool |
		[][]string |
		[][]int |
		[][]int64 |
		[][]uint64 |
		map[string]bool |
		map[string]string |
		map[string]int |
		map[string]int64 |
		map[string][]bool |
		map[string][]string |
		map[string][]int |
		map[string][]int64 |
		map[string][]uint64
}

type Safe interface {
	Basic | Seq |
		*big.Int |
		Datum |
		[]Datum |
		[][]Datum |
		map[string]Datum |
		map[Query]Datum |
		map[Query]bool |
		map[Query]string |
		map[Query]int |
		map[Query]int64 |
		map[Query]uint64 |
		map[string][]Datum |
		map[Query][]Datum |
		map[Query][]bool |
		map[Query][]string |
		map[Query][]int |
		map[Query][]int64 |
		map[Query][]uint64
}

type Implode = map[Query]Datum

// ============================================================================.

//nolint:recvcheck
type Datum struct {
	val   *any
	Debug gdatum.DatumDebug
	typ   Type
}

func (d Datum) SimpleString() (string, error) {
	switch {
	case d.typ.IsObj(), d.typ.IsArr():

	case d.typ.IsBln():
		switch {
		case M(d.Bln()):
			return "true", nil
		default:
			return "false", nil
		}

	case d.typ.IsFlt():
		return M(d.Flt()).String(), nil

	case d.typ.IsInt():
		return M(d.Int()).String(), nil

	case d.typ.IsStr():
		return M(d.Str()), nil
	}

	return "", EF("type does not support simple string formatting: %s", d.typ.String())
}

func (d Datum) String() string {
	if t11y.IsDebugToString() {
		return fmt.Sprintf("Dat[%s]", d.typ.String())
	}

	return d.string0()
}

func (d Datum) Pretty() string {
	if d.typ.isZero() {
		return ""
	}

	r, err := d.Raw()
	if err != nil {
		return getDataErrRepr()
	}

	rJ, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return getDataErrRepr()
	}

	return string(rJ)
}

func (d Datum) Raw() (any, error) {
	return d.raw()
}

func (d Datum) Plain() (any, error) {
	return d.plain()
}

func (d Datum) MarshalJSON() ([]byte, error) {
	n, err := d.raw()
	if err != nil {
		return nil, err
	}

	return gson.Marshal(n)
}

func (d Datum) MarshalJsonString() string {
	return string(M(gson.Marshal(M(d.raw()))))
}

func (d Datum) MarshalJSONTo(w io.Writer) error {
	n, err := d.raw()
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(n)
}

func (d *Datum) UnmarshalJSON(b []byte) error {
	var err error

	*d, err = ofJson(b)
	return err
}

func (d Datum) Type() Type {
	return d.typ
}

func (d Datum) Eq(other Datum) bool {
	return d.eq(other)
}

func (d Datum) Lt(other Datum) (bool, error) {
	cmp, err := d.cmp(other)
	if err != nil {
		return false, err
	}

	return cmp < 0, nil
}

func (d Datum) Gt(other Datum) (bool, error) {
	cmp, err := d.cmp(other)
	if err != nil {
		return false, err
	}

	return cmp > 0, nil
}

func (d Datum) Gte(other Datum) (bool, error) {
	cmp, err := d.cmp(other)
	if err != nil {
		return false, err
	}

	return cmp >= 0, nil
}

func (d Datum) Lte(other Datum) (bool, error) {
	cmp, err := d.cmp(other)
	if err != nil {
		return false, err
	}

	return cmp <= 0, nil
}

func (d Datum) Merge(
	right ...Datum,
) (Datum, error) {
	return d.merge(right)
}

func (d Datum) Set(
	query Query,
	value any,
) (Datum, error) {
	return d.set(query.impl(), value)
}

func (d Datum) Append(
	value any,
) (Datum, error) {
	return d.Set(
		Q(cmd.Append.String()),
		value,
	)
}

func (d Datum) Nest(
	q Query,
) (Datum, error) {
	return d.nest(q.impl())
}

func (d Datum) Get(
	q Query,
) (Datum, error) {
	return d.get(q.impl())
}

func (d Datum) Iter() (iter.Seq[Datum], error) {
	val, err := d.tryArr()
	if err != nil {
		return nil, err
	}

	return slices.Values(val), nil
}

// Iter2 TODO must unescape keys.
func (d Datum) Iter2() (iter.Seq2[string, Datum], error) {
	val, err := d.tryObj()
	if err != nil {
		return nil, err
	}

	return func(yield func(string, Datum) bool) {
		for k, v := range val {
			yield(k, v)
		}
	}, nil
}

func (d Datum) At(
	index int,
) (Datum, error) {
	k, err := GQParse(strconv.Itoa(index))
	if err != nil {
		return OfErr(), err
	}

	return d.get(k.impl())
}

func (d Datum) Has(
	query Query,
) (bool, error) {
	return d.has(query.impl())
}

func (d Datum) Tree() []Query {
	return z.Applied(d.tree(), func(it queryT) Query {
		return Query(it.String())
	})
}

func (d Datum) Keys() ([]string, error) {
	val, err := d.tryObj()
	if err != nil {
		return nil, err
	}

	return slices.Collect(maps.Keys(val)), nil
}

func (d Datum) Kv() (map[string]string, error) {
	keys, err := d.Keys()
	if err != nil {
		return nil, err
	}

	kv := make(map[string]string, len(keys))

	for _, k := range keys {
		str, err := d.Get(Q(internal.Escaped(k)))
		if err != nil {
			return nil, err
		}

		v, err := str.Str()
		if err != nil {
			return nil, err
		}

		kv[k] = v
	}

	return kv, nil
}

func (d Datum) Len() (int, error) {
	return d.tryLen()
}

func (d Datum) HasLen() bool {
	return d.typ.IsObj() || d.typ.IsArr()
}

// =====================================.

func (d Datum) Int() (*big.Int, error) {
	switch {
	case d.typ.IsNil():
		return nil, newNilError()

	case !d.typ.IsInt():
		return nil, newTypeCastError(d.typ, Int)

	default:
		cp := big.NewInt(0)
		cp.Set(cast[*big.Int](d))

		return cp, nil
	}
}

func (d Datum) Flt() (*big.Float, error) {
	switch {
	case d.typ.IsNil():
		return nil, newNilError()

	case !d.typ.IsFlt():
		return nil, newTypeCastError(d.typ, Flt)

	default:
		cp := big.NewFloat(0)
		cp.Set(cast[*big.Float](d))

		return cp, nil
	}
}

func (d Datum) Bln() (bool, error) {
	switch {
	case d.typ.IsNil():
		return false, newNilError()

	case !d.typ.IsBln():
		return false, newTypeCastError(d.typ, Bln)

	default:
		return cast[bool](d), nil
	}
}

func (d Datum) Str() (string, error) {
	switch {
	case d.typ.IsNil():
		return "", newNilError()

	case !d.typ.IsStr():
		return "", newTypeCastError(d.typ, Str)

	default:
		return cast[string](d), nil
	}
}

func (d Datum) FmtStr() (string, error) {
	switch {
	case d.typ.IsStr():
		return cast[string](d), nil

	case d.typ.IsBln():
		return strconv.FormatBool(cast[bool](d)), nil

	case d.typ.IsInt():
		return M(d.Int()).String(), nil

	default:
		return "", EF("cannot format datatype as simple string: %v", d.String())
	}
}

// =====================================.

func (d Datum) I08() (int8, error) {
	v, err := d.Int()
	if err != nil {
		return 0, err
	}

	vI := v.Int64()

	if !v.IsInt64() || vI < math.MinInt8 || math.MaxInt8 < vI {
		return 0, newDataReadIntegerOverflowError(reflected.TI8)
	}

	//nolint:gosec
	return int8(vI), nil
}

func (d Datum) I16() (int16, error) {
	v, err := d.Int()
	if err != nil {
		return 0, err
	}

	vI := v.Int64()

	if !v.IsInt64() || vI < math.MinInt16 || math.MaxInt16 < vI {
		return 0, newDataReadIntegerOverflowError(reflected.TI16)
	}

	//nolint:gosec
	return int16(vI), nil
}

func (d Datum) I32() (int32, error) {
	v, err := d.Int()
	if err != nil {
		return 0, err
	}

	vI := v.Int64()

	if !v.IsInt64() || vI < math.MinInt32 || math.MaxInt32 < vI {
		return 0, newDataReadIntegerOverflowError(reflected.TI32)
	}

	//nolint:gosec
	return int32(vI), nil
}

func (d Datum) I64() (int64, error) {
	v, err := d.Int()
	if err != nil {
		return 0, err
	}

	vI := v.Int64()

	if !v.IsInt64() || vI < math.MinInt64 || math.MaxInt64 < vI {
		return 0, newDataReadIntegerOverflowError(reflected.TI64)
	}

	return vI, nil
}

func (d Datum) U08() (uint8, error) {
	v, err := d.Int()
	if err != nil {
		return 0, err
	}

	vI := v.Int64()

	if !v.IsInt64() || vI < 0 || vI > math.MaxUint8 {
		return 0, newDataReadIntegerOverflowError(reflected.TU8)
	}

	return uint8(vI), nil
}

func (d Datum) U16() (uint16, error) {
	v, err := d.Int()
	if err != nil {
		return 0, err
	}

	vI := v.Int64()

	if !v.IsInt64() || vI < 0 || vI > math.MaxUint16 {
		return 0, newDataReadIntegerOverflowError(reflected.TU16)
	}

	return uint16(vI), nil
}

func (d Datum) U32() (uint32, error) {
	v, err := d.Int()
	if err != nil {
		return 0, err
	}

	vI := v.Int64()

	if !v.IsInt64() || vI < 0 || vI > math.MaxUint32 {
		return 0, newDataReadIntegerOverflowError(reflected.TU32)
	}

	return uint32(vI), nil
}

func (d Datum) U64() (uint64, error) {
	v, err := d.Int()
	if err != nil {
		return 0, err
	}

	vI := v.Int64()

	if !v.IsInt64() || vI < 0 {
		return 0, newDataReadIntegerOverflowError(reflected.TU64)
	}

	return uint64(vI), nil
}

func (d Datum) ISz() (int, error) {
	v, err := d.Int()
	if err != nil {
		return 0, err
	}

	vI := v.Int64()

	if !v.IsInt64() || vI > math.MaxInt {
		return 0, newDataReadIntegerOverflowError(reflected.TISize)
	}

	return int(vI), nil
}

func (d Datum) USz() (uint, error) {
	v, err := d.Int()
	if err != nil {
		return 0, err
	}

	vI := v.Int64()

	if !v.IsInt64() || vI < 0 {
		return 0, newDataReadIntegerOverflowError(reflected.TUSize)
	}

	return uint(vI), nil
}

func (d Datum) USzs() ([]uint, error) {
	return copyToArr(d, ToUsz)
}

func (d Datum) U64s() ([]uint64, error) {
	return copyToArr(d, ToU64)
}

func (d Datum) ISzs() ([]int, error) {
	return copyToArr(d, ToIsz)
}

func (d Datum) I64s() ([]int64, error) {
	return copyToArr(d, ToI64)
}

func (d Datum) Strs() ([]string, error) {
	return copyToArr(d, ToStr)
}

// =====================================.

func (d Datum) QUSz(q Query) (uint, error) {
	get, err := d.Get(q)
	if err != nil {
		return 0, err
	}

	return get.USz()
}

func (d Datum) QU08(q Query) (uint8, error) {
	get, err := d.Get(q)
	if err != nil {
		return 0, err
	}

	return get.U08()
}

func (d Datum) QU16(q Query) (uint16, error) {
	get, err := d.Get(q)
	if err != nil {
		return 0, err
	}

	return get.U16()
}

func (d Datum) QU32(q Query) (uint32, error) {
	get, err := d.Get(q)
	if err != nil {
		return 0, err
	}

	return get.U32()
}

func (d Datum) QU64(q Query) (uint64, error) {
	get, err := d.Get(q)
	if err != nil {
		return 0, err
	}

	return get.U64()
}

func (d Datum) QISz(q Query) (int, error) {
	get, err := d.Get(q)
	if err != nil {
		return 0, err
	}

	return get.ISz()
}

func (d Datum) QI08(q Query) (int8, error) {
	get, err := d.Get(q)
	if err != nil {
		return 0, err
	}

	return get.I08()
}

func (d Datum) QI16(q Query) (int16, error) {
	get, err := d.Get(q)
	if err != nil {
		return 0, err
	}

	return get.I16()
}

func (d Datum) QI32(q Query) (int32, error) {
	get, err := d.Get(q)
	if err != nil {
		return 0, err
	}

	return get.I32()
}

func (d Datum) QI64(q Query) (int64, error) {
	get, err := d.Get(q)
	if err != nil {
		return 0, err
	}

	return get.I64()
}

func (d Datum) QUSzs(q Query) ([]uint, error) {
	get, err := d.Get(q)
	if err != nil {
		return nil, err
	}

	return get.USzs()
}

func (d Datum) QU64s(q Query) ([]uint64, error) {
	get, err := d.Get(q)
	if err != nil {
		return nil, err
	}

	return get.U64s()
}

func (d Datum) QISzs(q Query) ([]int, error) {
	get, err := d.Get(q)
	if err != nil {
		return nil, err
	}

	return get.ISzs()
}

func (d Datum) QI64s(q Query) ([]int64, error) {
	get, err := d.Get(q)
	if err != nil {
		return nil, err
	}

	return get.I64s()
}

func (d Datum) QStrs(q Query) ([]string, error) {
	get, err := d.Get(q)
	if err != nil {
		return nil, err
	}

	return get.Strs()
}

func (d Datum) QInt(q Query) (*big.Int, error) {
	get, err := d.Get(q)
	if err != nil {
		return nil, err
	}

	return get.Int()
}

func (d Datum) QFlt(q Query) (*big.Float, error) {
	get, err := d.Get(q)
	if err != nil {
		return nil, err
	}

	return get.Flt()
}

func (d Datum) QBln(q Query) (bool, error) {
	get, err := d.Get(q)
	if err != nil {
		return false, err
	}

	return get.Bln()
}

func (d Datum) QStr(q Query) (string, error) {
	get, err := d.Get(q)
	if err != nil {
		return "", err
	}

	return get.Str()
}

func (d Datum) QFmtStr(q Query) (string, error) {
	get, err := d.Get(q)
	if err != nil {
		return "", err
	}

	return get.FmtStr()
}

func (d Datum) QKv(q Query) (map[string]string, error) {
	get, err := d.Get(q)
	if err != nil {
		return nil, err
	}

	return get.Kv()
}
