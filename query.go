package giraffe

import (
	"reflect"

	. "github.com/hkoosha/giraffe/internal/dot"
	"github.com/hkoosha/giraffe/internal/gquery"
)

//goland:noinspection GoUnusedConst
const (
	CmdOverwrite = string(gquery.CmdOverwrite)
	CmdMake      = string(gquery.CmdMake)
	CmdMaybe     = string(gquery.CmdMaybe)
	CmdAppend    = string(gquery.CmdAppend)
	CmdDelete    = string(gquery.CmdDelete)
	CmdSep       = string(gquery.CmdSep)
	CmdEscape    = string(gquery.CmdEscape)
	CmdAt        = string(gquery.CmdAt)
	CmdSelf      = string(gquery.CmdSelf)
)

func Escaped(
	spec string,
) string {
	return gquery.Escaped(spec)
}

func EscapedQ(
	spec string,
) Query {
	return Q(Escaped(spec))
}

func Q(
	spec string,
) Query {
	return M(Parse(spec))
}

func QErr() Query {
	return invalid
}

func Parse(
	spec string,
) (Query, error) {
	if _, err := gquery.Parse(spec); err != nil {
		return invalid, err
	}

	return Query(spec), nil
}

func Parser(
	prefix string,
) func(string) (Query, error) {
	if prefix != "" {
		prefix += CmdSep
		M(gquery.Parse(prefix + "dummy"))
	}

	return func(spec string) (Query, error) {
		return Parse(prefix + spec)
	}
}

func ParserMust(
	prefix string,
) func(string) Query {
	parser := Parser(prefix)

	return func(spec string) Query {
		return M(parser(spec))
	}
}

// =============================================================================.

// Query NEVER INSTANTIATE DIRECTLY. NEVER CAST TO. NEVER CAST FROM.
type Query string

func (q Query) String() string {
	return q.impl().String()
}

func (q Query) Plus(other Query) Query {
	return Query(q.impl().Plus(other.impl()).String())
}

func (q Query) Parser() func(string) (Query, error) {
	return Parser(q.impl().String())
}

func (q Query) ParserMust() func(string) Query {
	return ParserMust(q.impl().String())
}

func (q Query) WithMake() Query {
	return Query(q.impl().WithMake().String())
}

func (q Query) WithOverwrite() Query {
	return Query(q.impl().WithOverwrite().String())
}

func (q Query) WithoutOverwrite() Query {
	return Query(q.impl().WithoutOverwrite().String())
}

func (q Query) impl() gquery.Query {
	return M(gquery.Parse(string(q)))
}

var (
	invalid = Query("")
	tQuery  = reflect.TypeOf((*Query)(nil)).Elem()
)
