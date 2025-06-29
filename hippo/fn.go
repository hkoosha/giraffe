package hippo

import (
	"maps"
	"slices"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/hippo/internal/privnames"
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/typing"
)

type Fn_ struct {
	exe           Exe
	replicated    map[giraffe.Query][]giraffe.Query
	swapped       map[giraffe.Query]giraffe.Query
	scopedOut     giraffe.Query
	scopedIn      giraffe.Query
	name          string
	inputs        []giraffe.Query
	outputs       []giraffe.Query
	optionals     []giraffe.Query
	selected      []giraffe.Query
	noOverwriting bool
	typ           typing.Type
}

func (f *Fn_) Type() typing.Type {
	if f == nil {
		return typing.OfErr()
	}

	return f.typ
}

func (f *Fn_) IsValid() bool {
	return f != nil && f.exe != nil && f.typ.IsValid()
}

func (f *Fn_) AndReplicate(
	replicated map[giraffe.Query][]giraffe.Query,
) *Fn_ {
	f.ensure()

	m := maps.Clone(f.replicated)
	for k, v := range replicated {
		m[k] = slices.Clone(v)
	}

	clone := f.clone()
	clone.replicated = m
	return clone
}

func (f *Fn_) WithReplicated(
	replicated map[giraffe.Query][]giraffe.Query,
) *Fn_ {
	f.ensure()

	m := make(map[giraffe.Query][]giraffe.Query, len(replicated))
	for k, v := range replicated {
		m[k] = slices.Clone(v)
	}

	clone := f.clone()
	clone.replicated = m
	return clone
}

func (f *Fn_) AndSwapping(
	swapping map[giraffe.Query]giraffe.Query,
) *Fn_ {
	f.ensure()

	swapping = maps.Clone(swapping)
	maps.Copy(swapping, f.swapped)

	clone := f.clone()
	clone.swapped = swapping
	return clone
}

func (f *Fn_) WithSwapping(
	swapping map[giraffe.Query]giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.swapped = maps.Clone(swapping)
	return clone
}

func (f *Fn_) AndScopeOut(
	scope giraffe.Query,
) *Fn_ {
	f.ensure()

	return f.WithScopeOut(f.scopedOut.Plus(scope))
}

func (f *Fn_) WithScopeOut(
	scope giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.scopedOut = scope
	return clone
}

func (f *Fn_) AndInputs(
	inputs ...giraffe.Query,
) *Fn_ {
	return f.WithInput(append(inputs, f.inputs...)...)
}

func (f *Fn_) WithInput(
	inputs ...giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.inputs = inputs
	return clone
}

func (f *Fn_) AndOptionals(
	optionals ...giraffe.Query,
) *Fn_ {
	return f.WithOptional(append(optionals, f.optionals...)...)
}

func (f *Fn_) WithOptional(
	optionals ...giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.optionals = optionals
	return clone
}

func (f *Fn_) AndOutputs(
	outputs ...giraffe.Query,
) *Fn_ {
	return f.WithOutput(append(outputs, f.outputs...)...)
}

func (f *Fn_) WithOutput(
	outputs ...giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.outputs = slices.Clone(outputs)
	return clone
}

func (f *Fn_) WithNoOverwriting() *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.noOverwriting = true
	return clone
}

func (f *Fn_) WithoutNoOverwriting() *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.noOverwriting = false
	return clone
}

func (f *Fn_) AndSelect(
	select_ ...giraffe.Query,
) *Fn_ {
	return f.Select(append(select_, f.selected...)...)
}

func (f *Fn_) Select(
	select_ ...giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.selected = slices.Clone(select_)
	return clone
}

func (f *Fn_) SelectAll() *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.selected = nil
	return clone
}

func (f *Fn_) Named(
	name string,
) *Fn_ {
	f.ensure()

	if !privnames.SimpleName.MatchString(name) {
		panic(EF("invalid fn name: %s", name))
	}

	clone := f.clone()
	clone.name = name
	return clone
}

func (f *Fn_) Dump() *Fn_ {
	return f
}
