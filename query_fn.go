package giraffe

import (
	"reflect"

	. "github.com/hkoosha/giraffe/internal/dot0"
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

func Q[T QueryT](
	spec T,
) Query {
	//nolint:gocritic
	if asQ, ok := any(spec).(Query); ok {
		return M(Parse(asQ.impl().Reconstructed()))
	} else if asStr, ok := any(spec).(string); ok {
		return M(Parse(asStr))
	} else {
		panic("unknown type for Q: " + reflect.TypeOf(spec).String())
	}
}

type QueryT interface {
	Query | string
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
