package zcache

import (
	"slices"
	"strings"
	"time"

	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/glog"
	"github.com/hkoosha/giraffe/typing"
)

const (
	minNameLen = 1
	maxNameLen = 15

	namespaceKeySep     = ":"
	namespaceSep        = "_"
	maxNamespaceNesting = 3
	maxNamespaceLen     = 31
)

func NewConfig(
	lg glog.GLog,
	domain string,
	ttl time.Duration,
) *Config {
	lg = lg.Named("cache")

	cfg := &Config{
		lg:             lg,
		domain:         domain,
		keyPrefix:      "",
		namespaceParts: []string{},
		ttl:            ttl,
	}

	// Redo for validation.
	return cfg.WithTTL(ttl).Namespaced(domain)
}

type Config struct {
	lg             glog.GLog
	domain         string
	keyPrefix      string
	namespaceParts []string
	ttl            time.Duration
}

func (c *Config) clone() *Config {
	return &Config{
		lg:             c.lg,
		domain:         c.domain,
		keyPrefix:      c.keyPrefix,
		ttl:            c.ttl,
		namespaceParts: slices.Clone(c.namespaceParts),
	}
}

func (c *Config) Lg() glog.GLog {
	return c.lg
}

func (c *Config) Domain() string {
	return c.domain
}

func (c *Config) KeyPrefix() string {
	return c.keyPrefix
}

func (c *Config) TTL() time.Duration {
	return c.ttl
}

func (c *Config) WithTTL(
	ttl time.Duration,
) *Config {
	if ttl < 1*time.Minute || ttl > 365*24*time.Hour {
		panic(EF("invalid ttl"))
	}

	cp := c.clone()
	cp.ttl = ttl

	return cp
}

func (c *Config) Namespaced(
	namespace string,
) *Config {
	cp := c.clone()
	cp.keyPrefix, cp.namespaceParts = mkNamespace(namespace, cp.namespaceParts...)

	return cp
}

//nolint:nonamedreturns
func mkNamespace(
	namespace string,
	parts ...string,
) (nestedNamespace string, nestedNamespaceParts []string) {
	parts = append(parts, namespace)
	ns := strings.Join(parts, namespaceSep) + namespaceKeySep

	switch {
	case !typing.IsMachineReadableName(namespace, minNameLen, maxNameLen):
		panic(EF("%s", "invalid namespace: "+ns))

	case len(parts) > maxNamespaceNesting:
		panic(EF("%s", "namespace too deep: "+ns))

	case len(ns) > maxNamespaceLen:
		panic(EF("%s", "namespace too long: "+ns))
	}

	return ns, parts
}
