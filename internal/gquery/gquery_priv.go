package gquery

import (
	"regexp"
	"strings"

	. "github.com/hkoosha/giraffe/internal/dot"
)

// MaxQueryDepth must fit in the QFlag in the sequence part, i.e., 8 bits.
const MaxQueryDepth = 255

var uintRegex = regexp.MustCompile(`^\d+$`)

func newQuery(
	path *[]Query,
	ref string,
	flags QFlag,
) Query {
	return Query{
		Path:  path,
		ref:   ref,
		flags: flags,
		Debug: newQueryDebug(),
	}
}

func (q Query) at(
	seq int,
) Query {
	return (*q.Path)[seq]
}

func (q Query) reconstructedIn(
	sb *strings.Builder,
) {
	q.reconstructInAs(sb, q.flags)
}

func (q Query) reconstructedAs(
	flags QFlag,
) Query {
	sb := strings.Builder{}

	q.bef(&sb)
	q.reconstructInAs(&sb, flags)
	q.aft(&sb)

	flagged := M(Parse(sb.String()))

	return flagged.at(q.flags.Seq())
}

func (q Query) reconstructInAs(
	sb *strings.Builder,
	flags QFlag,
) {
	flags.reconstructPreModIn(sb)
	sb.WriteString(q.ref)
	flags.reconstructPostModIn(sb)
}

func (q Query) string0() string {
	sb := strings.Builder{}

	for j, p := range *q.Path {
		if j > 0 {
			sb.WriteByte(CmdSep)

			if q.flags.Seq() == j {
				sb.WriteByte(CmdAt)
			}
		}

		sb.WriteString(p.flags.reconstructPreMod())
		sb.WriteString(p.ref)
		sb.WriteString(q.flags.reconstructPostMod())
	}

	return sb.String()
}

func (q Query) bef(
	sb *strings.Builder,
) {
	if q.Flags().IsRoot() {
		return
	}

	path := *q.Path
	for i := range q.flags.Seq() {
		qI := path[i]
		sb.WriteString(qI.flags.reconstructPreMod())
		sb.WriteString(qI.ref)
		sb.WriteString(qI.flags.reconstructPostMod())
		sb.WriteByte(CmdSep)
	}
}

func (q Query) aft(
	sb *strings.Builder,
) {
	if q.Flags().IsLeaf() {
		return
	}

	path := *q.Path
	for i := q.flags.Seq() + 1; i < len(path); i++ {
		sb.WriteByte(CmdSep)

		qI := path[i]
		sb.WriteString(qI.flags.reconstructPreMod())
		sb.WriteString(qI.ref)
		sb.WriteString(qI.flags.reconstructPostMod())
	}
}
