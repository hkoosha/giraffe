package gquery

import (
	"slices"
	"strings"

	"github.com/hkoosha/giraffe/dialects"
	"github.com/hkoosha/giraffe/internal/queryerrors"
	"github.com/hkoosha/giraffe/internal/queryimpl"
	"github.com/hkoosha/giraffe/qcmd"
	"github.com/hkoosha/giraffe/qflag"
	. "github.com/hkoosha/giraffe/t11y/dot"
)

func newQuery(
	path *[]GiraffeQuery,
	ref string,
	flags qflag.QFlag,
) GiraffeQuery {
	return GiraffeQuery{
		// Debug: newDebug(),

		path:  path,
		ref:   ref,
		flags: flags,
	}
}

type GiraffeQuery struct {
	path  *[]GiraffeQuery
	ref   string
	flags qflag.QFlag
}

func (q GiraffeQuery) VisibleForTestingPath() []GiraffeQuery {
	return slices.Clone(*q.path)
}

func (q GiraffeQuery) Dialect() dialects.Dialect {
	return dialects.Giraffe1v1
}

func (q GiraffeQuery) Escaped() string {
	return Escaped(q.String())
}

func (q GiraffeQuery) Flags() qflag.QFlag {
	return q.flags
}

func (q GiraffeQuery) Attr() string {
	if !q.flags.IsObj() {
		return ""
	}

	return q.ref
}

func (q GiraffeQuery) Index() int {
	if !q.flags.IsArr() || q.flags.IsAppend() {
		return -1
	}

	return q.flags.Val()
}

func (q GiraffeQuery) Root() queryimpl.QueryImpl {
	return (*q.path)[0]
}

func (q GiraffeQuery) Leaf() queryimpl.QueryImpl {
	return (*q.path)[len(*q.path)-1]
}

func (q GiraffeQuery) Prev() queryimpl.QueryImpl {
	if seq := q.flags.Seq() - 1; seq >= 0 {
		prev := (*q.path)[seq]

		return prev
	}

	panic(EF("unreachable: no prev"))
}

func (q GiraffeQuery) Next() queryimpl.QueryImpl {
	if seq := q.flags.Seq() + 1; seq < len(*q.path) {
		return (*q.path)[seq]
	}

	panic(EF("unreachable: no next"))
}

func (q GiraffeQuery) String() string {
	return q.string0()
}

func (q GiraffeQuery) string0() string {
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

func (q GiraffeQuery) bef(
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

func (q GiraffeQuery) aft(
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

func (q GiraffeQuery) Reconstructed() string {
	sb := strings.Builder{}
	q.reconstructInAs(&sb, q.flags)
	return sb.String()
}

func (q GiraffeQuery) reconstructedIn(
	sb *strings.Builder,
) {
	q.reconstructInAs(sb, q.flags)
}

func (q GiraffeQuery) reconstructedAs(
	flags qflag.QFlag,
) GiraffeQuery {
	sb := strings.Builder{}

	q.bef(&sb)
	q.reconstructInAs(&sb, flags)
	q.aft(&sb)

	flagged := M(Parse(
		queryimpl.MaxDepth,
		sb.String(),
	))

	return flagged.at(q.flags.Seq())
}

func (q GiraffeQuery) reconstructInAs(
	sb *strings.Builder,
	flags qflag.QFlag,
) {
	flags.ReconstructPreModIn(sb)
	sb.WriteString(q.ref)
}

func (q GiraffeQuery) at(
	seq int,
) GiraffeQuery {
	return (*q.path)[seq]
}

// =====================================.

// UpTo TODO go through mem cache.
func (q GiraffeQuery) UpTo(withSelf bool) GiraffeQuery {
	if q.flags.IsSingle() {
		return q
	}

	sb := strings.Builder{}

	q.bef(&sb)

	if withSelf {
		q.reconstructedIn(&sb)
	}

	return M(Parse(
		queryimpl.MaxDepth,
		sb.String(),
	))
}

// Originating TODO go through mem cache.
func (q GiraffeQuery) Originating(withSelf bool) GiraffeQuery {
	if q.flags.IsSingle() {
		return q
	}

	sb := strings.Builder{}

	if withSelf {
		q.reconstructedIn(&sb)
	}
	q.aft(&sb)

	return M(Parse(
		queryimpl.MaxDepth,
		sb.String(),
	))
}

// =====================================.

func (q GiraffeQuery) WithWrite() queryimpl.QueryImpl {
	return q.reconstructedAs(q.flags | qflag.QModWrite)
}

func (q GiraffeQuery) WithMake() queryimpl.QueryImpl {
	return q.reconstructedAs(q.flags | qflag.QModeMake)
}

func (q GiraffeQuery) WithOverwrite() queryimpl.QueryImpl {
	return q.reconstructedAs(q.flags | qflag.QModOverwrit)
}

func (q GiraffeQuery) WithoutOverwrite() queryimpl.QueryImpl {
	return q.reconstructedAs(q.flags & ^qflag.QModOverwrit)
}

// PlusS panics if the resulting query is too deep, set by iface.MaxDepth.
func (q GiraffeQuery) PlusS(other string) queryimpl.QueryImpl {
	sb := strings.Builder{}

	q.bef(&sb)

	sb.WriteString(q.flags.ReconstructPreMod())
	sb.WriteString(q.ref)
	sb.WriteByte(qcmd.Sep.Byte())
	sb.WriteString(other)

	return M(Parse(
		queryimpl.MaxDepth,
		sb.String(),
	)).at(q.flags.Seq())
}

func (q GiraffeQuery) MustReadonly() error {
	if !q.flags.IsReadonly() {
		return nil
	}

	return queryerrors.NewNotWritableError(q.String())
}
