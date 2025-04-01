package giraffe

import (
	"math/bits"
	"strconv"
	"strings"

	"github.com/hkoosha/giraffe/g11y"
	. "github.com/hkoosha/giraffe/internal/dot"
)

// Types.
const (
	// ======================= Modifiers.

	// Nil must be exactly equal to 1, so foo%2==1 is nil and foo%2==0 is not.
	Nil Type = 0b0000_0001
	Arr Type = 0b0000_0010

	// =========================== Types.

	Obj Type = 0b0000_0010_0000_0000
	Int Type = 0b0000_0100_0000_0000
	Flt Type = 0b0000_1000_0000_0000
	Bln Type = 0b0001_0000_0000_0000
	Str Type = 0b0010_0000_0000_0000
)

// States.
const (
	types = Obj | Int | Flt | Bln | Str
	mods  = Nil | Arr
	valid = types | mods

	Err = Type((^uint64(0)) & (^uint64(valid)))
)

const (
	arrRepr = "[]"
	nilRepr = "?"
)

var dTypes = []Type{
	Obj,
	Int,
	Flt,
	Bln,
	Str,
}

//nolint:recvcheck
type Type uint64

//goland:noinspection GoMixedReceiverTypes
func (t Type) String() string {
	if g11y.IsDebugToString() {
		return "Typ=0b" + strings.TrimPrefix(strconv.FormatUint(uint64(t), 2), "0")
	}

	x := t
	if bits.OnesCount64(uint64(x)) > 1 {
		return "err@" + strings.TrimPrefix(strconv.FormatUint(uint64(t), 2), "0")
	}

	mod := ""
	if x.IsArr() {
		mod = arrRepr
		x &= ^Arr
	}

	if x.IsNil() {
		mod += nilRepr
		x &= ^Nil
	}

	var typ string

	switch {
	case x.Is(Obj):
		typ = "obj"
	case x.Is(Int):
		typ = "int"
	case x.Is(Flt):
		typ = "flt"
	case x.Is(Bln):
		typ = "bln"
	case x.Is(Str):
		typ = "str"
	case x > 0:
		return "err@" + strings.TrimPrefix(strconv.FormatUint(uint64(t), 2), "0")
	}

	return typ + mod
}

//goland:noinspection GoMixedReceiverTypes
func (t *Type) MarshalText() ([]byte, error) {
	if t == nil {
		return nil, newNilError()
	}

	s := t.String()

	return []byte(s), nil
}

//goland:noinspection GoMixedReceiverTypes
func (t *Type) UnmarshalJSON(b []byte) error {
	if t == nil {
		return newNilError()
	}

	str := string(b)
	parse := ParseType(str)

	if parse.IsErr() {
		return newTypeParseError(str)
	}

	*t = parse

	return nil
}

func ParseType(repr string) Type {
	x := Type(0)
	r := repr

	if strings.HasSuffix(r, nilRepr) {
		x |= Nil
		r = r[:len(r)-1]
	}

	if strings.HasSuffix(r, arrRepr) {
		x |= Arr
		r = r[:len(r)-2]
	}

	for _, t := range dTypes {
		if t.String() == r {
			return x | t
		}
	}

	return Err
}

// =====================================.

//goland:noinspection GoMixedReceiverTypes
func (t Type) WithArr() Type {
	return t | Arr
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) WithNil() Type {
	return t | Nil
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) IsNil() bool {
	return t&Nil == Nil
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) IsArr() bool {
	return t&Arr == Arr
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) IsErr() bool {
	return t&Err == Err
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) Is(o Type) bool {
	return (t&types) == (o&types) &&
		(t&Arr) == (o&Arr) &&
		((t&Nil) == 0) || ((o & Nil) == Nil)
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) isZero() bool {
	return t == 0
}

// =====================================.

//goland:noinspection GoMixedReceiverTypes
func (t Type) IsInt() bool {
	return t.Is(Int)
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) IsFlt() bool {
	return t.Is(Flt)
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) IsBln() bool {
	return t.Is(Bln)
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) IsStr() bool {
	return t.Is(Str)
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) IsObj() bool {
	return t.Is(Obj)
}

// =====================================.

func newTypeParseError(
	t string,
) error {
	return E(newGiraffeError(
		ErrCodeTypeParseError,
		"type parse error, type="+t,
	))
}
