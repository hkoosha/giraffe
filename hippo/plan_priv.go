package hippo

import (
	"regexp"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/hippo/internal"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
	"github.com/hkoosha/giraffe/typing"
)

func newNamedStep(
	name string,
	fn *Fn,
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
	}
}

type namedStep struct {
	fn   *Fn
	name string
}

func (fn *namedStep) String() string {
	const prefix = "NamedStep["
	const suffix = "]"

	value := "nil"
	if fn != nil {
		value = fn.name
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

// =============================================================================

var zeroPlan = &Plan{
	compensator: Compensator{
		comp: make([]compCondition, 0),
	},
	registry: &FnRegistry{
		scope:  nil,
		byType: nil,
	},
	steps: make([]namedStep, 0),
}

var zeroRegistry = &FnRegistry{
	scope:  nil,
	byType: make(map[typing.Type]regEntry),
}
