package gquery

import (
	"fmt"

	. "github.com/hkoosha/giraffe/internal/dot"
)

var ErrQ = newQuery(nil, "", QFlag(0))

//goland:noinspection GoNameStartsWithPackageName
type QueryError struct {
	Msg  string
	Code uint64
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("query error #%d: %s", e.Code, e.Msg)
}

func newQueryError(
	code uint64,
	msg string,
) error {
	return E(&QueryError{
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
