package setup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinalizerRegistry_Execute(t *testing.T) {
	t.Run("runs finalizers", func(t *testing.T) {
		f := NewFinalizerRegistry("datagen_test")

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
