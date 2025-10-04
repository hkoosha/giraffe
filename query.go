package giraffe

import (
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/queryimpl"
)

// Query NEVER INSTANTIATE DIRECTLY. NEVER CAST TO. NEVER CAST FROM.
type Query string

func (q Query) impl() queryimpl.DialecticalQuery {
	return M(queryimpl.Parse(string(q)))
}

func (q Query) String() string {
	return q.impl().String()
}

func (q Query) Plus(other Query) (Query, error) {
	sum, err := q.impl().Plus(other.impl())
	if err != nil {
		return "", err
	}

	return Query(sum.String()), nil
}

func (q Query) Parser() func(string) (Query, error) {
	return Parser(q.impl().String())
}

func (q Query) ParserMust() func(string) Query {
	return ParserMust(q.impl().String())
}

func (q Query) WithMake() Query {
	panic("todo")
	// return Query(q.impl().WithMake().String())
}

func (q Query) WithOverwrite() Query {
	panic("todo")
	// return Query(q.impl().WithOverwrite().String())
}

func (q Query) WithoutOverwrite() Query {
	panic("todo")
	// return Query(q.impl().WithoutOverwrite().String())
}
