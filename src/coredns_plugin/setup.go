package squawker

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register(PLUGIN_NAME, setup) }

func setup(c *caddy.Controller) error {
	config, err := parseArgs(c)
	if err != nil {
		return err
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return NewSquawker(next, *config)
	})

	return nil
}

func parseArgs(c *caddy.Controller) (*SquawkerConfig, error) {
	config := SquawkerConfig{}

	c.Next() // Ignore "squawker" and give us the next token.

	args := c.RemainingArgs()

	if len(args) != 1 {
		return nil, plugin.Error(PLUGIN_NAME, c.ArgErr())
	}

	config.address = args[0]

	return &config, nil
}
