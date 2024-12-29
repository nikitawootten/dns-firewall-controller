package firewall

import (
	"log"
)

type MockPollingFirewallBackend struct{}

func (m *MockPollingFirewallBackend) AddRule(rule FirewallRule) error {
	log.Printf("Rule added: %v", rule)
	return nil
}

func (m *MockPollingFirewallBackend) RemoveRule(rule FirewallRule) error {
	log.Printf("Rule removed: %v", rule)
	return nil
}

func NewMockBackend() FirewallBackend {
	return NewPollingFirewallAdapter(&MockPollingFirewallBackend{})
}
