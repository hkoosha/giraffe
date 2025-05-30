package hippo

import (
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
	"github.com/hkoosha/giraffe/hippo/internal/privnames"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

func newNamedStep(
	name string,
	fn *Fn_,
) namedStep {
	if name == "" {
		panic(EF("empty step name"))
	}
	if fn == nil {
		panic(EF("nil step fn"))
	}
	if !privnames.SimpleName.MatchString(name) {
		panic(hippoerr.NewPlanInvalidFnName(name))
	}

	return namedStep{
		fn:   fn,
		name: name,
	}
}

type namedStep struct {
	fn   *Fn_
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
