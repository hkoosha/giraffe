package hippo

import (
	"regexp"
	"slices"
	"strconv"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	"github.com/hkoosha/giraffe/zebra/z"
)

type Compensator struct {
	comp []compCondition
}

func (c Compensator) String() string {
	const prefix = "Compensator["
	const suffix = "]"

	value := strconv.Itoa(len(c.comp))

	return prefix + value + suffix
}

func (c Compensator) clone() Compensator {
	return Compensator{
		comp: slices.Clone(c.comp),
	}
}

func (c Compensator) compensate(
	ctx gtx.Context,
	sCtx *StepContext,
	err error,
) (giraffe.Datum, bool) {
	for _, comp := range c.comp {
		if comp.canCompensate(sCtx, err) {
			if next, err := comp.fn.call(ctx, mkCall(sCtx.stepName, sCtx.dat, sCtx.arg)); err == nil {
				return next, true
			}
		}
	}

	return dErr, false
}

func (c Compensator) For(
	msg *regexp.Regexp,
	name *regexp.Regexp,
	step int,
	with *Fn,
) Compensator {
	c.comp = z.Appended(c.comp, compCondition{
		onErr:  &*msg,
		onName: &*name,
		onStep: step,
		fn:     with,
	})

	return c
}

func (c Compensator) ForWith(
	msg *regexp.Regexp,
	name *regexp.Regexp,
	step int,
	with giraffe.Datum,
) Compensator {
	return c.For(msg, name, step, Static(with))
}

func (c Compensator) ForError(
	msg *regexp.Regexp,
	with *Fn,
) Compensator {
	c.comp = z.Appended(c.comp, compCondition{
		onErr:  &*msg,
		onName: nil,
		onStep: -1,
		fn:     with,
	})

	return c
}

func (c Compensator) ForErrorWith(
	msg *regexp.Regexp,
	with giraffe.Datum,
) Compensator {
	return c.ForError(msg, Static(with))
}

func (c Compensator) ForStep(
	step int,
	with *Fn,
) Compensator {
	c.comp = z.Appended(c.comp, compCondition{
		onErr:  nil,
		onName: nil,
		onStep: step,
		fn:     with,
	})

	return c
}

func (c Compensator) ForStepWith(
	step int,
	with giraffe.Datum,
) Compensator {
	return c.ForStep(step, Static(with))
}

func (c Compensator) ForNamed(
	name *regexp.Regexp,
	with *Fn,
	steps ...int,
) Compensator {
	if len(steps) > 0 {
		cp := slices.Clone(c.comp)

		for _, step := range steps {
			if step < 0 {
				panic(EF("invalid step: %d", step))
			}

			cp = append(cp, compCondition{
				onErr:  nil,
				onName: &*name,
				onStep: step,
				fn:     with,
			})
		}

		c.comp = cp

		return c
	}

	return Compensator{
		comp: z.Appended(c.comp, compCondition{
			onErr:  nil,
			onName: &*name,
			onStep: -1,
			fn:     with,
		}),
	}
}

func (c Compensator) ForNamedWith(
	name *regexp.Regexp,
	with giraffe.Datum,
	steps ...int,
) Compensator {
	return c.ForNamed(name, Static(with), steps...)
}
