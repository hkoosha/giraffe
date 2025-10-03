package queryimpl

import (
	"github.com/hkoosha/giraffe/dialects"
)

type Segment struct {
	dialect dialects.Dialect
	ref     string
}
