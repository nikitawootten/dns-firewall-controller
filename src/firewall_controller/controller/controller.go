package controller

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/nikitawootten/dns-firewall-controller/src/firewall_controller/firewall"
	pb "github.com/nikitawootten/dns-firewall-controller/src/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FirewallController struct {
	pb.UnimplementedFirewallControllerServer
	firewall firewall.Firewall
}

func NewFirewallControllerServer(firewall firewall.Firewall) pb.FirewallControllerServer {
	return &FirewallController{
		firewall: firewall,
	}
}

func (s *FirewallController) writeClientPolicy(client net.IP, allowed []net.IP, allow_until time.Time) error {
	log.Printf("Adding rule for %v: %v\n", client, allowed)
	err := s.firewall.AddRule(client, allowed)
	if err != nil {
		return fmt.Errorf("failed to add rule: %w", err)
	}

	go func() {
		time.Sleep(time.Until(allow_until))
		log.Printf("Removing rule for %v\n", client)
		err := s.firewall.RemoveRule(client)
		if err != nil {
			log.Printf("failed to remove rule: %v", err)
		}
	}()

	return nil
}

func (s *FirewallController) WriteClientPolicy(ctx context.Context, policy *pb.ClientPolicy) (*emptypb.Empty, error) {
	client, err := pb.FromProtoIpAddress(policy.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize client IP: %w", err)
	}

	allowed_ips, err := pb.FromProtoIpAddresses(policy.AllowedIps)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize allowed IPs: %w", err)
	}

	allow_until := policy.AllowUntil.AsTime()

	log.Printf("WriteClientPolicy ip: %v, allowed: %v, until: %v", client, allowed_ips, allow_until)

	err = s.writeClientPolicy(client, allowed_ips, allow_until)
	if err != nil {
		return nil, fmt.Errorf("failed to write client policy: %w", err)
	}

	return &emptypb.Empty{}, nil
}
