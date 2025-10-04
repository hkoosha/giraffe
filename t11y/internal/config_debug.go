package internal

import (
	"regexp"
	"sync/atomic"
)

var (
	skippedLines   = atomic.Pointer[[]*regexp.Regexp]{}
	collapsedLines = atomic.Pointer[[]*regexp.Regexp]{}
)

func init() {
	skippedLines.Store(&DefaultSkippedLines)
	collapsedLines.Store(&DefaultCollapsedLines)
}

func SetSkippedLine(
	re []*regexp.Regexp,
) {
	skippedLines.Store(&re)
}

func GetSkippedLine() []*regexp.Regexp {
	load := skippedLines.Load()
	if load == nil {
		return []*regexp.Regexp{}
	}

	return *load
}

func SetCollapsedLines(
	re []*regexp.Regexp,
) {
	collapsedLines.Store(&re)
}

func GetCollapsedLines() []*regexp.Regexp {
	load := collapsedLines.Load()
	if load == nil {
		return []*regexp.Regexp{}
	}

	return *load
}
