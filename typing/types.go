package typing

import (
	"reflect"
	"sync"
)

var (
	mu           = sync.RWMutex{}
	typeRegistry = make(map[reflect.Type]*Type_)
)

func optimistic(
	ty reflect.Type,
) (*Type_, bool) {
	mu.RLock()
	defer mu.RUnlock()

	typ, ok := typeRegistry[ty]

	return typ, ok
}

func pessimistic(
	ty reflect.Type,
) *Type_ {
	mu.Lock()
	defer mu.Unlock()

	typ, ok := typeRegistry[ty]

	if !ok {
		typ = &Type_{ty: ty}
		typeRegistry[ty] = typ
	}

	return typ
}

type Type_ struct {
	ty reflect.Type
}

func (t *Type_) String() string {
	return "Type[" + t.ty.Name() + "]"
}

type Type = *Type_

func TypeOf[T any]() Type {
	ty := reflect.TypeOf((*T)(nil)).Elem()

	if typ, ok := optimistic(ty); ok {
		return typ
	}

	return pessimistic(ty)
}
