package z

import (
	"iter"
)

func TryItApplied2[K comparable, V, U any](
	it iter.Seq2[K, V],
	fn func(K, V) (U, error),
) (map[K]U, error) {
	next, stop := iter.Pull2(it)
	defer stop()

	transformed := make(map[K]U)
	for k, v, ok := next(); ok; k, v, ok = next() {
		t, err := fn(k, v)
		if err != nil {
			//nolint:nilnil
			return transformed, err
		}

		transformed[k] = t
	}

	return transformed, nil
}

func TryApplied2[Map ~map[K]V, K comparable, V, U any](
	it Map,
	fn func(K, V) (U, error),
) (map[K]U, error) {
	transformed := make(map[K]U, len(it))
	for k, v := range it {
		t, err := fn(k, v)
		if err != nil {
			//nolint:nilnil
			return transformed, err
		}

		transformed[k] = t
	}

	return transformed, nil
}
