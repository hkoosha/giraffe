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

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/internal"
	. "github.com/hkoosha/giraffe/internal/dot"
	"github.com/hkoosha/giraffe/internal/gdatum"
	"github.com/hkoosha/giraffe/internal/gquery"
	"github.com/hkoosha/giraffe/zebra/z"
)

func Make[V any](v V) (Datum, error) {
	return of(v)
}

func Of[V SafeType1](v V) Datum {
	return M(Make(v))
}

func Of1[V SafeType1](
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

func OfEmpty() Datum {
	return emptyObj
}

func OfEmptyArr() Datum {
	return emptyArr
}

func OfErr() Datum {
	return errD
}

type SafeType0 interface {
	// The string type is not safe for tier 0.
	// Same for Query type.

	Datum |
		bool |
		int |
		uint |
		int8 |
		int16 |
		int32 |
		int64 |
		uint8 |
		uint16 |
		uint32 |
		uint64 |
		*big.Int
}

type SafeType1 interface {
	string |
		SafeType0 |
		[]Datum |
		[]bool |
		[]string |
		[]int |
		[]int64 |
		[]uint64 |
		[][]Datum |
		[][]bool |
		[][]string |
		[][]int |
		[][]int64 |
		[][]uint64 |
		map[string]Datum |
		map[string]bool |
		map[string]string |
		map[string]int |
		map[string]int64 |
		map[string]uint64 |
		map[Query]Datum |
		map[Query]bool |
		map[Query]string |
		map[Query]int |
		map[Query]int64 |
		map[Query]uint64 |
		map[string][]Datum |
		map[string][]bool |
		map[string][]string |
		map[string][]int |
		map[string][]int64 |
		map[string][]uint64 |
		map[Query][]Datum |
		map[Query][]bool |
		map[Query][]string |
		map[Query][]int |
		map[Query][]int64 |
		map[Query][]uint64
}

type Implode = map[Query]Datum

// ============================================================================.

type Datum struct {
	val   *any
	Debug gdatum.DatumDebug
	typ   Type
}

func (d Datum) String() string {
	if g11y.IsDebugToString() {
		return fmt.Sprintf("Dat[%s]", d.typ.String())
	}

	return d.String0()
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

func (d Datum) MarshalJSON() ([]byte, error) {
	n, err := d.raw()
	if err != nil {
		return nil, err
	}

	return json.Marshal(n)
}

func (d Datum) MarshalJSONTo(w io.Writer) error {
	n, err := d.raw()
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(n)
}

func (d Datum) UnmarshalJSON([]byte) error {
	panic(EF("unimplemented: unmarshal json"))
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
	right Datum,
) (Datum, error) {
	return d.merge(right, []string{})
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
		Q(CmdAppend),
		value,
	)
}

func (d Datum) Query(
	q string,
) (Datum, error) {
	query, err := Parse(q)
	if err != nil {
		return errD, err
	}

	return d.get(query.impl())
}

func (d Datum) Nest(
	q Query,
) (Datum, error) {
	return d.nest(q)
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
	k, err := Parse(strconv.Itoa(index))
	if err != nil {
		return OfErr(), err
	}

	return d.get(k.impl())
}

func (d Datum) Has(
	query Query,
) bool {
	return d.has(query.impl())
}

func (d Datum) Tree() []Query {
	return z.Applied(d.tree(), func(it gquery.Query) Query {
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

func (d Datum) Len() (int, error) {
	return d.tryLen()
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

// =====================================.

func (d Datum) I08() (int8, error) {
	v, err := d.Int()
	if err != nil {
		return 0, err
	}

	vI := v.Int64()

	if !v.IsInt64() || vI < math.MinInt8 || math.MaxInt8 < vI {
		return 0, newDataReadIntegerOverflowError(internal.TI8)
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
		return 0, newDataReadIntegerOverflowError(internal.TI16)
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
		return 0, newDataReadIntegerOverflowError(internal.TI32)
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
		return 0, newDataReadIntegerOverflowError(internal.TI64)
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
		return 0, newDataReadIntegerOverflowError(internal.TU8)
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
		return 0, newDataReadIntegerOverflowError(internal.TU16)
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
		return 0, newDataReadIntegerOverflowError(internal.TU32)
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
		return 0, newDataReadIntegerOverflowError(internal.TU64)
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
		return 0, newDataReadIntegerOverflowError(internal.TISize)
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
		return 0, newDataReadIntegerOverflowError(internal.TUSize)
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
