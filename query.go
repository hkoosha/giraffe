package giraffe

import (
	"github.com/hkoosha/giraffe/dialects"
	"github.com/hkoosha/giraffe/qflag"
)

// Query NEVER INSTANTIATE DIRECTLY. NEVER CAST TO. NEVER CAST FROM.
type Query string

func (q Query) Flags() qflag.QFlag {
	// TODO implement me
	panic("implement me")
}

func (q Query) Attr() string {
	return q.impl().Attr()
}

func (q Query) Index() int {
	return q.impl().Index()
}

func (q Query) Root() Query {
	return Query(q.impl().Root().String())
}

func (q Query) Leaf() Query {
	return Query(q.impl().Leaf().String())
}

func (q Query) Prev() Query {
	return Query(q.impl().Prev().String())
}

func (q Query) Next() Query {
	return Query(q.impl().Next().String())
}

func (q Query) String() string {
	return q.impl().String()
}

func (q Query) Dialect() dialects.Dialect {
	return q.impl().Dialect()
}

func (q Query) Parser() func(string) (Query, error) {
	return GQParser(q.impl().String())
}

func (q Query) ParserMust() func(string) Query {
	return GQParserMust(q.impl().String())
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
