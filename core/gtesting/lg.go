package gtesting

import (
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/hkoosha/giraffe/contrib/zap/gzapadapter"
	"github.com/hkoosha/giraffe/core/t11y/glog"
)

func Zap(t *testing.T) glog.Lg {
	t.Helper()

	return gzapadapter.Of(zaptest.NewLogger(t))
}
