package conn

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hkoosha/giraffe/conn/headers"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/zebra/z"
)

var (
	defaultHeaders = map[string]string{
		headers.UserAgent: UserAgent,
	}

	// Copy of http.DefaultTransport, but surely not messed up with.
	//nolint:exhaustruct
	defaultTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	defaultFilteredHeaders = z.Set[string]{
		headers.StrictTransportSecurity: z.None,
		headers.PermissionsPolicy:       z.None,
		headers.ReferrerPolicy:          z.None,
		headers.ExpectCt:                z.None,
		headers.ContentSecurityPolicy:   z.None,
		headers.CacheControl:            z.None,
		headers.XContentTypeOptions:     z.None,
		headers.XFrameOptions:           z.None,
	}

	defaultMaskedHeaders = z.Set[string]{
		headers.Authorization: z.None,
		headers.Cookie:        z.None,
		headers.SetCookie:     z.None,
		headers.Signature:     z.None,
	}
)

func defaultHeaderFilter(
	_ context.Context,
	_ Config,
	h string,
	_ string,
) bool {
	_, ok := defaultFilteredHeaders[h]
	return ok
}

func defaultHeaderMasked(
	_ context.Context,
	_ Config,
	h string,
	_ string,
) bool {
	_, ok := defaultMaskedHeaders[h]
	return ok
}

// =============================================================================.

//nolint:unused
type retryKeyT int

//nolint:unused
var retryKey retryKeyT

//nolint:unused
func getRetries(ctx context.Context) int {
	retries, ok := ctx.Value(retryKey).(*int)
	if !ok {
		panic(EF("retry key not set"))
	}

	return *retries
}

//nolint:unused
func incRetries(ctx context.Context) {
	retries, ok := ctx.Value(retryKey).(*int)
	if !ok {
		panic(EF("retry key not set"))
	}

	*retries++
}

// =============================================================================.

// type seal struct {
// 	sealed bool
// }

// =============================================================================.

func withBearerPrefix(token string) string {
	const prefix = "Bearer "
	const l = len(prefix)
	const prefixLower = "bearer "

	switch {
	case len(token) < l,
		strings.ToLower(token[:l]) != prefixLower:

		return prefix + token

	default:
		return token
	}
}

// =============================================================================.
// From [url] pkg.
func escape(s string) string {
	hexCount := 0
	for i := range len(s) {
		c := s[i]
		if !('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z') && !('0' <= c && c <= '9') {
			switch c {
			case '!', '$', '&', '\'', '(', ')', '*', '+', ',',
				';', '=', ':', '[', ']', '<', '>', '"',
				'-', '_', '.', '~':
				// Nothing.

			default:
				hexCount++
			}
		}
	}

	if hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte

	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	copy(t, s)
	for i := range len(s) {
		if s[i] == ' ' {
			t[i] = '+'
		}
	}
	return string(t)
}

// Copied from [url.URL.String].
//
//nolint:nestif
func partsOf(u *url.URL) (string, string) {
	const nn = len(":" + "//" + "//" + ":" + "@" + "/" + "./" + "?" + "#")

	if u.Opaque != "" {
		panic(EF("opaque url not supported: %s", u.Redacted()))
	}

	var buf0 strings.Builder

	n0 := len(u.Scheme)
	if !u.OmitHost && (u.Scheme != "" || u.Host != "" || u.User != nil) {
		username := u.User.Username()
		password, _ := u.User.Password()
		n0 += len(username) + len(password) + len(u.Host)
	}
	n0 += nn
	buf0.Grow(n0)

	var buf1 strings.Builder
	n1 := nn + len(u.Path) + len(u.RawQuery) + len(u.RawFragment)
	buf1.Grow(n1)

	if u.Scheme != "" {
		buf0.WriteString(u.Scheme)
		buf0.WriteByte(':')
	}

	if u.Scheme != "" || u.Host != "" || u.User != nil {
		if !(u.OmitHost && u.Host == "") || u.User != nil {
			if u.Host != "" || u.Path != "" || u.User != nil {
				buf0.WriteString("//")
			}
			if ui := u.User; ui != nil {
				buf0.WriteString(ui.String())
				buf0.WriteByte('@')
			}
			if h := u.Host; h != "" {
				buf0.WriteString(escape(h))
			}
		}
	}

	path := u.EscapedPath()
	if path != "" && path[0] != '/' && u.Host != "" {
		buf1.WriteByte('/')
	}
	if buf0.Len() == 0 {
		if segment, _, _ := strings.Cut(path, "/"); strings.Contains(segment, ":") {
			buf1.WriteString("./")
		}
	}
	buf1.WriteString(path)

	if u.ForceQuery || u.RawQuery != "" {
		buf1.WriteByte('?')
		buf1.WriteString(u.RawQuery)
	}
	if u.Fragment != "" {
		buf1.WriteByte('#')
		buf1.WriteString(u.EscapedFragment())
	}

	return buf0.String(), buf1.String()
}

type giraffeRT struct {
	cfg *config
	rt  http.RoundTripper
}

func (u giraffeRT) modify(
	ctx context.Context,
	orig *http.Request,
	req func() *http.Request,
) error {
	for h, v := range u.cfg.header.overwrite {
		req().Header.Set(h, v)
	}

	for h, fn := range u.cfg.header.overwriters {
		v := fn(ctx, u.cfg)
		req().Header.Set(h, v)
	}

	if u.cfg.http.endpoint != "" || u.cfg.http.pathPrefix != "" {
		endpoint, path := partsOf(orig.URL)

		if u.cfg.http.endpoint != "" {
			endpoint = u.cfg.http.endpoint
		}

		uParse, err := url.Parse(Join(endpoint, u.cfg.http.pathPrefix, path))
		if err != nil {
			return err
		}

		req().URL = uParse
	}

	return nil
}

func (u giraffeRT) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	cloned := false

	r := req

	//nolint:contextcheck
	err := u.modify(req.Context(), req, func() *http.Request {
		if !cloned {
			r = req.Clone(req.Context())
		}
		return r
	})
	if err != nil {
		return nil, err
	}

	return u.rt.RoundTrip(r)
}
