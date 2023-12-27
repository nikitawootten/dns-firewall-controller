package squawker

import (
	"context"
	"net"

	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/nikitawootten/dns-firewall-controller/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var PLUGIN_NAME = "squawker"

var log = clog.NewWithPlugin(PLUGIN_NAME)

type Squawker struct {
	client proto.FirewallControllerClient
	Next   plugin.Handler
}

type SquawkerConfig struct {
	address string
}

func NewSquawker(next plugin.Handler, config SquawkerConfig) plugin.Handler {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(config.address, opts...)
	if err != nil {
		plugin.Error(PLUGIN_NAME, err)
	}

	return &Squawker{
		Next:   next,
		client: proto.NewFirewallControllerClient(conn),
	}
}

func (s *Squawker) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	log.Info("Squawk!")

	state := request.Request{W: w, Req: r}
	client := net.ParseIP(state.IP())

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

	for_duration := time.Duration(max_ttl) * time.Second

	log.Infof("Squawker writing policy for %v: %v for %v seconds", client, allowed_ips, max_ttl)
	policy, err := proto.NewClientPolicy(client, allowed_ips, time.Now().Add(for_duration))
	_, err = s.client.WriteClientPolicy(ctx, policy)
	if err != nil {
		return rc, plugin.Error(PLUGIN_NAME, err)
	}

	return rc, err
}

func (s Squawker) Name() string { return PLUGIN_NAME }
