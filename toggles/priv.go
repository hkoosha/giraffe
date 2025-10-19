package toggles

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"slices"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/serdes/gson"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	"github.com/hkoosha/giraffe/toggles/internal"
)

func andOf(rest []Condition) Condition {
	uids := make([]string, 0, len(rest))
	filtered := make([]Condition, 0, len(rest))
	for _, it := range rest {
		if it.isNo() {
			return no_
		} else if !it.isYes() {
			filtered = append(filtered, it)
			uids = append(uids, it.uid())
		}
	}

	if len(rest) == 0 {
		return yes_
	} else if len(rest) == 1 {
		return rest[0]
	}

	return &and{
		rest: rest,
		cond: cond{
			Sealer: internal.Sealer{},
			name:   "",
			uid_:   fmt.Sprintf("and(%d, %s)", len(uids), uidOf(uids)),
		},
	}
}

func orOf(rest []Condition) Condition {
	uids := make([]string, 0, len(rest))
	filtered := make([]Condition, 0, len(rest))
	for _, it := range rest {
		if it.isYes() {
			return yes_
		} else if !it.isNo() {
			filtered = append(filtered, it)
			uids = append(uids, it.uid())
		}
	}

	if len(rest) == 0 {
		return yes_
	} else if len(rest) == 1 {
		return rest[0]
	}

	return &or{
		rest: rest,
		cond: cond{
			Sealer: internal.Sealer{},
			name:   "",
			uid_:   fmt.Sprintf("or(%d, %s)", len(uids), uidOf(uids)),
		},
	}
}

func condOf(
	name string,
	op string,
	args string,
) cond {
	return cond{
		Sealer: internal.Sealer{},
		name:   name,
		uid_:   fmt.Sprintf("%s(%s)", op, args),
	}
}

func uidOf(v any) string {
	sum := sha256.Sum256(gson.MustMarshal(v))
	uid := hex.EncodeToString(sum[:])
	return uid
}

func valueOf(
	name string,
	v any,
) Value {
	return &value{
		Sealer: internal.Sealer{},
		name:   name,
		val:    v,
	}
}

func inOf[V giraffe.Ord](
	name string,
	v []V,
) Condition {
	switch {
	case len(v) == 0:
		return no_

	case len(v) == 1:
		return eqOf(name, v[0])

	case len(v) < 8:
		//goland:noinspection GoPrintFunctions
		return &in[V]{
			val:  v,
			cond: condOf(name, "in", uidOf(v)),
		}

	default:
		//goland:noinspection GoPrintFunctions
		return &search[V]{
			val:  v,
			cond: condOf(name, "search", uidOf(v)),
		}
	}
}

func eqOf[V giraffe.Basic](
	name string,
	v V,
) Condition {
	return &eq{
		val:  v,
		cond: condOf(name, "eq", fmt.Sprintf("%s=%v", name, v)),
	}
}

func gtOf[V giraffe.Num](
	name string,
	v V,
) Condition {
	return &gt[V]{
		val:  v,
		cond: condOf(name, "gt", fmt.Sprintf("%s>%v", name, v)),
	}
}

func ltOf[V giraffe.Num](
	name string,
	v V,
) Condition {
	return &lt[V]{
		val:  v,
		cond: condOf(name, "lt", fmt.Sprintf("%s<%v", name, v)),
	}
}

// ============================================================================.

type condition interface {
	internal.Sealed

	test([]Value) bool
	test0(Value) bool
	uid() string
	isYes() bool
	isNo() bool
}

// ====================================.

type cond struct {
	internal.Sealer

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
	//goland:noinspection GoPrintFunctions
	return &not{
		rest: q,
		cond: condOf("", "not", q.uid_),
	}
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
	panic(EF("unreachable: cond.test0()"))
}

// ====================================.

type and struct {
	cond
	rest []Condition
}

func (q *and) test(values []Value) bool {
	if len(q.rest) == 0 {
		panic(EF("unreachable: empty conjunction"))
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
		panic(EF("unreachable: empty disjunction"))
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

var no_ Condition = &no{condOf("", "no", "")}

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
	return yes_
}

func (q *no) isNo() bool {
	return true
}

// ====================================.

var yes_ Condition = &yes{condOf("", "yes", "")}

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
	return no_
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
	if vv, ok := v.Value().(V); ok {
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
	if vv, ok := v.Value().(V); ok {
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
	if vv, ok := v.Value().(V); !ok || !slices.Contains(q.val, vv) {
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
	internal.Sealer

	val  any
	name string
}

func (q *value) Str() (string, bool) {
	if v, ok := q.val.(string); ok {
		return v, true
	}

	return "", false
}

func (q *value) Bln() (v, ok bool) { //nolint:nonamedreturns
	if v, ok = q.val.(bool); ok {
		return v, true
	}

	return false, false
}

func (q *value) I64() (int64, bool) {
	if v, ok := q.val.(int64); ok {
		return v, true
	}

	return 0, false
}

func (q *value) U64() (uint64, bool) {
	if v, ok := q.val.(uint64); ok {
		return v, true
	}

	return 0, false
}

func (q *value) Name() string {
	return q.name
}

func (q *value) Value() any {
	return q.val
}

// ====================================.

func newRouter(
	defaultCase Condition,
	togglers []Storage,
) Toggler {
	return &router{
		Sealer: internal.Sealer{},

		defaultCase: defaultCase,
		togglers:    togglers,
	}
}

type router struct {
	internal.Sealer
	defaultCase Condition
	togglers    []Storage
}

func (r *router) Query(
	ctx gtx.Context,
	name string,
	values ...Value,
) (bool, error) {
	var err error

	for _, t := range r.togglers {
		var en *bool
		if en, err = t.Get(ctx, name, slices.Clone(values)); err == nil && en != nil && *en {
			return true, nil
		}
	}

	return r.defaultCase.test(values), err
}
