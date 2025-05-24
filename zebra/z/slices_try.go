package z

import (
	"errors"
	"iter"
)

var errDuplicateKey = errors.New("duplicate key")

func TryItApplied[U, V any](
	it iter.Seq[U],
	fn func(U) (V, error),
) ([]V, error) {
	next, stop := iter.Pull(it)
	defer stop()

	transformed := make([]V, 0)
	for u, ok := next(); ok; u, ok = next() {
		t, err := fn(u)
		if err != nil {
			return transformed, err
		}

		transformed = append(transformed, t)
	}

	return transformed, nil
}

func TryApplied[Slice ~[]U, U, V any](
	it Slice,
	fn func(U) (V, error),
) ([]V, error) {
	transformed := make([]V, len(it))
	for i, u := range it {
		t, err := fn(u)
		if err != nil {
			if i > 0 {
				transformed = transformed[:i]
			} else {
				transformed = nil
			}

			return transformed, err
		}

		transformed[i] = t
	}

	return transformed, nil
}

func TryKeyBy[Slice ~[]U, U any, K comparable](
	it Slice,
	fn func(U) K,
) (map[K]U, error) {
	mapped := make(map[K]U, len(it))

	for _, each := range it {
		conv := fn(each)
		if _, ok := mapped[conv]; ok {
			return nil, errDuplicateKey
		}
		mapped[conv] = each
	}

	return mapped, nil
}

func TryKeyValBy[Slice ~[]U, U any, K comparable, V any](
	it Slice,
	kFn func(U) K,
	vFn func(U) V,
) (map[K]V, error) {
	mapped := make(map[K]V, len(it))

	for _, each := range it {
		k := kFn(each)
		v := vFn(each)
		if _, ok := mapped[k]; ok {
			return nil, errDuplicateKey
		}
		mapped[k] = v
	}

	return mapped, nil
}

func TryValBy[Slice ~[]U, U comparable, V any](
	it Slice,
	fn func(U) V,
) (map[U]V, error) {
	mapped := make(map[U]V, len(it))

	for _, each := range it {
		if _, ok := mapped[each]; !ok {
			mapped[each] = fn(each)
		} else {
			return nil, errDuplicateKey
		}
	}

	return mapped, nil
}

func TryMapBy[Slice ~[]U, U comparable, V any](
	it Slice,
	fn func(U) V,
) (map[U]V, error) {
	mapped := make(map[U]V, len(it))

	for _, k := range it {
		v := fn(k)
		if _, ok := mapped[k]; ok {
			return nil, errDuplicateKey
		} else {
			mapped[k] = v
		}
	}

	return mapped, nil
}
