package gquery

import (
	"math"
	"strings"
)

// TODO: no 'overwrite' and 'maybe' at the same time.

// TODO: better conflicting cmd checks.

var (
	ErrCodeQueryParseEmptyQuery         uint64 = math.MaxUint64
	ErrCodeQueryParseDuplicatedCmd      uint64 = math.MaxUint64
	ErrCodeQueryParseConflictingCmd     uint64 = math.MaxUint64
	ErrCodeQueryParseUnexpectedToken    uint64 = math.MaxUint64
	ErrCodeQueryParseUnexpectedSegments uint64 = math.MaxUint64
	ErrCodeQueryParseNestingTooDeep     uint64 = math.MaxUint64
	ErrCodeQueryParseNotWritable        uint64 = math.MaxUint64
)

func Escaped(
	spec string,
) string {
	sb := strings.Builder{}
	sb.Grow(len(spec))

	for _, c := range spec {
		if _, ok := commands[c]; ok {
			sb.WriteRune(CmdEscape)
		}

		sb.WriteRune(c)
	}

	return spec
}

func Parse(
	spec string,
) (Query, error) {
	return parse(spec)
}
