package hippo

import (
	"slices"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/g"
	"github.com/hkoosha/giraffe/typing"
	"github.com/hkoosha/giraffe/zebra/z"
)

var Plan_ = &Plan{
	compensator: Compensator{
		comp: make([]compCondition, 0),
	},
	registry: FnRegistry{
		scope:  nil,
		byName: make(map[string]*Fn),
		byType: make(map[typing.Type]*Fn),
	},
	steps: make([]namedStep, 0),
}

type Plan struct {
	compensator Compensator
	registry    FnRegistry
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

func (p *Plan) WithNext(
	name string,
	fn *Fn,
) *Plan {
	cp := p.clone()
	cp.steps = z.Appended(cp.steps, newNamedStep(name, fn))

	return cp
}

func (p *Plan) WithNextNamed(
	name string,
) (*Plan, error) {
	cp := p.clone()

	fn, err := cp.registry.Get(name)
	if err != nil {
		return nil, err
	}

	cp.steps = z.Appended(cp.steps, newNamedStep(name, fn))

	return cp, nil
}

func (p *Plan) WithNextTyped(
	ty typing.Type,
) (*Plan, error) {
	cp := p.clone()

	fn, err := cp.registry.GetTyped(ty)
	if err != nil {
		return nil, err
	}

	cp.steps = z.Appended(cp.steps, newNamedStep(ty.String(), fn))

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

func (p *Plan) AndRegistry(
	reg FnRegistry,
) *Plan {
	cp := p.clone()
	cp.registry = cp.registry.Merge(reg)

	return cp
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
		value = g.Join(p.Names()...)
	}

	return prefix + value + suffix
}
