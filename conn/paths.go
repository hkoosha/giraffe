package conn

import (
	"slices"
	"strings"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

var danger = []string{
	"//",
	"..",
}

func Join(
	parts ...string,
) string {
	return join(parts)
}

func join(
	parts []string,
) string {
	if parts == nil {
		return ""
	}

	orig := slices.Clone(parts)

	preprocess(&parts)

	parts, queries := splitQuery(parts)

	ensureParts(orig, parts)

	parts, prefix := processParts(parts)

	sb := &strings.Builder{}
	sb.WriteString(prefix)

	writeParts(sb, parts)
	writeQueries(sb, queries)

	fin := sb.String()

	ensure(orig, fin)

	return fin
}

func writeQueries(
	sb *strings.Builder,
	queries []string,
) {
	if len(queries) == 0 {
		return
	}

	sb.WriteByte('?')

	last := len(queries) - 1
	for i, p := range queries {
		sb.WriteString(p)
		if i != last {
			sb.WriteByte('&')
		}
	}
}

func writeParts(
	sb *strings.Builder,
	parts []string,
) {
	last := len(parts) - 1
	for i, v := range parts {
		sb.WriteString(v)
		if i != last {
			sb.WriteByte('/')
		}
	}
}

//nolint:nonamedreturns
func processParts(
	parts []string,
) (
	_ []string,
	prefix string,
) {
	// Check and see if any non-empty path originally started with slash,
	// before we remove all the slashes in the next block.
	prefix = ""
	if len(parts) > 0 && parts[0] != "" && parts[0][0] == '/' {
		prefix = "/"
	}

	// Remove extra slashes.
	for i := range parts {
		parts[i] = strings.Trim(parts[i], "/")
	}
	parts = slices.DeleteFunc(parts, func(s string) bool { return s == "" })

	return parts, prefix
}

func ensureParts(
	orig []string,
	parts []string,
) {
	if slices.Contains(parts, "") {
		panic(EF("empty path parts: %v", orig))
	}
}

//goland:noinspection HttpUrlsUsage
func ensure(
	orig []string,
	fin string,
) {
	var probe string
	switch {
	case strings.HasPrefix(fin, "https://"):
		probe = strings.TrimPrefix(fin, "https://")

	case strings.HasPrefix(fin, "http://"):
		probe = strings.TrimPrefix(fin, "http://")

	default:
		probe = fin
	}

	for _, d := range danger {
		if strings.Contains(probe, d) {
			panic(EF("illegal path: %v", orig))
		}
	}
}

func preprocess(parts *[]string) {
	p := *parts
	for i := range p {
		p[i] = strings.TrimSpace(p[i])
	}
	*parts = slices.DeleteFunc(p, func(it string) bool {
		return it == ""
	})
}

//nolint:nonamedreturns
func splitQuery(parts []string) (
	pathParts []string,
	queries []string,
) {
	i := slices.Index(parts, "?")

	if i < 0 {
		return parts, nil
	}

	return parts[0:i], parts[i+1:]
}
