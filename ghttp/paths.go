package ghttp

import (
	"slices"
	"strings"
)

//goland:noinspection HttpUrlsUsage
func Join(parts ...string) string {
	if len(parts) == 0 {
		panic("no parts to join")
	}

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	parts = slices.DeleteFunc(parts, func(s string) bool { return s == "" })

	// Check and see if any non-empty path originally started with slash,
	// before we remove all the slashes in the next block:.
	prefix := ""
	if len(parts) > 0 && parts[0] != "" && parts[0][0] == '/' {
		prefix = "/"
	}

	// Remove extra slashes:.
	for i := range parts {
		parts[i] = strings.Trim(parts[i], "/")
	}
	parts = slices.DeleteFunc(parts, func(s string) bool { return s == "" })

	fin := prefix + strings.Join(parts, "/")

	var probe string
	if strings.HasPrefix(fin, "http://") {
		probe = strings.TrimPrefix(fin, "http://")
	} else {
		probe = strings.TrimPrefix(fin, "https://")
	}

	if strings.Contains(probe, "//") {
		panic("http path has multiple consecutive fwd slashes: " + fin)
	}

	return fin
}
