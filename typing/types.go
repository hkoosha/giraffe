package typing

import (
	"fmt"
	"math"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/hkoosha/giraffe/core/t11y"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

const minId = 11

var (
	idCnt atomic.Uint64
	mu    = sync.RWMutex{}
	byTy  = make(map[reflect.Type]*info)
	byId  = make(map[uint64]*info)

	typeErr = Type{
		id: 1,
	}
)

type info struct {
	//nolint:unused
	ty   reflect.Type
	name string
	//nolint:unused
	clonedFrom []uint64
	typ        Type
	isVirtual  bool
}

func init() {
	idCnt.Add(minId - 1)
}

func newTypeInfo(
	ty reflect.Type,
	id uint64,
) *info {
	i := &info{
		ty:         ty,
		name:       fmt.Sprintf("Type[%s@%d]", ty.Name(), id),
		isVirtual:  false,
		clonedFrom: nil,
		typ: Type{
			id: id,
		},
	}

	return i
}

func newVirtualInfo(
	id uint64,
	clonedFrom []uint64,
) *info {
	i := &info{
		ty:         nil,
		name:       fmt.Sprintf("VirtualType[%d]", id),
		isVirtual:  true,
		clonedFrom: clonedFrom,
		typ: Type{
			id: id,
		},
	}

	return i
}

func getOrRegisterTy(
	ty reflect.Type,
) *info {
	t11y.NonNil(ty)

	inf := func(ty reflect.Type) *info {
		mu.RLock()
		defer mu.RUnlock()

		read, ok := byTy[ty]
		if !ok {
			return nil
		}

		return read
	}(ty)

	if inf == nil {
		inf = func(ty reflect.Type) *info {
			mu.Lock()
			defer mu.Unlock()

			read, ok := byTy[ty]
			if !ok {
				read = newTypeInfo(ty, idCnt.Add(1))
			}

			byTy[ty] = read
			byId[read.typ.id] = read

			return read
		}(ty)
	}

	return inf
}

func get(
	id uint64,
) (*info, bool) {
	mu.RLock()
	defer mu.RUnlock()

	read, ok := byId[id]
	if !ok {
		return nil, false
	}

	return read, true
}

func mustGet(
	id uint64,
) *info {
	inf, ok := get(id)
	Assertf(ok, "invalid type: %d", id)

	return inf
}

func registerVirtual(
	clonedFrom []uint64,
) *info {
	mu.Lock()
	defer mu.Unlock()

	if idCnt.Load() == math.MaxUint64 {
		panic(EF("type id pool exhausted"))
	}

	inf := newVirtualInfo(idCnt.Add(1), clonedFrom)
	byId[inf.typ.id] = inf

	return inf
}

type Type struct {
	id uint64
}

func (t Type) ensure() {
	Assertf(t.IsValid(), "invalid type")
}

func (t Type) String() string {
	if !t.IsValid() {
		return "Type[invalid]"
	}

	inf, ok := get(t.id)
	if !ok {
		return "Type[unknown]"
	}

	return inf.name
}

func (t Type) Id() uint64 {
	return t.id
}

func (t Type) IsValid() bool {
	// Id pool is private and id assignment is controlled by this pkg only,
	// hence we don't need to check for the existence of the type in type
	// registry.
	return t.id >= minId
}

func (t Type) IsVirtual() bool {
	t.ensure()
	return mustGet(t.id).isVirtual
}

func (t Type) Clone() Type {
	t.ensure()
	return registerVirtual([]uint64{t.id}).typ
}

// =============================================================================.

func Of[T any]() Type {
	ty := reflect.TypeOf((*T)(nil)).Elem()

	return getOrRegisterTy(ty).typ
}

func OfVirtual() Type {
	return registerVirtual(nil).typ
}

func OfErr() Type {
	return typeErr
}
