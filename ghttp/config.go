package ghttp

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/hkoosha/giraffe/glog"
)

const (
	UserAgent = "Giraffe/1.0"
)

type HeaderFilter = func(ctx context.Context, h, v string) bool

type HeaderProvider = func(context.Context, Config) string

type RetryIfFn = func(
	ctx context.Context,
	resp *http.Response,
	err error,
	attempt uint,
	cfg Config,
) (bool, error)

type ConfigRead interface {
	Ensure()
	Std() *http.Client

	Lg() glog.Lg
	IsLogged() bool
	IsPlainLog() bool
	LgHeaderFilter() HeaderFilter
	LgHeaderMask() HeaderFilter
	HeaderOverwrites() map[string]string
	HeaderOverwriters() map[string]HeaderProvider
	ExpectingStatusCode() int
	Endpoint() string
	PathPrefix() string
	Timeout() time.Duration
	TraceOptions() []otelhttp.Option
	IsTraced() bool
}

type ConfigRetryRead interface {
	IsRetryLog() bool
	RetryMax() uint
	RetryStatusCodes() []int
	RetryIf() RetryIfFn
	RetryBackoffDuration() time.Duration
}

type ConfigRetryWrite interface {
	WithRetryLogged() Config
	WithoutRetryLogged() Config
	SetRetryLogged(bool) Config
	WithMaxRetries(uint) Config
	WithRetryStatusCodes(...int) Config
	WithRetryIf(fn RetryIfFn) Config
	WithoutRetryIf() Config
	WithRetryBackoffDuration(time.Duration) Config
}

// Config
// Keep all setter methods prefixed with either of:
// - Is
// - With
// - Without
// - Set
// - And
// This makes refactoring easier and more consistent, since all the methods
// minus the getters are implemented in [Client] too.
type Config interface {
	ConfigRead
	ConfigRetryRead
	ConfigRetryWrite

	WithLg(
		glog.Lg,
	) Config

	WithLogged() Config
	WithoutLogged() Config
	SetLogged(bool) Config

	WithPlainLog() Config
	WithoutPlainLog() Config
	SetPlainLog(bool) Config

	WithLgFilteredHeaders(
		HeaderFilter,
	) Config
	WithLgMaskedHeaders(
		HeaderFilter,
	) Config

	WithBearerToken(string) Config
	WithBearerProvider(HeaderProvider) Config
	WithoutBearerProvider() Config

	WithExpectingStatusCode(int) Config
	WithoutExpectingStatusCode() Config

	WithEndpoint(string) Config
	WithoutEndpoint() Config
	WithPathPrefix(string) Config
	WithoutPathPrefix() Config
	AndPathPrefix(string) Config

	WithTimeout(time.Duration) Config

	WithTransport(http.RoundTripper) Config

	WithTraced() Config
	WithoutTraced() Config
	SetTraced(bool) Config
	WithTraceOptions(...otelhttp.Option) Config

	WithHeaderOverwrites(
		includeDefaults bool,
		h map[string]string,
	) Config
}
