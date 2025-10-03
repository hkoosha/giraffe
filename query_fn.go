package giraffe

import (
	"reflect"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/gquery"
	"github.com/hkoosha/giraffe/internal/gquery/gcmd"
)

//goland:noinspection GoUnusedConst
const (
	CmdOverwrite = string(gqcmd.CmdOverwrite)
	CmdMake      = string(gqcmd.CmdMake)
	CmdMaybe     = string(gqcmd.CmdMaybe)
	CmdAppend    = string(gqcmd.CmdAppend)
	CmdDelete    = string(gqcmd.CmdDelete)
	CmdSep       = string(gqcmd.CmdSep)
	CmdEscape    = string(gqcmd.CmdEscape)
	CmdAt        = string(gqcmd.CmdAt)
	CmdSelf      = string(gqcmd.CmdSelf)
)

func Escaped(
	spec string,
) string {
	return queryimpl.Escaped(spec)
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
	if _, err := queryimpl.Parse(spec); err != nil {
		return invalid, err
	}

	return Query(spec), nil
}

func Parser(
	prefix string,
) func(string) (Query, error) {
	if prefix != "" {
		prefix += CmdSep
		M(queryimpl.Parse(prefix + "dummy"))
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
