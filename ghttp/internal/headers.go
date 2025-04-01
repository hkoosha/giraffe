package internal

import (
	"maps"

	"github.com/hkoosha/giraffe/ghttp/headers"
	"github.com/hkoosha/giraffe/zebra/z"
)

var (
	masked = map[string]z.NA{
		headers.Authorization: z.None,
		headers.Cookie:        z.None,
		headers.SetCookie:     z.None,
		headers.Signature:     z.None,
	}

	filtered = map[string]z.NA{
		headers.StrictTransportSecurity: z.None,
		headers.PermissionsPolicy:       z.None,
		headers.ReferrerPolicy:          z.None,
		headers.ExpectCt:                z.None,
		headers.ContentSecurityPolicy:   z.None,
		headers.CacheControl:            z.None,
		headers.XContentTypeOptions:     z.None,
		headers.XFrameOptions:           z.None,
	}
)

func MaskedHeaders() map[string]z.NA {
	return maps.Clone(masked)
}

func FilteredHeaders() map[string]z.NA {
	return maps.Clone(filtered)
}
