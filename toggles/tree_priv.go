package toggles

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/zebra/z"
)

const (
	opAnd = "&&"
	opOr  = "||"
)

func join(
	rest []Condition,
	op string,
) Condition {
	if op == opAnd {
		rest = z.Filtered(rest, func(it Condition) bool {
			return !it.isYes()
		})
	} else if op == opOr {
		rest = z.Filtered(rest, func(it Condition) bool {
			return !it.isNo()
		})
	}

	if len(rest) == 0 {
		return &yes{}
	} else if len(rest) == 1 {
		return rest[0]
	}

	uid := z.Applied(rest, func(it Condition) string {
		return it.uid()
	})

	slices.Sort(uid)

	uid_ := "(" + strings.Join(uid, op) + ")"

	switch op {
	case opAnd:
		return &and{
			rest: rest,
			uid_: uid_,
		}

	case opOr:
		return &or{
			rest: rest,
			uid_: uid_,
		}

	default:
		panic(EF("unknown operation '%s'", op))
	}
}

func andOf(rest []Condition) Condition {
	return join(rest, "&&")
}

func orOf(rest []Condition) Condition {
	return join(rest, "||")
}

func notOf(c Condition) Condition {
	return &not{rest: c}
}

type seal struct{}

// ====================================.

func newAttr(
	name string,
	op Op,
	value any,
) Attr {
	sum := md5.Sum(M(json.Marshal(value)))
	return &attrImpl{
		name:  name,
		op:    op,
		value: value,
		uid_: fmt.Sprintf(
			"attr::%s::%s",
			name,
			hex.EncodeToString(sum[:]),
		),
	}
}

type attrImpl struct {
	name  string
	op    Op
	value any
	uid_  string
}

func (q *attrImpl) Name() string {
	return q.name
}

func (q *attrImpl) Test(attrs ...Attr) bool {
	myUid := q.uid()

	for _, attr := range attrs {
		if attr.uid() == myUid {
			return true
		}
	}

	return false
}

func (q *attrImpl) And(rest ...Condition) Condition {
	return andOf(append(rest, q))
}

func (q *attrImpl) Or(rest ...Condition) Condition {
	return orOf(append(rest, q))
}

func (q *attrImpl) Not() Condition {
	return notOf(q)
}

func (q *attrImpl) uid() string {
	return q.uid_
}

func (q *attrImpl) seal() seal {
	panic("do not call")
}

func (q *attrImpl) isYes() bool {
	return false
}

func (q *attrImpl) isNo() bool {
	return false
}

// ====================================.

type and struct {
	rest []Condition
	uid_ string
}

func (q *and) Test(attrs ...Attr) bool {
	if len(q.rest) == 0 {
		return true
	}

	for _, c := range q.rest {
		if !c.Test(attrs...) {
			return false
		}
	}

	return true
}

func (q *and) And(rest ...Condition) Condition {
	return andOf(append(rest, q))
}

func (q *and) Or(rest ...Condition) Condition {
	return orOf(append(rest, q))
}

func (q *and) Not() Condition {
	return notOf(q)
}

func (q *and) seal() seal {
	panic("do not call")
}

func (q *and) uid() string {
	return q.uid_
}

func (q *and) isYes() bool {
	return false
}

func (q *and) isNo() bool {
	return false
}

// ====================================.

type or struct {
	rest []Condition
	uid_ string
}

func (q *or) Test(attrs ...Attr) bool {
	if len(q.rest) == 0 {
		return true
	}

	for _, c := range q.rest {
		if c.Test(attrs...) {
			return true
		}
	}

	return false
}

func (q *or) And(rest ...Condition) Condition {
	return andOf(append(rest, q))
}

func (q *or) Or(rest ...Condition) Condition {
	return orOf(append(rest, q))
}

func (q *or) Not() Condition {
	return notOf(q)
}

func (q *or) seal() seal {
	panic("do not call")
}

func (q *or) uid() string {
	return q.uid_
}

func (q *or) isYes() bool {
	return false
}

func (q *or) isNo() bool {
	return false
}

// ====================================.

type not struct {
	rest Condition
}

func (q *not) Test(attrs ...Attr) bool {
	return !q.rest.Test(attrs...)
}

func (q *not) And(rest ...Condition) Condition {
	return andOf(append(rest, q))
}

func (q *not) Or(rest ...Condition) Condition {
	return orOf(append(rest, q))
}

func (q *not) Not() Condition {
	return q.rest
}

func (q *not) seal() seal {
	panic("do not call")
}

func (q *not) uid() string {
	return "not::" + q.rest.uid()
}

func (q *not) isYes() bool {
	return false
}

func (q *not) isNo() bool {
	return false
}

// ====================================.

type no struct {
}

func (q *no) Test(...Attr) bool {
	return false
}

func (q *no) And(...Condition) Condition {
	return q
}

func (q *no) Or(rest ...Condition) Condition {
	return orOf(rest)
}

func (q *no) Not() Condition {
	return &yes{}
}

func (q *no) seal() seal {
	panic("do not call")
}

func (q *no) uid() string {
	return "no"
}

func (q *no) isYes() bool {
	return false
}

func (q *no) isNo() bool {
	return true
}

// ====================================.

type yes struct {
}

func (q *yes) Test(...Attr) bool {
	return true
}

func (q *yes) And(rest ...Condition) Condition {
	return andOf(rest)
}

func (q *yes) Or(...Condition) Condition {
	return q
}

func (q *yes) Not() Condition {
	return &no{}
}

func (q *yes) seal() seal {
	panic("do not call")
}

func (q *yes) uid() string {
	return "yes"
}

func (q *yes) isYes() bool {
	return true
}

func (q *yes) isNo() bool {
	return false
}
