package pipelines

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/dot"
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

type Probe func(
	context.Context,
	*StepContext,
)

type AfterProbe func(
	context.Context,
	*StepContext,
	giraffe.Datum,
	error,
)

type StepContext struct {
	Fn       Fn
	StepName string
	dat      giraffe.Datum
	StepNo   int
}

func (s *StepContext) clone() *StepContext {
	cp := *s

	return &cp
}

func Runner(
	plan Plan,
) (*RunnerFn, error) {
	if len(plan.steps) == 0 {
		panic(EF("empty plan"))
	}

	return &RunnerFn{
		before: nil,
		after:  nil,
		plan:   plan,
	}, nil
}

type RunnerFn struct {
	before Probe
	after  Probe
	plan   Plan
}

func (r *RunnerFn) String() string {
	return fmt.Sprintf("RunnerFn[%s]", g.Join(r.plan.Names()...))
}

func (r *RunnerFn) clone() *RunnerFn {
	return &RunnerFn{
		plan:   r.plan,
		before: r.before,
		after:  r.after,
	}
}

func (r *RunnerFn) WithBefore(
	probe Probe,
) *RunnerFn {
	clone := r.clone()
	clone.before = probe

	return clone
}

func (r *RunnerFn) WithAfter(
	probe Probe,
) *RunnerFn {
	clone := r.clone()
	clone.after = probe

	return clone
}

func (r *RunnerFn) Ekran(
	ctx context.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	hist, hErr := history(dat)
	if hErr != nil {
		return dErr, hErr
	}

	for i, fn := range r.plan.steps {
		sCtx := StepContext{
			StepNo:   i,
			StepName: fn.name,
			Fn:       fn.fn,
			dat:      dat,
		}

		next, eErr := r.exe(ctx, &sCtx)
		if eErr != nil {
			return dErr, onFnErr(&sCtx, hist, eErr)
		}

		if cErr := chkResult(&sCtx, hist, next); cErr != nil {
			return dErr, cErr
		}

		merged, mErr := dat.Merge(next)
		if mErr != nil {
			return dErr, onMergeErr(&sCtx, hist, mErr, dat, next)
		}

		dat = merged
		hist = M(hist.Append(
			M(OfN(
				P(qName, sCtx.StepName),
				P(qState, dat),
			)),
		))
	}

	return Of0(giraffe.Implode{
		qFin:   dat,
		qSteps: Of0(hist),
	}), nil
}

func (r *RunnerFn) exe(
	ctx context.Context,
	sCtx *StepContext,
) (giraffe.Datum, error) {
	if r.before != nil {
		r.before(ctx, sCtx.clone())
	}

	next, err := sCtx.Fn(ctx, sCtx.dat)
	if err != nil {
		if fix, ok := r.plan.compensator.compensate(ctx, sCtx, err); ok {
			next = fix
			err = nil
		}
	}

	if r.after != nil {
		r.after(ctx, sCtx.clone())
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
	if err := g.As[*PipelineError](err); err != nil &&
		err.code == ErrCodeMissingKeys {
		return newPipelineError(
			ErrCodeMissingKeys,
			"missing keys",
			sCtx.StepName,
			sCtx.StepNo,
			history,
			err.queries...,
		)
	}

	return errors.Join(err, newPipelineError(
		ErrCodeFailedStep,
		"pipeline error",
		sCtx.StepName,
		sCtx.StepNo,
		Of0(history),
	))
}

func chkResult(
	sCtx *StepContext,
	history giraffe.Datum,
	next giraffe.Datum,
) error {
	if !next.Type().IsObj() {
		return newPipelineError(
			ErrCodeInvalidStepResult,
			"invalid step result. expecting an object, got: "+next.Type().String(),
			sCtx.StepName,
			sCtx.StepNo,
			Of0(history),
		)
	}

	return nil
}

func onMergeErr(
	sCtx *StepContext,
	history giraffe.Datum,
	err error,
	dat giraffe.Datum,
	next giraffe.Datum,
) error {
	msg := strings.Builder{}
	msg.WriteString("error on merging step data")
	msg.WriteString("\nstate: ")
	msg.WriteString(dat.Pretty())
	msg.WriteString("\nnext: ")
	msg.WriteString(next.Pretty())
	msg.WriteByte('\n')

	return errors.Join(err, newPipelineError(
		ErrCodeFailedStep,
		msg.String(),
		sCtx.StepName,
		sCtx.StepNo,
		Of0(history),
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
