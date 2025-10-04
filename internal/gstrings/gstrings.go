package gstrings

import (
	"fmt"
	"iter"
	"slices"
	"strings"
)

const sep = ", "

func Joined(
	str []string,
) string {
	for i := range str {
		str[i] = strings.TrimSpace(str[i])
	}

	return strings.Join(slices.DeleteFunc(
		str,
		func(it string) bool { return it == "" },
	), sep)
}

func JoinIt[V any](
	it iter.Seq[V],
) string {
	collect := slices.Collect(func(yield func(string) bool) {
		for v := range it {
			if !yield(fmt.Sprint(v)) {
				return
			}
		}
	})

	return strings.Join(collect, sep)
}
