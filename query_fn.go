package giraffe

import (
	"reflect"

	"github.com/hkoosha/giraffe/internal"
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/qcmd"
)

func Q[T interface{ Query | string }](
	spec T,
) Query {
	switch q := any(spec).(type) {
	case Query:
		return M(GQParse(q.impl().String()))

	case string:
		return M(GQParse(q))

	default:
		panic("unknown query type: " + reflect.TypeOf(spec).String())
	}
}

func GQErr() Query {
	return ""
}

func GQParse(
	spec string,
) (Query, error) {
	if _, err := internal.Parse(spec); err != nil {
		return "", err
	}

	return Query(spec), nil
}

func GQParser(
	prefix string,
) func(string) (Query, error) {
	if prefix != "" {
		prefix += qcmd.Sep.String()
		M(internal.Parse(prefix + "dummy"))
	}

	return func(spec string) (Query, error) {
		return GQParse(prefix + spec)
	}
}

func GQParserMust(
	prefix string,
) func(string) Query {
	parser := GQParser(prefix)

	return func(spec string) Query {
		return M(parser(spec))
	}
}
