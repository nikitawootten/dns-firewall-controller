package firewall

import (
	"net"
)

// NoopFirewall is a test firewall that does not enforce any rules.
// The rules are stored in memory and can be inspected.
type NoopFirewall struct {
	rules map[string][]net.IP
}

func NewNaiveFirewall() Firewall {
	return &NoopFirewall{
		rules: make(map[string][]net.IP),
	}
}

func (f *NoopFirewall) AddRule(client net.IP, allowed []net.IP) error {

	return nil
}

func (f *NoopFirewall) RemoveRule(client net.IP) error {
	return nil
}
