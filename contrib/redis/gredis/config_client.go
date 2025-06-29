package gredis

import (
	"github.com/hkoosha/giraffe/contrib/redis/gredis/internal"
)

type ClientConfig interface {
	internal.Sealed

	IsEnabled() bool
	IsTLS() bool
	IsInsecure() bool

	Host() string
	Port() uint16
	Auth() string
	DB() uint8
}

type ClientConfigWrite interface {
	internal.Sealed

	ClientConfig

	WithEnabled() ClientConfigWrite
	WithoutEnabled() ClientConfigWrite
	SetEnabled(bool) ClientConfigWrite

	WithTLS() ClientConfigWrite
	WithoutTLS() ClientConfigWrite
	SetTLS(bool) ClientConfigWrite

	WithInsecure() ClientConfigWrite
	WithoutInsecure() ClientConfigWrite
	SetInsecure(bool) ClientConfigWrite

	WithHost(string) ClientConfigWrite
	WithPort(uint16) ClientConfigWrite
	WithAuth(string) ClientConfigWrite
	WithDB(uint8) ClientConfigWrite
}

func NewClientConfig() ClientConfigWrite {
	return &config{}
}

type config struct {
	internal.Sealer
	isEnabled  bool
	isTLS      bool
	isInsecure bool
	host       string
	port       uint16
	auth       string
	db         uint8
}

func (c *config) shallow() *config {
	return &*c
}

// ============================================================================.

func (c *config) IsEnabled() bool {
	return c.isEnabled
}

func (c *config) IsTLS() bool {
	return c.isTLS
}

func (c *config) IsInsecure() bool {
	return c.isInsecure
}

func (c *config) Host() string {
	return c.host
}

func (c *config) Port() uint16 {
	return c.port
}

func (c *config) Auth() string {
	return c.auth
}

func (c *config) DB() uint8 {
	return c.db
}

func (c *config) WithEnabled() ClientConfigWrite {
	return c.SetEnabled(true)
}

func (c *config) WithoutEnabled() ClientConfigWrite {
	return c.SetEnabled(false)
}

func (c *config) SetEnabled(b bool) ClientConfigWrite {
	cp := c.shallow()
	cp.isEnabled = b
	return cp
}

func (c *config) WithTLS() ClientConfigWrite {
	return c.SetTLS(true)
}

func (c *config) WithoutTLS() ClientConfigWrite {
	return c.SetTLS(false)
}

func (c *config) SetTLS(b bool) ClientConfigWrite {
	cp := c.shallow()
	cp.isTLS = b
	return cp
}

func (c *config) WithInsecure() ClientConfigWrite {
	return c.SetInsecure(true)
}

func (c *config) WithoutInsecure() ClientConfigWrite {
	return c.SetInsecure(false)
}

func (c *config) SetInsecure(b bool) ClientConfigWrite {
	cp := c.shallow()
	cp.isInsecure = b
	return cp
}

func (c *config) WithHost(s string) ClientConfigWrite {
	cp := c.shallow()
	cp.host = s
	return cp
}

func (c *config) WithPort(u uint16) ClientConfigWrite {
	cp := c.shallow()
	cp.port = u
	return cp
}

func (c *config) WithAuth(s string) ClientConfigWrite {
	cp := c.shallow()
	cp.auth = s
	return cp
}

func (c *config) WithDB(u uint8) ClientConfigWrite {
	cp := c.shallow()
	cp.db = u
	return cp
}
