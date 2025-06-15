package gredis

import (
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/typing"
)

const (
	minNameLen = 1
	maxNameLen = 127

	namespaceKeySep     = ":"
	namespaceSep        = "_"
	maxNamespaceNesting = 5

	ttlMin = 2 * time.Second
)

var zero = &Config{
	namespace: "",
	nsParts:   []string{},
	ttl:       0,
	timeout:   2 * time.Second,
}

func NewConfig(
	namespace string,
	ttl time.Duration,
) *Config {
	return zero.
		WithTTL(ttl).
		Namespaced(namespace)
}

type Config struct {
	namespace string
	nsParts   []string
	ttl       time.Duration
	timeout   time.Duration
	rds       *redis.Client
}

func (c *Config) Ensure() *Config {
	switch {
	case c.ttl < ttlMin:
		panic(EF("ttl too low: %v", c.ttl))

	case c.timeout < 1*time.Millisecond:
		panic(EF("timeout: %v", c.timeout))

	case !typing.IsMachineReadableName(c.namespace, minNameLen, maxNameLen):
		panic(EF("invalid namespace: %s", c.namespace))

	case len(c.nsParts) > maxNamespaceNesting:
		panic(EF("namespace too deep: %s", c.namespace))
	}

	return c
}

func (c *Config) KeyPrefix() string {
	return c.namespace
}

func (c *Config) TTL() time.Duration {
	return c.ttl
}

func (c *Config) WithTTL(
	ttl time.Duration,
) *Config {
	cp := &*c
	cp.ttl = ttl
	return cp.Ensure()
}

func (c *Config) Timeout() time.Duration {
	return c.timeout
}

func (c *Config) WithTimeout(
	timeout time.Duration,
) *Config {
	cp := &*c
	cp.timeout = timeout
	return cp.Ensure()
}

func (c *Config) Namespaced(
	namespace string,
) *Config {
	cp := &*c
	cp.nsParts = append(cp.nsParts, namespace)
	cp.namespace = strings.Join(cp.nsParts, namespaceSep) + namespaceKeySep
	return cp.Ensure()
}
