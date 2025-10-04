package giraffe_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/t11y"
)

func TestMerge(t *testing.T) {
	t.Run("merge same values", func(t *testing.T) {
		t11y.TPreamble(t)

		defer func() {
			if r := recover(); r != nil {
				require.Nil(t, r, t11y.FmtStacktraceOf(r))
			}
		}()

		q := giraffe.Q("a.bb.c")

		d0 := giraffe.Of1(q, 123)
		d1 := giraffe.Of1(q, 123)
		dm, err := d0.Merge(d1)

		t11y.TNoError(t, err)
		assert.Equal(t, d0.Pretty(), dm.Pretty())
	})
}
