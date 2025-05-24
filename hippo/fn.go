package hippo

import (
	"fmt"
	"maps"
	"slices"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	"github.com/hkoosha/giraffe/hippo/internal"
	"github.com/hkoosha/giraffe/typing"
)

// TODO check duplicates.
// TODO check clashing.

type Exe = func(
	gtx.Context,
	Call,
) (giraffe.Datum, error)

func FnOf(
	exe Exe,
) *Fn {
	t := typing.OfVirtual()

	fn := &Fn{
		exe:          exe,
		inputs:       nil,
		scoped:       nil,
		outputs:      nil,
		optionals:    nil,
		copy:         nil,
		selected:     nil,
		skipOnExists: false,
		skipped:      false,
		skipWith:     nil,
		typ:          t,
		name:         "#" + t.String(),
		// args:      nil,
		// swapped:      nil,
	}

	if !fn.IsValid() {
		panic(E(errInvalidFn))
	}

	return fn
}

type Call interface {
	internal.Sealed

	Data() giraffe.Datum
	AndData(giraffe.Datum) (Call, error)
	WithData(giraffe.Datum) Call

	Args() giraffe.Datum
	AndArgs(giraffe.Datum) (Call, error)
	WithArgs(giraffe.Datum) Call
	WithoutArgs() Call

	Name() string
	WithName(string) Call

	CheckPresent(giraffe.Datum, []giraffe.Query) error
}

type Fn struct {
	exe          Exe
	copy         map[giraffe.Query]giraffe.Query
	name         string
	scoped       *giraffe.Query
	inputs       []giraffe.Query
	combine      map[giraffe.Query][]giraffe.Query
	outputs      []giraffe.Query
	optionals    []giraffe.Query
	selected     []giraffe.Query
	typ          typing.Type
	skipOnExists bool
	skipped      bool
	skipWith     *giraffe.Datum

	// swapped      map[giraffe.Query]giraffe.Query
	// args         []giraffe.Query
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

func (f *Fn) WithoutSkipOnExists() *Fn {
	return f.SetSkipOnExists(true)
}

func (f *Fn) WithSkipOnExists() *Fn {
	return f.SetSkipOnExists(true)
}

func (f *Fn) SetSkipOnExists(b bool) *Fn {
	f.ensure()

	clone := f.clone()
	clone.skipOnExists = b
	return clone
}

func (f *Fn) AndCopied(
	copied map[giraffe.Query]giraffe.Query,
) *Fn {
	f.ensure()

	copied = maps.Clone(copied)
	maps.Copy(copied, f.copy)

	clone := f.clone()
	clone.copy = copied
	return clone
}

func (f *Fn) WithCopied(
	copied map[giraffe.Query]giraffe.Query,
) *Fn {
	f.ensure()

	clone := f.clone()
	clone.copy = maps.Clone(copied)
	return clone
}

func (f *Fn) WithoutCopied() *Fn {
	f.ensure()

	clone := f.clone()
	clone.copy = nil
	return clone
}

/*func (f *Fn) AndSwapping(
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

func (f *Fn) WithoutSwapping() *Fn {
	f.ensure()

	clone := f.clone()
	clone.swapped = nil
	return clone
}
*/

func (f *Fn) WithScope(
	scope giraffe.Query,
) *Fn {
	f.ensure()

	cp := f.clone()
	cp.scoped = &scope
	return cp
}

func (f *Fn) WithoutScope() *Fn {
	f.ensure()

	cp := f.clone()
	cp.scoped = nil
	return cp
}

func (f *Fn) AndInputs(
	inputs ...giraffe.Query,
) *Fn {
	if len(inputs) == 0 {
		panic(EF("no input provided"))
	}

	return f.WithInputs(append(inputs, f.inputs...)...)
}

func (f *Fn) WithInputs(
	inputs ...giraffe.Query,
) *Fn {
	if len(inputs) == 0 {
		panic(EF("no input provided, use WithoutInput for this case"))
	}

	f.ensure()

	clone := f.clone()
	clone.inputs = inputs
	return clone
}

func (f *Fn) WithoutInputs() *Fn {
	f.ensure()

	clone := f.clone()
	clone.inputs = nil
	return clone
}

func (f *Fn) WithCombine(
	combine map[giraffe.Query][]giraffe.Query,
) *Fn {
	if len(combine) == 0 {
		panic(EF("no combine provided, use WithoutCombine for this case"))
	}

	f.ensure()

	clone := f.clone()
	clone.combine = cloneMapOfSlices(combine)
	return clone
}

func (f *Fn) WithoutCombine() *Fn {
	f.ensure()

	clone := f.clone()
	clone.combine = nil
	return clone
}

func (f *Fn) AndOptionals(
	optionals ...giraffe.Query,
) *Fn {
	if len(optionals) == 0 {
		panic(EF("no optionals provided"))
	}

	return f.WithOptional(append(optionals, f.optionals...)...)
}

func (f *Fn) WithOptional(
	optionals ...giraffe.Query,
) *Fn {
	if len(optionals) == 0 {
		panic(EF("no optionals provided, use WithoutOptionals in this case"))
	}

	f.ensure()

	clone := f.clone()
	clone.optionals = optionals
	return clone
}

func (f *Fn) WithoutOptional() *Fn {
	f.ensure()

	clone := f.clone()
	clone.optionals = nil
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

func (f *Fn) WithoutOutput() *Fn {
	f.ensure()

	clone := f.clone()
	clone.outputs = nil
	return clone
}

func (f *Fn) AndSelect(
	select_ ...giraffe.Query,
) *Fn {
	return f.WithSelect(append(select_, f.selected...)...)
}

func (f *Fn) WithSelect(
	select_ ...giraffe.Query,
) *Fn {
	if len(select_) == 0 {
		panic(EF("no select provided, use WithoutSelect in this case"))
	}

	f.ensure()

	clone := f.clone()
	clone.selected = slices.Clone(select_)
	return clone
}

func (f *Fn) WithoutSelect() *Fn {
	f.ensure()

	clone := f.clone()
	clone.selected = nil
	return clone
}

func (f *Fn) Skipped() bool {
	return f.skipped
}

func (f *Fn) WithSkipped() *Fn {
	return f.SetSkipped(true)
}

func (f *Fn) WithoutSkipped() *Fn {
	return f.SetSkipped(false)
}

func (f *Fn) SkippedWith() (giraffe.Datum, bool) {
	if f.skipWith == nil {
		return dErr, false
	}
	return *f.skipWith, true
}

func (f *Fn) WithSkippedWith(d giraffe.Datum) *Fn {
	cp := f.clone()
	cp.skipWith = &d
	return cp
}

func (f *Fn) WithoutSkippedWith() *Fn {
	cp := f.clone()
	cp.skipWith = nil
	return cp
}

func (f *Fn) SetSkipped(v bool) *Fn {
	cp := f.clone()
	cp.skipped = v
	return cp
}

func (f *Fn) Named(
	name string,
) *Fn {
	f.ensure()

	if !internal.ScopedSimpleName.MatchString(name) {
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
