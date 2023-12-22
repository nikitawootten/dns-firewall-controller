package firewall

import (
	"net"

	"github.com/coreos/go-iptables/iptables"
)

type IpTablesFirewall struct {
	ipt  *iptables.IPTables
	ip6t *iptables.IPTables
}

func NewIpTablesFirewall() (Firewall, error) {
	ipt, err := iptables.New()
	if err != nil {
		return nil, err
	}
	ip6t, err := iptables.NewWithProtocol(iptables.ProtocolIPv6)
	if err != nil {
		return nil, err
	}

	return &IpTablesFirewall{
		ipt:  ipt,
		ip6t: ip6t,
	}, nil
}

func (f *IpTablesFirewall) AddRule(client net.IP, allowed []net.IP) error {
	return nil
}

func (f *IpTablesFirewall) RemoveRule(client net.IP) error {
	return nil
}
