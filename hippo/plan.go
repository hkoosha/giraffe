package hippo

import (
	"fmt"
	"slices"

	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/internal/gstrings"
	"github.com/hkoosha/giraffe/typing"
	"github.com/hkoosha/giraffe/zebra/z"
)

func MkPlan() *Plan {
	return &Plan{
		compensator: Compensator{
			comp: make([]compCondition, 0),
		},
		registry: &FnRegistry{
			scope:  nil,
			byType: nil,
		},
		steps: make([]namedStep, 0),
	}
}

type Plan struct {
	compensator Compensator
	registry    *FnRegistry
	steps       []namedStep
}

func (p *Plan) clone() *Plan {
	return &Plan{
		compensator: p.compensator.clone(),
		registry:    p.registry.clone(),
		steps:       slices.Clone(p.steps),
	}
}

func (p *Plan) Names() []string {
	names := z.Applied(p.steps, func(it namedStep) string {
		return it.name
	})

	if names == nil {
		names = make([]string, 0)
	}

	return names
}

func (p *Plan) Dump() *Plan {
	return p
}

func (p *Plan) MustWithNext(
	name string,
	fn *Fn,
) *Plan {
	return M(p.WithNext(name, fn))
}

func (p *Plan) MustWithNextExe(
	name string,
	exe Exe,
) *Plan {
	return M(p.WithNext(name, FnOf(exe)))
}

func (p *Plan) WithNextExe(
	name string,
	exe Exe,
) (*Plan, error) {
	return p.WithNext(name, FnOf(exe))
}

func (p *Plan) WithNext(
	name string,
	fn *Fn,
) (*Plan, error) {
	cp := p.clone()

	if ok, err := p.registry.has(fn, name); err != nil {
		return nil, err
	} else if !ok {
		cp.registry = cp.registry.MustWith(fn, name)
	}

	cp.steps = z.Appended(cp.steps, newNamedStep(name, fn, nil))

	return cp, nil
}

func (p *Plan) WithSteps(
	fns ...FnConfig,
) (*Plan, error) {
	cp := p.clone()

	for i, f := range fns {
		if err := f.Validate(); err != nil {
			return nil, E(err...)
		}

		fn, err := p.registry.Named(f.Fn)
		if err != nil {
			return nil, err
		}

		fn, err = f.Configure(fn)
		if err != nil {
			return nil, err
		}

		stepName := fmt.Sprintf("%s#%04d", f.Fn, i)
		cp.steps = z.Appended(cp.steps, newNamedStep(stepName, fn, Clone(f.Args)))
	}

	return cp, nil
}

func (p *Plan) WithNextNamed(
	name string,
) (*Plan, error) {
	fn, err := p.registry.Named(name)
	if err != nil {
		return nil, err
	}

	cp := p.clone()
	cp.steps = z.Appended(cp.steps, newNamedStep(name, fn, nil))

	return cp, nil
}

func (p *Plan) WithNextTyped(
	ty typing.Type,
) (*Plan, error) {
	cp := p.clone()

	fn, err := cp.registry.Typed(ty)
	if err != nil {
		return nil, err
	}

	cp.steps = z.Appended(cp.steps, newNamedStep(ty.String(), fn, nil))

	return cp, nil
}

func (p *Plan) MustWithNextTyped(
	ty typing.Type,
) *Plan {
	return M(p.WithNextTyped(ty))
}

func (p *Plan) MustWithNextNamed(
	name string,
) *Plan {
	return M(p.WithNextNamed(name))
}

func (p *Plan) MustAndRegistry(
	reg *FnRegistry,
) *Plan {
	return M(p.AndRegistry(reg))
}

func (p *Plan) AndRegistry(
	reg *FnRegistry,
) (*Plan, error) {
	merged, err := p.registry.Merge(reg)
	if err != nil {
		return nil, err
	}

	cp := p.clone()
	cp.registry = merged

	return cp, nil
}

func (p *Plan) AndCompensator(
	c Compensator,
) *Plan {
	cp := p.clone()

	cp.compensator = Compensator{
		comp: z.Appended(cp.compensator.comp, c.comp...),
	}

	return cp
}

func (p *Plan) WithCompensator(
	c Compensator,
) *Plan {
	cp := p.clone()
	cp.compensator = c.clone()
	return cp
}

func (p *Plan) String() string {
	const prefix = "Plan["
	const suffix = "]"

	value := "nil"
	if p != nil {
		value = gstrings.Joined(p.Names())
	}

	return prefix + value + suffix
}
