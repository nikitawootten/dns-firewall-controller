package squawker

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
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
	log.Info("Squawk!")

	state := request.Request{W: w, Req: r}
	client := state.IP()

	rrw := dnstest.NewRecorder(w)
	rc, err := plugin.NextOrFailure(s.Name(), s.Next, ctx, rrw, r)
	if err != nil {
		return rc, err
	}

	allowed_ips := []net.IP{}
	var max_ttl uint32 = 0

	if rrw.Msg == nil {
		log.Info("Squawker found no answers")
		return rc, err
	}

	log.Infof("Squawker found %v answers", len(rrw.Msg.Answer))
	for _, answer := range rrw.Msg.Answer {
		switch record := answer.(type) {
		case *dns.A:
			log.Infof("Found A record: %v", answer)
			allowed_ips = append(allowed_ips, record.A)
		case *dns.AAAA:
			log.Infof("Found AAAA record: %v", answer)
			allowed_ips = append(allowed_ips, record.AAAA)
		}
		ttl := answer.Header().Ttl
		if ttl > max_ttl {
			max_ttl = ttl
		}
	}

	log.Infof("Squawker writing policy for %v: %v for %v seconds", client, allowed_ips, max_ttl)

	// TODO: Write policy to server

	return rc, err
}

func (s Squawker) Name() string { return PLUGIN_NAME }
