package gtesting

import (
	"testing"

	"github.com/hkoosha/giraffe/contrib/zap/gzapadapter"
	"github.com/hkoosha/giraffe/core/t11y/glog"
	"go.uber.org/zap/zaptest"
)

func Zap(t *testing.T) glog.Lg {
	t.Helper()

	return gzapadapter.Of(zaptest.NewLogger(t))
}
