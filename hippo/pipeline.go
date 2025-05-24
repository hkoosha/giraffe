package hippo

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/t11y"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	"github.com/hkoosha/giraffe/internal/gstrings"
)

type ProbeBefore = func(
	gtx.Context,
	*StepContext,
)

type ProbeAfter = func(
	gtx.Context,
	*StepContext,
	giraffe.Datum,
	error,
)

// ============================================================================.

type StepContext struct {
	fn       *Fn
	stepName string
	dat      giraffe.Datum
	arg      *giraffe.Datum
	stepNo   int
}

func (s *StepContext) clone() *StepContext {
	cp := *s
	return &cp
}

// ====================================.

func MkPipeline(
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
	clone := n.shallow()
	clone.before = probe

	return clone
}

func (n *PipelineFn) WithAfter(
	probe ProbeBefore,
) *PipelineFn {
	clone := n.shallow()
	clone.after = probe

	return clone
}

func (n *PipelineFn) Ekran(
	ctx gtx.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	return n.ekran(ctx, dat)
}
