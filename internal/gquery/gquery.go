package gquery

import (
	"strings"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/gquery/gqcmd"
	"github.com/hkoosha/giraffe/internal/gquery/gqflag"
)

// DebugImpl to enable debug, it is needed to switch this alias to QueryDebug.
// It will require recompilation but allows Query.Debug to be typed and
// have zero lengths when debug is not enabled (contrary to using pointers).
type DebugImpl = struct{} // = QueryDebug.

type Query struct {
	Debug DebugImpl
	Path  *[]Query
	ref   string
	flags gqflag.QFlag
}

func (q Query) Flags() gqflag.QFlag {
	return q.flags
}

func (q Query) Named() string {
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

func (q Query) Root() Query {
	return (*q.Path)[0]
}

func (q Query) Leaf() Query {
	return (*q.Path)[len(*q.Path)-1]
}

func (q Query) Prev() Query {
	if seq := q.flags.Seq() - 1; seq >= 0 {
		prev := (*q.Path)[seq]

		return prev
	}

	panic("unreachable: no prev")
}

func (q Query) Next() Query {
	if seq := q.flags.Seq() + 1; seq < len(*q.Path) {
		return (*q.Path)[seq]
	}

	panic("unreachable: no next")
}

func (q Query) String() string {
	return q.string0()
}

//nolint:nonamedreturns
func (q Query) Segments() (
	s0 Query,
	s1 Query,
	ok bool,
) {
	// TODO go through mem cache.

	mover := q.Root()
	for !mover.flags.IsLeaf() && !mover.flags.IsMover() {
		mover = mover.Next()
	}
	if !mover.flags.IsMover() {
		return ErrQ, ErrQ, false
	}

	s0 = mover.UpTo(true)
	s1 = mover.Originating(false)

	return s0, s1, true
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

	return M(Parse(sb.String()))
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

	return M(Parse(sb.String()))
}

// =====================================.

func (q Query) WithWrite() Query {
	return q.reconstructedAs(q.flags | gqflag.QModWrite)
}

func (q Query) WithMake() Query {
	return q.reconstructedAs(q.flags | gqflag.QModeMake)
}

func (q Query) WithOverwrite() Query {
	return q.reconstructedAs(q.flags | gqflag.QModOverwrit)
}

func (q Query) WithoutOverwrite() Query {
	return q.reconstructedAs(q.flags & ^gqflag.QModOverwrit)
}

// Plus panics if the resulting query is too deep, set by MaxQueryDepth.
func (q Query) Plus(other Query) Query {
	return q.PlusS(other.String())
}

// PlusS panics if the resulting query is too deep, set by MaxQueryDepth.
func (q Query) PlusS(other string) Query {
	sb := strings.Builder{}

	q.bef(&sb)

	sb.WriteString(q.flags.ReconstructPreMod())
	sb.WriteString(q.ref)
	sb.WriteString(q.flags.ReconstructPostMod())
	sb.WriteByte(gqcmd.Sep)
	sb.WriteString(other)

	return M(parse(sb.String())).at(q.flags.Seq())
}

// =====================================.

func (q Query) MustReadonly() error {
	if !q.flags.IsReadonly() {
		return nil
	}

	return newQueryNotWritableError(q)
}
