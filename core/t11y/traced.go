package t11y

import (
	"errors"
	"fmt"
	"regexp"
	"runtime/debug"
	"sync"

	"github.com/hkoosha/giraffe/core/t11y/internal"
)

func EnableDefaultTracer(
	skippedLines ...*regexp.Regexp,
) {
	EnableTracer()
	SetSkippedLines(skippedLines...)
	SetCollapsedLines()
}

func Traced(err error) error {
	if !internal.IsTracer.Load() {
		return err
	}

	var tE *tracedError
	if errors.As(err, &tE) {
		tE.Add(debug.Stack())
		return err
	}

	return errors.Join(
		err, &tracedError{
			Mu:          &sync.Mutex{},
			Stacktraces: [][]byte{debug.Stack()},
		},
	)
}

func TracedFmt(format string, v ...any) error {
	//nolint:err113
	return Traced(fmt.Errorf(format, v...))
}

func Must[A any](a A, err error) A {
	if err != nil {
		panic(Traced(err))
	}

	return a
}

func Ensure(err error) {
	if err != nil {
		panic(Traced(err))
	}
}
