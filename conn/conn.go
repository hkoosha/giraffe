package conn

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/hkoosha/giraffe/conn/internal"
	"github.com/hkoosha/giraffe/t11y/glog"
	"github.com/hkoosha/giraffe/zebra/serdes"
)

const UserAgent = "Giraffe/1.0"

//goland:noinspection GoUnusedConst
const (
	ReasonUnexpectedStatusCode FailureReason = 2
	ReasonEmptyResponse        FailureReason = 3
)

type FailureReason uint

// ============================================================================.

type FailedResponseError struct {
	Resp   any
	Reason FailureReason
}

func (e *FailedResponseError) Error() string {
	return "http request failed: " + strconv.FormatUint(uint64(e.Reason), 10)
}

// ============================================================================.

type HeaderFilter = func(
	context.Context,
	Config,
	string,
	string,
) bool

type HeaderProvider = func(
	context.Context,
	Config,
) string

type RetryIfFn = func(
	ctx context.Context,
	resp *http.Response,
	err error,
	attempt uint,
	cfg Config,
) (bool, error)

type ConfigLgRead interface {
	internal.Sealed

	Lg() glog.Lg
	IsLogged() bool
	IsPlainLog() bool
	LgHeaderFilter() HeaderFilter
	LgHeaderMask() HeaderFilter
}

type ConfigLgWrite interface {
	internal.Sealed

	WithLg(glog.Lg) Config

	WithLogged() Config
	WithoutLogged() Config
	SetLogged(bool) Config

	WithPlainLog() Config
	WithoutPlainLog() Config
	SetPlainLog(bool) Config

	WithLgFilteredHeaders(HeaderFilter) Config

	WithLgMaskedHeaders(HeaderFilter) Config
}

type ConfigRetryRead interface {
	internal.Sealed

	IsRetryLog() bool
	RetryMax() uint
	RetryStatusCodes() []int
	RetryIf() RetryIfFn
	RetryBackoffDuration() time.Duration
}

type ConfigRetryWrite interface {
	internal.Sealed

	WithRetryLogged() Config
	WithoutRetryLogged() Config
	SetRetryLogged(bool) Config
	WithMaxRetries(uint) Config
	WithRetryStatusCodes(...int) Config
	WithRetryIf(fn RetryIfFn) Config
	WithoutRetryIf() Config
	WithRetryBackoffDuration(time.Duration) Config
}

type ConfigRead interface {
	internal.Sealed

	ConfigLgRead
	ConfigRetryRead

	Ensure()
	Std() *http.Client

	HeaderOverwrites() map[string]string
	HeaderOverwriters() map[string]HeaderProvider
	ExpectingStatusCode() int
	Endpoint() string
	PathPrefix() string
	Method() string
	Timeout() time.Duration
	TraceOptions() []otelhttp.Option
	IsTraced() bool
}

type ConfigWrite interface {
	internal.Sealed

	ConfigLgWrite
	ConfigRetryWrite

	WithBearerToken(string) Config
	WithBearerProvider(HeaderProvider) Config
	WithoutBearerProvider() Config

	WithExpectingStatusCode(int) Config
	WithoutExpectingStatusCode() Config

	WithEndpoint(string) Config
	WithoutEndpoint() Config
	WithPathPrefix(string) Config
	WithMethod(string) Config
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
	AndHeaders(map[string]string) Config
	AndHeader(
		name string,
		value string,
	) Config
}

// Config
// Keep all setter methods prefixed with either of:
// - With
// - Without
// - Set
// - And
// This makes refactoring easier and more consistent, since all the methods
// minus the getters are implemented in [Client] too.
type Config interface {
	internal.Sealed

	Conn() Raw
	Serde() serdes.Serde[any]

	ConfigRead
	ConfigWrite
}

type Conn[R any] interface {
	internal.Sealed

	Std() *http.Client
	Cfg() Config
	Raw() Raw

	Call(
		ctx context.Context,
		body any,
		path ...string,
	) (R, error)

	CallAs(
		ctx context.Context,
		method string,
		body any,
		path ...string,
	) (R, error)

	Patch(
		ctx context.Context,
		body any,
		path ...string,
	) (R, error)

	Put(
		ctx context.Context,
		body any,
		path ...string,
	) (R, error)

	Post(
		ctx context.Context,
		body any,
		path ...string,
	) (R, error)

	PostForHeaders(
		ctx context.Context,
		body any,
		path ...string,
	) (http.Header, error)

	Get(
		ctx context.Context,
		path ...string,
	) (R, error)

	GetForHeaders(
		ctx context.Context,
		path ...string,
	) (http.Header, error)

	Delete(
		ctx context.Context,
		path ...string,
	) (R, error)
}

type Raw = Conn[[]byte]

// ============================================================================.

func MakeAny[R any](
	cfg Config,
	serde serdes.Serde[R],
) Conn[R] {
	cloned := cfgOf(cfg)

	return newConn[R](cloned, serde)
}

func MakeJson[R any](
	cfg Config,
) Conn[R] {
	return MakeAny[R](cfg, serdes.Json[R]())
}

// ====================================.

func OfAny(
	lg glog.Lg,
	timeout time.Duration,
	serde serdes.Serde[any],
) Config {
	return newConfig(lg, timeout, serde)
}

func OfJson(
	lg glog.Lg,
	timeout time.Duration,
) Config {
	return OfAny(lg, timeout, serdes.Json[any]())
}
