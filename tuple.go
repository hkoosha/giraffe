package giraffe

import (
	"fmt"
)

func TupleOf[V SafeType1](
	q Query,
	v V,
) Tuple {
	return Tuple{
		Query: q,
		Dat:   Of(v),
	}
}

type Tuple struct {
	Query Query
	Dat   Datum
}

func (t *Tuple) String() string {
	return fmt.Sprintf("Tuple2[%s, %s]", t.Query, t.Dat.String())
}

func (t *Tuple) Implode() Datum {
	return Of1(t.Query, t.Dat)
}

func (t *Tuple) Unpack() (Query, Datum) {
	return t.Query, t.Dat
}

func GetQuery(t Tuple) Query {
	return t.Query
}

func GetDatum(t Tuple) Datum {
	return t.Dat
}
