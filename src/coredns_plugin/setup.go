package squawker

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register(PLUGIN_NAME, setup) }

func setup(c *caddy.Controller) error {
	c.Next() // Ignore "squawker" and give us the next token.

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return NewSquawker(next)
	})

	return nil
}
