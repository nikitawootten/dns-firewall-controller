package firewall

import (
	"time"

	"github.com/nikitawootten/dns-firewall-controller/dns"
)

// FirewallRule represents a single rule allowing traffic from SourceIP to DestinationIP until Expires.
type FirewallRule struct {
	SourceIP      string
	DestinationIP string
}

func (rule FirewallRule) String() string {
	return rule.SourceIP + " -> " + rule.DestinationIP
}

type FirewallRuleWithExpiration struct {
	FirewallRule
	Expires time.Time
}

func (rule FirewallRuleWithExpiration) String() string {
	return rule.FirewallRule.String() + " (expires " + rule.Expires.String() + ")"
}

func FirewallRulesFromDNSResponse(response *dns.DNSResponse) []FirewallRuleWithExpiration {
	now := time.Now()
	rules := make([]FirewallRuleWithExpiration, 0, len(response.Records))
	for _, record := range response.Records {
		expires := time.Duration(record.TTL) * time.Second
		rules = append(rules, FirewallRuleWithExpiration{
			FirewallRule: FirewallRule{
				SourceIP:      response.SourceAddress.String(),
				DestinationIP: record.RecordIP.String(),
			},
			Expires: now.Add(expires),
		})
	}
	return rules
}

type FirewallBackend interface {
	Start() error
	Stop() error
	AddRule(rule FirewallRuleWithExpiration) error
	ListActiveRules() ([]FirewallRuleWithExpiration, error)
}
