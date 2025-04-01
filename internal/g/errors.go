package g

import (
	"errors"
)

func As[E error](err error) E {
	var cast E

	errors.As(err, &cast)

	return cast
}

func Is[E error](err error) bool {
	var cast E

	return errors.As(err, &cast)
}
