package setup_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hkoosha/giraffe/core/container/setup"
)

func TestOnce(t *testing.T) {
	t.Run("blocks second - level 1 nesting", func(t *testing.T) {
		setup.Finish(
			"giraffe_once_test1",
			"level1",
		)

		assert.Panics(t, func() {
			setup.Finish(
				"giraffe_once_test1",
				"level1",
			)
		})
	})

	t.Run("blocks second - level 2 nesting", func(t *testing.T) {
		setup.Finish(
			"giraffe_once_test2",
			"level1",
			"level2",
		)

		assert.Panics(t, func() {
			setup.Finish(
				"giraffe_once_test2",
				"level1",
				"level2",
			)
		})
	})

	t.Run("blocks second - level 3 nesting", func(t *testing.T) {
		setup.Finish(
			"giraffe_once_test3",
			"level1",
			"level2",
			"level3",
		)

		assert.Panics(t, func() {
			setup.Finish(
				"giraffe_once_test3",
				"level1",
				"level2",
				"level3",
			)
		})
	})

	t.Run("waterfall", func(t *testing.T) {
		setup.Finish(
			"giraffe_once_test_waterfall",
			"level1",
			"level2",
			"level3",
		)
		setup.EnsureDone(
			"giraffe_once_test_waterfall",
			"level1",
			"level2",
			"level3",
		)

		setup.Finish(
			"giraffe_once_test_waterfall",
			"level1",
			"level2",
		)
		setup.EnsureDone(
			"giraffe_once_test_waterfall",
			"level1",
			"level2",
		)

		setup.Finish(
			"giraffe_once_test_waterfall",
			"level1",
		)
		setup.EnsureDone(
			"giraffe_once_test_waterfall",
			"level1",
		)

		setup.Finish(
			"giraffe_once_test_waterfall",
		)
		setup.EnsureDone(
			"giraffe_once_test_waterfall",
		)
	})
}
