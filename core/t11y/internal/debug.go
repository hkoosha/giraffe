package internal

import (
	"regexp"
	"slices"
)

var (
	FnCall    = regexp.MustCompile(`^([^\t ]+/)(.*)(\(.*\))`)
	FileLine  = regexp.MustCompile(`^.*\.go:\d+ \+0x[0-9a-fA-F]*$`)
	LineNum   = regexp.MustCompile(`^.*\.go:\d+$`)
	GoSdkCode = regexp.MustCompile(
		`^.*/go/\d{1,3}(\.\d{1,3})*/src(/[a-zA-Z0-9_-]+)*/[a-zA-Z0-9_-]+\.go:\d+\s*$`,
	)

	DefaultSkippedLines = []*regexp.Regexp{
		// Go funcs.
		regexp.MustCompile(`go/pkg/mod/golang\.org/x/sync.*/errgroup/.*`),
		regexp.MustCompile(regexp.QuoteMeta("runtime/debug.Stack(")),
		regexp.MustCompile(`^panic\(`),
		regexp.MustCompile(regexp.QuoteMeta("golang.org/toolchain")),
		regexp.MustCompile(regexp.QuoteMeta("/giraffe/internal/queryerrors")),
		regexp.MustCompile(regexp.QuoteMeta("giraffe/errors.go")),

		// Giraffe packages.
		regexp.MustCompile(regexp.QuoteMeta("/giraffe/dot")),
		regexp.MustCompile("/giraffe/.*/dot"),
		regexp.MustCompile("/giraffe/t11y/.*"),
		regexp.MustCompile("/giraffe/.*/t11y.*"),

		// Hippo packages.
		regexp.MustCompile("giraffe/hippo/internal/hippoerr/.*"),

		// Third-party.
		regexp.MustCompile(regexp.QuoteMeta("go.uber.org/zap")),
	}

	DefaultCollapsedLines = []*regexp.Regexp{
		regexp.MustCompile(`errors?\.go`),
	}
)

func DeepCopyL1[S ~[]*E, E any](s S) S {
	if s == nil {
		return nil
	}

	cp := make(S, len(s))

	for i, v := range s {
		if v == nil {
			cp[i] = nil
		} else {
			vv := *v
			cp[i] = &vv
		}
	}

	return slices.Clip(cp)
}
