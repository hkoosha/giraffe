package g11y

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/hkoosha/giraffe/g11y/internal"
	"github.com/hkoosha/giraffe/internal/g"
)

func EnableDefaultTracer() {
	EnableTracer()
	SetSkippedLines(true, nil)
	SetCollapsedLines(true, nil)
}

func Traced(err error) error {
	switch {
	case !internal.IsTracer.Load():
		return err

	case g.Is[*TracedError](err):
		tE := g.As[*TracedError](err)
		tE.Add(debug.Stack())

		return err

	default:
		return errors.Join(
			err, &TracedError{
				Mu:          &sync.Mutex{},
				Stacktraces: [][]byte{debug.Stack()},
			},
		)
	}
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
