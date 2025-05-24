package conn

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn/internal"
	"github.com/hkoosha/giraffe/core/serdes"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/glog"
)

const (
	UserAgent      = "Giraffe/1.0"
	DefaultTimeout = 5 * time.Second
)

// ============================================================================.

type HeaderFilter = func(
	_ context.Context,
	_ Config,
	name string,
	value string,
) bool

type HeaderProvider = func(context.Context, Config) string

type RetryIfFn = func(
	_ context.Context,
	_ *http.Response,
	_ error,
	attempt uint,
	_ Config,
) (bool, error)

// ============================================================================.

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

type ConfigSerdeRead interface {
	internal.Sealed

	RxSerde() any
	TxSerde() any
}

type ConfigSerdeWrite interface {
	internal.Sealed

	WithSerdes(
		tx any,
		rx any,
	) Config

	WithRxSerde(any) Config
	WithTxSerde(any) Config

	WithJsonRxSerde() Config
	WithJsonTxSerde() Config
	WithJsonSerde() Config

	WithStringRxSerde() Config

	WithBytesRxSerde() Config
	WithBytesTxSerde() Config
	WithBytesSerde() Config
}

type ConfigRead interface {
	internal.Sealed

	ConfigLgRead
	ConfigRetryRead
	ConfigSerdeRead

	Ensure()
	Std() *http.Client

	HeaderOverwrites() map[string]string
	HeaderOverwriters() map[string]HeaderProvider
	ExpectingStatusCodes() []int
	Endpoint() string
	Endpoints() map[string]string
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
	ConfigSerdeWrite

	WithBearerToken(string) Config
	WithBearerProvider(HeaderProvider) Config
	WithoutBearerProvider() Config

	WithExpectingStatusCodes(...int) Config
	WithoutExpectingStatusCodes() Config

	AndEndpoint(name, addr string) Config
	WithEndpoints(map[string]string) Config
	WithoutEndpoints() Config
	WithEndpoint(string) Config
	WithoutEndpoint() Config
	WithPathPrefix(string) Config
	WithMethod(string) Config
	WithoutPathPrefix() Config
	AndPathPrefix(string) Config

	WithEndpointNamed(string) (Config, error)
	WithMustEndpointNamed(string) Config

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

	Raw() Raw
	Datum() Datum

	ConfigRead
	ConfigWrite
}

// ============================================================================.

type Configured interface {
	internal.Sealed

	Std() *http.Client
	Cfg() Config
}

type ToRaw interface {
	internal.Sealed

	Raw() Raw
}

type Headered[TX, RX any] interface {
	HCall(
		_ context.Context,
		_ *TX,
		path ...string,
	) (status int, headers map[string]string, _ RX, _ error)

	HPatch(
		_ context.Context,
		_ TX,
		path ...string,
	) (headers map[string]string, _ RX, _ error)

	HPut(
		_ context.Context,
		_ TX,
		path ...string,
	) (headers map[string]string, _ RX, _ error)

	HPost(
		_ context.Context,
		_ TX,
		path ...string,
	) (headers map[string]string, _ RX, _ error)

	HGet(
		_ context.Context,
		path ...string,
	) (headers map[string]string, _ RX, _ error)

	HDelete(
		_ context.Context,
		path ...string,
	) (headers map[string]string, _ RX, _ error)
}

type Headerless[TX, RX any] interface {
	Call(
		_ context.Context,
		_ *TX,
		path ...string,
	) (RX, error)

	Patch(
		_ context.Context,
		_ TX,
		path ...string,
	) (RX, error)

	Put(
		_ context.Context,
		_ TX,
		path ...string,
	) (RX, error)

	Post(
		_ context.Context,
		_ TX,
		path ...string,
	) (RX, error)

	Get(
		_ context.Context,
		path ...string,
	) (RX, error)

	Delete(
		_ context.Context,
		path ...string,
	) (RX, error)
}

type Conn[TX, RX any] interface {
	internal.Sealed

	IsExpected(
		_ context.Context,
		code int,
	) bool

	Configured
	ToRaw
	Headered[TX, RX]
	Headerless[TX, RX]
}

type Raw = Conn[[]byte, []byte]

type Datum = Conn[giraffe.Datum, giraffe.Datum]

// ============================================================================.

func MakeCfg(
	lg glog.Lg,
) Config {
	return newConfig(lg, DefaultTimeout)
}

func Make[TX, RX any](
	cfg Config,
	txSered serdes.Serde[TX],
	rxSered serdes.Serde[RX],
) Conn[TX, RX] {
	withSerde := cfgOf(cfg).withSerdes(txSered, rxSered)
	return newConn[TX, RX](withSerde)
}

func MakeJson[TX, RX any](
	cfg Config,
) Conn[TX, RX] {
	nonRaw := false

	var txSerde any
	var txSample TX
	if _, ok := any(txSample).([]byte); ok {
		txSerde = serdes.Bytes()
	} else {
		txSerde = serdes.Json[TX]
		nonRaw = true
	}

	var rxSerde any
	var rxSample RX
	if _, ok := any(rxSample).([]byte); ok {
		rxSerde = serdes.Bytes()
	} else {
		rxSerde = serdes.Json[RX]
		nonRaw = true
	}

	if !nonRaw {
		panic(EF("use Raw connection type instead of Json when both RX and TX are byte slices"))
	}

	withSerde := cfgOf(cfg).withSerdes(txSerde, rxSerde)

	return newConn[TX, RX](withSerde)
}
