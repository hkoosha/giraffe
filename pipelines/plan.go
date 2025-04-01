package pipelines

import (
	"fmt"
	"slices"

	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/internal/g"
	"github.com/hkoosha/giraffe/typing"
	"github.com/hkoosha/giraffe/zebra/z"
)

var (
	fnRegistryErr = FnRegistry{
		scope:  nil,
		byName: nil,
		byType: nil,
	}
	planErr = Plan{
		registry:    fnRegistryErr,
		steps:       nil,
		compensator: Compensator{comp: nil},
	}
)

func NewPlan() Plan {
	return Plan{
		compensator: NewCompensator(),
		registry:    NewFnRegistry(),
		steps:       nil,
	}
}

type Plan struct {
	compensator Compensator
	registry    FnRegistry
	steps       []namedStep
}

func (p Plan) clone() Plan {
	return Plan{
		compensator: p.compensator,
		registry:    p.registry,
		steps:       slices.Clone(p.steps),
	}
}

func (p Plan) Names() []string {
	return z.Applied(p.steps, func(it namedStep) string {
		return it.name
	})
}

func (p Plan) Dump() Plan {
	return p
}

func (p Plan) WithNext(
	name string,
	fn Fn,
) Plan {
	cp := p.clone()
	cp.steps = z.Appended(cp.steps, newNamedStep(name, fn))

	return cp
}

func (p Plan) WithNextNamed(
	name string,
) (Plan, error) {
	cp := p.clone()

	fn, err := cp.registry.Get(name)
	if err != nil {
		return planErr, err
	}

	cp.steps = z.Appended(cp.steps, newNamedStep(name, fn))

	return cp, nil
}

func (p Plan) WithNextTyped(
	ty typing.Type,
) (Plan, error) {
	cp := p.clone()

	fn, err := cp.registry.GetTyped(ty)
	if err != nil {
		return planErr, err
	}

	cp.steps = z.Appended(cp.steps, newNamedStep(ty.String(), fn))

	return cp, nil
}

func (p Plan) MustWithNextTyped(
	ty typing.Type,
) Plan {
	return M(p.WithNextTyped(ty))
}

func (p Plan) MustWithNextNamed(
	name string,
) Plan {
	return M(p.WithNextNamed(name))
}

func (p Plan) MergeRegistry(
	reg FnRegistry,
) Plan {
	cp := p.clone()
	cp.registry = cp.registry.Merge(reg)

	return cp
}

func (p Plan) WithCompensator(
	c Compensator,
) Plan {
	cp := p.clone()

	cp.compensator = Compensator{
		comp: z.Appended(cp.compensator.comp, c.comp...),
	}

	return cp
}

func (p Plan) String() string {
	return fmt.Sprintf("Plan[%s]", g.Join(p.Names()...))
}
