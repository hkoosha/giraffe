package gtesting

import (
	"testing"

	"github.com/hkoosha/giraffe/t11y"
)

func Preamble(t *testing.T) {
	t.Helper()

	{
		v := t11y.IsUnsafeError()
		t11y.EnableUnsafeError()
		t.Cleanup(func() {
			if v {
				t11y.EnableUnsafeError()
			} else {
				t11y.DisableUnsafeError()
			}
		})
	}

	{
		v := t11y.IsTracer()
		t11y.EnableTracer()
		t.Cleanup(func() {
			if v {
				t11y.EnableTracer()
			} else {
				t11y.DisableTracer()
			}
		})
	}

	{
		v := t11y.GetSkippedLines()
		t11y.SetSkippedLines(true)
		t.Cleanup(func() {
			t11y.SetSkippedLines(false, v...)
		})
	}

	{
		v := t11y.GetCollapsedLines()
		t11y.SetCollapsedLines(true)
		t.Cleanup(func() {
			t11y.SetCollapsedLines(false, v...)
		})
	}
}
