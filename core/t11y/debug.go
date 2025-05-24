package t11y

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"

	"github.com/hkoosha/giraffe/core/t11y/internal"
)

func GetSkippedLines() []*regexp.Regexp {
	r := internal.GetSkippedLine()
	return internal.DeepCopyL1(r)
}

func SetSkippedLines(
	lines ...*regexp.Regexp,
) {
	SetSkippedLines0(true, lines...)
}

func SetSkippedLines0(
	withDefaults bool,
	lines ...*regexp.Regexp,
) {
	s := internal.DeepCopyL1(lines)
	if withDefaults {
		s = append(s, internal.DefaultSkippedLines...)
	}

	internal.SetSkippedLine(s)
}

func GetCollapsedLines() []*regexp.Regexp {
	r := internal.GetCollapsedLines()
	return internal.DeepCopyL1(r)
}

func SetCollapsedLines(
	lines ...*regexp.Regexp,
) {
	SetCollapsedLines0(true, lines...)
}

func SetCollapsedLines0(
	withDefaults bool,
	lines ...*regexp.Regexp,
) {
	s := internal.DeepCopyL1(lines)
	if withDefaults {
		s = append(s, internal.DefaultCollapsedLines...)
	}

	internal.SetCollapsedLines(s)
}

// ==============================================================================.

func formatMsg(
	err any,
) string {
	msg := fmt.Sprint(err)

	if e, ok := err.(error); ok {
		var err *tracedError
		if errors.As(e, &err) {
			msg = strings.Replace(msg, tracedMsg+"\n", "", 1)
		}
	}

	return msg
}

func fmtStacktrace(
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

func fmtStacktraces(
	stacktraces []string,
) string {
	if len(stacktraces) == 0 {
		return "<missing trace>\n\n" + string(debug.Stack())
	}

	traces := make([][]string, len(stacktraces))
	for i, st := range stacktraces {
		traces[i] = fmtStacktrace(st)
	}

	isProper := func(smaller, bigger []string) bool {
		if len(smaller) >= len(bigger) {
			return false
		}

		sp := 0
		for i := range bigger {
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
		var tE *tracedError
		if errors.As(e, &tE) {
			st := tE.Get()
			str := make([]string, len(st))

			for i, s := range st {
				str[i] = string(s)
			}

			return fmtStacktraces(str)
		}
	}

	return fmt.Sprintf(
		"<missing trace>\n%s\n\n\n\n%s",
		strings.Join(fmtStacktrace(string(debug.Stack())), "\n"),
		string(debug.Stack()),
	)
}
