package g11y

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"strings"

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
) string {
	defer func() {
		if r := recover(); r != nil {
			_, _ = os.Stderr.WriteString("nested panic " + fmt.Sprintf("%+v", r))
			os.Exit(13)
		}
	}()

	fmtLine := func(
		sb *strings.Builder,
		line string,
		fn string,
	) {
		if internal.GoSdkCode.MatchString(line) {
			return
		}

		line = strings.TrimSpace(line)

		if internal.FileLine.MatchString(line) {
			line = strings.Split(line, " ")[0]
		}

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

		sb.WriteByte('\n')
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

	sb := strings.Builder{}
	fnName := ""

	for line := range strings.Lines(stacktrace) {
		switch {
		case skipped(line):
			// Do nothing.

		case internal.FnCall.MatchString(line):
			fnName = internal.FnCall.FindStringSubmatch(line)[2]

		default:
			fmtLine(&sb, line, fnName)
		}
	}

	return sb.String()
}

func FmtStacktraces(
	stacktraces []string,
) string {
	if len(stacktraces) == 0 {
		return "<missing trace>\n\n" + string(debug.Stack())
	}

	var sb strings.Builder

	for i, st := range stacktraces {
		if t := FmtStacktrace(st); strings.TrimSpace(t) != "" {
			sb.WriteString(fmt.Sprintf("trace %d:\n%s\n", i, t))
		}
	}

	if fin := sb.String(); fin != "" {
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

	return "<missing trace>\n" +
		FmtStacktrace(string(debug.Stack())) + "\n\n\n\n" + string(debug.Stack())
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
