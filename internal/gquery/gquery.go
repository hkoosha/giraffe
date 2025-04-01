package gquery

import (
	"strings"

	. "github.com/hkoosha/giraffe/internal/dot"
)

const (
	cmdSepStr = string(CmdSep)

	CmdOverwrite = '%'
	CmdMake      = '$'
	CmdMaybe     = '?'
	CmdAppend    = '+'
	CmdDelete    = '^'
	CmdSep       = '.'
	CmdEscape    = '\\'
	CmdAt        = '@'
	CmdSelf      = '#'
	CmdMove      = '>'

	CmdNonDeterministic = '~'
)

// DebugImpl to enable debug, it is needed to switch this alias to QueryDebug.
// It will require recompilation but allows Query.Debug to be typed and
// have zero lengths when debug is not enabled (contrary to using pointers).
type DebugImpl = struct{} // = QueryDebug.

type Query struct {
	Debug DebugImpl
	Path  *[]Query
	ref   string
	flags QFlag
}

func (q Query) Flags() QFlag {
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

// Segment0 TODO go through mem cache.
func (q Query) Segment0() (Query, bool) {
	mover := q.Root()
	for !mover.flags.IsLeaf() && !mover.flags.IsMover() {
		mover = mover.Next()
	}

	if !mover.flags.IsMover() {
		return ErrQ, false
	}

	return mover.UpTo(true), true
}

// Segment1 TODO go through mem cache.
func (q Query) Segment1() (Query, bool) {
	mover := q.Root()
	for !mover.flags.IsLeaf() && !mover.flags.IsMover() {
		mover = mover.Next()
	}

	if !mover.flags.IsMover() {
		return ErrQ, false
	}

	if mover.flags.IsMover() {
		return mover.Originating(false), true
	}

	return ErrQ, false
}

//nolint:nonamedreturns
func (q Query) Segments() (
	s0 Query,
	s1 Query,
	ok bool,
) {
	s0, ok = q.Segment0()
	if !ok {
		return ErrQ, ErrQ, false
	}

	s1, ok = q.Segment1()
	if !ok {
		return ErrQ, ErrQ, false
	}

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
	return q.reconstructedAs(q.flags | QModWrite)
}

func (q Query) WithMake() Query {
	return q.reconstructedAs(q.flags | QModeMake)
}

func (q Query) WithOverwrite() Query {
	return q.reconstructedAs(q.flags | QModOverwrit)
}

func (q Query) WithoutOverwrite() Query {
	return q.reconstructedAs(q.flags & ^QModOverwrit)
}

// Plus panics if the resulting query is too deep, set by MaxQueryDepth.
func (q Query) Plus(other Query) Query {
	return q.PlusS(other.String())
}

// PlusS panics if the resulting query is too deep, set by MaxQueryDepth.
func (q Query) PlusS(other string) Query {
	sb := strings.Builder{}

	q.bef(&sb)

	sb.WriteString(q.flags.reconstructPreMod())
	sb.WriteString(q.ref)
	sb.WriteString(q.flags.reconstructPostMod())
	sb.WriteByte(CmdSep)
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
