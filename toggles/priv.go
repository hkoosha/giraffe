package toggles

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/zebra/z"
)

const (
	opFirst binOp = iota + 1

	opEq
	opGt
	opLt

	opAnd
	opOr

	opIn

	opLast
)

type binOp uint

func (o binOp) String() string {
	switch o {
	case opEq:
		return "="

	case opGt:
		return ">"

	case opLt:
		return "<"

	case opAnd:
		return "&"

	case opOr:
		return "|"

	case opIn:
		return "~"

	default:
		return fmt.Sprintf("op[invalid=%d]", o)
	}
}

func join(
	rest []Condition,
	op binOp,
) Condition {
	if op == opAnd {
		rest = z.Filtered(rest, func(it Condition) bool {
			return !it.isYes()
		})
		if slices.ContainsFunc(rest, func(it Condition) bool { return it.isNo() }) {
			return no_
		}
	} else if op == opOr {
		rest = z.Filtered(rest, func(it Condition) bool {
			return !it.isNo()
		})
		if slices.ContainsFunc(rest, func(it Condition) bool { return it.isYes() }) {
			return yes_
		}
	}

	if len(rest) == 0 {
		return &yes{}
	} else if len(rest) == 1 {
		return rest[0]
	}

	var uid string
	{
		uids := z.Applied(rest, func(it Condition) string { return it.uid() })
		slices.Sort(uids)
		uid = "(" + strings.Join(uids, op.String()) + ")"
	}

	switch op {
	case opAnd:
		return &and{
			rest: rest,
			cond: cond{
				uid_: uid,
			},
		}

	case opOr:
		return &or{
			rest: rest,
			cond: cond{
				uid_: uid,
			},
		}

	default:
		panic(EF("unknown operation '%s'", op))
	}
}

func andOf(rest []Condition) Condition {
	return join(rest, opAnd)
}

func orOf(rest []Condition) Condition {
	return join(rest, opOr)
}

func notOf(c Condition) Condition {
	return &not{rest: c}
}

func (o binOp) ensure() binOp {
	if o <= opFirst || opLast <= o {
		panic(EF("invalid op: %s", o))
	}

	return o
}

func uidOf(v any) string {
	sum := md5.Sum(M(json.Marshal(v)))
	uid := hex.EncodeToString(sum[:])
	return uid
}

type seal struct{}

type sealed interface {
	seal() seal
}

type sealer struct{}

func (q *sealer) seal() seal {
	panic("do not call")
}

// ============================================================================.

type condition interface {
	sealed

	test([]Value) bool
	test0(Value) bool
	uid() string
	isYes() bool
	isNo() bool
}

// ====================================.

type cond struct {
	sealer

	name string
	uid_ string
}

func (q *cond) Name() string {
	return q.name
}

func (q *cond) And(rest ...Condition) Condition {
	return andOf(append(rest, q))
}

func (q *cond) Or(rest ...Condition) Condition {
	return orOf(append(rest, q))
}

func (q *cond) Not() Condition {
	return notOf(q)
}

func (q *cond) uid() string {
	return q.uid_
}

func (q *cond) isYes() bool {
	return false
}

func (q *cond) isNo() bool {
	return false
}

func (q *cond) test(rest []Value) bool {
	for _, r := range rest {
		if r.Name() == q.name && !q.test0(r) {
			return false
		}
	}

	return true
}

func (q *cond) test0(Value) bool {
	panic("unreachable: cond.test0()")
}

// ====================================.

type and struct {
	cond
	rest []Condition
}

func (q *and) test(values []Value) bool {
	if len(q.rest) == 0 {
		panic("unreachable: empty conjunction")
	}

	for _, c := range q.rest {
		if !c.test(values) {
			return false
		}
	}
	return true
}

// ====================================.

type or struct {
	cond
	rest []Condition
}

func (q *or) test(values []Value) bool {
	if len(q.rest) == 0 {
		panic("unreachable: empty disjunction")
	}

	for _, c := range q.rest {
		if c.test(values) {
			return true
		}
	}
	return false
}

// ====================================.

type not struct {
	rest Condition
	cond
}

func (q *not) test(values []Value) bool {
	return !q.rest.test(values)
}

func (q *not) Not() Condition {
	return q.rest
}

func (q *not) And(rest ...Condition) Condition {
	return andOf(append(rest, q))
}

func (q *not) Or(rest ...Condition) Condition {
	return orOf(append(rest, q))
}

// ====================================.

var no_ Condition = &no{cond{uid_: "no"}}

type no struct {
	cond
}

func (q *no) test([]Value) bool {
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

func (q *no) isNo() bool {
	return true
}

// ====================================.

var yes_ Condition = &yes{cond{uid_: "yes"}}

type yes struct {
	cond
}

func (q *yes) test([]Value) bool {
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

func (q *yes) isYes() bool {
	return true
}

// ====================================.

type eq struct {
	val any
	cond
}

func (q *eq) test0(v Value) bool {
	return reflect.DeepEqual(q.val, v)
}

// ====================================.

type gt[V giraffe.Ord] struct {
	val V
	cond
}

func (q *gt[V]) test0(v Value) bool {
	if vv, ok := v.value().(V); ok {
		return q.val > vv
	}

	return false
}

// ====================================.

type lt[V giraffe.Ord] struct {
	val V
	cond
}

func (q *lt[V]) test0(v Value) bool {
	if vv, ok := v.value().(V); ok {
		return q.val < vv
	}

	return false
}

// ====================================.

type in[V comparable] struct {
	cond
	val []V
}

func (q *in[V]) test0(v Value) bool {
	if vv, ok := v.value().(V); !ok || !slices.Contains(q.val, vv) {
		return false
	}

	return true
}

// ====================================.

type search[V giraffe.Ord] struct {
	cond
	val []V
}

func (q *search[V]) test0(v Value) bool {
	if vv, ok := v.(V); ok {
		_, ok = slices.BinarySearch(q.val, vv)
		return ok
	}

	return false
}

// ====================================.

type value struct {
	sealer
	val  any
	name string
}

func (q *value) Name() string {
	return q.name
}

func (q *value) value() any {
	return q.value
}
