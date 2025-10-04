package hippo

import (
	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/internal/dot1"
	"github.com/hkoosha/giraffe/internal/gstrings"
	"github.com/hkoosha/giraffe/t11y"
	"github.com/hkoosha/giraffe/t11y/gtx"
)

type ProbeBefore = func(
	Context,
	*StepContext,
)

type ProbeAfter = func(
	Context,
	*StepContext,
	giraffe.Datum,
	error,
)

// ============================================================================.

type StepContext struct {
	fn       *Fn
	stepName string
	dat      giraffe.Datum
	stepNo   int
}

func (s *StepContext) clone() *StepContext {
	cp := *s

	return &cp
}

// ====================================.

func Pipeline(
	plan *Plan,
) (*PipelineFn, error) {
	t11y.NonNil(plan)

	if len(plan.steps) == 0 {
		panic(EF("empty plan"))
	}

	return &PipelineFn{
		before: nil,
		after:  nil,
		plan:   plan,
	}, nil
}

type PipelineFn struct {
	before ProbeBefore
	after  ProbeBefore
	plan   *Plan
}

func (n *PipelineFn) String() string {
	const prefix = "PipelineFn["
	const suffix = "]"

	value := "nil"
	if n != nil {
		value = gstrings.Joined(n.plan.Names())
	}

	return prefix + value + suffix
}

func (n *PipelineFn) WithBefore(
	probe ProbeBefore,
) *PipelineFn {
	clone := n.clone()
	clone.before = probe

	return clone
}

func (n *PipelineFn) WithAfter(
	probe ProbeBefore,
) *PipelineFn {
	clone := n.clone()
	clone.after = probe

	return clone
}

func (n *PipelineFn) Ekran(
	ctx Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	htx := gtx.Of(ctx)

	return n.ekran(htx, dat)
}
