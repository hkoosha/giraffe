package dialectical

import (
	"github.com/hkoosha/giraffe/dialects"
	"github.com/hkoosha/giraffe/internal/queryimpl"
	"github.com/hkoosha/giraffe/qcmd"
	"github.com/hkoosha/giraffe/qflag"
)

func New() DialecticalQuery {
	return DialecticalQuery{
		dialect: dialects.Unknown,
		impl:    nil,
	}
}

//goland:noinspection GoNameStartsWithPackageName
type DialecticalQuery struct {
	dialect dialects.Dialect
	impl    queryimpl.QueryImpl
}

func (q DialecticalQuery) Flags() qflag.QFlag {
	return q.impl.Flags()
}

func (q DialecticalQuery) Attr() string {
	return q.impl.Attr()
}

func (q DialecticalQuery) Index() int {
	return q.impl.Index()
}

func (q DialecticalQuery) Plus(query DialecticalQuery) (DialecticalQuery, error) {
	sum, err := q.impl.Plus(query.String())
	if err != nil {
		return DialecticalQuery{}, err
	}

	return DialecticalQuery{
		dialect: q.dialect,
		impl:    sum,
	}, nil
}

func (q DialecticalQuery) Root() DialecticalQuery {
	return DialecticalQuery{
		dialect: q.dialect,
		impl:    q.impl.Root(),
	}
}

func (q DialecticalQuery) Leaf() DialecticalQuery {
	return DialecticalQuery{
		dialect: q.dialect,
		impl:    q.impl.Leaf(),
	}
}

func (q DialecticalQuery) Prev() DialecticalQuery {
	return DialecticalQuery{
		dialect: q.dialect,
		impl:    q.impl.Prev(),
	}
}

func (q DialecticalQuery) Next() DialecticalQuery {
	return DialecticalQuery{
		dialect: q.dialect,
		impl:    q.impl.Next(),
	}
}

func (q DialecticalQuery) WithMake() DialecticalQuery {
	return DialecticalQuery{
		dialect: q.dialect,
		impl:    q.impl.WithMake(),
	}
}

func (q DialecticalQuery) WithOverwrite() DialecticalQuery {
	return DialecticalQuery{
		dialect: q.dialect,
		impl:    q.impl.WithOverwrite(),
	}
}

func (q DialecticalQuery) WithoutOverwrite() DialecticalQuery {
	return DialecticalQuery{
		dialect: q.dialect,
		impl:    q.impl.WithoutOverwrite(),
	}
}

func (q DialecticalQuery) WithImpl(impl queryimpl.QueryImpl) DialecticalQuery {
	return DialecticalQuery{
		dialect: q.dialect,
		impl:    impl,
	}
}

func (q DialecticalQuery) WithDialect(d dialects.Dialect) DialecticalQuery {
	return DialecticalQuery{
		dialect: d,
		impl:    q.impl,
	}
}

func (q DialecticalQuery) Dialect() dialects.Dialect {
	return q.dialect
}

func (q DialecticalQuery) Impl() queryimpl.QueryImpl {
	return q.impl
}

func (q DialecticalQuery) String() string {
	return qcmd.Dialect.String() + q.dialect.String() + q.impl.String()
}
