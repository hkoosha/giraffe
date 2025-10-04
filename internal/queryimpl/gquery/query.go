package gquery

import (
	"strings"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/queryerrors"
	"github.com/hkoosha/giraffe/internal/queryimpl"
	"github.com/hkoosha/giraffe/qcmd"
	"github.com/hkoosha/giraffe/qflag"
)

func newQuery(
	path *[]Query,
	ref string,
	flags qflag.QFlag,
) Query {
	return Query{
		// Debug: newDebug(),

		path:  path,
		ref:   ref,
		flags: flags,
	}
}

type Query struct {
	path  *[]Query
	ref   string
	flags qflag.QFlag
}

func (q Query) Flags() qflag.QFlag {
	return q.flags
}

func (q Query) Attr() string {
	if !q.flags.IsObj() {
		return ""
	}

	return q.ref
}

func (q Query) Index() int {
	if !q.flags.IsArr() || q.flags.IsAppend() {
		return -1
	}

	return q.flags.Val()
}

func (q Query) Root() queryimpl.QueryImpl {
	return (*q.path)[0]
}

func (q Query) Leaf() queryimpl.QueryImpl {
	return (*q.path)[len(*q.path)-1]
}

func (q Query) Prev() queryimpl.QueryImpl {
	if seq := q.flags.Seq() - 1; seq >= 0 {
		prev := (*q.path)[seq]

		return prev
	}

	panic("unreachable: no prev")
}

func (q Query) Next() queryimpl.QueryImpl {
	if seq := q.flags.Seq() + 1; seq < len(*q.path) {
		return (*q.path)[seq]
	}

	panic("unreachable: no next")
}

// Plus panics if the resulting query is too deep, set by iface.MaxDepth.
func (q Query) Plus(other string) (queryimpl.QueryImpl, error) {
	return q.PlusS(other), nil
}

func (q Query) String() string {
	return q.string0()
}

func (q Query) string0() string {
	sb := strings.Builder{}

	for j, p := range *q.path {
		if j > 0 {
			sb.WriteByte(qcmd.Sep.Byte())

			if q.flags.Seq() == j {
				sb.WriteByte(qcmd.At.Byte())
			}
		}

		sb.WriteString(p.flags.ReconstructPreMod())
		sb.WriteString(p.ref)
	}

	return sb.String()
}

func (q Query) bef(
	sb *strings.Builder,
) {
	if q.Flags().IsRoot() {
		return
	}

	path := *q.path
	for i := range q.flags.Seq() {
		qI := path[i]
		sb.WriteString(qI.flags.ReconstructPreMod())
		sb.WriteString(qI.ref)
		sb.WriteByte(qcmd.Sep.Byte())
	}
}

func (q Query) aft(
	sb *strings.Builder,
) {
	if q.Flags().IsLeaf() {
		return
	}

	path := *q.path
	for i := q.flags.Seq() + 1; i < len(path); i++ {
		sb.WriteByte(qcmd.Sep.Byte())

		qI := path[i]
		sb.WriteString(qI.flags.ReconstructPreMod())
		sb.WriteString(qI.ref)
	}
}

func (q Query) Reconstructed() string {
	sb := strings.Builder{}
	q.reconstructInAs(&sb, q.flags)
	return sb.String()
}

func (q Query) reconstructedIn(
	sb *strings.Builder,
) {
	q.reconstructInAs(sb, q.flags)
}

func (q Query) reconstructedAs(
	flags qflag.QFlag,
) Query {
	sb := strings.Builder{}

	q.bef(&sb)
	q.reconstructInAs(&sb, flags)
	q.aft(&sb)

	flagged := M(parse(sb.String()))

	return flagged.at(q.flags.Seq())
}

func (q Query) reconstructInAs(
	sb *strings.Builder,
	flags qflag.QFlag,
) {
	flags.ReconstructPreModIn(sb)
	sb.WriteString(q.ref)
}

func (q Query) at(
	seq int,
) Query {
	return (*q.path)[seq]
}

// =====================================.

// UpTo TODO go through mem cache.
func (q Query) UpTo(withSelf bool) Query {
	if q.flags.IsSingle() {
		return q
	}

	sb := strings.Builder{}

	q.bef(&sb)

	if withSelf {
		q.reconstructedIn(&sb)
	}

	return M(parse(sb.String()))
}

// Originating TODO go through mem cache.
func (q Query) Originating(withSelf bool) Query {
	if q.flags.IsSingle() {
		return q
	}

	sb := strings.Builder{}

	if withSelf {
		q.reconstructedIn(&sb)
	}
	q.aft(&sb)

	return M(parse(sb.String()))
}

// =====================================.

func (q Query) WithWrite() Query {
	return q.reconstructedAs(q.flags | qflag.QModWrite)
}

func (q Query) WithMake() Query {
	return q.reconstructedAs(q.flags | qflag.QModeMake)
}

func (q Query) WithOverwrite() Query {
	return q.reconstructedAs(q.flags | qflag.QModOverwrit)
}

func (q Query) WithoutOverwrite() Query {
	return q.reconstructedAs(q.flags & ^qflag.QModOverwrit)
}

// PlusS panics if the resulting query is too deep, set by iface.MaxDepth.
func (q Query) PlusS(other string) Query {
	sb := strings.Builder{}

	q.bef(&sb)

	sb.WriteString(q.flags.ReconstructPreMod())
	sb.WriteString(q.ref)
	sb.WriteByte(qcmd.Sep.Byte())
	sb.WriteString(other)

	return M(parse(sb.String())).at(q.flags.Seq())
}

func (q Query) MustReadonly() error {
	if !q.flags.IsReadonly() {
		return nil
	}

	return queryerrors.NewNotWritableError(q.String())
}
