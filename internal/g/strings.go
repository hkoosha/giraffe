package g

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"

	"github.com/hkoosha/giraffe/zebra/z"
)

const sep = ", "

func Joined(
	str []string,
) string {
	return strings.Join(slices.DeleteFunc(z.Applied(str, strings.TrimSpace), IsEmpty), sep)
}

func JoinedFn[T any](
	values []T,
	toString func(T) string,
) string {
	return Joined(z.Applied(values, toString))
}

func Joiner[T fmt.Stringer](
	values []T,
) string {
	return Joined(z.Applied(values, ToString))
}

func Join(
	str ...string,
) string {
	return Joined(str)
}

func JoinIt[V any](
	it iter.Seq[V],
) string {
	return strings.Join(
		z.ItApplied(it, func(v V) string {
			return fmt.Sprint(v)
		}),
		sep,
	)
}

func JoinKeys[K comparable, V any](
	m map[K]V,
) string {
	return JoinIt(maps.Keys(m))
}

func IsEmpty(
	str string,
) bool {
	return str == "" || strings.TrimSpace(str) == ""
}

func ToString[T fmt.Stringer](s T) string {
	return s.String()
}
