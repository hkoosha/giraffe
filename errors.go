package giraffe

import (
	"fmt"
	"reflect"

	. "github.com/hkoosha/giraffe/internal/dot"
	"github.com/hkoosha/giraffe/internal/g"
	"github.com/hkoosha/giraffe/internal/gquery"
)

const (
	ErrCodeUnexpectedNil uint64 = iota + 1
	ErrCodeInvalidDatum
	ErrCodeCastError
	ErrCodeOverflowError

	ErrCodeQueryParseEmptyQuery
	ErrCodeQueryParseDuplicatedCmd
	ErrCodeQueryParseConflictingCmd
	ErrCodeQueryParseUnexpectedToken
	ErrCodeQueryParseUnexpectedSegments
	ErrCodeQueryParseNestingTooDeep
	ErrCodeQueryParseNotWritable

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

//nolint:reassign
func init() {
	gquery.ErrCodeQueryParseEmptyQuery = ErrCodeQueryParseEmptyQuery
	gquery.ErrCodeQueryParseDuplicatedCmd = ErrCodeQueryParseDuplicatedCmd
	gquery.ErrCodeQueryParseConflictingCmd = ErrCodeQueryParseConflictingCmd
	gquery.ErrCodeQueryParseUnexpectedToken = ErrCodeQueryParseUnexpectedToken
	gquery.ErrCodeQueryParseUnexpectedSegments = ErrCodeQueryParseUnexpectedSegments
	gquery.ErrCodeQueryParseNestingTooDeep = ErrCodeQueryParseNestingTooDeep
	gquery.ErrCodeQueryParseNotWritable = ErrCodeQueryParseNotWritable
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
	need gquery.QFlag,
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
