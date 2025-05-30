package hippoerr

import (
	"github.com/hkoosha/giraffe/typing"
)

type planErrorState struct {
	fnType typing.Type
	fnName string
}

func (e *planErrorState) String(
	hE *hippoError,
) string {
	_ = e.fnType
	_ = e.fnName
	return "TODO::planErrorState :: " + hE.msg
}

func NewPlanDuplicateFnError(
	name string,
) error {
	return NewHippoError(
		ErrCodeDuplicateFn,
		"duplicate fn: "+name,
		&planErrorState{
			fnType: nil,
			fnName: name,
		},
	)
}

func NewPlanInvalidStepName(
	name string,
) error {
	return NewHippoError(
		ErrCodeInvalidStepName,
		"bad fn name: "+name,
		&planErrorState{
			fnType: nil,
			fnName: name,
		},
	)
}

func NewPlanMissingFnError(
	ty typing.Type,
	name string,
) error {
	return NewHippoError(
		ErrCodeMissingFn,
		"missing fn: "+name,
		&planErrorState{
			fnName: name,
			fnType: ty,
		},
	)
}
