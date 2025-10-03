package queryimpl

import (
	"fmt"
	"math"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/qflag"
)

// ErrStart Keep in sync with errors.go in giraffe.
const ErrStart = uint64(math.MaxInt32)

const (
	ErrCodeQueryParseEmptyQuery = iota + ErrStart
	ErrCodeQueryParseDuplicatedCmd
	ErrCodeQueryParseConflictingCmd
	ErrCodeQueryParseUnexpectedToken
	ErrCodeQueryParseUnexpectedSegments
	ErrCodeQueryParseNestingTooDeep
	ErrCodeQueryParseNotWritable
	ErrCodeQueryParseUnclosedExtern
)

var ErrQ = newQuery(nil, "", qflag.QFlag(0))

//goland:noinspection GoNameStartsWithPackageName
type queryError struct {
	Msg  string
	Code uint64
}

func (e *queryError) Error() string {
	return fmt.Sprintf("query error #%d: %s", e.Code, e.Msg)
}

func newQueryError(
	code uint64,
	msg string,
) error {
	return E(&queryError{
		Code: code,
		Msg:  msg,
	})
}

func newQueryNotWritableError(
	q Query,
) error {
	return newQueryError(
		ErrCodeQueryParseNotWritable,
		fmt.Sprintf(
			"wrong query modifier: need=read, have=write, query=%s",
			q,
		),
	)
}
