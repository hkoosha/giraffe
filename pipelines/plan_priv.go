package pipelines

import (
	"fmt"
	"regexp"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/dot"
)

var stepNameRe = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9_]*`)

func newPipelineDuplicatedFnError(
	name string,
) error {
	return newPipelineError(
		ErrCodeDuplicateFn,
		"duplicate fn: "+name,
		"$constructor",
		-1,
		giraffe.OfErr(),
	)
}

func newPipelineInvalidStepName(
	name string,
) error {
	return newPipelineError(
		ErrCodeInvalidStepName,
		"bad fn name: "+name,
		"$constructor",
		-1,
		giraffe.OfErr(),
	)
}

// =============================================================================.

func newNamedStep(
	name string,
	fn Fn,
) namedStep {
	if name == "" {
		panic(EF("empty step name"))
	}
	if fn == nil {
		panic(EF("nil step fn"))
	}

	return namedStep{
		fn:   fn,
		name: name,
	}
}

type namedStep struct {
	fn   Fn
	name string
}

func (fn *namedStep) String() string {
	return fmt.Sprintf("namedStep[%s]", fn.name)
}
