package giraffe_test

import (
	"fmt"
	"io"
	"iter"
	"math/big"

	"github.com/hkoosha/giraffe"
)

var (
	//nolint:exhaustruct
	d = Datum{}

	_ Typed          = d
	_ Cast           = d
	_ CastArray      = d
	_ TypedQuery     = d
	_ CastQuery      = d
	_ CastArrayQuery = d
	_ Formatted      = d
	_ Queried        = d
	_ Modified       = d
	_ DatumPub       = d
	_ Dyn            = d
)

type (
	Query = giraffe.Query
	Datum = giraffe.Datum
	Type  = giraffe.Type
)

// ====

type Dyn interface {
	Has(Query) (bool, error)
	Set(Query, any) (Datum, error)
	Get(Query) (Datum, error)
}

type TypedQuery interface {
	QInt(Query) (*big.Int, error)
	QFlt(Query) (*big.Float, error)
	QBln(Query) (bool, error)
	QStr(Query) (string, error)
}

type CastQuery interface {
	QISz(Query) (int, error)
	QI08(Query) (int8, error)
	QI16(Query) (int16, error)
	QI32(Query) (int32, error)
	QI64(Query) (int64, error)
	QUSz(Query) (uint, error)
	QU08(Query) (uint8, error)
	QU16(Query) (uint16, error)
	QU32(Query) (uint32, error)
	QU64(Query) (uint64, error)
}

type CastArrayQuery interface {
	QStrs(Query) ([]string, error)
	QISzs(Query) ([]int, error)
	QI64s(Query) ([]int64, error)
	QUSzs(Query) ([]uint, error)
	QU64s(Query) ([]uint64, error)
}

type Typed interface {
	Type() Type

	Int() (*big.Int, error)
	Flt() (*big.Float, error)
	Bln() (bool, error)
	Str() (string, error)

	Raw() (any, error)
}

type Cast interface {
	ISz() (int, error)
	I08() (int8, error)
	I16() (int16, error)
	I32() (int32, error)
	I64() (int64, error)
	USz() (uint, error)
	U08() (uint8, error)
	U16() (uint16, error)
	U32() (uint32, error)
	U64() (uint64, error)
}

type CastArray interface {
	Strs() ([]string, error)
	ISzs() ([]int, error)
	USzs() ([]uint, error)
	I64s() ([]int64, error)
	U64s() ([]uint64, error)
}

type Queried interface {
	Tree() []Query
	Keys() ([]string, error)
	Kv() (map[string]string, error)
	Len() (int, error)
	HasLen() bool
	At(int) (Datum, error)
}

type Modified interface {
	Merge(...Datum) (Datum, error)
	Append(any) (Datum, error)
}

type Formatted interface {
	fmt.Stringer
	String() string
	Pretty() string
	Plain() (any, error)
	MarshalJSON() ([]byte, error)
	MarshalJSONTo(w io.Writer) error
}

type Eq interface {
	Eq(Datum) bool
}

type Rel interface {
	Gt(Datum) (bool, error)
	Lt(Datum) (bool, error)
	Gte(Datum) (bool, error)
	Lte(Datum) (bool, error)
}

type Iter interface {
	Iter() (iter.Seq[Datum], error)
	Iter2() (iter.Seq2[string, Datum], error)
}

type DatumPub interface {
	TypedQuery
	CastQuery
	CastArrayQuery
	Dyn

	Typed
	Cast
	CastArray
	Queried
	Modified
	Formatted
	Eq
	Rel
	Iter
}
