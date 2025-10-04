package giraffe

import (
	"reflect"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/queryimpl"
	"github.com/hkoosha/giraffe/qcmd"
)

func Q[T interface{ GQuery | string }](
	spec T,
) GQuery {
	switch q := any(spec).(type) {
	case GQuery:
		return M(GQParse(q.impl().String()))

	default:
		panic("unknown query type: " + reflect.TypeOf(spec).String())
	}
}

func GQErr() GQuery {
	return ""
}

func GQParse(
	spec string,
) (GQuery, error) {
	if _, err := queryimpl.Parse(spec); err != nil {
		return "", err
	}

	return GQuery(spec), nil
}

func GQParser(
	prefix string,
) func(string) (GQuery, error) {
	if prefix != "" {
		prefix += qcmd.Sep.String()
		M(queryimpl.Parse(prefix + "dummy"))
	}

	return func(spec string) (GQuery, error) {
		return GQParse(prefix + spec)
	}
}

func GQParserMust(
	prefix string,
) func(string) GQuery {
	parser := GQParser(prefix)

	return func(spec string) GQuery {
		return M(parser(spec))
	}
}
