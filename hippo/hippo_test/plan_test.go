package hippo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hkoosha/giraffe/hippo"
)

func TestPlan_Empty(t *testing.T) {
	t.Run("names of empty plan", func(t *testing.T) {
		p := hippo.Plan{}

		assert.Empty(t, p.Names())
		assert.NotNil(t, p.Names())
	})
}
