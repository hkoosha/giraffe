package pipeimpl

import (
	"slices"
	"strings"
	"time"

	"github.com/hkoosha/giraffe/cmd"
	"github.com/hkoosha/giraffe/core/inmem"
	"github.com/hkoosha/giraffe/internal"
	"github.com/hkoosha/giraffe/internal/queryimpl"
)

var cache = inmem.Make[PipeImpl](
	"github.com/hkoosha/giraffe|parse_pipe",
	7*24*time.Hour,
)

type PipeImpl struct {
	steps  []queryimpl.QueryImpl
	sSteps []string
}

func (p PipeImpl) String() string {
	return strings.Join(p.sSteps, " | ")
}

func (p PipeImpl) Steps() []queryimpl.QueryImpl {
	return slices.Clone(p.steps)
}

func (p PipeImpl) SSteps() []string {
	return slices.Clone(p.sSteps)
}

func parse(
	spec string,
) (PipeImpl, error) {
	queries := strings.Split(spec, cmd.Pipe.String())
	steps := make([]queryimpl.QueryImpl, len(queries))
	sSteps := make([]string, len(queries))
	for i, spec := range queries {
		parsed, err := internal.Parse(spec)
		if err != nil {
			return PipeImpl{}, err
		}
		steps[i] = parsed
		sSteps[i] = parsed.String()
	}

	return PipeImpl{
		steps:  steps,
		sSteps: sSteps,
	}, nil
}

func Parse(
	spec string,
) (PipeImpl, error) {
	cached, ok := cache.Get(spec)

	if !ok {
		pipe, err := parse(spec)
		cache.Set(spec, pipe, err)
		return pipe, err
	}

	return cached.Unpack()
}
