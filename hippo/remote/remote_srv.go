package remote

import (
	"io"
	"maps"
	"net/http"
	"regexp"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/serdes"
	"github.com/hkoosha/giraffe/core/t11y"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	"github.com/hkoosha/giraffe/hippo"
	"github.com/hkoosha/giraffe/hippo/internal"
)

type Server func(
	gtx.Context,
	io.Reader,
	io.Writer,
) error

func (s Server) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := gtx.Of(r.Context())

	if err := s(ctx, r.Body, w); err != nil {
		// TODO make error message more deterministic and safer

		msg := err.Error()
		if t11y.IsUnsafeError() {
			msg += "\n\n" + t11y.FmtStacktraceOf(err)
		}
		http.Error(w, msg, http.StatusBadRequest)
	}
}

type server struct {
	reg       *hippo.FnRegistry
	templates map[string]*hippo.Plan
	serde     serdes.Serde[Request]
}

func (s server) ekran(
	ctx gtx.Context,
	r io.Reader,
	w io.Writer,
) error {
	req, err := s.serde.StreamFrom(r)
	if err != nil {
		return newErrorParsingPayload(err)
	}

	init, err := giraffe.From(req.Init)
	if err != nil {
		return newErrorParsingPayload(err)
	}

	var compensator hippo.Compensator
	if req.Compensations == nil {
		req.Compensations = &[]requestCompensations{}
	}
	for _, comp := range *req.Compensations {
		var msgRe *regexp.Regexp
		if comp.OnErrRe != nil {
			msgRe, err = regexp.Compile(*comp.OnErrRe)
			if err != nil {
				return newErrorParsingPayload(err)
			}
		}

		var nameRe *regexp.Regexp
		if comp.OnNameRe != nil {
			nameRe, err = regexp.Compile(*comp.OnNameRe)
			if err != nil {
				return newErrorParsingPayload(err)
			}
		}

		step := -1
		if comp.OnStep != nil {
			step = *comp.OnStep
		}

		var with *hippo.Fn
		//nolint:nestif
		if comp.WithFn != "" {
			if with, err = s.reg.Named(comp.WithFn); err != nil {
				return newErrorMissingFn(comp.WithFn)
			}
		} else {
			withDatum, mkErr := giraffe.From(comp.With)
			if mkErr != nil {
				return newErrorParsingPayload(mkErr)
			}
			with = hippo.Static(withDatum)
		}

		// TODO msgRe & nameRe can be nil?
		compensator = compensator.For(msgRe, nameRe, step, with)
	}

	plan, ok := s.templates[req.Plan]
	if !ok {
		return newErrorMissingPlan(req.Plan)
	}

	plan = plan.AndCompensator(compensator)

	if len(plan.Names()) == 0 {
		return newErrorMissingPlan(req.Plan)
	}

	runner, err := hippo.MkPipeline(plan)
	if err != nil {
		return newUnknownError(err)
	}

	fin, err := runner.Ekran(ctx, init)
	if err != nil {
		return newErrorProcessingRequest(err)
	}

	if err := fin.MarshalJSONTo(w); err != nil {
		return newUnknownError(err)
	}

	return nil
}

func NewServer(
	reg *hippo.FnRegistry,
	templates map[string]*hippo.Plan,
) (Server, error) {
	for name, plan := range templates {
		if !internal.SimpleName.MatchString(name) {
			return nil, EF("invalid plan name: %s", name)
		}

		if plan == nil {
			return nil, EF("nil plan: %s", name)
		}
	}

	return server{
		reg:       reg,
		templates: maps.Clone(templates),
		serde:     RequestSerde(),
	}.ekran, nil
}
