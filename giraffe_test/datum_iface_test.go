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
)

type (
	QQQuery = giraffe.Query
	Datum   = giraffe.Datum
	Type    = giraffe.Type
)

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

type TypedQuery interface {
	QInt(QQQuery) (*big.Int, error)
	QFlt(QQQuery) (*big.Float, error)
	QBln(QQQuery) (bool, error)
	QStr(QQQuery) (string, error)
}

type CastQuery interface {
	QISz(QQQuery) (int, error)
	QI08(QQQuery) (int8, error)
	QI16(QQQuery) (int16, error)
	QI32(QQQuery) (int32, error)
	QI64(QQQuery) (int64, error)
	QUSz(QQQuery) (uint, error)
	QU08(QQQuery) (uint8, error)
	QU16(QQQuery) (uint16, error)
	QU32(QQQuery) (uint32, error)
	QU64(QQQuery) (uint64, error)
}

type CastArrayQuery interface {
	QStrs(QQQuery) ([]string, error)
	QISzs(QQQuery) ([]int, error)
	QI64s(QQQuery) ([]int64, error)
	QUSzs(QQQuery) ([]uint, error)
	QU64s(QQQuery) ([]uint64, error)
}

type Queried interface {
	Tree() []QQQuery
	Keys() ([]string, error)
	Len() (int, error)
	At(int) (Datum, error)
	Has(query QQQuery) bool
	Get(QQQuery) (Datum, error)
	Query(q string) (Datum, error)
}

type Modified interface {
	Set(QQQuery, any) (Datum, error)
	Merge(Datum) (Datum, error)
	Append(value any) (Datum, error)
}

type Formatted interface {
	fmt.Stringer
	String() string
	Pretty() string
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
	Typed
	TypedQuery
	Cast
	CastQuery
	CastArray
	CastArrayQuery
	Queried
	Modified
	Formatted
	Eq
	Rel
	Iter
}
