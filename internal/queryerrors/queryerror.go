package queryerrors

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

// ErrStart Keep in sync with errors.go in giraffe.
const ErrStart = uint64(math.MaxInt32)

const (
	ErrCodeEmpty = iota + ErrStart
	ErrCodeDuplicatedCmd
	ErrCodeConflictingCmd
	ErrCodeUnexpectedToken
	ErrCodeNestingTooDeep
	ErrCodeNotWritable
)

//goland:noinspection GoNameStartsWithPackageName
type qError struct {
	Msg  string
	Code uint64
}

func (e *qError) Error() string {
	return fmt.Sprintf("query error #%d: %s", e.Code, e.Msg)
}

func newError(
	code uint64,
	msg string,
) error {
	return E(&qError{
		Code: code,
		Msg:  msg,
	})
}

func NewNotWritableError(
	q string,
) error {
	return newError(
		ErrCodeNotWritable,
		"wrong modifier: need=read, have=write, query="+q,
	)
}

func NewParseError(
	code uint64,
	at int,
	spec string,
	msg string,
	extra ...string,
) error {
	sb := strings.Builder{}
	sb.Grow(len(spec) + len(msg) + 16)

	sb.WriteString("query parse error, ")
	sb.WriteString(msg)
	sb.WriteString(": at=")
	sb.WriteString(strconv.Itoa(at))
	sb.WriteString(", query=")
	sb.WriteString(spec)

	for _, e := range extra {
		sb.WriteString(", ")
		sb.WriteString(e)
	}

	return newError(code, sb.String())
}

func UnexpectedTokenError(
	at int,
	spec string,
	actual byte,
) error {
	return NewParseError(
		ErrCodeUnexpectedToken,
		at,
		spec,
		"expected token not seen",
		"actual="+string(actual),
	)
}

func ConflictingCmdError(
	at int,
	spec string,
	conflictWith byte,
) error {
	return NewParseError(
		ErrCodeConflictingCmd,
		at,
		spec,
		"conflicting command",
		"cmd="+string(conflictWith),
	)
}

func DuplicatedCmdError(
	at int,
	spec string,
	actual byte,
) error {
	return NewParseError(
		ErrCodeDuplicatedCmd,
		at,
		spec,
		"duplicated command",
		"cmd="+string(actual),
	)
}

func EmptyError(
	at int,
	spec string,
) error {
	return NewParseError(
		ErrCodeEmpty,
		at,
		spec,
		"query is empty",
	)
}

func NestingTooDeepError(
	at int,
	spec string,
) error {
	return NewParseError(
		ErrCodeNestingTooDeep,
		at,
		spec,
		"query nesting is too deep",
	)
}
