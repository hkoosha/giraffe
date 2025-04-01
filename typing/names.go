package typing

import (
	"math"
	"regexp"
	"strings"
)

// Contains only lowercase letters, numbers, and underscores; starts with
// character, and ends with a lowercase letter or number (not underscore).
var machineReadableNamePattern = regexp.MustCompile("^[a-z][a-z0-9_]*[^_]$")

func IsMachineReadableName(
	name string,
	minLenInclusive uint,
	maxLenInclusive uint,
) bool {
	if maxLenInclusive > math.MaxInt ||
		minLenInclusive > math.MaxInt ||
		maxLenInclusive < minLenInclusive {
		panic("invalid args for min/max length")
	}

	l := uint(len(name))

	return minLenInclusive <= l && l <= maxLenInclusive &&
		machineReadableNamePattern.MatchString(name) &&
		!strings.Contains(name, "__")
}

func IsSimpleMachineReadableName(
	name string,
	minLenInclusive uint,
	maxLenInclusive uint,
) bool {
	return IsMachineReadableName(name, minLenInclusive, maxLenInclusive) &&
		!strings.Contains(name, "_")
}
