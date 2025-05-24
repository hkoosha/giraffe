package hippo

import (
	"regexp"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/hippo/internal"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
)

func newNamedStep(
	name string,
	fn *Fn,
	arg *giraffe.Datum,
) namedStep {
	if name == "" {
		panic(EF("empty step name"))
	}
	if fn == nil {
		panic(EF("nil step fn"))
	}
	if !internal.SimpleName.MatchString(name) {
		panic(E(hippoerr.NewPlanInvalidFnName(name)))
	}

	return namedStep{
		fn:   fn,
		name: name,
		arg:  arg,
	}
}

type namedStep struct {
	fn   *Fn
	name string
	arg  *giraffe.Datum
}

func (fn *namedStep) String() string {
	const prefix = "NamedStep["
	const suffix = "]"

	value := "nil"
	if fn != nil {
		value = fn.name + "::" + fn.arg.String()
	}

	return prefix + value + suffix
}

// =============================================================================

type compCondition struct {
	onErr  *regexp.Regexp
	onName *regexp.Regexp
	fn     *Fn
	onStep int
}

func (c compCondition) canCompensate(
	sCtx *StepContext,
	err error,
) bool {
	if c.onStep >= 0 && c.onStep != sCtx.stepNo {
		return false
	}

	if c.onName != nil && !c.onName.MatchString(sCtx.stepName) {
		return false
	}

	if c.onErr != nil && !c.onErr.MatchString(err.Error()) {
		return false
	}

	return true
}
