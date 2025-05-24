package giraffe

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hkoosha/giraffe/cmd"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/internal/gstrings"
	"github.com/hkoosha/giraffe/internal/queryerrors"
)

//goland:noinspection GoUnusedConst
const (
	ErrCodeUnexpectedNil uint64 = iota + 1
	ErrCodeInvalidDatum
	ErrCodeInvalidJsonable
	ErrCodeCastError
	ErrCodeOverflowError

	ErrCodeDataMakeUnexpectedType
	ErrCodeDataMakeUnimplementedType
	ErrCodeDataMakeProhibitedType
	ErrCodeDataMakeInvalidType
	ErrCodeDataMakeSerializationFailure
	ErrCodeDataMakeDeserializationFailure
	ErrCodeDataMakeDuplicateKey

	ErrCodeDataMergeIncompatibleTypes
	ErrCodeDataMergeClashingKeys

	ErrCodeTypeParseError

	ErrCodeDataReadIndeterministicQuery
	ErrCodeDataReadIndexOutOfBounds
	ErrCodeDataReadMissingKey
	ErrCodeDataReadUnexpectedType
	ErrCodeDataReadOnly

	ErrCodeDataWriteMissingKey
	ErrCodeDataWriteImplicitOverwrite
	ErrCodeDataWriteUnexpectedValue
	ErrCodeDataWriteUnsegmentedQuery

	ErrCodeDataModifyOperationTakesNoValue
)

//goland:noinspection GoUnusedConst
const (
	ErrCodeQueryParseEmptyQuery      = queryerrors.ErrCodeEmpty
	ErrCodeQueryParseDuplicatedCmd   = queryerrors.ErrCodeDuplicatedCmd
	ErrCodeQueryParseConflictingCmd  = queryerrors.ErrCodeConflictingCmd
	ErrCodeQueryParseUnexpectedToken = queryerrors.ErrCodeUnexpectedToken
	ErrCodeQueryParseNestingTooDeep  = queryerrors.ErrCodeNestingTooDeep
	ErrCodeQueryParseNotWritable     = queryerrors.ErrCodeNotWritable
)

func init() {
}

func getDataErrRepr() string {
	const msg = "<err>"

	// X panic(TA(msg)).

	return msg
}

func newNilError() error {
	return newGiraffeError(
		ErrCodeUnexpectedNil,
		"unexpected nil",
	)
}

func newNotJsonableError() error {
	return newGiraffeError(
		ErrCodeInvalidJsonable,
		"given object cannot be transformed to of from json",
	)
}

func newInvalidDatumError() error {
	return E(newGiraffeError(
		ErrCodeInvalidDatum,
		"datum is in invalid state",
	))
}

// ==============================================================================.

//goland:noinspection GoNameStartsWithPackageName
type GiraffeError struct {
	msg  string
	code uint64
}

func (e *GiraffeError) Error() string {
	return e.msg
}

func (e *GiraffeError) Code() uint64 {
	return e.code
}

func newGiraffeError(
	code uint64,
	msg string,
) error {
	var err error = &GiraffeError{
		code: code,
		msg:  msg,
	}

	return E(err)
}

func newDataMakeError(
	code uint64,
	msg string,
) error {
	const prefix = "failed to make datum"

	return newGiraffeError(code, gstrings.Joined([]string{prefix, msg}))
}

func newTypeCastError(
	have Type,
	need Type,
	extra ...string,
) error {
	return newGiraffeError(
		ErrCodeCastError,
		fmt.Sprintf(
			"cannot cast: from=%s, to=%s%s",
			have,
			need,
			gstrings.Joined(extra),
		),
	)
}

func newQueryTypeCastError(
	have Type,
	need cmd.QFlag,
	extra ...string,
) error {
	needTyp := Arr
	if need.IsObj() {
		needTyp = Obj
	}

	return newTypeCastError(have, needTyp, extra...)
}

func newReflectiveTypeCastError(
	have Type,
	need any,
	extra ...string,
) error {
	return newGiraffeError(
		ErrCodeCastError,
		fmt.Sprintf(
			"cannot cast: from=%s, to=%s%s",
			have,
			reflect.TypeOf(need),
			gstrings.Joined(extra),
		),
	)
}

// =============================================================================

type MissingKeysError struct {
	keys []Query
}

func (e *MissingKeysError) Error() string {
	last := len(e.keys) - 1
	sb := strings.Builder{}
	sb.Grow(255 * (last + 1))

	sb.WriteString("missing keys: [")

	for i, k := range e.keys {
		if i != last {
			sb.WriteString(", ")
		}
		sb.WriteString(k.String())
	}
	sb.WriteByte(']')

	return sb.String()
}

func NewMissingKeyError(
	keys ...Query,
) error {
	return E(&MissingKeysError{
		keys: keys,
	})
}
