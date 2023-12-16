package squawker

import (
	"context"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

var PLUGIN_NAME = "squawker"

var log = clog.NewWithPlugin(PLUGIN_NAME)

type Squawker struct {
	Next plugin.Handler
}

func NewSquawker(next plugin.Handler) plugin.Handler {
	return &Squawker{Next: next}
}

func (s *Squawker) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	log.Debug("Received response")
	return plugin.NextOrFailure(s.Name(), s.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (s Squawker) Name() string { return PLUGIN_NAME }
