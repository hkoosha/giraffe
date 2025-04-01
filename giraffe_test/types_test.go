package giraffe_test

import (
	"math/bits"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hkoosha/giraffe"
)

var dTypes = []giraffe.Type{
	giraffe.Obj,
	giraffe.Int,
	giraffe.Flt,
	giraffe.Bln,
	giraffe.Str,
}

var dMods = []giraffe.Type{
	giraffe.Nil,
	giraffe.Arr,
}

func Test_Bits(t *testing.T) {
	all := append(slices.Clone(dTypes), dMods...)
	slices.Sort(all)

	t.Run("unique values", func(t *testing.T) {
		seen := make(map[uint64]any)
		for _, v := range all {
			_, dup := seen[uint64(v)]
			assert.Falsef(t, dup, "duplicated: %s", v)

			seen[uint64(v)] = nil
		}
	})

	t.Run("bit count", func(t *testing.T) {
		for _, v := range all {
			assert.Equalf(t, 1, bits.OnesCount64(uint64(v)), "not single-bit: %s", v)
		}
	})
}

func Test_Range(t *testing.T) {
	t.Run("modifier range", func(t *testing.T) {
		for _, v := range dMods {
			assert.GreaterOrEqual(t, uint64(v), uint64(1))
			assert.LessOrEqual(t, uint64(v), uint64(0b1111_1111))
		}
	})

	t.Run("type range", func(t *testing.T) {
		for _, v := range dTypes {
			assert.Greater(t, uint64(v), uint64(0b1111_1111))
		}
	})
}

func Test_FixedValues(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Equal(t, uint64(1), uint64(giraffe.Nil))
	})
}
