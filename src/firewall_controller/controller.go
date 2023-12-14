package firewall_controller

import (
	pb "github.com/nikitawootten/dns-firewall-controller/src/common"
)

type FirewallController struct {
	pb.UnimplementedFirewallControllerServer
}

func NewFirewallControllerServer() pb.FirewallControllerServer {
	return &FirewallController{}
}
