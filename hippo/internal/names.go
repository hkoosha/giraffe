package internal

import (
	"regexp"
)

var SimpleName = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9_]*`)
