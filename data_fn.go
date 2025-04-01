package giraffe

import (
	"math/big"
)

func ToU08(it Datum) (uint8, error) {
	return it.U08()
}

func ToU16(it Datum) (uint16, error) {
	return it.U16()
}

func ToU32(it Datum) (uint32, error) {
	return it.U32()
}

func ToU64(it Datum) (uint64, error) {
	return it.U64()
}

func ToI08(it Datum) (int8, error) {
	return it.I08()
}

func ToI16(it Datum) (int16, error) {
	return it.I16()
}

func ToI32(it Datum) (int32, error) {
	return it.I32()
}

func ToI64(it Datum) (int64, error) {
	return it.I64()
}

func ToIsz(it Datum) (int, error) {
	return it.ISz()
}

func ToUsz(it Datum) (uint, error) {
	return it.USz()
}

func ToInt(it Datum) (*big.Int, error) {
	return it.Int()
}

func ToFlt(it Datum) (*big.Float, error) {
	return it.Flt()
}

func ToStr(it Datum) (string, error) {
	return it.Str()
}

func ToBln(it Datum) (bool, error) {
	return it.Bln()
}

func ToArr(it Datum) (bool, error) {
	return it.Bln()
}
