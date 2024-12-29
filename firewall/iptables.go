package firewall

import (
	"github.com/coreos/go-iptables/iptables"
)

type IPTablesFirewallBackend struct {
	Table    string
	Chain    string
	IPTables *iptables.IPTables
}

func (i IPTablesFirewallBackend) buildRuleSpec(rule FirewallRule) []string {
	return []string{"-s", rule.SourceIP, "-d", rule.DestinationIP, "-j", "ACCEPT"}
}

func (i *IPTablesFirewallBackend) AddRule(rule FirewallRule) error {
	ruleSpec := i.buildRuleSpec(rule)
	if err := i.IPTables.AppendUnique(i.Table, i.Chain, ruleSpec...); err != nil {
		return err
	}
	return nil
}

func (i *IPTablesFirewallBackend) RemoveRule(rule FirewallRule) error {
	ruleSpec := i.buildRuleSpec(rule)
	if err := i.IPTables.DeleteIfExists(i.Table, i.Chain, ruleSpec...); err != nil {
		return err
	}
	return nil
}

func NewIPTablesBackend(pollingBackend *IPTablesFirewallBackend) FirewallBackend {
	return NewPollingFirewallAdapter(pollingBackend)
}
