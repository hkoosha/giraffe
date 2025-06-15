package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnce(t *testing.T) {
	t.Run("blocks second - level 1 nesting", func(t *testing.T) {
		Once("datagen", "setup_once_test", "level1")

		assert.Panics(t, func() {
			Once("datagen", "setup_once_test", "level1")
		})
	})

	t.Run("blocks second - level 2 nesting", func(t *testing.T) {
		Once("datagen", "setup_once_test", "level1", "level2")

		assert.Panics(t, func() {
			Once("datagen", "setup_once_test", "level1", "level2")
		})
	})

	t.Run("blocks second - level 3 nesting", func(t *testing.T) {
		Once("datagen", "setup_once_test", "level1", "level2", "level3")

		assert.Panics(t, func() {
			Once("datagen", "setup_once_test", "level1", "level2", "level3")
		})
	})

	t.Run("blocks second - level 4 nesting", func(t *testing.T) {
		Once("datagen", "setup_once_test", "level1", "level2", "level3", "level4")

		assert.Panics(t, func() {
			Once("datagen", "setup_once_test", "level1", "level2", "level3", "level4")
		})
	})

	t.Run("blocks second - level 5 nesting", func(t *testing.T) {
		Once("datagen", "setup_once_test", "level1", "level2", "level3", "level4", "level5")

		assert.Panics(t, func() {
			Once("datagen", "setup_once_test", "level1", "level2", "level3", "level4", "level5")
		})
	})
}
