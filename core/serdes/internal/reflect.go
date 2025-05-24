package internal

import (
	"errors"
	"maps"
	"reflect"
	"slices"
	"time"

	"github.com/hkoosha/giraffe/core/inmem"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

var cache = inmem.Make[bool](
	"github.com/hkoosha/giraffe|serdes|type_reflect",
	24*time.Hour,
)

var (
	tErr                 = reflect.TypeOf((*error)(nil)).Elem()
	errContainsNoMethods = errors.New("type does not have methods")
)

// TODO while exact generic type is erased, same generic among methods of same
// implementor can still be checked.
type methodSig struct {
	In         int
	Out        int
	ReturnsErr bool
}

func extractMethods(
	t reflect.Type,
) (map[string]methodSig, error) {
	if t.Kind() != reflect.Interface && t.Kind() != reflect.Struct {
		return nil, E(errContainsNoMethods)
	}

	methods := make(map[string]methodSig)
	for i := range t.NumMethod() {
		method := t.Method(i)
		out := method.Type.NumOut()
		methods[method.Name] = methodSig{
			In:         method.Type.NumIn(),
			Out:        out,
			ReturnsErr: out > 0 && method.Type.Out(out-1).Implements(tErr),
		}
	}

	return methods, nil
}

func implementsMethods(
	t map[string]methodSig,
	iface map[string]methodSig,
) (bool, error) {
	t = maps.Clone(t)
	for _, k := range slices.Collect(maps.Keys(t)) {
		if _, ok := iface[k]; !ok {
			delete(t, k)
		}
	}

	for name, iSig := range iface {
		tSig, ok := t[name]
		if !ok || iSig != tSig {
			return false, nil
		}
	}

	return true, nil
}

func implementsGenericErased(
	t reflect.Type,
	iface reflect.Type,
) (bool, error) {
	tMethods, err := extractMethods(t)
	if err != nil {
		return false, err
	}

	iMethods, err := extractMethods(iface)
	if err != nil {
		return false, err
	}

	return implementsMethods(tMethods, iMethods)
}

func ImplementsGenericErased(
	t reflect.Type,
	iface reflect.Type,
) (bool, error) {
	cached, err := cache.GetOr(
		t.String()+"|"+iface.String(),
		func() (bool, error) {
			return implementsGenericErased(t, iface)
		},
	)
	if err != nil {
		return false, err
	}
	return cached.V, nil
}
