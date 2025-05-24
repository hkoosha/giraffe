package t11y

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
)

//goland:noinspection GoUnusedExportedFunction
func DieIf(
	err any,
) {
	if err == nil {
		return
	}

	defer func() {
		_, _ = os.Stderr.WriteString("dying on error")
		os.Exit(13)
	}()

	_, err = fmt.Fprintf(
		os.Stderr,
		"%s\n\n[panic]\n%v\n\n%s\n\n",
		strings.Repeat("=", 80),
		formatMsg(err),
		FmtStacktraceOf(err),
	)
	if err != nil {
		panic(err)
	}
}

func Join(
	err any,
	with any,
) error {
	e0 := toErr(err)
	e1 := toErr(with)

	return errors.Join(e0, e1)
}

func Mix(
	err *atomic.Value,
	with any,
) {
	e0 := toErr(err.Load())
	e1 := toErr(with)
	err.Store(Join(e0, e1))
}

func MixAndGet(
	err *atomic.Value,
	with any,
) error {
	e0 := toErr(err.Load())
	e1 := toErr(with)
	join := Join(e0, e1)
	err.Store(join)
	return join
}

func toErr(a any) error {
	switch cast, ok := a.(error); {
	case ok:
		return cast

	case a != nil:
		//nolint:err113
		return fmt.Errorf("%v", a)

	default:
		return nil
	}
}
