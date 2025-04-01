package giraffe_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/g11y"
)

func setup() {
	g11y.EnableDefaultTracer()
}

func TestMerge(t *testing.T) {
	t.Run("merge same values", func(t *testing.T) {
		setup()
		defer func() {
			r := recover()
			require.Nil(
				t, r, "merge failed: %s", g11y.FmtStacktraceOf(r))
		}()

		d0 := giraffe.Of1("a.bb.c", 123)
		d1 := giraffe.Of1("a.bb.c", 123)

		dMerge, err := d0.Merge(d1)
		require.NoError(
			t, err, "merge failed: %s", g11y.FmtStacktraceOf(err))

		assert.Equal(t, d0.Pretty(), dMerge.Pretty())
	})
}
