package hippoerr

import (
	"fmt"
	"strings"

	"github.com/hkoosha/giraffe"
)

type PipelineErrorState struct {
	stepName string
	state    giraffe.Datum
	step     int
}

func (e *PipelineErrorState) String(
	hE *HippoError,
) string {
	sb := strings.Builder{}
	sb.WriteString("hippo error [")

	if e.stepName != "" {
		sb.WriteString(e.stepName)
	} else {
		sb.WriteString(strings.Repeat("?", 6))
	}

	sb.WriteByte('#')
	if e.step >= 0 {
		sb.WriteString(fmt.Sprintf("%02d", e.step))
	} else {
		sb.WriteString("??")
	}

	sb.WriteString("]: err=")
	sb.WriteString(hE.msg)

	sb.WriteString(" | state=")
	sb.WriteString(e.state.Pretty())

	return sb.String()
}

func NewPipelineStepError(
	stepName string,
	step int,
	fn string,
	state giraffe.Datum,
) error {
	return NewHippoError(
		ErrCodeFailedStep,
		"failed step: "+fn,
		&PipelineErrorState{
			stepName: stepName,
			step:     step,
			state:    state,
		},
	)
}
