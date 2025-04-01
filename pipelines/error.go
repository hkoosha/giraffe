package pipelines

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/internal/g"
)

const (
	ErrCodeFailedStep = iota + 1
	ErrCodeMissingKeys
	ErrCodeDuplicateFn
	ErrCodeInvalidStepName
	ErrCodeRemoteCallFailure
	ErrCodeInvalidStepResult
)

//goland:noinspection GoNameStartsWithPackageName
type PipelineError struct {
	msg      string
	stepName string
	state    giraffe.Datum
	queries  []giraffe.Query
	code     int
	step     int
}

func (e *PipelineError) Error() string {
	sb := strings.Builder{}
	sb.WriteString("runner error [")

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
	sb.WriteString(e.msg)

	if queries := g.Joiner(e.queries); queries != "" {
		sb.WriteString(" | queries=")
		sb.WriteString(queries)
	}

	sb.WriteString(" | state=")
	sb.WriteString(e.state.Pretty())

	return sb.String()
}

func (e *PipelineError) Step() int {
	return e.step
}

func (e *PipelineError) Queries() []giraffe.Query {
	return slices.Clone(e.queries)
}

func (e *PipelineError) State() giraffe.Datum {
	return e.state
}

func (e *PipelineError) StepName() string {
	return e.stepName
}

func newRemoteError(
	msg string,
	err error,
) error {
	return E(err, &PipelineError{
		code:     ErrCodeRemoteCallFailure,
		msg:      msg,
		stepName: "remote",
		step:     -1,
		state:    giraffe.OfEmpty(),
		queries:  nil,
	})
}

func newPipelineError(
	code int,
	msg string,
	stepName string,
	step int,
	state giraffe.Datum,
	queries ...giraffe.Query,
) error {
	return E(&PipelineError{
		code:     code,
		msg:      msg,
		stepName: stepName,
		step:     step,
		state:    state,
		queries:  queries,
	})
}

func newPipelineQueriesError(
	code int,
	msg string,
	queries ...giraffe.Query,
) error {
	return newPipelineError(
		code,
		msg,
		"",
		-1,
		giraffe.OfErr(),
		queries...,
	)
}
