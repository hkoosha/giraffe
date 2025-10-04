package g11y

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/hkoosha/giraffe/g11y/internal"
	"github.com/hkoosha/giraffe/internal/g"
)

func SetSkippedLines(
	withDefaults bool,
	lines []*regexp.Regexp,
) {
	s := internal.DeepCopyL1(lines)
	if withDefaults {
		s = append(s, internal.DefaultSkippedLines...)
	}

	internal.SetSkippedLine(s)
}

func SetCollapsedLines(
	withDefaults bool,
	lines []*regexp.Regexp,
) {
	s := internal.DeepCopyL1(lines)
	if withDefaults {
		s = append(s, internal.DefaultCollapsedLines...)
	}

	internal.SetCollapsedLines(s)
}

// ==============================================================================.

func FmtStacktrace(
	stacktrace string,
) []string {
	defer func() {
		if r := recover(); r != nil {
			_, _ = os.Stderr.WriteString("nested panic " + fmt.Sprintf("%+v", r))
			os.Exit(13)
		}
	}()

	fmtLine := func(
		line string,
		fn string,
	) string {
		line = strings.TrimSpace(line)

		if internal.FileLine.MatchString(line) {
			line = strings.Split(line, " ")[0]
		}

		sb := strings.Builder{}
		sb.Grow(len(line) + len(fn) + 3)

		if internal.LineNum.MatchString(line) {
			if s0, s1, cut := strings.Cut(line, ":"); cut {
				sb.WriteString(s0)
				sb.WriteByte(':')
				sb.WriteString(s1)
			} else {
				sb.WriteString("ERROR PARSING THIS LINE: ")
				sb.WriteString(line)
			}
		}

		if fn != "" {
			sb.WriteByte(' ')
			sb.WriteString(fn)
			sb.WriteString("()")
		}

		return sb.String()
	}

	var skipped func(string) bool
	{
		cList := internal.GetCollapsedLines()

		var lastCollapse *regexp.Regexp

		sList := internal.GetSkippedLine()
		skipping := false

		skipped = func(
			line string,
		) bool {
			if internal.GoSdkCode.MatchString(line) {
				return true
			}

			if skipping {
				skipping = false
				lastCollapse = nil

				return true
			}

			if strings.HasPrefix(line, "panic(") {
				skipping = true

				return true
			}

			for _, re := range sList {
				if re.MatchString(line) {
					skipping = true

					return true
				}
			}

			if lastCollapse != nil && lastCollapse.MatchString(line) {
				return true
			}

			for _, re := range cList {
				if re.MatchString(line) {
					lastCollapse = re

					break
				}
			}

			return false
		}
	}

	var lines []string

	fnName := ""
	for line := range strings.Lines(stacktrace) {
		switch {
		case skipped(line):
			// Do nothing.

		case internal.FnCall.MatchString(line):
			fnName = internal.FnCall.FindStringSubmatch(line)[2]

		default:
			l := fmtLine(line, fnName)
			lines = append(lines, l)
		}
	}

	return slices.DeleteFunc(lines, func(it string) bool {
		return strings.TrimSpace(it) == ""
	})
}

func FmtStacktraces(
	stacktraces []string,
) string {
	if len(stacktraces) == 0 {
		return "<missing trace>\n\n" + string(debug.Stack())
	}

	var traces [][]string
	for _, st := range stacktraces {
		traces = append(traces, FmtStacktrace(st))
	}

	isProper := func(smaller, bigger []string) bool {
		if len(smaller) >= len(bigger) {
			return false
		}

		sp := 0
		for i := 0; i < len(bigger); i++ {
			if bigger[i] == smaller[sp] {
				sp++
			} else {
				sp = 0
			}
			if sp == len(smaller) {
				return true
			}
		}

		return false
	}

	drop := map[int]struct{}{}
	for _, big := range traces {
		for i, small := range traces {
			if isProper(small, big) {
				drop[i] = struct{}{}
			}
		}
	}

	j := 0
	for i := range traces {
		if _, dropped := drop[i]; !dropped && len(traces[i]) > 0 {
			traces[j] = traces[i]
			j++
		}
	}
	traces = traces[:j]

	var sb strings.Builder
	for i, t := range traces {
		if len(traces) > 1 {
			sb.WriteString("trace ")
			sb.WriteString(strconv.Itoa(i))
			sb.WriteString(":\n")
		}
		for _, tl := range t {
			sb.WriteString(tl)
			sb.WriteString("\n")
		}
	}

	fin := sb.String()
	if fin != "" {
		return fin
	}

	return "<missing trace>\n\n" + string(debug.Stack())
}

func FmtStacktraceOf(
	err any,
) string {
	if e, ok := err.(error); ok {
		if tE := g.As[*TracedError](e); tE != nil {
			st := tE.Get()
			str := make([]string, len(st))

			for i, s := range st {
				str[i] = string(s)
			}

			return FmtStacktraces(str)
		}
	}

	return fmt.Sprintf(
		"<missing trace>\n%s\n\n\n\n%s",
		strings.Join(FmtStacktrace(string(debug.Stack())), "\n"),
		string(debug.Stack()),
	)
}

func FmtMsg(
	err any,
) string {
	msg := fmt.Sprint(err)

	if e, ok := err.(error); ok {
		var err *TracedError
		if errors.As(e, &err) {
			msg = strings.Replace(msg, tracedMsg+"\n", "", 1)
		}
	}

	return msg
}

func TNoError(
	t *testing.T,
	err error,
) {
	t.Helper()

	if err == nil {
		return
	}

	content := []struct{ label, content string }{
		{"G11y Trace", FmtStacktraceOf(err)},
		{"Error", fmt.Sprintf("Received unexpected error:\n%+v", err)},
		{"Test", t.Name()},
	}

	longest := len("Error Trace")

	var sb0 strings.Builder
	for _, v := range content {
		sb1 := strings.Builder{}
		for i, scanner := 0, bufio.NewScanner(strings.NewReader(v.content)); scanner.Scan(); i++ {
			if i != 0 {
				sb1.WriteString("\n\t" + strings.Repeat(" ", longest+1) + "\t")
			}
			sb1.WriteString(scanner.Text())
		}

		sb0.WriteByte('\t')
		sb0.WriteString(v.label)
		sb0.WriteByte(':')
		sb0.WriteString(strings.Repeat(" ", longest-len(v.label)))
		sb0.WriteByte('\t')
		sb0.WriteString(sb1.String())
		sb0.WriteByte('\n')
	}

	t.Errorf("\n%s", ""+sb0.String())
	t.FailNow()
}

func TPreamble(
	t *testing.T,
) {
	t.Helper()

	// TODO undo on cleanup?
	EnableDefaultTracer()
}
