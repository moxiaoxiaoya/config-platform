package client

import (
	"config-platform/resolver"
)

type Client struct {
	subSystem string

	discover resolver.Discovery
}

func NewClient(conf ClientConfig, opts ...Option) *Client {
	c := &Client{
		subSystem: conf.SubSystem,
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(c)
		}
	}
	return c
}

type Option func(*Client)

func WithDiscovery(r resolver.Discovery) Option {
	return func(c *Client) {
		c.discover = r
	}
}

func (c *Client) Start() {

}
