package hippo

import (
	"fmt"
	"maps"
	"slices"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/hippo/internal"
	"github.com/hkoosha/giraffe/typing"
)

// TODO check duplicates.
// TODO check clashing.

type Exe = func(
	ctx Context,
	dat giraffe.Datum,
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
		exe:          exe,
		inputs:       nil,
		scoped:       "",
		outputs:      nil,
		optionals:    nil,
		replicated:   nil,
		selected:     nil,
		swapped:      nil,
		skipOnExists: false,
		typ:          t,
		name:         "#" + t.String(),
		args:         nil,
	}

	var err error = nil
	if !fn.IsValid() {
		err = E(errInvalidFn)
	}

	//nolint:nilnil
	return fn, err
}

// FnConfig
//
// TODO mix as private field in Fn itself.
//
//nolint:lll
type FnConfig struct {
	Require      *[]giraffe.Query                 `json:"require,omitempty"        yaml:"require,omitempty"`
	Select       *[]giraffe.Query                 `json:"select,omitempty"         yaml:"select,omitempty"`
	Replicate    *map[giraffe.Query]giraffe.Query `json:"replicate,omitempty"      yaml:"replicate,omitempty"`
	Swap         *map[giraffe.Query]giraffe.Query `json:"swap,omitempty"           yaml:"swap,omitempty"`
	Scoped       *giraffe.Query                   `json:"scoped,omitempty"         yaml:"scoped,omitempty"`
	SkipOnExists *bool                            `json:"skip_on_exists,omitempty" yaml:"require,omitempty"`
	Args         *giraffe.Datum                   `json:"args,omitempty"           yaml:"args,omitempty"`
	Fn           string                           `json:"fn"                       yaml:"fn"`
}

func (f *FnConfig) Clone() *FnConfig {
	require := f.Require
	if require != nil {
		require = Ref(slices.Clone(*require))
	}

	select_ := f.Select
	if select_ != nil {
		select_ = Ref(slices.Clone(*select_))
	}

	replicate := f.Replicate
	if replicate != nil {
		replicate = Ref(maps.Clone(*replicate))
	}

	swap := f.Swap
	if swap != nil {
		swap = Ref(maps.Clone(*swap))
	}

	scoped := f.Scoped
	if scoped != nil {
		scoped = Ref(*scoped)
	}

	skipOnExists := f.SkipOnExists
	if skipOnExists != nil {
		skipOnExists = Ref(*skipOnExists)
	}

	args := f.Args
	if args != nil {
		args = Ref(*args)
	}

	return &FnConfig{
		Fn:           f.Fn,
		Require:      require,
		Select:       select_,
		Replicate:    replicate,
		Swap:         swap,
		Scoped:       scoped,
		SkipOnExists: skipOnExists,
		Args:         args,
	}
}

func (f *FnConfig) Configure(
	fn *Fn,
) (*Fn, error) {
	if f.Require != nil {
		fn = fn.WithInput(*f.Require...)
	}

	if f.Select != nil {
		fn = fn.Select(*f.Select...)
	}

	if f.Replicate != nil {
		fn = fn.WithReplicated(*f.Replicate)
	}

	if f.Swap != nil {
		fn = fn.WithSwapping(*f.Swap)
	}

	if f.Scoped != nil {
		fn = fn.WithScope(*f.Scoped)
	}

	if f.SkipOnExists != nil {
		fn = fn.SetSkipOnExists(*f.SkipOnExists)
	}

	if f.Args != nil {
		fn = fn.WithArgs(*f.Args)
	}

	return fn, nil
}

type Fn struct {
	exe          Exe
	args         *giraffe.Datum
	replicated   map[giraffe.Query]giraffe.Query
	swapped      map[giraffe.Query]giraffe.Query
	name         string
	scoped       giraffe.Query
	inputs       []giraffe.Query
	outputs      []giraffe.Query
	optionals    []giraffe.Query
	selected     []giraffe.Query
	typ          typing.Type
	skipOnExists bool
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

	cp := f.clone()
	cp.scoped = scope
	return cp
}

func (f *Fn) AndArgs(
	args giraffe.Datum,
) (*Fn, error) {
	f.ensure()

	cp := f.clone()

	if cp.args == nil {
		cp.args = &args
	} else {
		merged, err := cp.args.Merge(args)
		if err != nil {
			return nil, err
		}
		cp.args = &merged
	}
	return cp, nil
}

func (f *Fn) WithArgs(
	args giraffe.Datum,
) *Fn {
	f.ensure()

	cp := f.clone()
	cp.args = &args
	return cp
}

func (f *Fn) WithoutArgs() *Fn {
	f.ensure()

	cp := f.clone()
	cp.args = nil
	return cp
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
