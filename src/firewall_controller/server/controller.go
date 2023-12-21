package firewall_controller

import (
	"context"

	"github.com/coredns/coredns/plugin/pkg/log"
	pb "github.com/nikitawootten/dns-firewall-controller/src/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FirewallController struct {
	pb.UnimplementedFirewallControllerServer
}

func NewFirewallControllerServer() pb.FirewallControllerServer {
	return &FirewallController{}
}

func (s *FirewallController) WriteClientPolicy(ctx context.Context, policy *pb.ClientPolicy) (*emptypb.Empty, error) {
	log.Info("WriteClientPolicy called:")

	client, _ := pb.FromProtoIpAddress(policy.Client)
	log.Infof("\tClient: %v", client)

	allowed_ips, _ := pb.FromProtoIpAddresses(policy.AllowedIps)
	log.Infof("\tAllowed IPs: %v", allowed_ips)

	allow_until := policy.AllowUntil.AsTime()
	log.Infof("\tAllow until: %v", allow_until)

	return &emptypb.Empty{}, nil
}
