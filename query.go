package giraffe

import (
	"github.com/hkoosha/giraffe/internal"
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/queryimpl/dialectical"
)

// GQuery NEVER INSTANTIATE DIRECTLY. NEVER CAST TO. NEVER CAST FROM.
type GQuery string

func (q GQuery) impl() dialectical.DialecticalQuery {
	return M(internal.Parse(string(q)))
}

func (q GQuery) String() string {
	return q.impl().String()
}

func (q GQuery) Plus(other GQuery) (GQuery, error) {
	sum, err := q.impl().Plus(other.impl())
	if err != nil {
		return "", err
	}

	return GQuery(sum.String()), nil
}

func (q GQuery) Parser() func(string) (GQuery, error) {
	return GQParser(q.impl().String())
}

func (q GQuery) ParserMust() func(string) GQuery {
	return GQParserMust(q.impl().String())
}

func (q GQuery) WithMake() GQuery {
	return GQuery(q.impl().WithMake().String())
}

func (q GQuery) WithOverwrite() GQuery {
	return GQuery(q.impl().WithOverwrite().String())
}

func (q GQuery) WithoutOverwrite() GQuery {
	return GQuery(q.impl().WithoutOverwrite().String())
}
