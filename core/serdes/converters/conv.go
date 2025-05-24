package converters

type Conv[T any, U any] interface {
	Write(T) (U, error)

	Read(U) (T, error)
}

// =============================================================================

func Bytes() Conv[[]byte, []byte] {
	return bytesConv{}
}

func String() Conv[string, string] {
	return stringConv{}
}

func Json[T any]() Conv[T, []byte] {
	return jsonConv[T]{}
}

func JsonStr[T any]() Conv[T, string] {
	return jsonStr[T]{}
}
