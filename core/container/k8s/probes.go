package k8s

import (
	"context"
	"maps"
	"net/http"

	"github.com/hkoosha/giraffe/core/t11y/glog"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

const (
	readinessPath = "/readiness"
	livenessPath  = "/liveness"
)

type HealthCheckFn = func(context.Context) error

type probes struct {
	lg         glog.Lg
	checks     map[string]HealthCheckFn
	okMsg      []byte
	failureMsg []byte
}

func (p *probes) say(
	w http.ResponseWriter,
	ok bool,
) {
	status := http.StatusServiceUnavailable
	msg := p.failureMsg
	if ok {
		msg = p.okMsg
		status = http.StatusOK
	}

	w.WriteHeader(status)

	if _, err := w.Write(msg); err != nil {
		p.lg.Error("failed to write probe response", err)
	}
}

func (p *probes) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	var name string
	var err error
	var check HealthCheckFn
	for name, check = range p.checks {
		if err = check(r.Context()); err != nil {
			break
		}
	}

	if err == nil {
		p.say(w, true)
	} else {
		p.lg.Warn("probe failed", N("probe", name), err)
		p.say(w, false)
	}
}

func externalChecker(
	lg glog.Lg,
	checks map[string]HealthCheckFn,
) http.Handler {
	return &probes{
		lg:         lg,
		okMsg:      []byte("OK"),
		failureMsg: []byte(http.StatusText(http.StatusServiceUnavailable)),
		checks:     maps.Clone(checks),
	}
}

func RegisterProbes(
	lg glog.Lg,
	mux *http.ServeMux,
	readiness map[string]HealthCheckFn,
	liveness map[string]HealthCheckFn,
) {
	rR := lg.Named("readiness")
	mux.Handle(readinessPath, externalChecker(rR, readiness))

	rL := lg.Named("liveness")
	mux.Handle(livenessPath, externalChecker(rL, liveness))
}
