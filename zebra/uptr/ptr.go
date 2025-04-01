package uptr

func To[T any](t T) *T {
	return &t
}

func From[V any](v *V) V {
	if v == nil {
		var res V

		return res
	}

	return *v
}

func FromOr[V any](v *V, or V) V {
	if v == nil {
		return or
	}

	return *v
}

func ToSlice[V any](vs []V) []*V {
	result := make([]*V, len(vs))
	for i, v := range vs {
		result[i] = &v
	}

	return result
}

func FromSlice[V any](vs []*V) []V {
	result := make([]V, len(vs))
	for i, v := range vs {
		result[i] = *v
	}

	return result
}

func MustCopy[T any](v *T) *T {
	if v == nil {
		return nil
	}

	return &(*v)
}
