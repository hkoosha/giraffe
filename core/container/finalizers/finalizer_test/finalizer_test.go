package finalizer_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe/core/container/finalizers"
	"github.com/hkoosha/giraffe/core/container/setup"
)

func run(
	ctx context.Context,
	lg func(string),
	ch chan<- []string,
) {
	var tl []string

	lg("creating")
	f := finalizers.NewFinalizer(setup.New())

	lg("add: f0")
	f.Add00(func() {
		lg("called: f0")
		tl = append(tl, "f0")
	})

	lg("add: f1")
	f.Add00(func() {
		lg("called: f1")
		tl = append(tl, "f1")
	})

	lg("add: f2")
	f.Add00(func() {
		lg("called: f2")
		tl = append(tl, "f2")
	})

	lg("finalize")
	f.Finalize(ctx)

	lg("fin")
	ch <- tl
}

func TestFinalizerRegistry_Execute(t *testing.T) {
	t.Parallel()

	timeout := 2 * time.Second
	lg := func(s string) { t.Log(s) }

	t.Run("runs finalizers", func(t *testing.T) {
		t.Parallel()

		timeline := make(chan []string, 1)
		go run(t.Context(), lg, timeline)

		select {
		case fin := <-timeline:
			require.Equal(t, []string{"f2", "f1", "f0"}, fin)

		case <-time.After(timeout):
			t.Error("timed out waiting for finalizers to run")
		}
	})
}
