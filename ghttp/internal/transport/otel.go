package transport

import (
	"errors"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	"go.opentelemetry.io/otel/trace"

	"github.com/hkoosha/giraffe/g11y"
)

var tracer = otel.Tracer("datagen/http")

var errUnsuccessfulRequest = errors.New("http request failed")

type OtelTransport struct {
	rt http.RoundTripper
}

type OtelTransportOption func(*OtelTransport)

func NewOtelTransport(
	base http.RoundTripper,
) *OtelTransport {
	g11y.NonNil(base)

	return &OtelTransport{
		rt: otelhttp.NewTransport(base),
	}
}

func (ot *OtelTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, span := tracer.Start(req.Context(), "HTTP "+req.Method+" "+req.URL.Host)
	defer span.End()

	req = req.WithContext(ctx)

	reqAttr := requestAttributes(req)
	span.SetAttributes(reqAttr...)

	res, err := ot.rt.RoundTrip(req)

	switch {
	case err != nil:
		span.RecordError(err, trace.WithAttributes(reqAttr...))
		span.SetStatus(codes.Error, err.Error())

	case res != nil:
		resAttrs := responseAttributes(res)
		span.SetAttributes(resAttrs...)
		span.SetStatus(httpconv.ClientStatus(res.StatusCode))

		if res.StatusCode >= 400 {
			span.RecordError(
				errUnsuccessfulRequest,
				trace.WithAttributes(resAttrs...),
			)
		}
	}

	//nolint:nilnil
	return res, err
}

func requestAttributes(req *http.Request) []attribute.KeyValue {
	attrs := httpconv.ClientRequest(req)
	if req.URL != nil {
		attrs = append(
			attrs,
			attribute.String("http.host", req.URL.Host),
			attribute.String("http.path", req.URL.Path),
			attribute.String("http.scheme", req.URL.Scheme),
			attribute.String("http.raw_query", req.URL.RawQuery),
		)
	}

	return attrs
}

func responseAttributes(res *http.Response) []attribute.KeyValue {
	if res == nil {
		return nil
	}

	attrs := httpconv.ClientResponse(res)
	if res.Status != "" {
		attrs = append(
			attrs,
			attribute.String("http.status_text", res.Status),
		)
	}

	return attrs
}
