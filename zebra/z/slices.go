package z

import (
	"iter"
	"slices"
)

func ItFlatten[Slice ~[]U, U any](
	it iter.Seq[Slice],
) iter.Seq[U] {
	return func(yield func(U) bool) {
		for vs := range it {
			for _, v := range vs {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func ItFlattened[Slice ~[]U, U any](
	it iter.Seq[Slice],
) Slice {
	return slices.Collect(ItFlatten(it))
}

func Flatten[Slice ~[]U, U any](
	it []Slice,
) iter.Seq[U] {
	return ItFlatten(slices.Values(it))
}

func Flattened[Slice ~[]U, U any](
	it []Slice,
) Slice {
	return slices.Collect(Flatten(it))
}

// =============================================================================.

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

// =============================================================================.

func ItFilter[V any](
	it iter.Seq[V],
	fn func(V) bool,
) iter.Seq[V] {
	return func(yield func(V) bool) {
		for u := range it {
			if fn(u) {
				if !yield(u) {
					return
				}
			}
		}
	}
}

func ItFiltered[V any](
	it iter.Seq[V],
	fn func(V) bool,
) []V {
	return slices.Collect(ItFilter(it, fn))
}

func Filter[Slice ~[]V, V any](
	it Slice,
	fn func(V) bool,
) iter.Seq[V] {
	return ItFilter(slices.Values(it), fn)
}

func Filtered[Slice ~[]V, V any](
	it Slice,
	fn func(V) bool,
) []V {
	return slices.Collect(Filter(it, fn))
}

// =============================================================================.

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

func GroupByIf[Slice ~[]U, U any, K comparable](
	it Slice,
	fn func(U) (K, bool),
) map[K][]U {
	mapped := make(map[K][]U, len(it))

	for _, each := range it {
		if conv, ok0 := fn(each); ok0 {
			if items, ok1 := mapped[conv]; ok1 {
				mapped[conv] = append(items, each)
			} else {
				mapped[conv] = []U{each}
			}
		}
	}

	return mapped
}

func GroupByKey[Slice ~[]U, U comparable, V any](
	it Slice,
	fn func(U) V,
) map[U][]V {
	mapped := make(map[U][]V, len(it))

	for _, k := range it {
		v := fn(k)
		if items, ok := mapped[k]; ok {
			mapped[k] = append(items, v)
		} else {
			mapped[k] = []V{v}
		}
	}

	return mapped
}

func GroupByWith[Slice ~[]U, U any, K comparable, V any](
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

func GroupByWithIf[Slice ~[]U, U any, K comparable, V any](
	it Slice,
	fn func(U) (K, V, bool),
) map[K][]V {
	mapped := make(map[K][]V, len(it))

	for _, each := range it {
		if k, v, ok0 := fn(each); ok0 {
			if items, ok1 := mapped[k]; ok1 {
				mapped[k] = append(items, v)
			} else {
				mapped[k] = []V{v}
			}
		}
	}

	return mapped
}

func Appended[S ~[]E, E any](s S, e ...E) S {
	return append(slices.Clone(s), e...)
}
