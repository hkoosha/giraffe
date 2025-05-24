package httpmethod

import (
	"net/http"

	"github.com/hkoosha/giraffe/core/t11y"
)

//goland:noinspection GoUnusedConst
const (
	GET    T = http.MethodGet
	POST   T = http.MethodPost
	PUT    T = http.MethodPut
	PATCH  T = http.MethodPatch
	DELETE T = http.MethodDelete
)

type T string

func (t T) String() string {
	return string(t)
}

func (t T) HasBody() bool {
	return t != GET && t != DELETE
}

func Of(v string) (T, bool) {
	switch v {
	case GET.String():
		return GET, true

	case POST.String():
		return POST, true

	case PUT.String():
		return PUT, true

	case PATCH.String():
		return PATCH, true

	case DELETE.String():
		return DELETE, true
	}

	return T(""), false
}

func MustOf(v string) T {
	switch v {
	case GET.String():
		return GET

	case POST.String():
		return POST

	case PUT.String():
		return PUT

	case PATCH.String():
		return PATCH

	case DELETE.String():
		return DELETE
	}

	panic(t11y.TracedFmt("unsupported http method: %s", v))
}

func All() []T {
	return []T{
		GET,
		POST,
		PUT,
		PATCH,
		DELETE,
	}
}
