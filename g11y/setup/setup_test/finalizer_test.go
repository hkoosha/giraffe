package setup_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hkoosha/giraffe/g11y/setup"
)

func TestFinalizerRegistry_Execute(t *testing.T) {
	t.Run("runs finalizers", func(t *testing.T) {
		f := setup.NewFinalizerRegistry("giraffe_test")

		var timeline []string

		f.Add(func(context.Context) {
			timeline = append(timeline, "f0")
		})

		f.Add(func(context.Context) {
			timeline = append(timeline, "f1")
		})

		f.Add(func(context.Context) {
			timeline = append(timeline, "f2")
		})

		f.Execute(t.Context())

		assert.Equal(t, []string{"f0", "f1", "f2"}, timeline)
	})
}
