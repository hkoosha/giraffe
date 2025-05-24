package internal

import "go.opentelemetry.io/otel/attribute"

func Join(attrs, extra []attribute.KeyValue) []attribute.KeyValue {
	cp := make([]attribute.KeyValue, len(attrs)+len(extra))
	copy(cp, attrs)
	copy(cp[len(attrs):], extra)

	return cp
}
