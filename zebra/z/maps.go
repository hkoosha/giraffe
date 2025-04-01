package z

import (
	"iter"
	"maps"
	"slices"
)

func ItApply2[K comparable, V, U any](
	it iter.Seq2[K, V],
	fn func(K, V) U,
) iter.Seq2[K, U] {
	return func(yield func(K, U) bool) {
		for k, v := range it {
			if !yield(k, fn(k, v)) {
				return
			}
		}
	}
}

func ItApplied2[K comparable, V, U any](
	it iter.Seq2[K, V],
	fn func(K, V) U,
) map[K]U {
	return maps.Collect(ItApply2(it, fn))
}

func Apply2[Map ~map[K]V, K comparable, V, U any](
	it Map,
	fn func(K, V) U,
) iter.Seq2[K, U] {
	return ItApply2(maps.All(it), fn)
}

func Applied2[Map ~map[K]V, K comparable, V, U any](
	it Map,
	fn func(K, V) U,
) map[K]U {
	return maps.Collect(Apply2(it, fn))
}

func Applied2V[Map ~map[K]V, K comparable, V, U any](
	it Map,
	fn func(V) U,
) map[K]U {
	fnV := func(_ K, v V) U {
		return fn(v)
	}

	return maps.Collect(Apply2(it, fnV))
}

func Applies2Ref[Map ~map[K]V, K comparable, V any](
	it Map,
) map[K]*V {
	m := make(map[K]*V, len(it))
	for k, v := range it {
		m[k] = &v
	}

	return m
}

func Values[Map ~map[K]V, K comparable, V any](
	it Map,
) []V {
	return slices.Collect(maps.Values(it))
}

func MapEq[K, V comparable](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k1, v1 := range m1 {
		if v2, ok := m2[k1]; !ok || v1 != v2 {
			return false
		}
	}

	return true
}

func UnionLeft[K, V comparable](m1, m2 map[K]V) map[K]V {
	union := maps.Clone(m1)

	for k, v := range m2 {
		if _, ok := union[k]; !ok {
			union[k] = v
		}
	}

	return union
}
