package zio

import (
	"bytes"
	"io"
)

func Tee(r io.ReadCloser) ([]byte, io.ReadCloser, error) {
	if r == nil {
		return nil, nil, nil
	}

	defer func(r io.ReadCloser) {
		_ = r.Close()
	}(r)

	read, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}

	return read, io.NopCloser(bytes.NewReader(read)), nil
}

func ReadAll(r io.ReadCloser) ([]byte, error) {
	if r == nil {
		return nil, nil
	}

	defer func(r io.ReadCloser) {
		_ = r.Close()
	}(r)

	read, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return read, nil
}
