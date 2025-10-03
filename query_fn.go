package giraffe

import (
	"reflect"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/queryimpl"
	"github.com/hkoosha/giraffe/qcmd"
)

func Escaped(
	spec string,
) string {
	return queryimpl.Escaped(spec)
}

func Q[T interface{ Query | string }](
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

func QErr() Query {
	return ""
}

func Parse(
	spec string,
) (Query, error) {
	if _, err := queryimpl.Parse(spec); err != nil {
		return "", err
	}

	return Query(spec), nil
}

func Parser(
	prefix string,
) func(string) (Query, error) {
	if prefix != "" {
		prefix += qcmd.Sep.String()
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
