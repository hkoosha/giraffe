package hippo

import (
	"fmt"
	"maps"
	"slices"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/hippo/internal"
	. "github.com/hkoosha/giraffe/internal/dot1"
	"github.com/hkoosha/giraffe/typing"
)

// TODO check duplicates.
// TODO check clashing

type ExeCtx = func(
	Context,
	giraffe.Datum,
) (giraffe.Datum, error)

type Exe = func(
	Context,
	giraffe.Datum,
) (giraffe.Datum, error)

// =============================================================================.

func FnOf(
	exe Exe,
) *Fn {
	return M(TryFnOf(exe))
}

func TryFnOf(
	exe Exe,
) (*Fn, error) {
	t := typing.OfVirtual()

	fn := &Fn{
		exe:        exe,
		scoped:     "",
		inputs:     nil,
		outputs:    nil,
		optionals:  nil,
		replicated: nil,
		selected:   nil,
		swapped:    nil,
		typ:        t,
		name:       "#" + t.String(),
	}

	var err error = nil
	if !fn.IsValid() {
		err = E(errInvalidFn)
	}

	//nolint:nilnil
	return fn, err
}

func MustFnCtxOf(
	exe ExeCtx,
) *Fn {
	return M(FnCtxOf(exe))
}

func FnCtxOf(
	exeCtx ExeCtx,
) (*Fn, error) {
	exe := func(
		ctx Context,
		dat giraffe.Datum,
	) (giraffe.Datum, error) {
		return exeCtx(ctx, dat)
	}

	return TryFnOf(exe)
}

type Fn struct {
	exe        Exe
	replicated map[giraffe.Query]giraffe.Query
	swapped    map[giraffe.Query]giraffe.Query
	scoped     giraffe.Query
	name       string
	inputs     []giraffe.Query
	outputs    []giraffe.Query
	optionals  []giraffe.Query
	selected   []giraffe.Query
	typ        typing.Type
}

func (f *Fn) ensure() {
	if !f.IsValid() {
		panic(EF("invalid fn"))
	}
}

func (f *Fn) Type() typing.Type {
	if f == nil {
		return typing.OfErr()
	}

	return f.typ
}

func (f *Fn) IsValid() bool {
	return f != nil && f.exe != nil && f.typ.IsValid()
}

func (f *Fn) AndReplicate(
	replicated map[giraffe.Query]giraffe.Query,
) *Fn {
	f.ensure()

	replicated = maps.Clone(replicated)
	maps.Copy(replicated, f.replicated)

	clone := f.clone()
	clone.replicated = replicated
	return clone
}

func (f *Fn) WithReplicated(
	replicated map[giraffe.Query]giraffe.Query,
) *Fn {
	f.ensure()

	clone := f.clone()
	clone.replicated = maps.Clone(replicated)
	return clone
}

func (f *Fn) AndSwapping(
	swapping map[giraffe.Query]giraffe.Query,
) *Fn {
	f.ensure()

	swapping = maps.Clone(swapping)
	maps.Copy(swapping, f.swapped)

	clone := f.clone()
	clone.swapped = swapping
	return clone
}

func (f *Fn) WithSwapping(
	swapping map[giraffe.Query]giraffe.Query,
) *Fn {
	f.ensure()

	clone := f.clone()
	clone.swapped = maps.Clone(swapping)
	return clone
}

func (f *Fn) WithScope(
	scope giraffe.Query,
) *Fn {
	f.ensure()

	clone := f.clone()
	clone.scoped = scope
	return clone
}

func (f *Fn) AndInputs(
	inputs ...giraffe.Query,
) *Fn {
	return f.WithInput(append(inputs, f.inputs...)...)
}

func (f *Fn) WithInput(
	inputs ...giraffe.Query,
) *Fn {
	f.ensure()

	clone := f.clone()
	clone.inputs = inputs
	return clone
}

func (f *Fn) AndOptionals(
	optionals ...giraffe.Query,
) *Fn {
	return f.WithOptional(append(optionals, f.optionals...)...)
}

func (f *Fn) WithOptional(
	optionals ...giraffe.Query,
) *Fn {
	f.ensure()

	clone := f.clone()
	clone.optionals = optionals
	return clone
}

func (f *Fn) AndOutputs(
	outputs ...giraffe.Query,
) *Fn {
	return f.WithOutput(append(outputs, f.outputs...)...)
}

func (f *Fn) WithOutput(
	outputs ...giraffe.Query,
) *Fn {
	f.ensure()

	clone := f.clone()
	clone.outputs = slices.Clone(outputs)
	return clone
}

func (f *Fn) AndSelect(
	select_ ...giraffe.Query,
) *Fn {
	return f.Select(append(select_, f.selected...)...)
}

func (f *Fn) Select(
	select_ ...giraffe.Query,
) *Fn {
	f.ensure()

	clone := f.clone()
	clone.selected = slices.Clone(select_)
	return clone
}

func (f *Fn) SelectAll() *Fn {
	f.ensure()

	clone := f.clone()
	clone.selected = nil
	return clone
}

func (f *Fn) Named(
	name string,
) *Fn {
	f.ensure()

	if !internal.SimpleName.MatchString(name) {
		panic(EF("invalid fn name: %s", name))
	}

	clone := f.clone()
	clone.name = name
	return clone
}

func (f *Fn) Dump() *Fn {
	return f
}

func (f *Fn) String() string {
	return fmt.Sprintf("Fn[%s][%s]", f.typ, f.name)
}
