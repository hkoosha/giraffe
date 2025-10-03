package giraffe

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
