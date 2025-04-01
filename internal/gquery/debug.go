package gquery

import (
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
)

const debugNilId = "<nil>"

var (
	debugQueryId = atomic.Uint64{}

	//nolint:exhaustruct
	hasDebug = reflect.TypeOf(Query{}) == reflect.TypeOf(QueryDebug{})
)

func init() {
	debugQueryId.Add(11)
}

func newQueryDebug() DebugImpl {
	if !hasDebug {
		return DebugImpl{}
	}

	db := QueryDebug{
		ID:   ref(debugQueryId.Add(1)),
		Root: ref("?"),
		Str:  ref("?"),
		Seq:  ref(-1),
		All:  nil,
	}

	cast, ok := any(db).(DebugImpl)
	if !ok {
		panic("unreachable: query debug is not a DebugImpl: " +
			reflect.TypeOf(db).String())
	}

	return cast
}

//goland:noinspection GoNameStartsWithPackageName
type QueryDebug struct {
	All  *[]string `json:"all"`
	Root *string   `json:"root"`
	Str  *string   `json:"str"`
	ID   *uint64   `json:"id"`
	Seq  *int      `json:"seq"`
}

func (d *QueryDebug) String() string {
	if d == nil {
		return debugNilId
	}

	props := []string{
		"id=" + strconv.FormatUint(*d.ID, 10),
		"root=" + *d.Root,
		"str=" + *d.Str,
		"seq=" + strconv.Itoa(*d.Seq),
	}

	return strings.Join(props, " | ")
}

func (d *QueryDebug) populate(q Query) string {
	d.Root = ref(q.Root().String())
	d.Str = ref(q.String())
	d.Seq = ref(q.Flags().Seq())

	return d.String()
}

func debugPopulateQueries(
	path []Query,
) {
	if !hasDebug {
		return
	}

	all := make([]string, len(path))

	for i, q := range path {
		if dDebug, ok := any(&path[i].Debug).(*QueryDebug); ok {
			all[i] = dDebug.populate(q)
		} else {
			panic("unreachable: debug not enabled")
		}
	}

	for i := range path {
		if dDebug, ok := any(&path[i].Debug).(*QueryDebug); ok {
			dDebug.All = &all
		} else {
			panic("unreachable: debug not enabled")
		}
	}
}

func ref[T any](t T) *T {
	return &t
}
