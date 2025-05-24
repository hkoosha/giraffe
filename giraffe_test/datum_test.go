package giraffe_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/gtesting"
)

func TestMerge(t *testing.T) {
	t.Run("merge same values", func(t *testing.T) {
		gtesting.Preamble(t)

		defer func() {
			if r := recover(); r != nil {
				gtesting.NoError(t, r)
			}
		}()

		q := giraffe.Q("a.bb.c")

		d0 := giraffe.Of1(q, 123)
		d1 := giraffe.Of1(q, 123)
		dm, err := d0.Merge(d1)

		gtesting.NoError(t, err)
		assert.Equal(t, d0.Pretty(), dm.Pretty())
	})
}
