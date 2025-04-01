package g11y

import (
	"fmt"
	"os"
	"strings"
)

//goland:noinspection GoUnusedExportedFunction
func DieIf(
	err any,
) {
	if err == nil {
		return
	}

	defer func() {
		_, _ = os.Stderr.WriteString("will die on panic")
		os.Exit(13)
	}()

	_, err = fmt.Fprintf(
		os.Stderr,
		"%s\n\n[panic]\n%v\n\n%s\n\n",
		strings.Repeat("=", 80),
		FmtMsg(err),
		FmtStacktraceOf(err),
	)
	if err != nil {
		panic(err)
	}
}
