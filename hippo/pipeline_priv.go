package hippo

import (
	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
)

var (
	qFin   = giraffe.Q("fin")
	qSteps = giraffe.Q("steps")
	qName  = giraffe.Q("name")
	qState = giraffe.Q("state")

	stepInit = "init"

	dErr = giraffe.OfErr()
)

func (n *PipelineFn) ekran(
	ctx gtx.Context,
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
			arg:      fn.arg,
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
			M(giraffe.OfN(
				giraffe.TupleOf(qName, sCtx.stepName),
				giraffe.TupleOf(qState, dat),
			)),
		))
	}

	return giraffe.Of(giraffe.Implode{
		qFin:   dat,
		qSteps: giraffe.Of(hist),
	}), nil
}

func (n *PipelineFn) exe(
	ctx gtx.Context,
	sCtx *StepContext,
) (giraffe.Datum, error) {
	if n.before != nil {
		n.before(ctx, sCtx.clone())
	}

	next, err := sCtx.fn.call(ctx, mkCall(sCtx.stepName, sCtx.dat, sCtx.arg))
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

func (n *PipelineFn) shallow() *PipelineFn {
	return &PipelineFn{
		plan:   n.plan,
		before: n.before,
		after:  n.after,
	}
}

func onFnErr(
	sCtx *StepContext,
	history giraffe.Datum,
	err error,
) error {
	// Var msg string
	// if cast := g.As[*hippoerr.HippoError](err); cast != nil {
	// 	if cast.Code() == hippoerr.ErrCodeMissingKeys {
	// 		cast.State().String()
	// 	}
	// }.

	return E(err, hippoerr.NewPipelineStepError(
		sCtx.stepName,
		sCtx.stepNo,
		sCtx.fn.String(),
		history,
	))
}

func history(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	ini, err := giraffe.OfN(
		giraffe.TupleOf(qName, stepInit),
		giraffe.TupleOf(qState, dat),
	)
	if err != nil {
		return dErr, err
	}

	return giraffe.Of([]giraffe.Datum{ini}), err
}
