package hippoerr

import (
	"github.com/hkoosha/giraffe/typing"
)

type PlanErrorState struct {
	fnName string
	fnType typing.Type
}

func (e *PlanErrorState) String(
	hE *HippoError,
) string {
	_ = e.fnType
	_ = e.fnName
	return "TODO::planErrorState :: " + hE.msg
}

func NewRegDuplicateFnError(
	name string,
) error {
	return NewHippoError(
		ErrCodeDuplicateFn,
		"duplicate fn: "+name,
		&PlanErrorState{
			fnType: typing.OfErr(),
			fnName: name,
		},
	)
}

func NewPlanInvalidFnName(
	name string,
) error {
	return NewHippoError(
		ErrCodeInvalidStepName,
		"bad fn name: "+name,
		&PlanErrorState{
			fnType: typing.OfErr(),
			fnName: name,
		},
	)
}

func NewRegMissingFnError(
	ty typing.Type,
	name string,
) error {
	return NewHippoError(
		ErrCodeMissingFn,
		"missing fn: "+name,
		&PlanErrorState{
			fnName: name,
			fnType: ty,
		},
	)
}
