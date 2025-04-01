package pipelines

import (
	"context"
	"encoding/json"
	"io"
	"maps"
	"net/http"
	"regexp"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/g11y"
)

type Server func(context.Context, io.Reader, io.Writer) error

func (s Server) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := s(r.Context(), r.Body, w); err != nil {
		msg := err.Error()
		if g11y.IsUnsafeError() {
			msg += "\n\n" + g11y.FmtStacktraceOf(err)
		}
		http.Error(w, msg, http.StatusBadRequest)

		return
	}
}

type server struct {
	reg       FnRegistry
	templates map[string]Plan
}

func (s *server) ekran(
	ctx context.Context,
	r io.Reader,
	w io.Writer,
) error {
	var req Request
	dec := json.NewDecoder(r)
	dec.UseNumber()
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		return newErrorParsingPayload(err)
	}

	init, err := giraffe.Make(req.Init)
	if err != nil {
		return newErrorParsingPayload(err)
	}

	var compensator Compensator
	if req.Compensations == nil {
		req.Compensations = &[]RequestCompensations{}
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

		var with Fn
		//nolint:nestif
		if comp.WithFn != "" {
			if with, err = s.reg.Get(comp.WithFn); err != nil {
				return newErrorMissingFn(comp.WithFn)
			}
		} else {
			withDatum, mkErr := giraffe.Make(comp.With)
			if mkErr != nil {
				return newErrorParsingPayload(mkErr)
			}
			with = Static(withDatum)
		}

		compensator = compensator.For(msgRe, nameRe, step, with)
	}

	plan, ok := s.templates[req.Plan]
	if !ok {
		return newErrorMissingPlan(req.Plan)
	}

	plan = plan.WithCompensator(compensator)

	if len(plan.Names()) == 0 {
		return newErrorMissingPlan(req.Plan)
	}

	runner, err := Runner(plan)
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
	reg FnRegistry,
	templates map[string]Plan,
) Server {
	s := server{
		reg:       reg,
		templates: maps.Clone(templates),
	}

	return s.ekran
}
