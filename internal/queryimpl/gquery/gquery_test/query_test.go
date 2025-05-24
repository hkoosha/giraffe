package gquery_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe/core/gtesting"
	"github.com/hkoosha/giraffe/internal/queryimpl/gquery"
)

func TestParse(t *testing.T) {
	t.Run("simple_1", func(t *testing.T) {
		gtesting.Preamble(t)

		spec := "k0"

		first, err := gquery.Parse(spec)
		gtesting.NoError(t, err)

		path := first.VisibleForTestingPath()
		require.Len(t, path, 1)

		actual := path[0]

		assert.Equal(t, "k0", actual.String())
		assert.Equal(t, "k0", actual.Attr())
		assert.Equal(t, 0, actual.Flags().Seq())
		assert.Equal(t, actual, actual.Root())
		assert.Equal(t, actual, actual.Leaf())
		assert.True(t, actual.Flags().IsSingle())
		assert.True(t, actual.Flags().IsLeaf())
		assert.True(t, actual.Flags().IsRoot())

		assert.Panics(t, func() { actual.Prev() })
		assert.Panics(t, func() { actual.Next() })
	})

	t.Run("simple_2", func(t *testing.T) {
		gtesting.Preamble(t)

		spec := "k0.k1"

		first, err := gquery.Parse(spec)
		gtesting.NoError(t, err)

		path := first.VisibleForTestingPath()
		require.Len(t, path, 2)

		k0 := path[0]
		k1 := path[1]

		assert.Equal(t, "k0.k1", k0.String())
		assert.Equal(t, "k0", k0.Attr())
		assert.Equal(t, 0, k0.Flags().Seq())
		assert.Equal(t, k0, k0.Root())
		assert.Equal(t, k1, k0.Leaf())
		assert.False(t, k0.Flags().IsSingle())
		assert.False(t, k0.Flags().IsLeaf())
		assert.True(t, k0.Flags().IsRoot())
		assert.Panics(t, func() { k0.Prev() })
		require.NotNil(t, k0.Next())
		assert.Equal(t, k1, k0.Next())

		assert.Equal(t, "k0.@k1", k1.String())
		assert.Equal(t, "k1", k1.Attr())
		assert.Equal(t, 1, k1.Flags().Seq())
		assert.Equal(t, k0, k1.Root())
		assert.Equal(t, k1, k1.Leaf())
		assert.False(t, k1.Flags().IsSingle())
		assert.True(t, k1.Flags().IsLeaf())
		assert.False(t, k1.Flags().IsRoot())
		assert.Panics(t, func() { k1.Next() })
		require.NotNil(t, k1.Prev())

		p := k1.Prev()
		assert.Equal(t, k0, p)
	})
}

func TestNext(t *testing.T) {
	t.Run("test next", func(t *testing.T) {
		gtesting.Preamble(t)

		spec := "dynamic.static.thingy.foo"

		q, err := gquery.Parse(spec)
		gtesting.NoError(t, err)

		assert.Equal(t, "dynamic", q.Attr())
		assert.Equal(t, "static", q.Next().Attr())
		assert.Equal(t, "thingy", q.Next().Next().Attr())
		assert.Equal(t, "foo", q.Next().Next().Next().Attr())

		require.Panics(t, func() { q.Next().Next().Next().Next() })
	})
}

func TestQuery_ToString(t *testing.T) {
	t.Run("to string", func(t *testing.T) {
		gtesting.Preamble(t)

		q, err := gquery.Parse("k0.k1.k2")
		gtesting.NoError(t, err)

		str := q.Next().String()
		require.Equal(t, "k0.@k1.k2", str)
	})
}
