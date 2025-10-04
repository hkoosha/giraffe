package giraffe

import (
	"fmt"
)

func TupleOf[V Safe](
	q GQuery,
	v V,
) Tuple {
	return Tuple{
		Query: q,
		Dat:   Of(v),
	}
}

type Tuple struct {
	Query GQuery
	Dat   Datum
}

func (t *Tuple) String() string {
	return fmt.Sprintf("Tuple2[%s, %s]", t.Query, t.Dat.String())
}

func (t *Tuple) Implode() Datum {
	return Of1(t.Query, t.Dat)
}

func (t *Tuple) Unpack() (GQuery, Datum) {
	return t.Query, t.Dat
}

func GetQuery(t Tuple) GQuery {
	return t.Query
}

func GetDatum(t Tuple) Datum {
	return t.Dat
}
