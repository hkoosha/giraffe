package pipelines

import (
	"errors"
	"maps"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/typing"
)

var errNoSuchFn = errors.New("no such function")

func NewFnRegistry() FnRegistry {
	return FnRegistry{
		scope:  nil,
		byName: make(map[string]Fn),
		byType: make(map[typing.Type]Fn),
	}
}

type FnRegistry struct {
	scope  *giraffe.Query
	byName map[string]Fn
	byType map[typing.Type]Fn
}

func (r FnRegistry) clone() FnRegistry {
	cp := FnRegistry{
		scope:  r.scope,
		byName: maps.Clone(r.byName),
		byType: maps.Clone(r.byType),
	}

	if cp.byName == nil {
		cp.byName = make(map[string]Fn)
	}
	if cp.byType == nil {
		cp.byType = make(map[typing.Type]Fn)
	}

	return cp
}

func (r FnRegistry) WithNamed(
	name string,
	fn Fn,
) (FnRegistry, error) {
	if fn == nil {
		panic(EF("nil fn"))
	}

	if !stepNameRe.MatchString(name) {
		return fnRegistryErr, newPipelineInvalidStepName(name)
	}

	cp := r.clone()
	if _, ok := cp.byName[name]; ok {
		return fnRegistryErr, newPipelineDuplicatedFnError(name)
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

	if !stepNameRe.MatchString(name) {
		return fnRegistryErr, newPipelineInvalidStepName(name)
	}

	cp := r.clone()
	if _, ok := cp.byName[name]; ok {
		return fnRegistryErr, newPipelineDuplicatedFnError(name)
	}

	cp.byName[name] = fn

	return cp, nil
}

func (r FnRegistry) WithTyped(
	ty typing.Type,
	fn Fn,
) (FnRegistry, error) {
	if fn == nil {
		panic(EF("nil fn"))
	}

	cp := r.clone()
	if _, ok := cp.byType[ty]; ok {
		return fnRegistryErr, newPipelineDuplicatedFnError(ty.String())
	}

	cp.byType[ty] = fn

	return cp, nil
}

func (r FnRegistry) MustWithNamed(
	name string,
	fn Fn,
) FnRegistry {
	return M(r.WithNamed(name, fn))
}

func (r FnRegistry) MustWithTyped(
	ty typing.Type,
	fn Fn,
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
) (Fn, error) {
	if r.byName == nil {
		return nil, E(errNoSuchFn)
	}

	if fn, ok := r.byName[name]; ok {
		return fn, nil
	}

	return nil, E(errNoSuchFn)
}

func (r FnRegistry) GetTyped(
	ty typing.Type,
) (Fn, error) {
	if r.byType == nil {
		return nil, E(errNoSuchFn)
	}

	if fn, ok := r.byType[ty]; ok {
		return fn, nil
	}

	return nil, E(errNoSuchFn)
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
