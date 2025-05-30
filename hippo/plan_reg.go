package hippo

import (
	"maps"
	"strings"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
	"github.com/hkoosha/giraffe/hippo/internal/privnames"
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/g"
	"github.com/hkoosha/giraffe/typing"
	"github.com/hkoosha/giraffe/zebra/z"
)

var fnRegistryErr = FnRegistry{
	scope:  nil,
	byName: nil,
	byType: nil,
}

//goland:noinspection GoUnusedGlobalVariable
var FnRegistry_ = &FnRegistry{
	scope:  nil,
	byName: make(map[string]*Fn),
	byType: make(map[typing.Type]*Fn),
}

type FnRegistry struct {
	scope  *giraffe.Query
	byName map[string]*Fn
	byType map[typing.Type]*Fn
}

func (r FnRegistry) String() string {
	const prefix = "FnRegistry["
	const suffix = "]"

	value := strings.Builder{}

	value.WriteString("scope=")
	if r.scope != nil {
		value.WriteString(r.scope.String())
	} else {
		value.WriteString("nil")
	}

	value.WriteString(", names=[")
	value.WriteString(g.JoinIt(maps.Keys(r.byName)))
	value.WriteString("]")

	value.WriteString(", types=[")
	value.WriteString(g.Join(z.ItApplied(maps.Keys(r.byType), func(it typing.Type) string {
		return it.String()
	})...))
	value.WriteString("]")

	return prefix + value.String() + suffix
}

func (r FnRegistry) clone() FnRegistry {
	cp := FnRegistry{
		scope:  r.scope,
		byName: maps.Clone(r.byName),
		byType: maps.Clone(r.byType),
	}

	if cp.byName == nil {
		cp.byName = make(map[string]*Fn)
	}
	if cp.byType == nil {
		cp.byType = make(map[typing.Type]*Fn)
	}

	return cp
}

func (r FnRegistry) WithNamed(
	name string,
	fn *Fn,
) (FnRegistry, error) {
	if fn == nil {
		panic(EF("nil fn"))
	}

	if !privnames.SimpleName.MatchString(name) {
		return fnRegistryErr, hippoerr.NewPlanInvalidStepName(name)
	}

	cp := r.clone()
	if _, ok := cp.byName[name]; ok {
		return fnRegistryErr, hippoerr.NewPlanDuplicateFnError(name)
	}

	cp.byName[name] = fn

	return cp, nil
}

func (r FnRegistry) WithTypeNamedAs(
	ty typing.Type,
	name string,
) (FnRegistry, error) {
	fn, err := r.GetTyped(ty)
	if err != nil {
		return fnRegistryErr, err
	}

	if !privnames.SimpleName.MatchString(name) {
		return fnRegistryErr, hippoerr.NewPlanInvalidStepName(name)
	}

	cp := r.clone()
	if _, ok := cp.byName[name]; ok {
		return fnRegistryErr, hippoerr.NewPlanDuplicateFnError(name)
	}

	cp.byName[name] = fn

	return cp, nil
}

func (r FnRegistry) WithTyped(
	ty typing.Type,
	fn *Fn,
) (FnRegistry, error) {
	if fn == nil {
		panic(EF("nil fn"))
	}

	cp := r.clone()
	if _, ok := cp.byType[ty]; ok {
		return fnRegistryErr, hippoerr.NewPlanDuplicateFnError(ty.String())
	}

	cp.byType[ty] = fn

	return cp, nil
}

func (r FnRegistry) MustWithNamed(
	name string,
	fn *Fn,
) FnRegistry {
	return M(r.WithNamed(name, fn))
}

func (r FnRegistry) MustWithTyped(
	ty typing.Type,
	fn *Fn,
) FnRegistry {
	return M(r.WithTyped(ty, fn))
}

func (r FnRegistry) MustWithTypeNamedAs(
	ty typing.Type,
	name string,
) FnRegistry {
	return M(r.WithTypeNamedAs(ty, name))
}

func (r FnRegistry) Get(
	name string,
) (*Fn, error) {
	if r.byName == nil {
		return nil, hippoerr.NewPlanMissingFnError(nil, name)
	}

	if fn, ok := r.byName[name]; ok {
		return fn, nil
	}

	return nil, hippoerr.NewPlanMissingFnError(nil, name)
}

func (r FnRegistry) GetTyped(
	ty typing.Type,
) (*Fn, error) {
	if r.byType == nil {
		return nil, hippoerr.NewPlanMissingFnError(ty, "")
	}

	if fn, ok := r.byType[ty]; ok {
		return fn, nil
	}

	return nil, hippoerr.NewPlanMissingFnError(ty, "")
}

func (r FnRegistry) Merge(
	other FnRegistry,
) FnRegistry {
	cp := r.clone()

	for k, v := range other.byName {
		if _, ok := cp.byName[k]; !ok {
			cp.byName[k] = v
		}
	}

	for k, v := range other.byType {
		if _, ok := cp.byType[k]; !ok {
			cp.byType[k] = v
		}
	}

	return cp
}

func (r FnRegistry) Dump() FnRegistry {
	return r
}
