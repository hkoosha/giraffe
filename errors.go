package giraffe

import (
	"fmt"
	"reflect"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/g"
	"github.com/hkoosha/giraffe/internal/queryerrors"
	"github.com/hkoosha/giraffe/qflag"
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

//nolint:reassign
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

	return newGiraffeError(code, g.Join(prefix, msg))
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
			g.Joined(extra),
		),
	)
}

func newQueryTypeCastError(
	have Type,
	need qflag.QFlag,
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
			g.Joined(extra),
		),
	)
}
