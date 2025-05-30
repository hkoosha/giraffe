package hippo

import (
	"context"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
	. "github.com/hkoosha/giraffe/internal/dot0"
	. "github.com/hkoosha/giraffe/internal/dot1"
	"github.com/hkoosha/giraffe/internal/g"
)

var (
	qFin   = Q("fin")
	qSteps = Q("steps")
	qName  = Q("name")
	qState = Q("state")

	stepInit = "init"

	dErr = giraffe.OfErr()
)

type ProbeBefore = func(
	context.Context,
	*StepContext,
)

type ProbeAfter = func(
	context.Context,
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
	g11y.NonNil(plan)

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
		value = g.Join(n.plan.Names()...)
	}

	return prefix + value + suffix
}

func (n *PipelineFn) clone() *PipelineFn {
	return &PipelineFn{
		plan:   n.plan,
		before: n.before,
		after:  n.after,
	}
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
	ctx context.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	hist, hErr := history(dat)
	if hErr != nil {
		return dErr, hErr
	}

	for i, fn := range n.plan.steps {
		sCtx := StepContext{
			stepNo:   i,
			stepName: fn.name,
			fn:       fn.fn,
			dat:      dat,
		}

		next, eErr := n.exe(ctx, &sCtx)
		if eErr != nil {
			return dErr, onFnErr(&sCtx, hist, eErr)
		}

		merged, mErr := dat.Merge(next)
		if mErr != nil {
			return dErr, onFnErr(&sCtx, hist, mErr)
		}

		dat = merged
		hist = M(hist.Append(
			M(OfN(
				P(qName, sCtx.stepName),
				P(qState, dat),
			)),
		))
	}

	return Of0(giraffe.Implode{
		qFin:   dat,
		qSteps: Of0(hist),
	}), nil
}

func (n *PipelineFn) exe(
	ctx context.Context,
	sCtx *StepContext,
) (giraffe.Datum, error) {
	if n.before != nil {
		n.before(ctx, sCtx.clone())
	}

	next, err := sCtx.fn.call(ctx, sCtx.dat)
	if err != nil {
		if fix, ok := n.plan.compensator.compensate(ctx, sCtx, err); ok {
			next = fix
			err = nil
		}
	}

	if n.after != nil {
		n.after(ctx, sCtx.clone())
	}

	if err != nil {
		return dErr, err
	}

	return next, nil
}

func onFnErr(
	sCtx *StepContext,
	history giraffe.Datum,
	err error,
) error {
	return E(err, hippoerr.NewPipelineStepError(
		sCtx.stepName,
		sCtx.stepNo,
		history,
	))
}

func history(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	ini, err := giraffe.OfN(
		P(qName, stepInit),
		P(qState, dat),
	)
	if err != nil {
		return dErr, err
	}

	return Of0([]giraffe.Datum{ini}), err
}
