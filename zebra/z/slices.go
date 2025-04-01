package z

import (
	"iter"
	"slices"
)

func ItApply[U, V any](
	it iter.Seq[U],
	fn func(U) V,
) iter.Seq[V] {
	return func(yield func(V) bool) {
		for u := range it {
			if !yield(fn(u)) {
				return
			}
		}
	}
}

func ItApplied[U, V any](
	it iter.Seq[U],
	fn func(U) V,
) []V {
	return slices.Collect(ItApply(it, fn))
}

func Apply[Slice ~[]U, U, V any](
	it Slice,
	fn func(U) V,
) iter.Seq[V] {
	return ItApply(slices.Values(it), fn)
}

func Applied[Slice ~[]U, U, V any](
	it Slice,
	fn func(U) V,
) []V {
	return slices.Collect(Apply(it, fn))
}

func GroupBy[Slice ~[]U, U any, K comparable](
	it Slice,
	fn func(U) K,
) map[K][]U {
	mapped := make(map[K][]U, len(it))

	for _, each := range it {
		conv := fn(each)
		if items, ok := mapped[conv]; ok {
			mapped[conv] = append(items, each)
		} else {
			mapped[conv] = []U{each}
		}
	}

	return mapped
}

func GroupByKeyVal[Slice ~[]U, U any, K comparable, V any](
	it Slice,
	fn func(U) (K, V),
) map[K][]V {
	mapped := make(map[K][]V, len(it))

	for _, each := range it {
		k, v := fn(each)
		if items, ok := mapped[k]; ok {
			mapped[k] = append(items, v)
		} else {
			mapped[k] = []V{v}
		}
	}

	return mapped
}

func Appended[S ~[]E, E any](s S, e ...E) S {
	return append(slices.Clone(s), e...)
}
