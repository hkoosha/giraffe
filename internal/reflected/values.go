package reflected

func Zero[T any]() T {
	var t T
	return t
}
