package internal

import (
	"regexp"
)

var (
	SimpleName       = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9_]*`)
	ScopedSimpleName = regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9_]*)(/[a-zA-Z][a-zA-Z0-9_]*)?`)
)
